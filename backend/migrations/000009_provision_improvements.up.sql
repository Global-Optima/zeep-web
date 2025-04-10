ALTER TABLE provisions
    ADD COLUMN default_expiration_in_minutes INT NOT NULL CHECK (default_expiration_in_minutes > 0);

ALTER TABLE store_provisions DROP CONSTRAINT store_provisions_volume_check;
