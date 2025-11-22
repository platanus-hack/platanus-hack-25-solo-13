-- Script de Limpieza de Datos de Prueba
-- Ejecutar ANTES de crear migraciones de producción
-- Uso: docker compose exec postgres psql -U admin -d hackathon -f /path/to/cleanup_test_data.sql

-- ===========================================================================
-- IMPORTANTE: Este script elimina datos de prueba/desarrollo
-- NO ejecutar en producción si hay usuarios reales activos
-- ===========================================================================

BEGIN;

-- 1. Eliminar OA de Matemáticas de prueba (ID 1)
--    Esto eliminará en cascada:
--    - oa_bloom_objectives asociados (6 registros)
--    - questions asociadas (~30 registros)
--    - diagnostic_answers si existen
--    - student_oa_progress si existe
--    - diagnostic_results si existen
DELETE FROM objetivos_aprendizaje
WHERE id = 1
  AND codigo = 'OA02'
  AND materia_id = 1
  AND descripcion LIKE '%ecuaciones de primer grado%';

-- Verificar eliminación
SELECT 'OAs eliminados' as operacion, COUNT(*) as registros
FROM objetivos_aprendizaje WHERE materia_id = 1;

-- 2. Eliminar usuarios de prueba
--    Esto eliminará en cascada:
--    - student_profiles
--    - diagnostic_sessions
--    - diagnostic_answers
--    - student_oa_progress
--    - student_oa_history
--    - user_gamification
--    - user_inventory
--    - user_equipment
--    - unlock_notifications
DELETE FROM users
WHERE email IN (
  'alice@example.com',
  'bob@example.com',
  'charlie@example.com',
  'diana@example.com',
  'eve@example.com',
  'test@lumera.com',
  'teacher@lumera.com',
  'testquestions@lumera.com',
  'testgamif@lumera.com'
);

-- OPCIONAL: Descomentar si jon@than.cl también es de prueba
-- DELETE FROM users WHERE email = 'jon@than.cl';

-- Verificar usuarios restantes
SELECT 'Usuarios restantes' as operacion, COUNT(*) as registros FROM users;
SELECT email, name, created_at FROM users ORDER BY id;

-- 3. Limpiar relación curso-materia de Matemáticas si no tiene OAs
--    (Mantener solo la de Lengua que tiene contenido real)
DELETE FROM curso_materias
WHERE materia_id = 1
  AND NOT EXISTS (
    SELECT 1 FROM objetivos_aprendizaje WHERE materia_id = 1
  );

-- Verificar relaciones restantes
SELECT 'Relaciones curso-materia restantes' as operacion, COUNT(*) as registros
FROM curso_materias;

-- 4. Verificar integridad final
SELECT
  'Datos de Producción Finales' as resumen,
  (SELECT COUNT(*) FROM cursos) as cursos,
  (SELECT COUNT(*) FROM materias) as materias,
  (SELECT COUNT(*) FROM curso_materias) as relaciones,
  (SELECT COUNT(*) FROM objetivos_aprendizaje) as oas,
  (SELECT COUNT(*) FROM oa_bloom_objectives) as bloom_objectives,
  (SELECT COUNT(*) FROM questions) as preguntas,
  (SELECT COUNT(*) FROM users) as usuarios;

-- 5. Mostrar distribución final de OAs por materia
SELECT
  m.nombre as materia,
  COUNT(oa.id) as total_oas,
  COUNT(DISTINCT obo.id) as bloom_objectives,
  COUNT(DISTINCT q.id) as preguntas
FROM materias m
LEFT JOIN objetivos_aprendizaje oa ON oa.materia_id = m.id
LEFT JOIN oa_bloom_objectives obo ON obo.oa_id = oa.id
LEFT JOIN questions q ON q.oa_bloom_objective_id = obo.id
GROUP BY m.id, m.nombre
ORDER BY m.id;

COMMIT;

-- ===========================================================================
-- Resultado Esperado:
-- - 0 OAs de Matemáticas
-- - 23 OAs de Lengua y Literatura
-- - 138 bloom objectives
-- - 672 preguntas
-- - 0-2 usuarios (dependiendo si jon@than.cl se mantiene)
-- ===========================================================================
