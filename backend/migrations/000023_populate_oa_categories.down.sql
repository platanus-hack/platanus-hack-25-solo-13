-- Revert all OA categories back to 'General'
UPDATE objetivos_aprendizaje
SET categoria = 'General'
WHERE codigo IN (
    'OA-2', 'OA-3', 'OA-4', 'OA-5', 'OA-6', 'OA-7', 'OA-8', 'OA-9', 'OA-10', 'OA-11', 'OA-12',
    'OA-13', 'OA-14', 'OA-15', 'OA-16', 'OA-17', 'OA-18', 'OA-19',
    'OA-20', 'OA-21', 'OA-22', 'OA-23',
    'OA-24'
);
