-- Create profile_history table for tracking profile changes over time
CREATE TABLE IF NOT EXISTS profile_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Profile snapshot at that moment
    snapshot JSONB NOT NULL,

    -- Event metadata
    evento VARCHAR(100),  -- "diagnostico_inicial", "post_leccion", "actualizacion_manual", "post_evaluacion"

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index to query user history
CREATE INDEX idx_profile_history_user_id ON profile_history(user_id);
CREATE INDEX idx_profile_history_created_at ON profile_history(created_at DESC);

-- Comments for documentation
COMMENT ON TABLE profile_history IS 'History of changes in student profiles';
COMMENT ON COLUMN profile_history.snapshot IS 'Complete profile snapshot at that moment';
COMMENT ON COLUMN profile_history.evento IS 'Type of event that generated the snapshot';
