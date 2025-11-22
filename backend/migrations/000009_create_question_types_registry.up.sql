-- Create question_types registry (catalog of allowed question types)
CREATE TABLE IF NOT EXISTS question_types (
    tipo VARCHAR(30) PRIMARY KEY,
    nombre_display VARCHAR(100) NOT NULL,
    descripcion TEXT,
    schema_example JSONB,
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Seed with 9 existing activity types from frontend
INSERT INTO question_types (tipo, nombre_display, descripcion, activo) VALUES
('multiple_choice', 'Seleccion Multiple', 'Pregunta con opciones A/B/C/D - Bloom: Recordar, Comprender, Aplicar', true),
('true_false', 'Verdadero o Falso', 'Afirmacion para evaluar - Bloom: Recordar, Comprender', true),
('fill_blanks', 'Completar Espacios', 'Texto con espacios en blanco - Bloom: Recordar, Comprender', true),
('drag_drop_matching', 'Relacionar Terminos', 'Emparejar terminos con definiciones - Bloom: Comprender, Aplicar', true),
('sequencing', 'Ordenar Secuencia', 'Ordenar items cronologicamente - Bloom: Comprender, Aplicar', true),
('compare_contrast', 'Comparar y Contrastar', 'Clasificar caracteristicas - Bloom: Analizar', true),
('open_ended', 'Respuesta Abierta', 'Ensayo o texto libre con rubrica - Bloom: Analizar, Evaluar, Crear', true),
('criteria_evaluation', 'Evaluacion por Criterios', 'Evaluar usando rubrica - Bloom: Evaluar', true),
('concept_map', 'Mapa Conceptual', 'Crear diagrama de conceptos - Bloom: Crear', true)
ON CONFLICT (tipo) DO NOTHING;

COMMENT ON TABLE question_types IS 'Catalog of allowed question types for the question bank';
COMMENT ON COLUMN question_types.tipo IS 'Unique identifier for question type (used in questions.tipo)';
COMMENT ON COLUMN question_types.schema_example IS 'Optional JSON schema example for question_data structure';
