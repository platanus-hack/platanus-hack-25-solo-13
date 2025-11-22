-- Rollback de Migración 17: Eliminar Cursos, Materias y Relaciones
-- ADVERTENCIA: Esto eliminará en cascada todos los OAs, bloom objectives y preguntas asociadas

BEGIN;

-- Eliminar relaciones curso-materia
DELETE FROM curso_materias WHERE id = 1;

-- Eliminar materias (cascadeará a objetivos_aprendizaje, oa_bloom_objectives, questions)
DELETE FROM materias WHERE id IN (1, 2);

-- Eliminar curso
DELETE FROM cursos WHERE id = 1;

-- Verificación
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM cursos WHERE id = 1) OR
       EXISTS (SELECT 1 FROM materias WHERE id IN (1, 2)) OR
       EXISTS (SELECT 1 FROM curso_materias WHERE id = 1) THEN
        RAISE EXCEPTION 'Error: No se eliminaron correctamente los registros';
    END IF;

    RAISE NOTICE 'Rollback completado: eliminados curso, materias y relaciones';
END $$;

COMMIT;
