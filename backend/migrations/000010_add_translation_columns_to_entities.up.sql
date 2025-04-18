CREATE TABLE app_translations (
  id SERIAL PRIMARY KEY,
  translation_id BIGINT NOT NULL,  
  language_code VARCHAR(10) NOT NULL,
  translated_text TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ,
  CONSTRAINT uq_translation_language UNIQUE (translation_id, language_code)
);

CREATE INDEX idx_translation_id ON app_translations(translation_id);

ALTER TABLE products 
  ADD COLUMN name_translation_id BIGINT NULL,
  ADD COLUMN description_translation_id BIGINT NULL;

ALTER TABLE additives
  ADD COLUMN name_translation_id BIGINT NULL,
  ADD COLUMN description_translation_id BIGINT NULL;

ALTER TABLE ingredients
  ADD COLUMN name_translation_id BIGINT NULL;

ALTER TABLE product_categories
  ADD COLUMN name_translation_id BIGINT NULL,
  ADD COLUMN description_translation_id BIGINT NULL;

ALTER TABLE additive_categories
  ADD COLUMN name_translation_id BIGINT NULL,
  ADD COLUMN description_translation_id BIGINT NULL;

ALTER TABLE ingredient_categories
  ADD COLUMN name_translation_id BIGINT NULL,
  ADD COLUMN description_translation_id BIGINT NULL;

ALTER TABLE units
  ADD COLUMN name_translation_id BIGINT NULL;