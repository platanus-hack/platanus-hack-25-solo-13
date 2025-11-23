-- Populate categories for existing OAs based on their content
-- This migration classifies all OAs from Lengua y Literatura - Primero Medio

-- LECTURA (11 OAs) - Comprensión de textos literarios y no literarios
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-2' AND titulo LIKE '%Sintetizar, registrar y ordenar%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-3' AND titulo LIKE '%Leer habitualmente%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-4' AND titulo LIKE '%Reflexionar sobre las diferentes dimensiones%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-5' AND titulo LIKE '%Analizar las narraciones%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-6' AND titulo LIKE '%Analizar los poemas%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-7' AND titulo LIKE '%Analizar los textos dramáticos%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-8' AND titulo LIKE '%Comprender la visión de mundo%tragedias%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-9' AND titulo LIKE '%Comprender la relevancia de las obras del Romanticismo%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-10' AND titulo LIKE '%Analizar y evaluar textos con finalidad argumentativa%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-11' AND titulo LIKE '%Analizar y evaluar textos de los medios de comunicación%';
UPDATE objetivos_aprendizaje SET categoria = 'Lectura' WHERE codigo = 'OA-12' AND titulo LIKE '%Leer y comprender textos no literarios%';

-- ESCRITURA (7 OAs) - Producción de textos escritos
UPDATE objetivos_aprendizaje SET categoria = 'Escritura' WHERE codigo = 'OA-13' AND titulo LIKE '%Aplicar flexiblemente y creativamente las habilidades de escritura%';
UPDATE objetivos_aprendizaje SET categoria = 'Escritura' WHERE codigo = 'OA-14' AND titulo LIKE '%Escribir, con el propósito de explicar un tema%';
UPDATE objetivos_aprendizaje SET categoria = 'Escritura' WHERE codigo = 'OA-15' AND titulo LIKE '%Escribir, con el propósito de persuadir%';
UPDATE objetivos_aprendizaje SET categoria = 'Escritura' WHERE codigo = 'OA-16' AND titulo LIKE '%Planificar, escribir, revisar, reescribir y editar%';
UPDATE objetivos_aprendizaje SET categoria = 'Escritura' WHERE codigo = 'OA-17' AND titulo LIKE '%Usar consistentemente el estilo directo%';
UPDATE objetivos_aprendizaje SET categoria = 'Escritura' WHERE codigo = 'OA-18' AND titulo LIKE '%Usar en sus textos recursos de correferencia léxica%';
UPDATE objetivos_aprendizaje SET categoria = 'Escritura' WHERE codigo = 'OA-19' AND titulo LIKE '%Escribir correctamente para facilitar la comprensión%';

-- COMUNICACIÓN ORAL (4 OAs) - Expresión y comprensión oral
UPDATE objetivos_aprendizaje SET categoria = 'Comunicación Oral' WHERE codigo = 'OA-20' AND titulo LIKE '%Comprender, comparar y evaluar textos orales%';
UPDATE objetivos_aprendizaje SET categoria = 'Comunicación Oral' WHERE codigo = 'OA-21' AND titulo LIKE '%Resumir un discurso argumentativo escuchado%';
UPDATE objetivos_aprendizaje SET categoria = 'Comunicación Oral' WHERE codigo = 'OA-22' AND titulo LIKE '%Dialogar constructivamente%';
UPDATE objetivos_aprendizaje SET categoria = 'Comunicación Oral' WHERE codigo = 'OA-23' AND titulo LIKE '%Expresarse frente a una audiencia%';

-- INVESTIGACIÓN (1 OA) - Análisis crítico y reflexión
UPDATE objetivos_aprendizaje SET categoria = 'Investigación' WHERE codigo = 'OA-24' AND titulo LIKE '%Formular una interpretación de los textos literarios%';

-- Verify classification (should return 23 rows)
-- SELECT categoria, COUNT(*) as total FROM objetivos_aprendizaje GROUP BY categoria ORDER BY categoria;
