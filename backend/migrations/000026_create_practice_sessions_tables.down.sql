-- Drop practice tables
DROP INDEX IF EXISTS idx_practice_answers_question_id;
DROP INDEX IF EXISTS idx_practice_answers_session_id;
DROP TABLE IF EXISTS practice_answers;

DROP INDEX IF EXISTS idx_practice_sessions_estado;
DROP INDEX IF EXISTS idx_practice_sessions_oa_id;
DROP INDEX IF EXISTS idx_practice_sessions_user_id;
DROP TABLE IF EXISTS practice_sessions;
