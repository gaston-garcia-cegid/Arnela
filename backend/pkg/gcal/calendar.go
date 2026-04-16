package gcal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type CalendarService struct {
	svc        *calendar.Service
	calendarID string
}

// NewCalendarService creates a Google Calendar service using a service account
// credentials JSON file. Set GOOGLE_CALENDAR_CREDENTIALS to the path of the
// JSON key file and GOOGLE_CALENDAR_ID to the target calendar.
func NewCalendarService() (*CalendarService, error) {
	credFile := os.Getenv("GOOGLE_CALENDAR_CREDENTIALS")
	calID := os.Getenv("GOOGLE_CALENDAR_ID")

	if credFile == "" || calID == "" {
		return nil, fmt.Errorf("GOOGLE_CALENDAR_CREDENTIALS and GOOGLE_CALENDAR_ID must be set")
	}

	data, err := os.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	conf, err := google.JWTConfigFromJSON(data, calendar.CalendarEventsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	ctx := context.Background()
	svc, err := calendar.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return &CalendarService{svc: svc, calendarID: calID}, nil
}

func IsConfigured() bool {
	return os.Getenv("GOOGLE_CALENDAR_CREDENTIALS") != "" && os.Getenv("GOOGLE_CALENDAR_ID") != ""
}

func (cs *CalendarService) CreateEvent(ctx context.Context, appt *domain.Appointment) (string, error) {
	event := appointmentToEvent(appt)

	created, err := cs.svc.Events.Insert(cs.calendarID, event).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create calendar event: %w", err)
	}

	log.Printf("Google Calendar event created: %s", created.Id)
	return created.Id, nil
}

func (cs *CalendarService) UpdateEvent(ctx context.Context, eventID string, appt *domain.Appointment) error {
	event := appointmentToEvent(appt)

	_, err := cs.svc.Events.Update(cs.calendarID, eventID, event).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to update calendar event: %w", err)
	}

	log.Printf("Google Calendar event updated: %s", eventID)
	return nil
}

func (cs *CalendarService) DeleteEvent(ctx context.Context, eventID string) error {
	err := cs.svc.Events.Delete(cs.calendarID, eventID).Context(ctx).Do()
	if err == nil {
		log.Printf("Google Calendar event deleted: %s", eventID)
		return nil
	}
	var gerr *googleapi.Error
	if errors.As(err, &gerr) && gerr.Code == 404 {
		log.Printf("Google Calendar event already absent: %s", eventID)
		return nil
	}
	return fmt.Errorf("failed to delete calendar event: %w", err)
}

func appointmentToEvent(appt *domain.Appointment) *calendar.Event {
	summary := appt.Title
	description := ""

	if appt.Client != nil {
		summary = fmt.Sprintf("%s - %s %s", appt.Title, appt.Client.FirstName, appt.Client.LastName)
		description = fmt.Sprintf("Cliente: %s %s\nEmail: %s\nTeléfono: %s",
			appt.Client.FirstName, appt.Client.LastName,
			appt.Client.Email, appt.Client.Phone)
	}

	if appt.Notes.Valid {
		description += "\n\nNotas: " + appt.Notes.String
	}

	location := ""
	if appt.Room != "" {
		location = fmt.Sprintf("Sala %s - Arnela Gabinete, C/ García Barbón 30, Vigo", string(appt.Room))
	}

	end := appt.StartTime.Add(time.Duration(appt.DurationMinutes) * time.Minute)

	return &calendar.Event{
		Summary:     summary,
		Description: description,
		Location:    location,
		Start: &calendar.EventDateTime{
			DateTime: appt.StartTime.Format(time.RFC3339),
			TimeZone: "Europe/Madrid",
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Europe/Madrid",
		},
		Reminders: &calendar.EventReminders{
			UseDefault:      false,
			ForceSendFields: []string{"UseDefault"},
			Overrides: []*calendar.EventReminder{
				{Method: "popup", Minutes: 30},
				{Method: "email", Minutes: 24 * 60}, // Google Calendar email reminder (Phase 7)
			},
		},
	}
}
