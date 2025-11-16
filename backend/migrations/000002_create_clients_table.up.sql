-- Create clients table
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    dni VARCHAR(20) NOT NULL,
    date_of_birth DATE,
    address TEXT,
    city VARCHAR(100),
    postal_code VARCHAR(10),
    province VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_visit TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_clients_email ON clients(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_dni ON clients(dni) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_phone ON clients(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_user_id ON clients(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_clients_last_name ON clients(last_name) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_is_active ON clients(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_clients_deleted_at ON clients(deleted_at) WHERE deleted_at IS NOT NULL;

-- Create unique constraint for email (excluding soft-deleted records)
CREATE UNIQUE INDEX idx_clients_email_unique ON clients(email) WHERE deleted_at IS NULL;

-- Create unique constraint for DNI (excluding soft-deleted records)
CREATE UNIQUE INDEX idx_clients_dni_unique ON clients(dni) WHERE deleted_at IS NULL;

-- Trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_clients_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_clients_updated_at
    BEFORE UPDATE ON clients
    FOR EACH ROW
    EXECUTE FUNCTION update_clients_updated_at();

-- Add comment to table
COMMENT ON TABLE clients IS 'Stores client information for the CRM system';
COMMENT ON COLUMN clients.user_id IS 'Optional reference to users table if client has login access';
COMMENT ON COLUMN clients.dni IS 'Spanish identification number (DNI/NIE/NIF)';
COMMENT ON COLUMN clients.deleted_at IS 'Soft delete timestamp';
