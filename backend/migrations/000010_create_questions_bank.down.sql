-- Drop indexes
DROP INDEX IF EXISTS idx_questions_activa;
DROP INDEX IF EXISTS idx_questions_question_data;
DROP INDEX IF EXISTS idx_questions_tags;
DROP INDEX IF EXISTS idx_questions_tipo_uso;
DROP INDEX IF EXISTS idx_questions_tipo;
DROP INDEX IF EXISTS idx_questions_oa_bloom;

-- Drop table
DROP TABLE IF EXISTS questions;
