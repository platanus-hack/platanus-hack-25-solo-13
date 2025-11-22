-- MigraciÃ³n 17: Seed de Cursos, Materias y Relaciones (Datos de ProducciÃ³n)
-- DescripciÃ³n: Carga inicial de estructura educativa para Primero Medio
-- Fecha: 2025-11-22
-- Datos: 1 curso, 2 materias, 1 relaciÃ³n curso-materia

-- ============================================================================
-- CURSOS
-- ============================================================================

INSERT INTO cursos (id, nombre, codigo, nivel_educativo, descripcion, activo, created_at, updated_at)
VALUES (
    1,
    'Primero Medio',
    '1M',
    'EnseÃ±anza Media',
    'Primer aÃ±o de educaciÃ³n media',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    nombre = EXCLUDED.nombre,
    codigo = EXCLUDED.codigo,
    nivel_educativo = EXCLUDED.nivel_educativo,
    descripcion = EXCLUDED.descripcion,
    activo = EXCLUDED.activo,
    updated_at = NOW();

-- ============================================================================
-- MATERIAS
-- ============================================================================

INSERT INTO materias (id, nombre, codigo, descripcion, color, activo, created_at, updated_at)
VALUES
    (
        1,
        'MatemÃ¡ticas',
        'MAT',
        'MatemÃ¡ticas y resoluciÃ³n de problemas',
        '#EF4444',
        true,
        NOW(),
        NOW()
    ),
    (
        2,
        'Lengua y Literatura',
        'LYL',
        'Lenguaje, comunicaciÃ³n y literatura',
        '#8B5CF6',
        true,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    nombre = EXCLUDED.nombre,
    codigo = EXCLUDED.codigo,
    descripcion = EXCLUDED.descripcion,
    color = EXCLUDED.color,
    activo = EXCLUDED.activo,
    updated_at = NOW();

-- ============================================================================
-- RELACIONES CURSO-MATERIA
-- ============================================================================

-- Solo Lengua y Literatura tiene OAs cargados (ver migraciÃ³n 18)
-- MatemÃ¡ticas estÃ¡ disponible como materia pero sin contenido aÃºn

INSERT INTO curso_materias (id, curso_id, materia_id, horas_semanales, created_at)
VALUES (
    1,
    1,  -- Primero Medio
    2,  -- Lengua y Literatura
    6,  -- Horas semanales
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    curso_id = EXCLUDED.curso_id,
    materia_id = EXCLUDED.materia_id,
    horas_semanales = EXCLUDED.horas_semanales;

-- ============================================================================
-- VERIFICACIÃN
-- ============================================================================

-- Verificar que se crearon correctamente
DO $$
DECLARE
    cursos_count INTEGER;
    materias_count INTEGER;
    relaciones_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO cursos_count FROM cursos WHERE id = 1;
    SELECT COUNT(*) INTO materias_count FROM materias WHERE id IN (1, 2);
    SELECT COUNT(*) INTO relaciones_count FROM curso_materias WHERE id = 1;

    IF cursos_count != 1 OR materias_count != 2 OR relaciones_count != 1 THEN
        RAISE EXCEPTION 'Error en verificaciÃ³n de datos seed: cursos=%, materias=%, relaciones=%',
            cursos_count, materias_count, relaciones_count;
    END IF;

    RAISE NOTICE 'Seed completado: 1 curso, 2 materias, 1 relaciÃ³n';
END $$;
