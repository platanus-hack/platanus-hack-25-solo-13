-- Add authentication fields to users table
ALTER TABLE users
ADD COLUMN password_hash VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';

-- Create index on role for faster role-based queries
CREATE INDEX idx_users_role ON users(role);
