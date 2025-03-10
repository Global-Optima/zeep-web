SET timezone TO 'UTC';

CREATE DOMAIN valid_phone AS VARCHAR(16)
CHECK (VALUE ~ '^\+[1-9]\d{1,14}$');

-- FacilityAddress Table
CREATE TABLE
	 facility_addresses (
		id SERIAL PRIMARY KEY,
		address VARCHAR(255) UNIQUE NOT NULL,
		longitude DECIMAL(9, 6) UNIQUE,
		latitude DECIMAL(9, 6) UNIQUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_facility_coordinates
    ON facility_addresses (longitude, latitude)
    WHERE deleted_at IS NULL AND longitude IS NOT NULL AND latitude IS NOT NULL;

-- Units Table
CREATE TABLE
	 units (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		conversion_factor DECIMAL(10, 4) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_unit_name ON units (name) WHERE deleted_at IS NULL;

-- ProductCategory Table
CREATE TABLE
	 product_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_product_category_name ON product_categories (name) WHERE deleted_at IS NULL;

-- AdditiveCategory Table
CREATE TABLE
	 additive_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		is_multiple_select BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unique_additive_category_name ON additive_categories (name) WHERE deleted_at IS NULL;

-- IngredientCategory Table
CREATE TABLE  ingredient_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unique_ingredient_category_name ON ingredient_categories (name) WHERE deleted_at IS NULL;

-- StockMaterialCategory Table
CREATE TABLE  stock_material_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unique_stock_material_category_name ON stock_material_categories (name) WHERE deleted_at IS NULL;

-- Product Table
CREATE TABLE
	 products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		image_url VARCHAR(2048),
		video_url VARCHAR(2048),
		category_id INT NOT NULL REFERENCES product_categories (id) ON UPDATE CASCADE ON DELETE RESTRICT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_product_name ON products (name) WHERE deleted_at IS NULL;

-- RecipeStep Table
CREATE TABLE
	 recipe_steps (
		id SERIAL PRIMARY KEY,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		step INT NOT NULL,
		name VARCHAR(100),
		description TEXT,
		image_url VARCHAR(2048),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Prevent duplicate step numbers for the same product
CREATE UNIQUE INDEX unique_recipe_step_number
    ON recipe_steps (product_id, step)
    WHERE deleted_at IS NULL;

-- ProductSize Table
CREATE TABLE
	 product_sizes (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		unit_id INT NOT NULL REFERENCES units (id) ON DELETE RESTRICT,
		base_price DECIMAL(10, 2) NOT NULL CHECK (base_price > 0),
        size DECIMAL(10, 2) NOT NULL CHECK (product_sizes.size > 0),
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		discount_id INT,
		machine_id VARCHAR(40) NOT NULL UNIQUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Additive Table
CREATE TABLE
	 additives (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		base_price DECIMAL(10, 2) NOT NULL CHECK (base_price > 0),
        size DECIMAL(10, 2) NOT NULL CHECK (additives.size > 0),
        unit_id INT NOT NULL REFERENCES units(id) ON DELETE RESTRICT,
		additive_category_id INT NOT NULL REFERENCES additive_categories (id) ON UPDATE CASCADE ON DELETE RESTRICT,
		image_url VARCHAR(2048),
        machine_id VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_additive_name ON additives (name) WHERE deleted_at IS NULL;

-- Franchisees Table
CREATE TABLE  franchisees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

-- Ensure franchisee names are unique for non-deleted rows
CREATE UNIQUE INDEX unique_franchisee_name
    ON franchisees (name)
    WHERE deleted_at IS NULL;

	-- Regions Table
CREATE TABLE  regions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

-- Warehouses Table
CREATE TABLE
	 warehouses (
		id SERIAL PRIMARY KEY,
		facility_address_id INT NOT NULL REFERENCES facility_addresses (id) ON DELETE RESTRICT,
        region_id INT NOT NULL REFERENCES regions (id) ON DELETE RESTRICT,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Store Table
CREATE TABLE
	 stores (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		facility_address_id INT REFERENCES facility_addresses (id),
		franchisee_id INT REFERENCES franchisees (id) ON DELETE CASCADE,
		warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
		is_active BOOLEAN DEFAULT TRUE,
		contact_phone valid_phone,
		contact_email VARCHAR(255),
		store_hours VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_store_name ON stores (name) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX unique_store_contact_phone ON stores (contact_phone) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX unique_store_contact_email ON stores (contact_email) WHERE deleted_at IS NULL;

-- StoreAdditive Table
CREATE TABLE
	 store_additives (
		id SERIAL PRIMARY KEY,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
    store_price DECIMAL(10, 2) CHECK (store_price > 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_store_additive
    ON store_additives (store_id, additive_id)
    WHERE deleted_at IS NULL;

-- StoreProduct Table
CREATE TABLE
	 store_products (
		id SERIAL PRIMARY KEY,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		is_available BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_store_product
    ON store_products (store_id, product_id)
    WHERE deleted_at IS NULL;

-- StoreProductSize Table
CREATE TABLE
     store_product_sizes (
    id SERIAL PRIMARY KEY,
    product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
    store_product_id INT NOT NULL REFERENCES store_products (id),
    store_price DECIMAL(10, 2) CHECK (store_price > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_store_product_size
    ON store_product_sizes (store_product_id, product_size_id)
    WHERE deleted_at IS NULL;

-- ProductAdditive Table
CREATE TABLE
	 product_size_additives (
		id SERIAL PRIMARY KEY,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
        is_default BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_product_size_additive
    ON product_size_additives (product_size_id, additive_id)
    WHERE deleted_at IS NULL;

-- Ingredient Table
CREATE TABLE
	 ingredients (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		calories DECIMAL(5, 2) CHECK (calories >= 0),
		fat DECIMAL(5, 2) CHECK (fat >= 0),
		carbs DECIMAL(5, 2) CHECK (carbs >= 0),
		proteins DECIMAL(5, 2) CHECK (proteins >= 0),
        expiration_in_days INT CHECK (expiration_in_days >= 0),
        unit_id INT NOT NULL REFERENCES units(id) ON DELETE RESTRICT,
        category_id INT NOT NULL REFERENCES ingredient_categories(id) ON DELETE RESTRICT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_ingredient_name ON ingredients (name) WHERE deleted_at IS NULL;

-- ProductIngredient Table
CREATE TABLE
	 product_size_ingredients (
		id SERIAL PRIMARY KEY,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_product_size_ingredient
    ON product_size_ingredients (product_size_id, ingredient_id)
    WHERE deleted_at IS NULL;

-- AdditiveIngredients Table
CREATE TABLE
	 additive_ingredients (
		id SERIAL PRIMARY KEY,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_additive_ingredient
    ON additive_ingredients (ingredient_id, additive_id)
    WHERE deleted_at IS NULL;

-- Store Stocks Table (Previously StoreWarehouseStocks)
CREATE TABLE
 	store_stocks (
		id SERIAL PRIMARY KEY,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		low_stock_threshold DECIMAL(10,2) NOT NULL CHECK (low_stock_threshold > 0),
		quantity DECIMAL(10,2) NOT NULL CHECK (quantity >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unique_store_stock
    ON store_stocks (store_id, ingredient_id)
    WHERE deleted_at IS NULL;

-- Customer Table
CREATE TABLE
	 customers (
		id SERIAL PRIMARY KEY,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		phone valid_phone,
		is_verified BOOLEAN DEFAULT FALSE,
		is_banned BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_customer_phone ON customers (phone) WHERE deleted_at IS NULL;

-- Employee Table
CREATE TABLE
	 employees (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
		phone valid_phone,
		email VARCHAR(255),
		hashed_password VARCHAR(255) NOT NULL,
		is_active BOOLEAN NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_employee_phone ON employees (phone) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX unique_employee_email ON employees (email) WHERE deleted_at IS NULL;

CREATE TYPE store_employee_role AS ENUM ('STORE_MANAGER', 'BARISTA');

-- StoreEmployee Table
CREATE TABLE
     store_employees (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
    store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
    role store_employee_role NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_store_employee ON store_employees (employee_id, deleted_at) WHERE deleted_at IS NULL;

CREATE TYPE warehouse_employee_role AS ENUM ('WAREHOUSE_MANAGER', 'WAREHOUSE_EMPLOYEE');

-- WarehouseEmployee Table
CREATE TABLE
     warehouse_employees (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
    warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
    role warehouse_employee_role NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_warehouse_employee ON warehouse_employees (employee_id, deleted_at) WHERE deleted_at IS NULL;

CREATE TYPE region_employee_role AS ENUM ('REGION_WAREHOUSE_MANAGER');

-- Region Managers Table
CREATE TABLE  region_employees (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
    region_id INT NOT NULL REFERENCES regions (id) ON DELETE CASCADE,
    role region_employee_role NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_region_employee ON region_employees (employee_id, deleted_at) WHERE deleted_at IS NULL;

CREATE TYPE franchisee_employee_role AS ENUM ('FRANCHISEE_MANAGER', 'FRANCHISEE_OWNER');

-- Franchisee Employees Table
CREATE TABLE  franchisee_employees (
    id SERIAL PRIMARY KEY,
    franchisee_id INT NOT NULL REFERENCES franchisees (id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
    role franchisee_employee_role NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_franchisee_employee ON franchisee_employees (employee_id, deleted_at) WHERE deleted_at IS NULL;

CREATE TYPE admin_role AS ENUM ('ADMIN', 'OWNER');

-- Admin table
CREATE TABLE  admin_employees (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    role admin_role NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_admin_employee ON admin_employees (employee_id, deleted_at) WHERE deleted_at IS NULL;

-- EmployeeWorkTrack Table
CREATE TABLE
     employee_work_tracks (
    id SERIAL PRIMARY KEY,
    start_work_at TIMESTAMPTZ,
    end_work_at TIMESTAMPTZ,
    employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE TYPE http_method AS ENUM ('GET', 'POST', 'PUT', 'PATCH', 'DELETE');
CREATE TYPE operation_type AS ENUM ('GET', 'CREATE', 'UPDATE', 'DELETE');
CREATE TYPE component_name AS ENUM (
    'FRANCHISEE',
    'REGION',
    'PRODUCT',
    'PRODUCT_CATEGORY',
    'STORE_PRODUCT',
    'EMPLOYEE',
    'STORE_EMPLOYEE',
    'WAREHOUSE_EMPLOYEE',
    'FRANCHISEE_EMPLOYEE',
    'REGION_EMPLOYEE',
    'ADMIN_EMPLOYEE',
    'ADDITIVE',
    'ADDITIVE_CATEGORY',
    'STORE_ADDITIVE',
    'PRODUCT_SIZE',
    'RECIPE_STEPS',
    'STORE',
    'WAREHOUSE',
    'STORE_STOCK',
    'INGREDIENT',
    'INGREDIENT_CATEGORY',
    'STOCK_REQUESTS',
    'STOCK_MATERIAL',
    'STOCK_MATERIAL_CATEGORY',
    'WAREHOUSE_STOCK',
    'SUPPLIER',
    'UNIT',
    'ORDER'
    );


CREATE TABLE employee_audits (
    id SERIAL PRIMARY KEY,
    employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
    operation_type operation_type NOT NULL,
    component_name VARCHAR(255) NOT NULL,
    details JSONB,
    ip_address VARCHAR(45) NOT NULL,
    resource_url TEXT NOT NULL,
    method http_method NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_employee_audits_timestamp ON employee_audits(created_at);
CREATE INDEX idx_employee_audits_employee_id ON employee_audits(employee_id);

-- EmployeeWorkday Table
CREATE TABLE
	 employee_workdays (
		id SERIAL PRIMARY KEY,
		day VARCHAR(15) NOT NULL,
		start_at TIME NOT NULL,
		end_at TIME NOT NULL,
		employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_employee_workday
    ON employee_workdays (employee_id, day)
    WHERE deleted_at IS NULL;

-- Table: employee_notifications
CREATE TABLE employee_notifications (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    details JSONB DEFAULT '{}'::JSONB,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ NULL
);

-- Table: employee_notification_recipients
CREATE TABLE employee_notification_recipients (
    id SERIAL PRIMARY KEY,
    notification_id INT NOT NULL REFERENCES employee_notifications(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT FALSE NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ NULL
);

-- Indexes for employee_notification_recipients
CREATE INDEX idx_notification_id ON employee_notification_recipients (notification_id);
CREATE INDEX idx_employee_id ON employee_notification_recipients (employee_id);

-- Referral Table
CREATE TABLE
	 referrals (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		referee_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_referrals
    ON referrals (customer_id, referee_id)
    WHERE deleted_at IS NULL;

-- VerificationCode Table
CREATE TABLE
	 verification_codes (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- CustomerAddress Table
CREATE TABLE
	 customer_addresses (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		address VARCHAR(255) NOT NULL,
		longitude VARCHAR(20),
		latitude VARCHAR(20),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Prevent duplicate addresses for the same customer
CREATE UNIQUE INDEX unique_customer_address
    ON customer_addresses (customer_id, LOWER(address))
    WHERE deleted_at IS NULL;

-- Bonus Table
CREATE TABLE
	 bonuses (
		id SERIAL PRIMARY KEY,
		bonuses DECIMAL(10, 2) CHECK (bonuses >= 0),
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		expires_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Orders Table
CREATE TABLE
	 orders (
		id SERIAL PRIMARY KEY,
		customer_id INT REFERENCES customers (id) ON DELETE SET NULL,
		display_number INT NOT NULL,
		customer_name VARCHAR(255) NOT NULL,
		store_employee_id INT REFERENCES store_employees (id) ON DELETE SET NULL,
		store_id INT REFERENCES stores (id) NOT NULL,
		delivery_address_id INT REFERENCES customer_addresses (id) ON DELETE SET NULL,
		status VARCHAR(50) NOT NULL,
		total DECIMAL(10, 2) NOT NULL CHECK (total >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE INDEX idx_orders_store_display
ON orders (store_id, display_number);

-- SubOrders Table
CREATE TABLE
	 suborders (
		id SERIAL PRIMARY KEY,
		order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
		store_product_size_id INT NOT NULL REFERENCES store_product_sizes (id) ON DELETE RESTRICT,
		price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- SubOrdersAdditives Table
CREATE TABLE
	 suborder_additives (
		id SERIAL PRIMARY KEY,
		suborder_id INT NOT NULL REFERENCES suborders (id) ON DELETE CASCADE,
		store_additive_id INT NOT NULL REFERENCES store_additives (id) ON DELETE CASCADE,
		price DECIMAL(10, 2) NOT NULL CHECK (price > 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StockRequests Table
CREATE TABLE
	 stock_requests (
		id SERIAL PRIMARY KEY,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
		status VARCHAR(50) NOT NULL,
		details JSONB,
		store_comment TEXT,
		warehouse_comment TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StockMaterials Table
CREATE TABLE  stock_materials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
	ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    safety_stock DECIMAL(10,2) NOT NULL CHECK (safety_stock >= 0),
    unit_id INT NOT NULL REFERENCES units(id) ON DELETE RESTRICT,
	size DECIMAL(10, 2) NOT NULL,
	category_id INT NOT NULL REFERENCES stock_material_categories(id) ON DELETE RESTRICT,
    barcode VARCHAR(255),
    expiration_period_in_days INT NOT NULL DEFAULT 1095, -- Default 3 years
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unique_stock_material_barcode ON stock_materials (barcode) WHERE deleted_at IS NULL;

-- StockRequestIngredients Table
CREATE TABLE
	 stock_request_ingredients (
		id SERIAL PRIMARY KEY,
		stock_request_id INT NOT NULL REFERENCES stock_requests (id) ON DELETE CASCADE,
		stock_material_id INT NOT NULL REFERENCES stock_materials(id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		delivered_date TIMESTAMPTZ,
		expiration_date TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Prevent duplicate ingredients in the same stock request
CREATE UNIQUE INDEX unique_stock_request_ingredient
    ON stock_request_ingredients (stock_request_id, stock_material_id)
    WHERE deleted_at IS NULL;

-- Suppliers Table
CREATE TABLE
	 suppliers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		contact_email VARCHAR(255),
		contact_phone valid_phone UNIQUE,
		city VARCHAR(100) NOT NULL,
		address VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Create supplier_warehouse_deliveries Table
CREATE TABLE  supplier_warehouse_deliveries (
    id SERIAL PRIMARY KEY,
    supplier_id INT NOT NULL REFERENCES suppliers (id) ON DELETE CASCADE,
    warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
	delivery_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Create supplier_warehouse_delivery_materials Table
CREATE TABLE  supplier_warehouse_delivery_materials (
    id SERIAL PRIMARY KEY,
    delivery_id INT NOT NULL REFERENCES supplier_warehouse_deliveries (id) ON DELETE CASCADE,
    stock_material_id INT NOT NULL REFERENCES stock_materials (id) ON DELETE CASCADE,
    barcode VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
    expiration_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE
	 warehouse_stocks (
		id SERIAL PRIMARY KEY,
		warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
		stock_material_id INT NOT NULL REFERENCES stock_materials (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_warehouse_stock
    ON warehouse_stocks (warehouse_id, stock_material_id)
    WHERE deleted_at IS NULL;

CREATE TABLE
	 supplier_materials (
		id SERIAL PRIMARY KEY,
		stock_material_id INT NOT NULL REFERENCES stock_materials (id) ON DELETE CASCADE,
		supplier_id INT NOT NULL REFERENCES suppliers (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_supplier_material
    ON supplier_materials (supplier_id, stock_material_id)
    WHERE deleted_at IS NULL;

-- Create the supplier_prices table
CREATE TABLE  supplier_prices (
    id SERIAL PRIMARY KEY,
    supplier_material_id INT NOT NULL REFERENCES supplier_materials(id) ON DELETE CASCADE,
    base_price DECIMAL(10, 2) NOT NULL CHECK (base_price > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    order_id INT NOT NULL REFERENCES orders(id),
    bin VARCHAR(20) NOT NULL,
    transaction_id VARCHAR(50) UNIQUE NOT NULL,
    process_id VARCHAR(50) UNIQUE NULL,
    payment_method VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    currency CHAR(3) NOT NULL,
    qr_number VARCHAR(50),
    card_mask VARCHAR(16),
    icc VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);