package email

import (
	"context"
	"fmt"
	"log"

	"github.com/gaston-garcia-cegid/arnela/backend/pkg/queue"
)

// NewEmailTaskHandler returns a queue.TaskHandler that sends emails via SMTP.
func NewEmailTaskHandler(mailer *Mailer) queue.TaskHandler {
	return func(ctx context.Context, task *queue.Task) error {
		to, _ := task.Payload["to"].(string)
		subject, _ := task.Payload["subject"].(string)
		body, _ := task.Payload["body"].(string)

		if to == "" || subject == "" || body == "" {
			return fmt.Errorf("missing required email fields (to=%q, subject=%q)", to, subject)
		}

		log.Printf("EMAIL HANDLER: Sending to=%s subject=%q", to, subject)

		if err := mailer.Send(to, subject, body); err != nil {
			log.Printf("EMAIL HANDLER: Failed to send to %s: %v", to, err)
			return err
		}

		log.Printf("EMAIL HANDLER: Successfully sent to %s", to)
		return nil
	}
}

// EnqueueAppointmentConfirmation enqueues an appointment confirmation email.
func EnqueueAppointmentConfirmation(wp *queue.WorkerPool, data AppointmentEmailData, recipientEmail string) error {
	body, err := RenderTemplate(AppointmentConfirmationTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render confirmation template: %w", err)
	}

	return wp.EnqueueTask(queue.TaskTypeSendEmail, map[string]interface{}{
		"to":      recipientEmail,
		"subject": fmt.Sprintf("Cita confirmada: %s - %s", data.Title, data.Date),
		"body":    body,
		"type":    "appointment_confirmation",
	})
}

// EnqueueAppointmentCancellation enqueues an appointment cancellation email.
func EnqueueAppointmentCancellation(wp *queue.WorkerPool, data AppointmentEmailData, recipientEmail string) error {
	body, err := RenderTemplate(AppointmentCancellationTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render cancellation template: %w", err)
	}

	return wp.EnqueueTask(queue.TaskTypeSendEmail, map[string]interface{}{
		"to":      recipientEmail,
		"subject": fmt.Sprintf("Cita cancelada: %s - %s", data.Title, data.Date),
		"body":    body,
		"type":    "appointment_cancellation",
	})
}
