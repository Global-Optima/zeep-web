ALTER TABLE provisions
    ADD COLUMN default_expiration_in_minutes INT NOT NULL CHECK (default_expiration_in_minutes > 0) DEFAULT 180;

ALTER TABLE store_provisions DROP CONSTRAINT store_provisions_volume_check;

ALTER TABLE store_provisions
    ADD COLUMN initial_volume DECIMAL(10,2);

UPDATE store_provisions
SET initial_volume = volume;

ALTER TABLE store_provisions
    ALTER COLUMN initial_volume SET NOT NULL,
ADD CONSTRAINT check_initial_volume_positive CHECK (initial_volume > 0);