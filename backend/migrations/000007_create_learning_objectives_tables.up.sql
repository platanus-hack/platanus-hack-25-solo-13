-- Create objetivos_aprendizaje table (Learning Objectives)
CREATE TABLE IF NOT EXISTS objetivos_aprendizaje (
    id SERIAL PRIMARY KEY,
    materia_id INTEGER NOT NULL REFERENCES materias(id) ON DELETE CASCADE,

    codigo VARCHAR(20) NOT NULL,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT,

    orden INTEGER,
    activo BOOLEAN DEFAULT true,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(materia_id, codigo)
);

CREATE INDEX idx_oa_materia ON objetivos_aprendizaje(materia_id);

-- Create oa_bloom_objectives table (Objectives broken down by Bloom level)
CREATE TABLE IF NOT EXISTS oa_bloom_objectives (
    id SERIAL PRIMARY KEY,
    oa_id INTEGER NOT NULL REFERENCES objetivos_aprendizaje(id) ON DELETE CASCADE,
    bloom_level_id INTEGER NOT NULL REFERENCES bloom_levels(id) ON DELETE CASCADE,

    objetivo_especifico TEXT NOT NULL,
    indicadores_logro TEXT[],
    tipo_actividad_sugerida VARCHAR(50),
    complejidad_estimada INTEGER,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(oa_id, bloom_level_id),
    CONSTRAINT complejidad_check CHECK (complejidad_estimada BETWEEN 1 AND 10)
);

CREATE INDEX idx_oa_bloom_oa ON oa_bloom_objectives(oa_id);
CREATE INDEX idx_oa_bloom_level ON oa_bloom_objectives(bloom_level_id);

-- Comments
COMMENT ON TABLE objetivos_aprendizaje IS 'Base learning objectives for each subject';
COMMENT ON TABLE oa_bloom_objectives IS 'Specific objectives for each OA broken down by Bloom taxonomy level';
COMMENT ON COLUMN oa_bloom_objectives.objetivo_especifico IS 'Specific objective for this Bloom level';
COMMENT ON COLUMN oa_bloom_objectives.indicadores_logro IS 'Success criteria for this objective';
