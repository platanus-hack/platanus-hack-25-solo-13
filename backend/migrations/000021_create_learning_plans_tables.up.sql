-- Table: learning_plans
-- Stores personalized learning plans generated for users by learning objective
CREATE TABLE IF NOT EXISTS learning_plans (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    oa_bloom_objective_id INTEGER NOT NULL REFERENCES oa_bloom_objectives(id) ON DELETE CASCADE,
    titulo VARCHAR(500) NOT NULL,
    descripcion TEXT,
    tiempo_estimado_minutos INTEGER NOT NULL DEFAULT 0,
    estado VARCHAR(50) NOT NULL DEFAULT 'generando' CHECK (estado IN ('generando', 'generado', 'error')),
    error_mensaje TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Constraint: one plan per user per OA
    CONSTRAINT unique_user_oa UNIQUE (user_id, oa_bloom_objective_id)
);

CREATE INDEX idx_learning_plans_user_id ON learning_plans(user_id);
CREATE INDEX idx_learning_plans_oa_bloom_objective_id ON learning_plans(oa_bloom_objective_id);
CREATE INDEX idx_learning_plans_estado ON learning_plans(estado);

-- Table: learning_plan_components
-- Individual teaching components (slides) that make up a learning plan
CREATE TABLE IF NOT EXISTS learning_plan_components (
    id SERIAL PRIMARY KEY,
    learning_plan_id INTEGER NOT NULL REFERENCES learning_plans(id) ON DELETE CASCADE,
    orden INTEGER NOT NULL,
    tipo_componente VARCHAR(100) NOT NULL,
    objetivo_especifico TEXT NOT NULL,
    tiempo_estimado_minutos INTEGER NOT NULL DEFAULT 0,
    estado VARCHAR(50) NOT NULL DEFAULT 'pendiente' CHECK (estado IN ('pendiente', 'generando', 'generado', 'error')),
    contenido_props JSONB,
    error_mensaje TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Constraint: unique order within a plan
    CONSTRAINT unique_plan_orden UNIQUE (learning_plan_id, orden)
);

CREATE INDEX idx_learning_plan_components_plan_id ON learning_plan_components(learning_plan_id);
CREATE INDEX idx_learning_plan_components_estado ON learning_plan_components(estado);
CREATE INDEX idx_learning_plan_components_tipo ON learning_plan_components(tipo_componente);

-- Comments for documentation
COMMENT ON TABLE learning_plans IS 'Personalized learning plans generated for students by learning objective';
COMMENT ON COLUMN learning_plans.estado IS 'Status: generando (generating), generado (complete), error (failed)';
COMMENT ON COLUMN learning_plans.tiempo_estimado_minutos IS 'Estimated time to complete the entire learning plan';

COMMENT ON TABLE learning_plan_components IS 'Individual teaching slides/components that compose a learning plan';
COMMENT ON COLUMN learning_plan_components.tipo_componente IS 'Component type: GrammarConceptSlide, ReadingStrategySlide, etc.';
COMMENT ON COLUMN learning_plan_components.contenido_props IS 'JSON props required by the Svelte component';
COMMENT ON COLUMN learning_plan_components.estado IS 'Status: pendiente (not generated), generando (generating), generado (complete), error (failed)';
