-- Drop indexes
DROP INDEX IF EXISTS idx_customization_items_rarity;
DROP INDEX IF EXISTS idx_customization_items_type;
DROP INDEX IF EXISTS idx_unlock_triggers_lookup;

-- Drop tables (triggers table first due to FK)
DROP TABLE IF EXISTS item_unlock_triggers;
DROP TABLE IF EXISTS customization_items;
