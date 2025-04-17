ALTER TABLE provisions
    DROP COLUMN IF EXISTS default_expiration_in_minutes;

ALTER TABLE store_provisions
    ADD CONSTRAINT store_provisions_volume_check CHECK (volume > 0);

ALTER TABLE store_provisions
    DROP CONSTRAINT IF EXISTS check_initial_volume_positive;

ALTER TABLE store_provisions
    DROP COLUMN IF EXISTS initial_volume;
