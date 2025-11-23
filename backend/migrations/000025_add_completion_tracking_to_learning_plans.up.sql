-- Add completion tracking fields to learning_plans
-- These fields allow tracking user progress in the plan

-- Add tracking fields
ALTER TABLE learning_plans
ADD COLUMN completado BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN fecha_inicio TIMESTAMP,
ADD COLUMN fecha_completado TIMESTAMP,
ADD COLUMN progreso_actual INTEGER NOT NULL DEFAULT 0,
ADD COLUMN total_slides INTEGER NOT NULL DEFAULT 0;

-- Comments to document the fields
COMMENT ON COLUMN learning_plans.completado IS 'Indicates if the user completed the entire learning plan';
COMMENT ON COLUMN learning_plans.fecha_inicio IS 'Date and time when the user started the plan';
COMMENT ON COLUMN learning_plans.fecha_completado IS 'Date and time when the user completed the plan';
COMMENT ON COLUMN learning_plans.progreso_actual IS 'Number of slides completed by the user';
COMMENT ON COLUMN learning_plans.total_slides IS 'Total number of slides in the plan';

-- Index for queries of completed plans by user
CREATE INDEX idx_learning_plans_completado ON learning_plans(user_id, completado);
CREATE INDEX idx_learning_plans_fecha_completado ON learning_plans(fecha_completado) WHERE fecha_completado IS NOT NULL;
