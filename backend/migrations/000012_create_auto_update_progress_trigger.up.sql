-- Function to auto-update student_oa_progress from diagnostic_results
CREATE OR REPLACE FUNCTION actualizar_progreso_desde_diagnostico()
RETURNS TRIGGER AS $$
BEGIN
    -- Insert or update student_oa_progress based on diagnostic_results
    INSERT INTO student_oa_progress (
        user_id,
        oa_bloom_objective_id,
        estado,
        porcentaje_logro,
        notas,
        ultima_actividad_fecha
    )
    SELECT
        ds.user_id,
        oabo.id,
        CASE
            WHEN NEW.porcentaje_aciertos >= 80 THEN 'logrado'
            WHEN NEW.porcentaje_aciertos >= 60 THEN 'en_proceso'
            ELSE 'no_iniciado'
        END,
        NEW.porcentaje_aciertos,
        'Diagnostico sesion ' || NEW.session_id || ' - Intento ' || ds.numero_intento,
        NOW()
    FROM diagnostic_sessions ds
    JOIN oa_bloom_objectives oabo ON oabo.oa_id = NEW.oa_id
        AND oabo.bloom_level_id = NEW.nivel_bloom_dominado
    WHERE ds.id = NEW.session_id
    ON CONFLICT (user_id, oa_bloom_objective_id)
    DO UPDATE SET
        estado = EXCLUDED.estado,
        porcentaje_logro = EXCLUDED.porcentaje_logro,
        notas = student_oa_progress.notas || E'\n' || EXCLUDED.notas,
        ultima_actividad_fecha = NOW(),
        updated_at = NOW();

    -- Also create history entry
    INSERT INTO student_oa_history (
        user_id,
        oa_bloom_objective_id,
        estado,
        porcentaje_logro,
        tipo_evento,
        notas
    )
    SELECT
        ds.user_id,
        oabo.id,
        CASE
            WHEN NEW.porcentaje_aciertos >= 80 THEN 'logrado'
            WHEN NEW.porcentaje_aciertos >= 60 THEN 'en_proceso'
            ELSE 'no_iniciado'
        END,
        NEW.porcentaje_aciertos,
        'diagnostico',
        'Diagnostico sesion ' || NEW.session_id || ': ' || NEW.preguntas_correctas || '/' || NEW.preguntas_respondidas
    FROM diagnostic_sessions ds
    JOIN oa_bloom_objectives oabo ON oabo.oa_id = NEW.oa_id
        AND oabo.bloom_level_id = NEW.nivel_bloom_dominado
    WHERE ds.id = NEW.session_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger on diagnostic_results
CREATE TRIGGER trigger_actualizar_progreso
    AFTER INSERT ON diagnostic_results
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_progreso_desde_diagnostico();

COMMENT ON FUNCTION actualizar_progreso_desde_diagnostico IS 'Auto-updates student_oa_progress and creates history entry when diagnostic results are created';
