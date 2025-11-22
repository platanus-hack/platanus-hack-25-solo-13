-- Drop tables in reverse order
DROP INDEX IF EXISTS idx_oa_history_tipo_evento;
DROP INDEX IF EXISTS idx_oa_history_date;
DROP INDEX IF EXISTS idx_oa_history_oa_bloom;
DROP INDEX IF EXISTS idx_oa_history_user;
DROP TABLE IF EXISTS student_oa_history;

DROP INDEX IF EXISTS idx_student_progress_estado;
DROP INDEX IF EXISTS idx_student_progress_oa_bloom;
DROP INDEX IF EXISTS idx_student_progress_user;
DROP TABLE IF EXISTS student_oa_progress;
