-- Drop trigger
DROP TRIGGER IF EXISTS trigger_actualizar_progreso ON diagnostic_results;

-- Drop function
DROP FUNCTION IF EXISTS actualizar_progreso_desde_diagnostico();
