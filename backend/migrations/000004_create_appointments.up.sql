-- Create appointment_status enum type
CREATE TYPE appointment_status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed', 'rescheduled');

-- Create appointments table
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    therapist_id VARCHAR(100) NOT NULL, -- Mock ID (therapist-1, therapist-2, therapist-3)
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration_minutes INTEGER NOT NULL DEFAULT 60, -- 45 or 60 minutes
    status appointment_status NOT NULL DEFAULT 'pending',
    notes TEXT, -- Internal notes (admin only)
    cancellation_reason TEXT,
    google_calendar_event_id VARCHAR(255), -- For future Google Calendar integration
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_appointments_client_id ON appointments(client_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_appointments_therapist_id ON appointments(therapist_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_appointments_start_time ON appointments(start_time) WHERE deleted_at IS NULL;
CREATE INDEX idx_appointments_status ON appointments(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_appointments_date_range ON appointments(start_time, end_time) WHERE deleted_at IS NULL;

-- Constraint: end_time must be after start_time
ALTER TABLE appointments ADD CONSTRAINT check_appointment_time_order 
    CHECK (end_time > start_time);

-- Constraint: duration must be 45 or 60 minutes
ALTER TABLE appointments ADD CONSTRAINT check_appointment_duration 
    CHECK (duration_minutes IN (45, 60));

-- Grant permissions
GRANT ALL PRIVILEGES ON TABLE appointments TO arnela_user;

-- Comment for documentation
COMMENT ON TABLE appointments IS 'Stores all client appointments with therapists';
COMMENT ON COLUMN appointments.therapist_id IS 'Mock therapist ID until employees feature is implemented';
COMMENT ON COLUMN appointments.duration_minutes IS 'Appointment duration: 45 or 60 minutes';
COMMENT ON COLUMN appointments.notes IS 'Internal notes visible only to admin/employees';
