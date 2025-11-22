-- Create student_profiles table with JSONB for flexible profile data
CREATE TABLE IF NOT EXISTS student_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

    -- Basic metadata (relational for fast queries)
    edad INTEGER,
    curso_actual VARCHAR(50),

    -- Complete profile in JSONB (flexible, ML-ready)
    profile_data JSONB NOT NULL DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_student_profiles_user_id ON student_profiles(user_id);
CREATE INDEX idx_student_profiles_profile_data ON student_profiles USING GIN(profile_data);

-- Comments for documentation
COMMENT ON TABLE student_profiles IS 'Student profiles for adaptive learning';
COMMENT ON COLUMN student_profiles.profile_data IS 'Profile data in JSON: prior knowledge, cognitive profile, preferences, motivation, self-efficacy, autonomy, interests';
