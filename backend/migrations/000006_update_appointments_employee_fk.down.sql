-- Rollback: Restore mock therapist_id column

-- Step 1: Add back therapist_id column
ALTER TABLE appointments 
ADD COLUMN therapist_id VARCHAR(100);

-- Step 2: Create index for therapist_id
CREATE INDEX idx_appointments_therapist_id ON appointments(therapist_id) WHERE deleted_at IS NULL;

-- Step 3: Drop employee_id column and its constraints
DROP INDEX IF EXISTS idx_appointments_employee_id;
ALTER TABLE appointments DROP CONSTRAINT IF EXISTS fk_appointments_employee;
ALTER TABLE appointments DROP COLUMN employee_id;

-- Restore comment
COMMENT ON COLUMN appointments.therapist_id IS 'Mock therapist ID until employees feature is implemented';
