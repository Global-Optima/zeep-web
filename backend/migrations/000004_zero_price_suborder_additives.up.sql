-- Drop the existing check constraint (if it exists)
ALTER TABLE suborder_additives
    DROP CONSTRAINT IF EXISTS suborder_additives_price_check;

-- Add a new check constraint that allows a zero price
ALTER TABLE suborder_additives
    ADD CONSTRAINT suborder_additives_price_check CHECK (price >= 0);