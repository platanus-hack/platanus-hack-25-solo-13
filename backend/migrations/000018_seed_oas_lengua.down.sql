-- Rollback de Migración 18: Eliminar OAs de Lengua y Literatura
-- ADVERTENCIA: Esto eliminará en cascada bloom objectives y preguntas asociadas

BEGIN;

-- Eliminar todos los OAs de Lengua y Literatura (materia_id = 2)
-- Cascadeará a oa_bloom_objectives y questions
DELETE FROM objetivos_aprendizaje WHERE materia_id = 2;

-- Verificación
DO $$
DECLARE
    oas_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO oas_count FROM objetivos_aprendizaje WHERE materia_id = 2;

    IF oas_count != 0 THEN
        RAISE EXCEPTION 'Error: OAs de Lengua no fueron eliminados correctamente (quedan %)', oas_count;
    END IF;

    RAISE NOTICE 'Rollback completado: 23 OAs de Lengua eliminados';
END $$;

COMMIT;
