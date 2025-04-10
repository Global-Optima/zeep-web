ALTER TABLE provisions
    DROP COLUMN default_expiration_in_minutes;

ALTER TABLE store_provisions
    ADD CONSTRAINT store_provisions_volume_check CHECK (volume > 0);
