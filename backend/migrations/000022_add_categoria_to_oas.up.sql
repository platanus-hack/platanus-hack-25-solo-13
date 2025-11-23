-- Add categoria column to objetivos_aprendizaje table
ALTER TABLE objetivos_aprendizaje
ADD COLUMN categoria VARCHAR(50) DEFAULT 'General';

-- Classify existing OAs into categories based on title keywords
UPDATE objetivos_aprendizaje
SET categoria = CASE
    WHEN LOWER(titulo) LIKE '%lectura%' OR LOWER(titulo) LIKE '%leer%' OR LOWER(titulo) LIKE '%texto%narrativo%' OR LOWER(titulo) LIKE '%comprender%' OR LOWER(titulo) LIKE '%interpretar%' THEN 'Lectura'
    WHEN LOWER(titulo) LIKE '%escrib%' OR LOWER(titulo) LIKE '%redact%' OR LOWER(titulo) LIKE '%produci%' THEN 'Escritura'
    WHEN LOWER(titulo) LIKE '%oral%' OR LOWER(titulo) LIKE '%habla%' OR LOWER(titulo) LIKE '%exposi%' OR LOWER(titulo) LIKE '%diálogo%' OR LOWER(titulo) LIKE '%dialogo%' OR LOWER(titulo) LIKE '%conversa%' THEN 'Comunicación Oral'
    WHEN LOWER(titulo) LIKE '%investiga%' OR LOWER(titulo) LIKE '%indaga%' OR LOWER(titulo) LIKE '%búsqueda%' OR LOWER(titulo) LIKE '%busqueda%' OR LOWER(titulo) LIKE '%fuente%' THEN 'Investigación'
    ELSE 'General'
END;

-- Create index for faster filtering by categoria
CREATE INDEX idx_objetivos_aprendizaje_categoria ON objetivos_aprendizaje(categoria);
