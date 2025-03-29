-- Add machine_category column with default value 'TEA'
ALTER TABLE product_categories
ADD COLUMN machine_category VARCHAR(20) NOT NULL DEFAULT 'TEA';
