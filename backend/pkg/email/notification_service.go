package email

import (
	"log"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/queue"
)

// QueueNotificationService implements service.NotificationService
// by enqueueing emails into the Redis task queue.
type QueueNotificationService struct {
	wp *queue.WorkerPool
}

func NewQueueNotificationService(wp *queue.WorkerPool) *QueueNotificationService {
	return &QueueNotificationService{wp: wp}
}

func (s *QueueNotificationService) NotifyAppointmentConfirmed(appt *domain.Appointment, clientEmail string) error {
	data := appointmentToEmailData(appt)

	if err := EnqueueAppointmentConfirmation(s.wp, data, clientEmail); err != nil {
		log.Printf("Failed to enqueue confirmation email for %s: %v", clientEmail, err)
		return err
	}
	return nil
}

func (s *QueueNotificationService) NotifyAppointmentCancelled(appt *domain.Appointment, clientEmail, reason string) error {
	data := appointmentToEmailData(appt)
	data.Reason = reason

	if err := EnqueueAppointmentCancellation(s.wp, data, clientEmail); err != nil {
		log.Printf("Failed to enqueue cancellation email for %s: %v", clientEmail, err)
		return err
	}
	return nil
}

func appointmentToEmailData(appt *domain.Appointment) AppointmentEmailData {
	data := AppointmentEmailData{
		Title: appt.Title,
		Date:  appt.StartTime.Format("02/01/2006"),
		Time:  appt.StartTime.Format("15:04"),
		Room:  string(appt.Room),
	}

	if appt.Client != nil {
		data.ClientName = appt.Client.FirstName + " " + appt.Client.LastName
	}
	if appt.Employee != nil {
		data.EmployeeName = appt.Employee.FirstName + " " + appt.Employee.LastName
	}

	return data
}
