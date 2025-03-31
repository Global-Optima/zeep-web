ALTER TABLE store_products
    DROP COLUMN is_out_of_stock;

ALTER TABLE store_additives
    DROP COLUMN is_out_of_stock;

ALTER TABLE product_size_additives
    DROP COLUMN is_hidden BOOLEAN;

ALTER TABLE product_sizes
    DROP COLUMN additives_updated_at;