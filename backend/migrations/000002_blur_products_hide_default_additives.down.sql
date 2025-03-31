ALTER TABLE store_products
    DROP COLUMN is_out_of_stock;

ALTER TABLE store_additives
    DROP COLUMN is_out_of_stock,
    DROP CONSTRAINT store_additives_price_check,
    ADD CONSTRAINT store_additives_price_check CHECK (base_price > 0);

ALTER TABLE additives
    DROP CONSTRAINT additives_base_price_check,
    ADD CONSTRAINT additives_base_price_check CHECK (base_price > 0);

ALTER TABLE product_size_additives
    DROP COLUMN is_hidden BOOLEAN;

ALTER TABLE product_sizes
    DROP COLUMN additives_updated_at;