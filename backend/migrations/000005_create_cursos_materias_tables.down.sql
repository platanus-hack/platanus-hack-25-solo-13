-- Drop tables in reverse order (respecting foreign keys)
DROP INDEX IF EXISTS idx_curso_materias_materia;
DROP INDEX IF EXISTS idx_curso_materias_curso;
DROP TABLE IF EXISTS curso_materias;
DROP TABLE IF EXISTS materias;
DROP TABLE IF EXISTS cursos;
