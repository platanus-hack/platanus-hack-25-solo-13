-- Create flexible questions bank
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    oa_bloom_objective_id INTEGER NOT NULL
        REFERENCES oa_bloom_objectives(id) ON DELETE CASCADE,

    -- Question type with foreign key to registry
    tipo VARCHAR(30) NOT NULL REFERENCES question_types(tipo),

    -- Usage context
    tipo_uso VARCHAR(20) NOT NULL DEFAULT 'all'
        CHECK (tipo_uso IN ('diagnostico', 'practica', 'evaluacion', 'all')),

    -- Flexible content in JSONB (structure varies by tipo)
    question_data JSONB NOT NULL,
    validation_data JSONB NOT NULL,

    -- Basic validations (don't change with new types)
    CONSTRAINT question_data_is_object CHECK (
        jsonb_typeof(question_data) = 'object'
    ),
    CONSTRAINT validation_data_is_object CHECK (
        jsonb_typeof(validation_data) = 'object'
    ),

    -- Metadata
    dificultad_relativa INTEGER DEFAULT 3
        CHECK (dificultad_relativa BETWEEN 1 AND 5),
    veces_usada INTEGER DEFAULT 0,
    activa BOOLEAN DEFAULT true,
    tags TEXT[],

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_questions_oa_bloom ON questions(oa_bloom_objective_id);
CREATE INDEX idx_questions_tipo ON questions(tipo);
CREATE INDEX idx_questions_tipo_uso ON questions(tipo_uso);
CREATE INDEX idx_questions_tags ON questions USING GIN(tags);
CREATE INDEX idx_questions_question_data ON questions USING GIN(question_data);
CREATE INDEX idx_questions_activa ON questions(activa) WHERE activa = true;

-- Comments
COMMENT ON TABLE questions IS 'Flexible question bank supporting 9 activity types';
COMMENT ON COLUMN questions.tipo IS 'Question type from question_types registry';
COMMENT ON COLUMN questions.question_data IS 'Question content (structure varies by tipo)';
COMMENT ON COLUMN questions.validation_data IS 'Validation criteria (structure varies by tipo)';
COMMENT ON COLUMN questions.dificultad_relativa IS '1=very easy, 5=very hard (within same Bloom level)';
COMMENT ON COLUMN questions.tags IS 'Search tags for categorization';
