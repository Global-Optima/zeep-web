-- Schema setup

-- Additive and Category Tables
CREATE TABLE IF NOT EXISTS additive_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS additives (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price DECIMAL(10, 2) DEFAULT 0,
    additive_category_id INT REFERENCES additive_categories (id) ON DELETE SET NULL,
    image_url VARCHAR(2048)
);

-- Product and Related Tables
CREATE TABLE IF NOT EXISTS product_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(2048),
    base_price DECIMAL(10, 2),
    category_id INT REFERENCES product_categories (id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS product_sizes (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    base_price DECIMAL(10, 2) NOT NULL,
    measure VARCHAR(50),
    size INT NOT NULL,
    is_default BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS product_additives (
    id SERIAL PRIMARY KEY,
    product_size_id INT REFERENCES product_sizes(id) ON DELETE CASCADE,
    additive_id INT REFERENCES additives(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recipe_steps (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    step INT NOT NULL,
    description TEXT,
    image_url VARCHAR(2048)
);

-- Store and Related Tables
CREATE TABLE IF NOT EXISTS facility_addresses (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    longitude DECIMAL(9, 6),
    latitude DECIMAL(9, 6)
);

CREATE TABLE IF NOT EXISTS stores (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    facility_address_id INT REFERENCES facility_addresses(id) ON UPDATE CASCADE ON DELETE SET NULL,
    is_franchise BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS store_additives (
    id SERIAL PRIMARY KEY,
    additive_id INT REFERENCES additives(id) ON DELETE CASCADE,
    store_id INT REFERENCES stores(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) DEFAULT 0
);

CREATE TABLE IF NOT EXISTS store_products (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    store_id INT REFERENCES stores(id) ON DELETE CASCADE,
    is_available BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS store_product_sizes (
    id SERIAL PRIMARY KEY,
    product_size_id INT REFERENCES product_sizes(id) ON DELETE CASCADE,
    store_id INT REFERENCES stores(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) DEFAULT 0
);

-- Employee and Role Tables
CREATE TABLE IF NOT EXISTS employee_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(15),
    email VARCHAR(255),
    role_id INT REFERENCES employee_roles(id) ON UPDATE CASCADE ON DELETE SET NULL,
    store_id INT REFERENCES stores(id) ON UPDATE CASCADE ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE
);

-- Populate tables with test data

-- Universal Test Data Setup

-- Additive Categories and Additives
INSERT INTO additive_categories (id, name, description) VALUES
(1, 'Sweeteners', 'Additives for sweetening'),
(2, 'Dairy', 'Dairy-based additives');

INSERT INTO additives (id, name, description, base_price, image_url, additive_category_id) VALUES
(1, 'Sugar', 'Sweet sugar', 0.50, 'https://example.com/sugar.jpg', 1),
(2, 'Honey', 'Natural honey', 1.00, 'https://example.com/honey.jpg', 1),
(3, 'Cream', 'Rich cream', 1.50, 'https://example.com/cream.jpg', 2);

-- Product Categories and Products
INSERT INTO product_categories (id, name, description) VALUES
(1, 'Beverages', 'Category for drinks'),
(2, 'Snacks', 'Category for snacks');

INSERT INTO products (id, name, description, image_url, base_price, category_id) VALUES
(1, 'Espresso', 'Strong coffee', 'https://example.com/espresso.jpg', 2.99, 1),
(2, 'Latte', 'Coffee with milk', 'https://example.com/latte.jpg', 3.99, 1),
(3, 'Chips', 'Crispy snack', 'https://example.com/chips.jpg', 1.50, 2);

-- Product Sizes and Additive Links
INSERT INTO product_sizes (id, product_id, name, base_price, measure, size, is_default) VALUES
(1, 1, 'Small', 2.99, 'oz', 8, TRUE),
(2, 1, 'Large', 3.99, 'oz', 16, FALSE),
(3, 2, 'Medium', 4.50, 'oz', 12, TRUE);

INSERT INTO product_additives (id, product_size_id, additive_id) VALUES
(1, 1, 1),  -- Small Espresso with Sugar
(2, 1, 2),  -- Small Espresso with Honey
(3, 3, 3);  -- Medium Latte with Cream

-- Recipe Steps for Products
INSERT INTO recipe_steps (id, product_id, step, description, image_url) VALUES
(1, 1, 1, 'Brew the espresso', 'https://example.com/brew.jpg'),
(2, 2, 1, 'Add steamed milk', 'https://example.com/steamed_milk.jpg');

-- Facility Addresses and Stores
INSERT INTO facility_addresses (id, address, longitude, latitude) VALUES
(1, '123 Coffee St', 12.3456, -98.7654),
(2, '456 Snack Ave', 34.0522, -118.2437);

INSERT INTO stores (id, name, facility_address_id, is_franchise) VALUES
(1, 'Downtown Coffee', 1, TRUE),
(2, 'Uptown Snacks', 2, FALSE);

-- Store-Specific Additive Prices
INSERT INTO store_additives (id, additive_id, store_id, price) VALUES
(1, 1, 1, 0.55),
(2, 3, 2, 1.60);

-- Store-Specific Product Availability and Prices
INSERT INTO store_products (id, product_id, store_id, is_available) VALUES
(1, 1, 1, TRUE),
(2, 2, 1, TRUE),
(3, 3, 2, TRUE);

INSERT INTO store_product_sizes (id, product_size_id, store_id, price) VALUES
(1, 1, 1, 3.10),
(2, 3, 1, 4.60);

-- Employee Roles and Employees for Stores
INSERT INTO employee_roles (id, name) VALUES
(1, 'Manager'),
(2, 'Barista');

INSERT INTO employees (id, name, phone, email, is_active, role_id, store_id) VALUES
(1, 'Alice Smith', '+1234567890', 'alice@example.com', TRUE, 1, 1),
(2, 'Bob Johnson', '+0987654321', 'bob@example.com', TRUE, 2, 1),
(3, 'Carol White', '+1122334455', 'carol@example.com', TRUE, 2, 2);
