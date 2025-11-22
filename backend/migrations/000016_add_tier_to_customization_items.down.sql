-- Rollback: Remove tier field from customization_items
ALTER TABLE customization_items
DROP COLUMN IF EXISTS tier;
