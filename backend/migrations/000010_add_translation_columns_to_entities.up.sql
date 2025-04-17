-- ========================================================
-- 0) Create or Recreate the translations table with new design.
-- ========================================================
DROP TABLE IF EXISTS app_translations;
CREATE TABLE app_translations (
  id SERIAL PRIMARY KEY,
  translation_id BIGINT NOT NULL,  -- translation group ID; many rows can share this value.
  language_code VARCHAR(10) NOT NULL,
  translated_text TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  CONSTRAINT uq_translation_language UNIQUE (translation_id, language_code)
);

-- Create index on translation_id
CREATE INDEX idx_translation_id ON app_translations(translation_id);

-- ========================================================
-- 1) Update the 'products' table
-- ========================================================
-- Drop pre-existing constraints and columns if they exist.
ALTER TABLE products 
  DROP CONSTRAINT IF EXISTS fk_products_name_translation,
  DROP CONSTRAINT IF EXISTS fk_products_description_translation;

ALTER TABLE products 
  DROP COLUMN IF EXISTS name_translation_id,
  DROP COLUMN IF EXISTS description_translation_id;

-- Add new columns (no FK constraint, because translation_id is not unique by itself).
ALTER TABLE products 
  ADD COLUMN name_translation_id BIGINT NULL,
  ADD COLUMN description_translation_id BIGINT NULL;

-- (Optional: If you wish, you can create a CHECK or application-level trigger to validate that the given translation_id exists in app_translations.)

-- ========================================================
-- 2) Update the 'additives' table
-- ========================================================
ALTER TABLE additives 
  DROP CONSTRAINT IF EXISTS fk_additives_name_translation,
  DROP CONSTRAINT IF EXISTS fk_additives_description_translation;

ALTER TABLE additives 
  DROP COLUMN IF EXISTS name_translation_id,
  DROP COLUMN IF EXISTS description_translation_id;

ALTER TABLE additives
  ADD COLUMN name_translation_id BIGINT NULL,
  ADD COLUMN description_translation_id BIGINT NULL;

-- ========================================================
-- 3) Update the 'ingredients' table
-- ========================================================
ALTER TABLE ingredients 
  DROP CONSTRAINT IF EXISTS fk_ingredients_name_translation;

ALTER TABLE ingredients 
  DROP COLUMN IF EXISTS name_translation_id;

ALTER TABLE ingredients
  ADD COLUMN name_translation_id BIGINT NULL;
