package gcal

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/queue"
	"github.com/google/uuid"
)

type fakeEventSync struct {
	created, updated, deleted int
	nextID                    string
}

func (f *fakeEventSync) CreateEvent(ctx context.Context, appt *domain.Appointment) (string, error) {
	f.created++
	if f.nextID == "" {
		return "gcal-new-id", nil
	}
	return f.nextID, nil
}

func (f *fakeEventSync) UpdateEvent(ctx context.Context, eventID string, appt *domain.Appointment) error {
	f.updated++
	return nil
}

func (f *fakeEventSync) DeleteEvent(ctx context.Context, eventID string) error {
	f.deleted++
	return nil
}

type memApptRepo struct {
	appt *domain.Appointment
}

func (m *memApptRepo) GetByIDWithRelations(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	_ = ctx
	_ = id
	return m.appt, nil
}

func (m *memApptRepo) Update(ctx context.Context, appointment *domain.Appointment) error {
	_ = ctx
	m.appt = appointment
	return nil
}

func TestSyncAppointmentCalendar_CreateWhenConfirmed(t *testing.T) {
	id := uuid.New()
	appt := &domain.Appointment{
		ID:          id,
		Status:      domain.AppointmentStatusConfirmed,
		Title:       "Sesión",
		StartTime:   time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
		DurationMinutes: 45,
		Room:        domain.RoomGabinete01,
	}
	f := &fakeEventSync{}
	repo := &memApptRepo{appt: appt}

	if err := syncAppointmentCalendar(context.Background(), f, repo, id); err != nil {
		t.Fatal(err)
	}
	if f.created != 1 || f.updated != 0 || f.deleted != 0 {
		t.Fatalf("unexpected calls: created=%d updated=%d deleted=%d", f.created, f.updated, f.deleted)
	}
	if !repo.appt.GoogleCalendarEventID.Valid || repo.appt.GoogleCalendarEventID.String == "" {
		t.Fatal("expected google calendar id persisted")
	}
}

func TestSyncAppointmentCalendar_UpdateWhenLinked(t *testing.T) {
	id := uuid.New()
	appt := &domain.Appointment{
		ID:     id,
		Status: domain.AppointmentStatusConfirmed,
		GoogleCalendarEventID: domain.NullableString{
			NullString: sql.NullString{String: "existing", Valid: true},
		},
		Title:             "Sesión",
		StartTime:         time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
		DurationMinutes:   45,
		Room:              domain.RoomGabinete01,
	}
	f := &fakeEventSync{}
	repo := &memApptRepo{appt: appt}

	if err := syncAppointmentCalendar(context.Background(), f, repo, id); err != nil {
		t.Fatal(err)
	}
	if f.created != 0 || f.updated != 1 || f.deleted != 0 {
		t.Fatalf("unexpected calls: created=%d updated=%d deleted=%d", f.created, f.updated, f.deleted)
	}
}

func TestSyncAppointmentCalendar_DeleteWhenCancelled(t *testing.T) {
	id := uuid.New()
	appt := &domain.Appointment{
		ID:     id,
		Status: domain.AppointmentStatusCancelled,
		GoogleCalendarEventID: domain.NullableString{
			NullString: sql.NullString{String: "evt-99", Valid: true},
		},
	}
	f := &fakeEventSync{}
	repo := &memApptRepo{appt: appt}

	if err := syncAppointmentCalendar(context.Background(), f, repo, id); err != nil {
		t.Fatal(err)
	}
	if f.deleted != 1 || f.created != 0 || f.updated != 0 {
		t.Fatalf("unexpected calls: created=%d updated=%d deleted=%d", f.created, f.updated, f.deleted)
	}
	if repo.appt.GoogleCalendarEventID.Valid {
		t.Fatal("expected google calendar id cleared")
	}
}

func TestCalendarTaskHandler_Upsert(t *testing.T) {
	id := uuid.New()
	appt := &domain.Appointment{
		ID:              id,
		Status:          domain.AppointmentStatusCompleted,
		Title:           "X",
		StartTime:       time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC),
		DurationMinutes: 60,
		Room:            domain.RoomGabinete02,
	}
	f := &fakeEventSync{}
	repo := &memApptRepo{appt: appt}
	h := NewCalendarTaskHandler(f, repo)

	err := h(context.Background(), &queue.Task{
		Payload: map[string]interface{}{
			"action":           "upsert",
			"appointment_id": id.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if f.created != 1 {
		t.Fatalf("expected one create, got %d", f.created)
	}
}
