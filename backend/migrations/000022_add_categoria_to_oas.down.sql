-- Drop index
DROP INDEX IF EXISTS idx_objetivos_aprendizaje_categoria;

-- Drop categoria column from objetivos_aprendizaje table
ALTER TABLE objetivos_aprendizaje
DROP COLUMN IF EXISTS categoria;
