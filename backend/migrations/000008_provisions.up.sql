-- Provisions Table
CREATE TABLE provisions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    absolute_volume DECIMAL(10,2) NOT NULL CHECK (absolute_volume > 0),
    net_cost DECIMAL(10,2) NOT NULL CHECK (net_cost > 0),
    unit_id INT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    preparation_in_minutes INT NOT NULL CHECK (preparation_in_minutes > 0),
    limit_per_day INT NOT NULL CHECK (limit_per_day > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unique_provision_name ON provisions (name) WHERE deleted_at IS NULL;

-- ProvisionIngredients Table
CREATE TABLE provision_ingredients (
    id SERIAL PRIMARY KEY,
    ingredient_id INT NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    provision_id INT NOT NULL REFERENCES provisions(id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    ingredients_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- AdditiveProvisions Table
CREATE TABLE additive_provisions (
    id SERIAL PRIMARY KEY,
    additive_id INT NOT NULL REFERENCES additives(id) ON DELETE CASCADE,
    provision_id INT NOT NULL REFERENCES provisions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- ProductSizeProvisions Table
CREATE TABLE product_size_provisions (
    id SERIAL PRIMARY KEY,
    product_size_id INT NOT NULL REFERENCES product_sizes(id) ON DELETE CASCADE,
    provision_id INT NOT NULL REFERENCES provisions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- StoreProvisions Table
CREATE TABLE store_provisions (
    id SERIAL PRIMARY KEY,
    provision_id INT NOT NULL REFERENCES provisions(id) ON DELETE CASCADE,
    volume DECIMAL(10,2) NOT NULL CHECK (volume > 0),
    status INT NOT NULL,
    store_id INT NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    expiration_in_hours INT NOT NULL,
    completed_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- StoreProvisionIngredients Table
CREATE TABLE store_provision_ingredients (
    id SERIAL PRIMARY KEY,
    ingredient_id INT NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    store_provision_id INT NOT NULL REFERENCES store_provisions(id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE additives
    ADD COLUMN provisions_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE product_sizes
    ADD COLUMN provisions_updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP;