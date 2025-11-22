-- Create customization_items table (global catalog)
CREATE TABLE IF NOT EXISTS customization_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('avatar', 'frame')),
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('common', 'rare', 'epic', 'legendary')),
    image_url TEXT NOT NULL,
    description TEXT,
    unlock_condition JSONB NOT NULL,  -- For display/reference
    base_coins_cost INTEGER NOT NULL DEFAULT 0 CHECK (base_coins_cost >= 0),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create item_unlock_triggers table (event index)
CREATE TABLE IF NOT EXISTS item_unlock_triggers (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES customization_items(id) ON DELETE CASCADE,
    trigger_type VARCHAR(50) NOT NULL,  -- 'oa_complete', 'bloom_mastery', 'streak', 'diagnostic_score', 'coins', 'default'
    trigger_key VARCHAR(100) NOT NULL,  -- 'oa_5', 'bloom_4', 'streak_7', 'materia_1_score_80', 'default'
    additional_data JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for fast event lookups
CREATE INDEX idx_unlock_triggers_lookup ON item_unlock_triggers(trigger_type, trigger_key);
CREATE INDEX idx_customization_items_type ON customization_items(type);
CREATE INDEX idx_customization_items_rarity ON customization_items(rarity);

-- ============================================
-- SEED DATA: AVATARS (10 items)
-- ============================================

-- Default avatars (3 - common, unlocked from start)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Avatar Estudiante', 'avatar', 'common', '/assets/avatars/student.svg', 'Avatar basico de estudiante', '{"type": "default", "value": true}', 0, true),
('Avatar Lechuza', 'avatar', 'common', '/assets/avatars/owl.svg', 'La sabiduria de la lechuza', '{"type": "default", "value": true}', 0, true),
('Avatar Libro', 'avatar', 'common', '/assets/avatars/book.svg', 'El conocimiento en forma de libro', '{"type": "default", "value": true}', 0, true);

-- OA completion avatars (2 - common/rare)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Avatar Matematico', 'avatar', 'common', '/assets/avatars/math.svg', 'Domina un OA de Matematicas para desbloquearlo', '{"type": "oa_complete", "materia": "Matematicas"}', 0, false),
('Avatar Literato', 'avatar', 'rare', '/assets/avatars/literature.svg', 'Completa 3 OAs de Lengua y Literatura', '{"type": "oa_count", "materia": "Lengua y Literatura", "count": 3}', 0, false);

-- Bloom mastery avatars (2 - rare/epic)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Avatar Analizador', 'avatar', 'rare', '/assets/avatars/analyzer.svg', 'Domina 3 OAs en nivel Bloom Analizar', '{"type": "bloom_mastery", "bloom_level": 4, "count": 3}', 0, false),
('Avatar Creador', 'avatar', 'epic', '/assets/avatars/creator.svg', 'Domina 2 OAs en nivel Bloom Crear', '{"type": "bloom_mastery", "bloom_level": 6, "count": 2}', 0, false);

-- Streak avatars (2 - rare/epic)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Avatar Constante', 'avatar', 'rare', '/assets/avatars/streak7.svg', 'Manten una racha de 7 dias consecutivos', '{"type": "streak", "days": 7}', 0, false),
('Avatar Dedicado', 'avatar', 'epic', '/assets/avatars/streak30.svg', 'Manten una racha de 30 dias consecutivos', '{"type": "streak", "days": 30}', 0, false);

-- Legendary avatar (1 - multiple conditions)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Avatar Maestro', 'avatar', 'legendary', '/assets/avatars/master.svg', 'Alcanza nivel 10 y completa 10 OAs', '{"type": "multiple", "conditions": [{"type": "level", "value": 10}, {"type": "oa_count_total", "count": 10}]}', 0, false);

-- ============================================
-- SEED DATA: FRAMES (6 items)
-- ============================================

-- Default frames (2 - common)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Marco Basico', 'frame', 'common', '/assets/frames/basic.svg', 'Marco simple y elegante', '{"type": "default", "value": true}', 0, true),
('Marco Clasico', 'frame', 'common', '/assets/frames/classic.svg', 'Marco clasico para tu avatar', '{"type": "default", "value": true}', 0, true);

-- Purchasable with coins (2 - rare)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Marco Dorado', 'frame', 'rare', '/assets/frames/gold.svg', 'Marco dorado de prestigio', '{"type": "coins", "amount": 500}', 500, false),
('Marco Brillante', 'frame', 'epic', '/assets/frames/shiny.svg', 'Marco con efecto de brillo', '{"type": "coins", "amount": 1000}', 1000, false);

-- Achievement frames (2 - epic/legendary)
INSERT INTO customization_items (name, type, rarity, image_url, description, unlock_condition, base_coins_cost, is_default)
VALUES
('Marco Evaluador', 'frame', 'epic', '/assets/frames/evaluator.svg', 'Completa 5 diagnosticos con 80%+ de acierto', '{"type": "diagnostic_count", "count": 5, "min_score": 80}', 0, false),
('Marco Leyenda', 'frame', 'legendary', '/assets/frames/legend.svg', 'Alcanza nivel 15', '{"type": "level", "value": 15}', 0, false);

-- ============================================
-- SEED DATA: TRIGGERS (Event Index)
-- ============================================

-- Default items triggers (3 avatars + 2 frames = 5)
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
VALUES
(1, 'default', 'default', '{}'),
(2, 'default', 'default', '{}'),
(3, 'default', 'default', '{}'),
(11, 'default', 'default', '{}'),
(12, 'default', 'default', '{}');

-- OA completion triggers
-- Avatar Matematico: Trigger on ANY Matematicas OA completion
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
SELECT 4, 'oa_complete', 'oa_' || id,
       jsonb_build_object('materia_id', materia_id, 'materia_nombre', 'Matematicas')
FROM objetivos_aprendizaje
WHERE materia_id = (SELECT id FROM materias WHERE nombre = 'Matematicas');

-- Avatar Literato: Trigger on ANY Lengua OA completion (check count in validator)
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
SELECT 5, 'oa_complete', 'oa_' || id,
       jsonb_build_object('materia_id', materia_id, 'materia_nombre', 'Lengua y Literatura', 'requires_count', 3)
FROM objetivos_aprendizaje
WHERE materia_id = (SELECT id FROM materias WHERE nombre = 'Lengua y Literatura');

-- Bloom mastery triggers
-- Avatar Analizador: Trigger when ANY OA with Bloom level 4 (Analizar) is completed
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
SELECT 6, 'bloom_mastery', 'bloom_4_oa_' || oa_id,
       jsonb_build_object('bloom_level', 4, 'bloom_name', 'Analizar', 'requires_count', 3)
FROM oa_bloom_objectives
WHERE bloom_level_id = 4;

-- Avatar Creador: Trigger when ANY OA with Bloom level 6 (Crear) is completed
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
SELECT 7, 'bloom_mastery', 'bloom_6_oa_' || oa_id,
       jsonb_build_object('bloom_level', 6, 'bloom_name', 'Crear', 'requires_count', 2)
FROM oa_bloom_objectives
WHERE bloom_level_id = 6;

-- Streak triggers
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
VALUES
(8, 'streak', 'streak_7', '{"days": 7}'),
(9, 'streak', 'streak_30', '{"days": 30}');

-- Legendary avatar (multiple conditions - will check in validator)
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
VALUES
(10, 'level_up', 'level_10', '{"level": 10, "also_requires_oa_count": 10}');

-- Purchasable frames (coins)
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
VALUES
(13, 'coins', 'coins_500', '{"amount": 500}'),
(14, 'coins', 'coins_1000', '{"amount": 1000}');

-- Achievement frames
INSERT INTO item_unlock_triggers (item_id, trigger_type, trigger_key, additional_data)
VALUES
(15, 'diagnostic_achievement', 'diagnostic_5_score_80', '{"count": 5, "min_score": 80}'),
(16, 'level_up', 'level_15', '{"level": 15}');

COMMENT ON TABLE customization_items IS 'Global catalog of all avatars, frames, and future customization items';
COMMENT ON TABLE item_unlock_triggers IS 'Event index mapping trigger events to items that should be unlocked';
COMMENT ON COLUMN item_unlock_triggers.trigger_key IS 'Specific key for fast lookups, e.g., oa_5, bloom_4_oa_10, streak_7';
