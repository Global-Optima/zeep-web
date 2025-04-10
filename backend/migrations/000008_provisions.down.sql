-- Rollback ALTERs first
ALTER TABLE product_sizes
    DROP COLUMN IF EXISTS provisions_updated_at;

ALTER TABLE additives
    DROP COLUMN IF EXISTS provisions_updated_at;

-- Then drop child tables first to avoid foreign key issues
DROP TABLE IF EXISTS store_provision_ingredients;
DROP TABLE IF EXISTS store_provisions;
DROP TABLE IF EXISTS product_size_provisions;
DROP TABLE IF EXISTS additive_provisions;
DROP TABLE IF EXISTS provision_ingredients;

-- Then drop the main table
DROP INDEX IF EXISTS unique_provision_name;
DROP TABLE IF EXISTS provisions;
