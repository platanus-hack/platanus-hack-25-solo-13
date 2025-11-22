-- Drop triggers first
DROP TRIGGER IF EXISTS trigger_update_longest_streak ON user_gamification;
DROP TRIGGER IF EXISTS trigger_update_level ON user_gamification;

-- Drop functions
DROP FUNCTION IF EXISTS update_longest_streak();
DROP FUNCTION IF EXISTS update_level_on_xp_change();
DROP FUNCTION IF EXISTS calculate_level_from_xp(INTEGER);

-- Drop indexes
DROP INDEX IF EXISTS idx_user_gamification_level;
DROP INDEX IF EXISTS idx_user_gamification_xp;

-- Drop table
DROP TABLE IF EXISTS user_gamification;
