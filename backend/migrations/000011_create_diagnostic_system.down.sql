-- Drop diagnostic tables in reverse order (foreign keys)

-- Drop diagnostic_results first
DROP INDEX IF EXISTS idx_diagnostic_results_oa;
DROP INDEX IF EXISTS idx_diagnostic_results_session;
DROP TABLE IF EXISTS diagnostic_results;

-- Drop diagnostic_answers
DROP INDEX IF EXISTS idx_diagnostic_answers_oa_bloom;
DROP INDEX IF EXISTS idx_diagnostic_answers_question;
DROP INDEX IF EXISTS idx_diagnostic_answers_session;
DROP TABLE IF EXISTS diagnostic_answers;

-- Drop diagnostic_sessions last
DROP INDEX IF EXISTS idx_diagnostic_sessions_user_materia;
DROP INDEX IF EXISTS idx_diagnostic_sessions_estado;
DROP INDEX IF EXISTS idx_diagnostic_sessions_materia;
DROP INDEX IF EXISTS idx_diagnostic_sessions_user;
DROP TABLE IF EXISTS diagnostic_sessions;
