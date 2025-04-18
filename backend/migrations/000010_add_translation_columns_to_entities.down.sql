ALTER TABLE products 
  DROP COLUMN name_translation_id,
  DROP COLUMN description_translation_id;

ALTER TABLE additives 
  DROP COLUMN name_translation_id,
  DROP COLUMN description_translation_id;

ALTER TABLE ingredients 
  DROP COLUMN name_translation_id;

ALTER TABLE additive_categories
  DROP COLUMN name_translation_id,
  DROP COLUMN description_translation_id;
  
ALTER TABLE product_categories
  DROP COLUMN name_translation_id,
  DROP COLUMN description_translation_id;

ALTER TABLE ingredient_categories
  DROP COLUMN name_translation_id,
  DROP COLUMN description_translation_id;

ALTER TABLE units 
  DROP COLUMN name_translation_id;

-- Drop the app_translations table.
DROP TABLE app_translations;
