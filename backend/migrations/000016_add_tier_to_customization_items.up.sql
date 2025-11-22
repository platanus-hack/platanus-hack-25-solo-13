-- Add tier field to customization_items for numerical tier representation (1-5 stars)
ALTER TABLE customization_items
ADD COLUMN tier INTEGER CHECK (tier BETWEEN 1 AND 5);

-- Map existing rarities to tiers
UPDATE customization_items
SET tier = CASE rarity
    WHEN 'common' THEN 1
    WHEN 'rare' THEN 3
    WHEN 'epic' THEN 4
    WHEN 'legendary' THEN 5
    ELSE 1
END;

-- Make tier NOT NULL after populating existing rows
ALTER TABLE customization_items
ALTER COLUMN tier SET NOT NULL;

-- Set default for future inserts
ALTER TABLE customization_items
ALTER COLUMN tier SET DEFAULT 1;
