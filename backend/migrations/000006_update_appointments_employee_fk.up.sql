-- Migration: Update appointments to use real employees instead of mock therapists
-- Changes therapist_id from VARCHAR(100) to employee_id UUID with FK constraint

-- Step 1: Add new employee_id column (nullable for migration)
ALTER TABLE appointments 
ADD COLUMN employee_id UUID;

-- Step 2: Add foreign key constraint to employees table
ALTER TABLE appointments
ADD CONSTRAINT fk_appointments_employee
FOREIGN KEY (employee_id) REFERENCES employees(id)
ON DELETE RESTRICT ON UPDATE CASCADE;

-- Step 3: Create index for employee_id
CREATE INDEX idx_appointments_employee_id ON appointments(employee_id) WHERE deleted_at IS NULL;

-- Step 4: Drop old therapist_id column and its index
DROP INDEX IF EXISTS idx_appointments_therapist_id;
ALTER TABLE appointments DROP COLUMN therapist_id;

-- Update comment
COMMENT ON COLUMN appointments.employee_id IS 'Reference to employee (therapist) assigned to this appointment';
