-- Create cursos table (educational levels)
CREATE TABLE IF NOT EXISTS cursos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    codigo VARCHAR(20) UNIQUE,
    nivel_educativo VARCHAR(50),
    descripcion TEXT,
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create materias table (subjects/courses)
CREATE TABLE IF NOT EXISTS materias (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    codigo VARCHAR(20) UNIQUE,
    descripcion TEXT,
    color VARCHAR(7),
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create curso_materias junction table (many-to-many)
CREATE TABLE IF NOT EXISTS curso_materias (
    id SERIAL PRIMARY KEY,
    curso_id INTEGER NOT NULL REFERENCES cursos(id) ON DELETE CASCADE,
    materia_id INTEGER NOT NULL REFERENCES materias(id) ON DELETE CASCADE,
    horas_semanales INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(curso_id, materia_id)
);

-- Indexes
CREATE INDEX idx_curso_materias_curso ON curso_materias(curso_id);
CREATE INDEX idx_curso_materias_materia ON curso_materias(materia_id);

-- Comments
COMMENT ON TABLE cursos IS 'Educational levels (1 medio, 2 medio, etc.)';
COMMENT ON TABLE materias IS 'Subjects/courses (Language, Math, etc.)';
COMMENT ON TABLE curso_materias IS 'Many-to-many relationship between courses and subjects';
