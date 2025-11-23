-- Create practice_sessions table
CREATE TABLE IF NOT EXISTS practice_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    oa_id INTEGER NOT NULL REFERENCES objetivos_aprendizaje(id) ON DELETE CASCADE,
    oa_bloom_objective_id INTEGER NOT NULL REFERENCES oa_bloom_objectives(id) ON DELETE CASCADE,
    bloom_level_inicial INTEGER NOT NULL,
    bloom_level_final INTEGER,
    numero_preguntas INTEGER NOT NULL DEFAULT 10,
    preguntas_respondidas INTEGER NOT NULL DEFAULT 0,
    preguntas_correctas INTEGER NOT NULL DEFAULT 0,
    estado VARCHAR(20) NOT NULL DEFAULT 'en_progreso',
    estrategia JSONB,
    resultado JSONB,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for practice_sessions
CREATE INDEX idx_practice_sessions_user_id ON practice_sessions(user_id);
CREATE INDEX idx_practice_sessions_oa_id ON practice_sessions(oa_id);
CREATE INDEX idx_practice_sessions_estado ON practice_sessions(estado);

-- Create practice_answers table
CREATE TABLE IF NOT EXISTS practice_answers (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES practice_sessions(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    bloom_level_id INTEGER NOT NULL REFERENCES bloom_levels(id) ON DELETE CASCADE,
    user_answer JSONB,
    is_correct BOOLEAN,
    score DECIMAL(5,2),
    tiempo_segundos INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for practice_answers
CREATE INDEX idx_practice_answers_session_id ON practice_answers(session_id);
CREATE INDEX idx_practice_answers_question_id ON practice_answers(question_id);
