-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_equipment_timestamp ON user_equipment;
DROP FUNCTION IF EXISTS update_equipment_timestamp();

-- Drop indexes
DROP INDEX IF EXISTS idx_unlock_notifications_user_unread;
DROP INDEX IF EXISTS idx_user_inventory_item;
DROP INDEX IF EXISTS idx_user_inventory_user;

-- Drop tables (order matters due to FKs)
DROP TABLE IF EXISTS unlock_notifications;
DROP TABLE IF EXISTS user_equipment;
DROP TABLE IF EXISTS user_inventory;
