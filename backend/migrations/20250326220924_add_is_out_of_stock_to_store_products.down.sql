ALTER TABLE store_products
DROP COLUMN is_out_of_stock;

ALTER TABLE store_additives
DROP COLUMN is_out_of_stock;

ALTER TABLE additives
DROP CONSTRAINT additives_base_price_check;

ALTER TABLE additives
ADD CONSTRAINT additives_base_price_check CHECK (base_price > 0);

ALTER TABLE product_size_additives
DROP COLUMN is_hidden BOOLEAN;