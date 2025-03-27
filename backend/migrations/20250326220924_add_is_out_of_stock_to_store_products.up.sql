ALTER TABLE store_products
    ADD COLUMN is_out_of_stock BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE store_additives
    ADD COLUMN is_out_of_stock BOOLEAN NOT NULL DEFAULT FALSE;