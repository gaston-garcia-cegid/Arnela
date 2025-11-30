-- Drop index
DROP INDEX IF EXISTS idx_appointments_room_time;

-- Drop room column
ALTER TABLE appointments DROP COLUMN IF EXISTS room;

-- Drop enum type
DROP TYPE IF EXISTS room_type;
