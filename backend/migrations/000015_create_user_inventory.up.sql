-- Create user_inventory table (items unlocked per user)
CREATE TABLE IF NOT EXISTS user_inventory (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES customization_items(id) ON DELETE CASCADE,
    acquired_at TIMESTAMP NOT NULL DEFAULT NOW(),
    acquisition_type VARCHAR(20) NOT NULL DEFAULT 'auto_unlock' CHECK (acquisition_type IN ('auto_unlock', 'purchase', 'gift', 'default')),
    PRIMARY KEY (user_id, item_id)
);

-- Create user_equipment table (currently equipped items)
-- Note: Type validation (avatar/frame) enforced in application layer
CREATE TABLE IF NOT EXISTS user_equipment (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    equipped_avatar_id INTEGER REFERENCES customization_items(id) ON DELETE SET NULL,
    equipped_frame_id INTEGER REFERENCES customization_items(id) ON DELETE SET NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create unlock_notifications table (for "New item unlocked!" UI)
CREATE TABLE IF NOT EXISTS unlock_notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES customization_items(id) ON DELETE CASCADE,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_user_inventory_user ON user_inventory(user_id);
CREATE INDEX idx_user_inventory_item ON user_inventory(item_id);
CREATE INDEX idx_unlock_notifications_user_unread ON unlock_notifications(user_id, is_read) WHERE is_read = false;

-- Trigger to auto-update user_equipment.updated_at
CREATE OR REPLACE FUNCTION update_equipment_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_equipment_timestamp
BEFORE UPDATE ON user_equipment
FOR EACH ROW
EXECUTE FUNCTION update_equipment_timestamp();

-- ============================================
-- SEED DATA: Give default items to all users
-- ============================================

-- Insert all default items into inventory for existing users
INSERT INTO user_inventory (user_id, item_id, acquired_at, acquisition_type)
SELECT
    u.id AS user_id,
    ci.id AS item_id,
    NOW() AS acquired_at,
    'default' AS acquisition_type
FROM users u
CROSS JOIN customization_items ci
WHERE ci.is_default = true
ON CONFLICT (user_id, item_id) DO NOTHING;

-- Create equipment records for all users with first default avatar and frame
INSERT INTO user_equipment (user_id, equipped_avatar_id, equipped_frame_id, updated_at)
SELECT
    u.id AS user_id,
    (SELECT id FROM customization_items WHERE type = 'avatar' AND is_default = true ORDER BY id LIMIT 1) AS equipped_avatar_id,
    (SELECT id FROM customization_items WHERE type = 'frame' AND is_default = true ORDER BY id LIMIT 1) AS equipped_frame_id,
    NOW() AS updated_at
FROM users u
WHERE u.id NOT IN (SELECT user_id FROM user_equipment);

COMMENT ON TABLE user_inventory IS 'Stores which items each user has unlocked';
COMMENT ON TABLE user_equipment IS 'Stores which avatar and frame each user currently has equipped';
COMMENT ON TABLE unlock_notifications IS 'Tracks unread unlock notifications for user UI';
