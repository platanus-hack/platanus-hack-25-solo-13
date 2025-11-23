-- Rollback: Remove completion tracking fields from learning_plans

-- Drop indexes
DROP INDEX IF EXISTS idx_learning_plans_fecha_completado;
DROP INDEX IF EXISTS idx_learning_plans_completado;

-- Drop columns
ALTER TABLE learning_plans
DROP COLUMN IF EXISTS total_slides,
DROP COLUMN IF EXISTS progreso_actual,
DROP COLUMN IF EXISTS fecha_completado,
DROP COLUMN IF EXISTS fecha_inicio,
DROP COLUMN IF EXISTS completado;
