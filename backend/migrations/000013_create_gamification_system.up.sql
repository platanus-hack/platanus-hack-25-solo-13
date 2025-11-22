-- Create user_gamification table
CREATE TABLE IF NOT EXISTS user_gamification (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    level INTEGER NOT NULL DEFAULT 1 CHECK (level >= 1),
    xp INTEGER NOT NULL DEFAULT 0 CHECK (xp >= 0),
    coins INTEGER NOT NULL DEFAULT 100 CHECK (coins >= 0),
    current_streak INTEGER NOT NULL DEFAULT 0 CHECK (current_streak >= 0),
    longest_streak INTEGER NOT NULL DEFAULT 0 CHECK (longest_streak >= 0),
    last_activity_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for leaderboard queries
CREATE INDEX idx_user_gamification_xp ON user_gamification(xp DESC);
CREATE INDEX idx_user_gamification_level ON user_gamification(level DESC);

-- Function to calculate level from XP
-- Formula: level = floor(sqrt(xp / 100)) + 1
-- Level 1: 0-99 XP
-- Level 2: 100-399 XP
-- Level 3: 400-899 XP
-- Level 4: 900-1599 XP
-- etc.
CREATE OR REPLACE FUNCTION calculate_level_from_xp(xp_amount INTEGER)
RETURNS INTEGER AS $$
BEGIN
    IF xp_amount < 0 THEN
        RETURN 1;
    END IF;
    RETURN FLOOR(SQRT(xp_amount::FLOAT / 100.0)) + 1;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Trigger to auto-update level when XP changes
CREATE OR REPLACE FUNCTION update_level_on_xp_change()
RETURNS TRIGGER AS $$
BEGIN
    NEW.level := calculate_level_from_xp(NEW.xp);
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_level
BEFORE INSERT OR UPDATE OF xp ON user_gamification
FOR EACH ROW
EXECUTE FUNCTION update_level_on_xp_change();

-- Trigger to update longest_streak when current_streak increases
CREATE OR REPLACE FUNCTION update_longest_streak()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.current_streak > NEW.longest_streak THEN
        NEW.longest_streak := NEW.current_streak;
    END IF;
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_longest_streak
BEFORE INSERT OR UPDATE OF current_streak ON user_gamification
FOR EACH ROW
EXECUTE FUNCTION update_longest_streak();

-- Seed: Create gamification record for all existing users
INSERT INTO user_gamification (user_id, level, xp, coins, current_streak, longest_streak, last_activity_date)
SELECT
    id,
    1,                    -- Starting level
    0,                    -- Starting XP
    100,                  -- Starting coins (welcome bonus)
    0,                    -- No streak yet
    0,                    -- No longest streak yet
    CURRENT_DATE          -- Set today as last activity
FROM users
WHERE id NOT IN (SELECT user_id FROM user_gamification);

COMMENT ON TABLE user_gamification IS 'Stores gamification stats per user: XP, level, coins, and streak tracking';
COMMENT ON FUNCTION calculate_level_from_xp IS 'Calculates user level based on XP using sqrt formula';
