package gcal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/queue"
	"github.com/google/uuid"
)

// NewCalendarTaskHandler returns a queue.TaskHandler that syncs appointment
// rows to Google Calendar (create / update / delete) using the latest DB state.
func NewCalendarTaskHandler(es EventSync, repo AppointmentCalendarRepository) queue.TaskHandler {
	return func(ctx context.Context, task *queue.Task) error {
		action, _ := task.Payload["action"].(string)
		switch action {
		case "delete":
			eventID, _ := task.Payload["event_id"].(string)
			if eventID == "" {
				return fmt.Errorf("missing event_id for delete action")
			}
			if err := es.DeleteEvent(ctx, eventID); err != nil {
				log.Printf("CALENDAR: delete failed for event %s: %v", eventID, err)
				return err
			}
			return nil
		case "upsert", "":
			idStr, _ := task.Payload["appointment_id"].(string)
			if idStr == "" {
				return fmt.Errorf("missing appointment_id for upsert action")
			}
			apptID, err := uuid.Parse(idStr)
			if err != nil {
				return fmt.Errorf("invalid appointment_id: %w", err)
			}
			return syncAppointmentCalendar(ctx, es, repo, apptID)
		default:
			log.Printf("CALENDAR: unknown action %q", action)
			return nil
		}
	}
}

func syncAppointmentCalendar(ctx context.Context, es EventSync, repo AppointmentCalendarRepository, apptID uuid.UUID) error {
	appt, err := repo.GetByIDWithRelations(ctx, apptID)
	if err != nil {
		return fmt.Errorf("load appointment: %w", err)
	}

	eventID := strings.TrimSpace(appt.GoogleCalendarEventID.String)

	// Cancelled: remove remote event and clear linkage.
	if appt.Status == domain.AppointmentStatusCancelled {
		if appt.GoogleCalendarEventID.Valid && eventID != "" {
			if err := es.DeleteEvent(ctx, eventID); err != nil {
				log.Printf("CALENDAR: delete event %s for cancelled appointment %s: %v", eventID, apptID, err)
				return err
			}
			appt.GoogleCalendarEventID = domain.NullableString{}
			if err := repo.Update(ctx, appt); err != nil {
				return fmt.Errorf("clear google_calendar_event_id: %w", err)
			}
		}
		return nil
	}

	// Pending: do not publish; strip orphan linkage if any.
	if appt.Status == domain.AppointmentStatusPending {
		if appt.GoogleCalendarEventID.Valid && eventID != "" {
			if err := es.DeleteEvent(ctx, eventID); err != nil {
				log.Printf("CALENDAR: delete orphan event %s for pending appointment %s: %v", eventID, apptID, err)
				return err
			}
			appt.GoogleCalendarEventID = domain.NullableString{}
			if err := repo.Update(ctx, appt); err != nil {
				return fmt.Errorf("clear google_calendar_event_id for pending: %w", err)
			}
		}
		return nil
	}

	// Confirmed, rescheduled, completed: ensure calendar reflects current times/title.
	if appt.Status != domain.AppointmentStatusConfirmed &&
		appt.Status != domain.AppointmentStatusRescheduled &&
		appt.Status != domain.AppointmentStatusCompleted {
		log.Printf("CALENDAR: skip unsupported status %q for appointment %s", appt.Status, apptID)
		return nil
	}

	if appt.GoogleCalendarEventID.Valid && eventID != "" {
		if err := es.UpdateEvent(ctx, eventID, appt); err != nil {
			log.Printf("CALENDAR: update event %s: %v", eventID, err)
			return err
		}
		return nil
	}

	newID, err := es.CreateEvent(ctx, appt)
	if err != nil {
		log.Printf("CALENDAR: create event for appointment %s: %v", apptID, err)
		return err
	}
	appt.GoogleCalendarEventID = domain.NullableString{
		NullString: sql.NullString{String: newID, Valid: true},
	}
	if err := repo.Update(ctx, appt); err != nil {
		return fmt.Errorf("persist google_calendar_event_id: %w", err)
	}
	return nil
}
