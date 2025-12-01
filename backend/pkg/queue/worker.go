package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	QueueName       = "arnela:tasks"
	ProcessingQueue = "arnela:tasks:processing"
	DeadLetterQueue = "arnela:tasks:failed"
	MaxRetries      = 3
)

// TaskType defines the type of task to execute
type TaskType string

const (
	TaskTypeSendEmail    TaskType = "send_email"
	TaskTypeSendSMS      TaskType = "send_sms"
	TaskTypeSendWhatsApp TaskType = "send_whatsapp"
	TaskTypeSyncCalendar TaskType = "sync_calendar"
)

// Task represents a job to be processed
type Task struct {
	ID        string                 `json:"id"`
	Type      TaskType               `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Retries   int                    `json:"retries"`
	CreatedAt time.Time              `json:"created_at"`
}

// TaskHandler is a function that processes a task
type TaskHandler func(ctx context.Context, task *Task) error

// WorkerPool manages a pool of workers that process tasks from Redis queue
type WorkerPool struct {
	redisClient *redis.Client
	numWorkers  int
	handlers    map[TaskType]TaskHandler
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	stats       WorkerStats
}

// WorkerStats tracks worker pool statistics
type WorkerStats struct {
	TasksProcessed int64 `json:"tasks_processed"`
	TasksFailed    int64 `json:"tasks_failed"`
	ActiveWorkers  int   `json:"active_workers"`
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(redisClient *redis.Client, numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		redisClient: redisClient,
		numWorkers:  numWorkers,
		handlers:    make(map[TaskType]TaskHandler),
		ctx:         ctx,
		cancel:      cancel,
		stats: WorkerStats{
			ActiveWorkers: numWorkers,
		},
	}

	// Register default handlers
	pool.RegisterHandler(TaskTypeSendEmail, defaultEmailHandler)
	pool.RegisterHandler(TaskTypeSendSMS, defaultSMSHandler)
	pool.RegisterHandler(TaskTypeSendWhatsApp, defaultWhatsAppHandler)
	pool.RegisterHandler(TaskTypeSyncCalendar, defaultCalendarHandler)

	return pool
}

// RegisterHandler registers a task handler for a specific task type
func (wp *WorkerPool) RegisterHandler(taskType TaskType, handler TaskHandler) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.handlers[taskType] = handler
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	wp.cancel()
	wp.wg.Wait()
	log.Println("✓ Worker pool stopped")
}

// worker is the main worker loop
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		default:
			// Try to fetch a task (blocking pop with timeout)
			result, err := wp.redisClient.BRPopLPush(
				wp.ctx,
				QueueName,
				ProcessingQueue,
				5*time.Second,
			).Result()

			if err != nil {
				if err == redis.Nil {
					// No tasks available, continue
					continue
				}
				log.Printf("Worker %d: Error fetching task: %v", id, err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Process the task
			wp.processTask(id, result)
		}
	}
}

// processTask processes a single task
func (wp *WorkerPool) processTask(workerID int, taskJSON string) {
	var task Task
	if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
		log.Printf("Worker %d: Failed to unmarshal task: %v", workerID, err)
		wp.moveToDeadLetter(taskJSON, fmt.Sprintf("unmarshal error: %v", err))
		return
	}

	log.Printf("Worker %d: Processing task %s (type: %s, attempt: %d/%d)",
		workerID, task.ID, task.Type, task.Retries+1, MaxRetries)

	// Get handler for task type
	wp.mu.RLock()
	handler, exists := wp.handlers[task.Type]
	wp.mu.RUnlock()

	if !exists {
		log.Printf("Worker %d: No handler for task type %s", workerID, task.Type)
		wp.moveToDeadLetter(taskJSON, fmt.Sprintf("no handler for type: %s", task.Type))
		return
	}

	// Execute handler with timeout
	ctx, cancel := context.WithTimeout(wp.ctx, 30*time.Second)
	defer cancel()

	if err := handler(ctx, &task); err != nil {
		log.Printf("Worker %d: Task %s failed: %v", workerID, task.ID, err)
		wp.handleTaskFailure(task, taskJSON, err)
		return
	}

	// Task succeeded - remove from processing queue
	wp.redisClient.LRem(wp.ctx, ProcessingQueue, 1, taskJSON)

	wp.mu.Lock()
	wp.stats.TasksProcessed++
	wp.mu.Unlock()

	log.Printf("Worker %d: Task %s completed successfully", workerID, task.ID)
}

// handleTaskFailure handles a failed task (retry or move to DLQ)
func (wp *WorkerPool) handleTaskFailure(task Task, taskJSON string, err error) {
	task.Retries++

	if task.Retries >= MaxRetries {
		// Max retries exceeded - move to dead letter queue
		wp.moveToDeadLetter(taskJSON, fmt.Sprintf("max retries exceeded: %v", err))
		return
	}

	// Retry with exponential backoff
	backoff := time.Duration(task.Retries*task.Retries) * time.Second
	log.Printf("Retrying task %s in %v (attempt %d/%d)", task.ID, backoff, task.Retries+1, MaxRetries)

	// Update task and re-queue
	updatedTaskJSON, _ := json.Marshal(task)

	// Remove from processing queue
	wp.redisClient.LRem(wp.ctx, ProcessingQueue, 1, taskJSON)

	// Re-queue after backoff
	time.AfterFunc(backoff, func() {
		wp.redisClient.RPush(wp.ctx, QueueName, updatedTaskJSON)
	})
}

// moveToDeadLetter moves a task to the dead letter queue
func (wp *WorkerPool) moveToDeadLetter(taskJSON string, reason string) {
	log.Printf("Moving task to dead letter queue: %s", reason)

	// Add to dead letter queue with metadata
	dlqEntry := map[string]interface{}{
		"task":      taskJSON,
		"reason":    reason,
		"failed_at": time.Now().Format(time.RFC3339),
	}
	dlqJSON, _ := json.Marshal(dlqEntry)

	wp.redisClient.RPush(wp.ctx, DeadLetterQueue, dlqJSON)
	wp.redisClient.LRem(wp.ctx, ProcessingQueue, 1, taskJSON)

	wp.mu.Lock()
	wp.stats.TasksFailed++
	wp.mu.Unlock()
}

// EnqueueTask enqueues a new task to the queue
func (wp *WorkerPool) EnqueueTask(taskType TaskType, payload map[string]interface{}) error {
	task := Task{
		ID:        generateTaskID(),
		Type:      taskType,
		Payload:   payload,
		Retries:   0,
		CreatedAt: time.Now(),
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	if err := wp.redisClient.RPush(wp.ctx, QueueName, taskJSON).Err(); err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Task %s enqueued (type: %s)", task.ID, task.Type)
	return nil
}

// GetStats returns current worker pool statistics
func (wp *WorkerPool) GetStats() WorkerStats {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.stats
}

// generateTaskID generates a unique task ID
func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

// Default handlers (placeholders - implement actual logic)

func defaultEmailHandler(ctx context.Context, task *Task) error {
	log.Printf("EMAIL HANDLER: Sending email - %+v", task.Payload)
	// TODO: Implement actual email sending via SendGrid/SMTP
	time.Sleep(500 * time.Millisecond) // Simulate work
	return nil
}

func defaultSMSHandler(ctx context.Context, task *Task) error {
	log.Printf("SMS HANDLER: Sending SMS - %+v", task.Payload)
	// TODO: Implement actual SMS sending via Twilio
	time.Sleep(500 * time.Millisecond) // Simulate work
	return nil
}

func defaultWhatsAppHandler(ctx context.Context, task *Task) error {
	log.Printf("WHATSAPP HANDLER: Sending WhatsApp message - %+v", task.Payload)
	// TODO: Implement actual WhatsApp sending via Twilio
	time.Sleep(500 * time.Millisecond) // Simulate work
	return nil
}

func defaultCalendarHandler(ctx context.Context, task *Task) error {
	log.Printf("CALENDAR HANDLER: Syncing calendar event - %+v", task.Payload)
	// TODO: Implement actual Google Calendar sync
	time.Sleep(500 * time.Millisecond) // Simulate work
	return nil
}
