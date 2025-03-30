ALTER TABLE store_additives
  DROP CONSTRAINT store_additives_store_price_check,
  ADD CONSTRAINT store_additives_store_price_check CHECK (store_price >= 0);

ALTER TABLE additives
  DROP CONSTRAINT additives_base_price_check,
  ADD CONSTRAINT additives_base_price_check CHECK (base_price >= 0);
