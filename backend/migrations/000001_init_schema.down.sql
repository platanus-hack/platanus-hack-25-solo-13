-- Rollback de la migración inicial

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
