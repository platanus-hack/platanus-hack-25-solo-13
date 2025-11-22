-- Create student_oa_progress table (Current progress on each OA-Bloom objective)
CREATE TABLE IF NOT EXISTS student_oa_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    oa_bloom_objective_id INTEGER NOT NULL REFERENCES oa_bloom_objectives(id) ON DELETE CASCADE,

    estado VARCHAR(20) NOT NULL,
    porcentaje_logro INTEGER DEFAULT 0,

    intentos INTEGER DEFAULT 0,
    ultima_actividad_fecha TIMESTAMP,
    notas TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, oa_bloom_objective_id),
    CONSTRAINT porcentaje_check CHECK (porcentaje_logro BETWEEN 0 AND 100),
    CONSTRAINT estado_check CHECK (estado IN ('no_iniciado', 'en_proceso', 'logrado', 'dominado'))
);

CREATE INDEX idx_student_progress_user ON student_oa_progress(user_id);
CREATE INDEX idx_student_progress_oa_bloom ON student_oa_progress(oa_bloom_objective_id);
CREATE INDEX idx_student_progress_estado ON student_oa_progress(estado);

-- Create student_oa_history table (Complete history of progress changes)
CREATE TABLE IF NOT EXISTS student_oa_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    oa_bloom_objective_id INTEGER NOT NULL REFERENCES oa_bloom_objectives(id) ON DELETE CASCADE,

    estado VARCHAR(20) NOT NULL,
    porcentaje_logro INTEGER,

    tipo_evento VARCHAR(50),
    puntaje_obtenido DECIMAL(5,2),
    puntaje_maximo DECIMAL(5,2),

    notas TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT history_estado_check CHECK (estado IN ('no_iniciado', 'en_proceso', 'logrado', 'dominado'))
);

CREATE INDEX idx_oa_history_user ON student_oa_history(user_id);
CREATE INDEX idx_oa_history_oa_bloom ON student_oa_history(oa_bloom_objective_id);
CREATE INDEX idx_oa_history_date ON student_oa_history(created_at DESC);
CREATE INDEX idx_oa_history_tipo_evento ON student_oa_history(tipo_evento);

-- Comments
COMMENT ON TABLE student_oa_progress IS 'Current progress of each student on OA-Bloom objectives';
COMMENT ON TABLE student_oa_history IS 'Complete history of all progress changes and evaluations';
COMMENT ON COLUMN student_oa_progress.estado IS 'no_iniciado | en_proceso | logrado | dominado';
COMMENT ON COLUMN student_oa_history.tipo_evento IS 'evaluacion | practica | diagnostico | repaso';
