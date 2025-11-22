-- Drop profile_history table
DROP INDEX IF EXISTS idx_profile_history_created_at;
DROP INDEX IF EXISTS idx_profile_history_user_id;
DROP TABLE IF EXISTS profile_history;
