package gcal

import (
	"fmt"

	"github.com/gaston-garcia-cegid/arnela/backend/pkg/queue"
	"github.com/google/uuid"
)

// EnqueueAppointmentCalendarSync enqueues a worker task to upsert or remove the
// Google Calendar event for this appointment (worker loads latest row from DB).
func EnqueueAppointmentCalendarSync(wp *queue.WorkerPool, appointmentID uuid.UUID) error {
	if wp == nil {
		return fmt.Errorf("worker pool is nil")
	}
	return wp.EnqueueTask(queue.TaskTypeSyncCalendar, map[string]interface{}{
		"action":         "upsert",
		"appointment_id": appointmentID.String(),
	})
}
