-- ========================================================
-- Down Migration: Remove the new columns and drop the translations table.
-- ========================================================

-- For 'products' table:
ALTER TABLE products 
  DROP COLUMN IF EXISTS name_translation_id,
  DROP COLUMN IF EXISTS description_translation_id;

-- For 'additives' table:
ALTER TABLE additives 
  DROP COLUMN IF EXISTS name_translation_id,
  DROP COLUMN IF EXISTS description_translation_id;

-- For 'ingredients' table:
ALTER TABLE ingredients 
  DROP COLUMN IF EXISTS name_translation_id;

-- Drop the app_translations table.
DROP TABLE IF EXISTS app_translations;
