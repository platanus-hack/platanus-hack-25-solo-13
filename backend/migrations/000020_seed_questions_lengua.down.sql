-- Rollback de Migración 20: Eliminar Preguntas de Lengua
-- ADVERTENCIA: Esto eliminará 672 preguntas del banco

BEGIN;

-- Eliminar preguntas asociadas a bloom objectives de Lengua
DELETE FROM questions
WHERE oa_bloom_objective_id IN (
    SELECT id FROM oa_bloom_objectives
    WHERE oa_id IN (SELECT id FROM objetivos_aprendizaje WHERE materia_id = 2)
);

-- Verificación
DO $$
DECLARE
    questions_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO questions_count
    FROM questions
    WHERE oa_bloom_objective_id IN (
        SELECT id FROM oa_bloom_objectives
        WHERE oa_id IN (SELECT id FROM objetivos_aprendizaje WHERE materia_id = 2)
    );

    IF questions_count != 0 THEN
        RAISE EXCEPTION 'Error: Preguntas no fueron eliminadas (quedan %)', questions_count;
    END IF;

    RAISE NOTICE 'Rollback completado: 672 preguntas eliminadas';
END $$;

COMMIT;
