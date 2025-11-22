-- ============================================================================
-- DIAGNOSTIC SESSIONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS diagnostic_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    materia_id INTEGER NOT NULL REFERENCES materias(id) ON DELETE CASCADE,

    numero_intento INTEGER NOT NULL DEFAULT 1,
    estado VARCHAR(20) NOT NULL DEFAULT 'en_progreso'
        CHECK (estado IN ('en_progreso', 'completado', 'abandonado')),

    -- Adaptive strategy in JSONB for flexibility
    estrategia JSONB NOT NULL DEFAULT '{
        "nivel_bloom_actual": 2,
        "oas_evaluados": [],
        "aciertos_consecutivos": 0,
        "fallos_consecutivos": 0,
        "patron_respuestas": []
    }'::jsonb,

    preguntas_totales INTEGER DEFAULT 0,
    preguntas_correctas INTEGER DEFAULT 0,

    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP
);

CREATE INDEX idx_diagnostic_sessions_user ON diagnostic_sessions(user_id);
CREATE INDEX idx_diagnostic_sessions_materia ON diagnostic_sessions(materia_id);
CREATE INDEX idx_diagnostic_sessions_estado ON diagnostic_sessions(estado);
CREATE INDEX idx_diagnostic_sessions_user_materia ON diagnostic_sessions(user_id, materia_id);

COMMENT ON TABLE diagnostic_sessions IS 'Adaptive diagnostic sessions for initial student assessment';
COMMENT ON COLUMN diagnostic_sessions.estrategia IS 'Adaptive algorithm state in JSON format';
COMMENT ON COLUMN diagnostic_sessions.numero_intento IS 'Diagnostic attempt number (allows re-diagnostics)';

-- ============================================================================
-- DIAGNOSTIC ANSWERS
-- ============================================================================

CREATE TABLE IF NOT EXISTS diagnostic_answers (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES diagnostic_sessions(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,

    -- Denormalized for fast queries
    oa_bloom_objective_id INTEGER NOT NULL,
    bloom_level_id INTEGER NOT NULL,

    -- Flexible answer in JSONB (structure varies by question type)
    user_answer JSONB NOT NULL,

    -- Validation results
    is_correct BOOLEAN,
    score DECIMAL(5,2),  -- 0-100 for partial scoring

    tiempo_segundos INTEGER,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_diagnostic_answers_session ON diagnostic_answers(session_id);
CREATE INDEX idx_diagnostic_answers_question ON diagnostic_answers(question_id);
CREATE INDEX idx_diagnostic_answers_oa_bloom ON diagnostic_answers(oa_bloom_objective_id);

COMMENT ON TABLE diagnostic_answers IS 'Individual answers during diagnostic sessions';
COMMENT ON COLUMN diagnostic_answers.user_answer IS 'Student answer in JSON (structure varies by question type)';
COMMENT ON COLUMN diagnostic_answers.score IS 'Partial score for questions with gradual scoring';

-- ============================================================================
-- DIAGNOSTIC RESULTS
-- ============================================================================

CREATE TABLE IF NOT EXISTS diagnostic_results (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES diagnostic_sessions(id) ON DELETE CASCADE,
    oa_id INTEGER NOT NULL REFERENCES objetivos_aprendizaje(id) ON DELETE CASCADE,

    nivel_bloom_dominado INTEGER NOT NULL CHECK (nivel_bloom_dominado BETWEEN 1 AND 6),
    nivel_bloom_nombre VARCHAR(20) NOT NULL,

    preguntas_respondidas INTEGER NOT NULL DEFAULT 0,
    preguntas_correctas INTEGER NOT NULL DEFAULT 0,
    porcentaje_aciertos INTEGER NOT NULL CHECK (porcentaje_aciertos BETWEEN 0 AND 100),

    recomendacion TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(session_id, oa_id)
);

CREATE INDEX idx_diagnostic_results_session ON diagnostic_results(session_id);
CREATE INDEX idx_diagnostic_results_oa ON diagnostic_results(oa_id);

COMMENT ON TABLE diagnostic_results IS 'Consolidated results per OA after completing diagnostic';
COMMENT ON COLUMN diagnostic_results.nivel_bloom_dominado IS 'Highest Bloom level mastered by student (1-6)';
COMMENT ON COLUMN diagnostic_results.recomendacion IS 'Personalized recommendation for learning path';
