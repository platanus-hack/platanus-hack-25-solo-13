-- Remove authentication fields from users table
DROP INDEX IF EXISTS idx_users_role;

ALTER TABLE users
DROP COLUMN IF EXISTS role,
DROP COLUMN IF EXISTS password_hash;
