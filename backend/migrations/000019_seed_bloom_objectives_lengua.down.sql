-- Rollback de Migración 19: Eliminar Bloom Objectives de Lengua
-- ADVERTENCIA: Esto eliminará en cascada las preguntas asociadas

BEGIN;

-- Eliminar bloom objectives asociados a OAs de Lengua
-- Cascadeará a questions
DELETE FROM oa_bloom_objectives
WHERE oa_id IN (SELECT id FROM objetivos_aprendizaje WHERE materia_id = 2);

-- Verificación
DO $$
DECLARE
    bloom_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO bloom_count
    FROM oa_bloom_objectives
    WHERE oa_id IN (SELECT id FROM objetivos_aprendizaje WHERE materia_id = 2);

    IF bloom_count != 0 THEN
        RAISE EXCEPTION 'Error: Bloom objectives no fueron eliminados (quedan %)', bloom_count;
    END IF;

    RAISE NOTICE 'Rollback completado: 134 bloom objectives eliminados';
END $$;

COMMIT;
