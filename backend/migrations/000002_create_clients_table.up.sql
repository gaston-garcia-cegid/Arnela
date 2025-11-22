-- Create clients table
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    dni VARCHAR(20) UNIQUE,
    nif VARCHAR(20) UNIQUE,
    address_street VARCHAR(255),
    address_city VARCHAR(100),
    address_province VARCHAR(100),
    address_postal_code VARCHAR(10),
    address_country VARCHAR(100),
    notes TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_clients_user_id ON clients(user_id);
CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email);
CREATE INDEX IF NOT EXISTS idx_clients_dni ON clients(dni);
CREATE INDEX IF NOT EXISTS idx_clients_nif ON clients(nif);
CREATE INDEX IF NOT EXISTS idx_clients_is_active ON clients(is_active);

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
