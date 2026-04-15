package service

import "github.com/gaston-garcia-cegid/arnela/backend/internal/domain"

// NotificationService enqueues user-facing notifications.
type NotificationService interface {
	NotifyAppointmentConfirmed(appointment *domain.Appointment, clientEmail string) error
	NotifyAppointmentCancelled(appointment *domain.Appointment, clientEmail, reason string) error
}
