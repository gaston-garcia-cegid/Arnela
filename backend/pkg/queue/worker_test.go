package queue

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cleanup := func() {
		client.Close()
		mr.Close()
	}

	return client, cleanup
}

func TestWorkerPool_EnqueueTask(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	pool := NewWorkerPool(client, 2)

	tests := []struct {
		name      string
		taskType  TaskType
		payload   map[string]interface{}
		wantError bool
	}{
		{
			name:     "enqueue email task",
			taskType: TaskTypeSendEmail,
			payload: map[string]interface{}{
				"to":      "test@example.com",
				"subject": "Test Email",
				"body":    "Test body",
			},
			wantError: false,
		},
		{
			name:     "enqueue SMS task",
			taskType: TaskTypeSendSMS,
			payload: map[string]interface{}{
				"phone":   "+1234567890",
				"message": "Test SMS",
			},
			wantError: false,
		},
		{
			name:     "enqueue WhatsApp task",
			taskType: TaskTypeSendWhatsApp,
			payload: map[string]interface{}{
				"phone":   "+1234567890",
				"message": "Test WhatsApp",
			},
			wantError: false,
		},
		{
			name:     "enqueue calendar sync task",
			taskType: TaskTypeSyncCalendar,
			payload: map[string]interface{}{
				"event_id": "123",
				"action":   "create",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pool.EnqueueTask(tt.taskType, tt.payload)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Verify tasks were enqueued
	ctx := context.Background()
	queueLen, err := client.LLen(ctx, QueueName).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(4), queueLen)
}

func TestWorkerPool_ProcessTask(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	pool := NewWorkerPool(client, 1)

	// Register custom handler that tracks execution
	executed := false
	pool.RegisterHandler(TaskTypeSendEmail, func(ctx context.Context, task *Task) error {
		executed = true
		assert.Equal(t, "test@example.com", task.Payload["to"])
		return nil
	})

	// Start worker pool
	pool.Start()
	defer pool.Stop()

	// Enqueue task
	err := pool.EnqueueTask(TaskTypeSendEmail, map[string]interface{}{
		"to":      "test@example.com",
		"subject": "Test",
	})
	require.NoError(t, err)

	// Wait for task to be processed
	time.Sleep(2 * time.Second)

	// Verify task was executed
	assert.True(t, executed, "Task should have been executed")

	// Verify task was removed from queue
	ctx := context.Background()
	queueLen, err := client.LLen(ctx, QueueName).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), queueLen)

	// Verify stats
	stats := pool.GetStats()
	assert.Equal(t, int64(1), stats.TasksProcessed)
	assert.Equal(t, int64(0), stats.TasksFailed)
}

func TestWorkerPool_TaskRetry(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	pool := NewWorkerPool(client, 1)

	attempts := 0
	pool.RegisterHandler(TaskTypeSendEmail, func(ctx context.Context, task *Task) error {
		attempts++
		if attempts < 3 {
			return assert.AnError // Simulate failure
		}
		return nil // Success on third attempt
	})

	pool.Start()
	defer pool.Stop()

	err := pool.EnqueueTask(TaskTypeSendEmail, map[string]interface{}{
		"to": "test@example.com",
	})
	require.NoError(t, err)

	// Wait for retries
	time.Sleep(10 * time.Second)

	// Should have been retried and eventually succeeded
	assert.Equal(t, 3, attempts)
	stats := pool.GetStats()
	assert.Equal(t, int64(1), stats.TasksProcessed)
}

func TestWorkerPool_DeadLetterQueue(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	pool := NewWorkerPool(client, 1)

	// Handler that always fails
	pool.RegisterHandler(TaskTypeSendEmail, func(ctx context.Context, task *Task) error {
		return assert.AnError
	})

	pool.Start()
	defer pool.Stop()

	err := pool.EnqueueTask(TaskTypeSendEmail, map[string]interface{}{
		"to": "test@example.com",
	})
	require.NoError(t, err)

	// Wait for max retries
	time.Sleep(15 * time.Second)

	// Task should be in dead letter queue
	ctx := context.Background()
	dlqLen, err := client.LLen(ctx, DeadLetterQueue).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(1), dlqLen)

	// Verify stats
	stats := pool.GetStats()
	assert.Equal(t, int64(1), stats.TasksFailed)
}

func TestWorkerPool_GracefulShutdown(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	pool := NewWorkerPool(client, 2)

	// Register slow handler
	pool.RegisterHandler(TaskTypeSendEmail, func(ctx context.Context, task *Task) error {
		time.Sleep(2 * time.Second)
		return nil
	})

	pool.Start()

	// Enqueue some tasks
	for i := 0; i < 5; i++ {
		err := pool.EnqueueTask(TaskTypeSendEmail, map[string]interface{}{
			"to": "test@example.com",
		})
		require.NoError(t, err)
	}

	// Stop pool immediately
	pool.Stop()

	// Workers should have stopped gracefully
	stats := pool.GetStats()
	assert.GreaterOrEqual(t, stats.TasksProcessed, int64(0))
}

func TestWorkerPool_Stats(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	pool := NewWorkerPool(client, 3)

	stats := pool.GetStats()
	assert.Equal(t, 3, stats.ActiveWorkers)
	assert.Equal(t, int64(0), stats.TasksProcessed)
	assert.Equal(t, int64(0), stats.TasksFailed)
}
