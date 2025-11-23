-- Drop appointments table and related objects
DROP INDEX IF EXISTS idx_appointments_date_range;
DROP INDEX IF EXISTS idx_appointments_status;
DROP INDEX IF EXISTS idx_appointments_start_time;
DROP INDEX IF EXISTS idx_appointments_therapist_id;
DROP INDEX IF EXISTS idx_appointments_client_id;

DROP TABLE IF EXISTS appointments;

DROP TYPE IF EXISTS appointment_status;
