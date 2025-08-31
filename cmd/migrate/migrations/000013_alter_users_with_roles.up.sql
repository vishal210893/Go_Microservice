-- Create roles table first
CREATE TABLE IF NOT EXISTS roles
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default roles
INSERT INTO roles (name)
VALUES ('user'),
       ('admin')
ON CONFLICT (name) DO NOTHING;

-- Add role_id column to users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role_id INT REFERENCES roles (id) DEFAULT 1;

-- Update existing users to have 'user' role
UPDATE users
SET role_id = (SELECT id FROM roles WHERE name = 'user')
WHERE role_id IS NULL;

-- Remove default and set NOT NULL constraint
ALTER TABLE users
    ALTER COLUMN role_id DROP DEFAULT;

ALTER TABLE users
    ALTER COLUMN role_id SET NOT NULL;
