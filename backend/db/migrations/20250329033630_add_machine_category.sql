-- migrate:up
-- Add machine_category column with default value 'TEA'
ALTER TABLE product_categories
ADD COLUMN machine_category VARCHAR(20) NOT NULL DEFAULT 'TEA';


-- migrate:down
-- Remove machine_category column
ALTER TABLE product_categories
DROP COLUMN machine_category;
