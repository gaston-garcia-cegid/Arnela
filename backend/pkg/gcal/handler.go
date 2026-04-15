package gcal

import (
	"context"
	"fmt"
	"log"

	"github.com/gaston-garcia-cegid/arnela/backend/pkg/queue"
)

// NewCalendarTaskHandler returns a queue.TaskHandler that syncs appointment
// events to Google Calendar.
func NewCalendarTaskHandler(cs *CalendarService) queue.TaskHandler {
	return func(ctx context.Context, task *queue.Task) error {
		action, _ := task.Payload["action"].(string)
		eventID, _ := task.Payload["event_id"].(string)

		switch action {
		case "delete":
			if eventID == "" {
				return fmt.Errorf("missing event_id for delete action")
			}
			if err := cs.DeleteEvent(ctx, eventID); err != nil {
				log.Printf("CALENDAR HANDLER: Failed to delete event %s: %v", eventID, err)
				return err
			}
		default:
			log.Printf("CALENDAR HANDLER: Unknown action %q", action)
		}

		return nil
	}
}
