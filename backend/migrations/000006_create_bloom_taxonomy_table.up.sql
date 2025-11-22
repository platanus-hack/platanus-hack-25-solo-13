-- Create bloom_levels table (Bloom's Taxonomy)
CREATE TABLE IF NOT EXISTS bloom_levels (
    id SERIAL PRIMARY KEY,
    nivel INTEGER NOT NULL UNIQUE,
    nombre VARCHAR(50) NOT NULL,
    nombre_en VARCHAR(50),
    descripcion TEXT,
    verbos_accion TEXT[],
    color VARCHAR(7),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT bloom_nivel_check CHECK (nivel BETWEEN 1 AND 6)
);

-- Seed data: The 6 levels of Bloom's Revised Taxonomy
INSERT INTO bloom_levels (nivel, nombre, nombre_en, descripcion, verbos_accion, color) VALUES
(1, 'Recordar', 'Remember',
 'Recuperar informacion relevante de la memoria de largo plazo',
 ARRAY['identificar', 'listar', 'nombrar', 'reconocer', 'recordar', 'seleccionar', 'definir'],
 '#EF4444'),

(2, 'Comprender', 'Understand',
 'Construir significado a partir de mensajes orales, escritos y graficos',
 ARRAY['explicar', 'interpretar', 'resumir', 'parafrasear', 'clasificar', 'comparar', 'ejemplificar'],
 '#F59E0B'),

(3, 'Aplicar', 'Apply',
 'Usar informacion en situaciones nuevas o concretas',
 ARRAY['aplicar', 'ejecutar', 'implementar', 'usar', 'demostrar', 'resolver', 'construir'],
 '#EAB308'),

(4, 'Analizar', 'Analyze',
 'Descomponer material en partes y determinar como las partes se relacionan',
 ARRAY['analizar', 'comparar', 'contrastar', 'diferenciar', 'organizar', 'distinguir', 'examinar'],
 '#22C55E'),

(5, 'Evaluar', 'Evaluate',
 'Hacer juicios basados en criterios y estandares',
 ARRAY['evaluar', 'juzgar', 'criticar', 'verificar', 'valorar', 'justificar', 'argumentar'],
 '#3B82F6'),

(6, 'Crear', 'Create',
 'Juntar elementos para formar un todo coherente o funcional; reorganizar en un nuevo patron',
 ARRAY['crear', 'disenar', 'producir', 'planificar', 'generar', 'construir', 'desarrollar'],
 '#8B5CF6');

-- Comments
COMMENT ON TABLE bloom_levels IS 'Bloom Revised Taxonomy levels for learning objectives';
COMMENT ON COLUMN bloom_levels.verbos_accion IS 'Action verbs typical of this cognitive level';
