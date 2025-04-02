ALTER TABLE suborder_additives
    DROP CONSTRAINT IF EXISTS suborder_additives_price_check;

ALTER TABLE suborder_additives
    ADD CONSTRAINT suborder_additives_price_check CHECK (price > 0);