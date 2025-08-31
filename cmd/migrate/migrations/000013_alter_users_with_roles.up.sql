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
