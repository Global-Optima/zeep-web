ALTER TABLE store_products
    ADD COLUMN is_out_of_stock BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE store_additives
    ADD COLUMN is_out_of_stock BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE additives
    DROP CONSTRAINT additives_base_price_check;

ALTER TABLE additives
    ADD CONSTRAINT additives_base_price_check CHECK (base_price >= 0);

ALTER TABLE product_sizes
    ADD COLUMN additives_updated_at TIMESTAMPTZ;

ALTER TABLE product_size_additives
    ADD COLUMN is_hidden BOOLEAN NOT NULL DEFAULT FALSE;