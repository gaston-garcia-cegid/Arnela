-- Create employees table
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    dni VARCHAR(20) UNIQUE NOT NULL,
    date_of_birth DATE,
    position VARCHAR(100), -- e.g., "Psicólogo", "Terapeuta", "Coordinador"
    specialties TEXT[], -- Array of specializations
    is_active BOOLEAN DEFAULT true,
    hire_date DATE,
    notes TEXT,
    avatar_color VARCHAR(7) DEFAULT '#6366F1', -- For UI display
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_employees_user_id ON employees(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_email ON employees(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_dni ON employees(dni) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_is_active ON employees(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_position ON employees(position) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_deleted_at ON employees(deleted_at);

-- Auto-update trigger for updated_at
CREATE TRIGGER update_employees_updated_at 
BEFORE UPDATE ON employees
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at_column();
