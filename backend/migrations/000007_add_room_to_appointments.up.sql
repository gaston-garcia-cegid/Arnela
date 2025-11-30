-- Add room column to appointments table
CREATE TYPE room_type AS ENUM ('gabinete_01', 'gabinete_02', 'gabinete_externo');

ALTER TABLE appointments 
ADD COLUMN room room_type NOT NULL DEFAULT 'gabinete_01';

-- Create index for room availability queries
CREATE INDEX idx_appointments_room_time ON appointments(room, start_time, end_time) WHERE deleted_at IS NULL;

COMMENT ON COLUMN appointments.room IS 'Room/office where the appointment takes place';
COMMENT ON TYPE room_type IS 'Available rooms: Gabinete 01, Gabinete 02, or External office';
