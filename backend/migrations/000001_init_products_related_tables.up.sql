CREATE DOMAIN valid_phone AS VARCHAR(16)
    CHECK (VALUE ~ '^\+[1-9]\d{1,14}$');

-- FacilityAddress Table
CREATE TABLE
	IF NOT EXISTS facility_addresses (
		id SERIAL PRIMARY KEY,
		address VARCHAR(255) UNIQUE NOT NULL,
		longitude DECIMAL(9, 6) UNIQUE,
		latitude DECIMAL(9, 6) UNIQUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Units Table
CREATE TABLE
	IF NOT EXISTS units (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		conversion_factor DECIMAL(10, 4) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- ProductCategory Table
CREATE TABLE
	IF NOT EXISTS product_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- AdditiveCategory Table
CREATE TABLE
	IF NOT EXISTS additive_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		description TEXT,
		is_multiple_select BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);
-- IngredientCategory Table
CREATE TABLE IF NOT EXISTS ingredient_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

-- Product Table
CREATE TABLE
	IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		description TEXT,
		image_url VARCHAR(2048),
		video_url VARCHAR(2048),
		category_id INT REFERENCES product_categories (id) ON UPDATE CASCADE ON DELETE SET NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- RecipeStep Table
CREATE TABLE
	IF NOT EXISTS recipe_steps (
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

-- ProductSize Table
CREATE TABLE
	IF NOT EXISTS product_sizes (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		measure VARCHAR(50),
		base_price DECIMAL(10, 2) NOT NULL,
		size INT NOT NULL,
		is_default BOOLEAN DEFAULT FALSE,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		discount_id INT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Additive Table
CREATE TABLE
	IF NOT EXISTS additives (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		description TEXT,
		base_price DECIMAL(10, 2) DEFAULT 0,
		size VARCHAR(200),
		additive_category_id INT REFERENCES additive_categories (id),
		image_url VARCHAR(2048),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Store Table
CREATE TABLE
	IF NOT EXISTS stores (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		facility_address_id INT REFERENCES facility_addresses (id),
		is_franchise BOOLEAN DEFAULT FALSE,
		status VARCHAR(20) DEFAULT 'ACTIVE',
		contact_phone valid_phone UNIQUE,
		contact_email VARCHAR(255) UNIQUE,
		store_hours VARCHAR(255),
		admin_id INT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
    );

-- StoreAdditive Table
CREATE TABLE
	IF NOT EXISTS store_additives (
		id SERIAL PRIMARY KEY,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		price DECIMAL(10, 2) DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE UNIQUE INDEX unique_store_additive
    ON store_additives (store_id, additive_id)
    WHERE deleted_at IS NULL;

-- StoreProduct Table
CREATE TABLE
	IF NOT EXISTS store_products (
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
    IF NOT EXISTS store_product_sizes (
    id SERIAL PRIMARY KEY,
    product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
    store_product_id INT NOT NULL REFERENCES store_products (id),
    price DECIMAL(10, 2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
    );

CREATE UNIQUE INDEX unique_store_product_size
    ON store_product_sizes (store_product_id, product_size_id)
    WHERE deleted_at IS NULL;

-- ProductAdditive Table
CREATE TABLE
	IF NOT EXISTS product_size_additives (
		id SERIAL PRIMARY KEY,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
        is_default BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Ingredient Table
CREATE TABLE
	IF NOT EXISTS ingredients (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		calories DECIMAL(5, 2) CHECK (calories >= 0),
		fat DECIMAL(5, 2) CHECK (fat >= 0),
		carbs DECIMAL(5, 2) CHECK (carbs >= 0),
		proteins DECIMAL(5, 2) CHECK (proteins >= 0),
		expires_at TIMESTAMPTZ,
    	unit_id INT NOT NULL REFERENCES units(id) ON DELETE SET NULL,
		ingredient_category_id INT NOT NULL REFERENCES ingredient_categories(id) ON DELETE SET NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- ProductIngredient Table
CREATE TABLE
	IF NOT EXISTS product_size_ingredients (
		id SERIAL PRIMARY KEY,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- AdditiveIngredients Table 
CREATE TABLE
	IF NOT EXISTS additive_ingredients (
		id SERIAL PRIMARY KEY,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Warehouses Table
CREATE TABLE
	IF NOT EXISTS warehouses (
		id SERIAL PRIMARY KEY,
		facility_address_id INT NOT NULL REFERENCES facility_addresses (id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StoreWarehouses Table
CREATE TABLE
	IF NOT EXISTS store_warehouses (
		id SERIAL PRIMARY KEY,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StoreWarehouseStock Table
CREATE TABLE
	IF NOT EXISTS store_warehouse_stocks (
		id SERIAL PRIMARY KEY,
		store_warehouse_id INT NOT NULL REFERENCES store_warehouses (id) ON DELETE CASCADE,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		low_stock_threshold DECIMAL(10, 2) NOT NULL CHECK (low_stock_threshold > 0),
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Customer Table
CREATE TABLE
	IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		phone valid_phone UNIQUE,
		is_verified BOOLEAN DEFAULT FALSE,
		is_banned BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Employee Table
CREATE TABLE
	IF NOT EXISTS employees (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
		phone valid_phone UNIQUE,
		email VARCHAR(255) UNIQUE,
		hashed_password VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL,
		type VARCHAR(50) NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StoreEmployee Table
CREATE TABLE
	IF NOT EXISTS store_employees (
		id SERIAL PRIMARY KEY,
		employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		is_franchise BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- WarehouseEmployee Table
CREATE TABLE
	IF NOT EXISTS warehouse_employees (
		id SERIAL PRIMARY KEY,
		employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
		warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- EmployeeAudit Table
CREATE TABLE
	IF NOT EXISTS employee_audits (
		id SERIAL PRIMARY KEY,
		start_work_at TIMESTAMPTZ,
		end_work_at TIMESTAMPTZ,
		employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- EmployeeWorkday Table
CREATE TABLE
	IF NOT EXISTS employee_workdays (
		id SERIAL PRIMARY KEY,
		day VARCHAR(15) NOT NULL,
		start_at TIME NOT NULL,
		end_at TIME NOT NULL,
		employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Referral Table
CREATE TABLE
	IF NOT EXISTS referrals (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		referee_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- VerificationCode Table
CREATE TABLE
	IF NOT EXISTS verification_codes (
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
	IF NOT EXISTS customer_addresses (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		address VARCHAR(255) NOT NULL,
		longitude VARCHAR(20),
		latitude VARCHAR(20),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Bonus Table
CREATE TABLE
	IF NOT EXISTS bonuses (
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
	IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		customer_id INT REFERENCES customers (id) ON DELETE SET NULL,
		customer_name VARCHAR(255) NOT NULL,
		employee_id INT REFERENCES employees (id) ON DELETE SET NULL,
		store_id INT REFERENCES stores (id) NOT NULL,
		delivery_address_id INT REFERENCES customer_addresses (id) ON DELETE SET NULL,
		status VARCHAR(50) NOT NULL,
		total DECIMAL(10, 2) NOT NULL CHECK (total >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- SubOrders Table
CREATE TABLE
	IF NOT EXISTS suborders (
		id SERIAL PRIMARY KEY,
		order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE SET NULL,
		price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- SubOrdersAdditives Table
CREATE TABLE
	IF NOT EXISTS suborder_additives (
		id SERIAL PRIMARY KEY,
		suborder_id INT NOT NULL REFERENCES suborders (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StockRequests Table
CREATE TABLE
	IF NOT EXISTS stock_requests (
		id SERIAL PRIMARY KEY,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
		status VARCHAR(50) NOT NULL,
		request_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StockMaterials Table
CREATE TABLE IF NOT EXISTS stock_materials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
	ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    safety_stock DECIMAL(10,2) NOT NULL CHECK (safety_stock >= 0),
    expiration_flag BOOLEAN NOT NULL,
    unit_id INT NOT NULL REFERENCES units(id) ON DELETE SET NULL,
    category VARCHAR(255),
    barcode VARCHAR(255) UNIQUE,
    expiration_period_in_days INT NOT NULL DEFAULT 1095, -- Default 3 years
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- StockRequestIngredients Table
CREATE TABLE
	IF NOT EXISTS stock_request_ingredients (
		id SERIAL PRIMARY KEY,
		stock_request_id INT NOT NULL REFERENCES stock_requests (id) ON DELETE CASCADE,
		stock_material_id INT NOT NULL REFERENCES stock_materials(id) ON DELETE CASCADE,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		delivered_date TIMESTAMPTZ,
		expiration_date TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Suppliers Table
CREATE TABLE
	IF NOT EXISTS suppliers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		contact_email VARCHAR(255),
		contact_phone valid_phone,
		address VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);



-- Packages Table
CREATE TABLE IF NOT EXISTS stock_material_packages (
    id SERIAL PRIMARY KEY,
    stock_material_id INT NOT NULL REFERENCES stock_materials(id) ON DELETE CASCADE,
    size DECIMAL(10,2) NOT NULL,
    unit_id INT NOT NULL REFERENCES units(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Deliveries Table
CREATE TABLE
	IF NOT EXISTS supplier_warehouse_deliveries (
		id SERIAL PRIMARY KEY,
		stock_material_id INT NOT NULL REFERENCES stock_materials (id) ON DELETE CASCADE,
		supplier_id INT NOT NULL,
		warehouse_id INT NOT NULL,
		barcode VARCHAR(255) NOT NULL,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		delivery_date TIMESTAMPTZ NOT NULL,
		expiration_date TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE TABLE
	IF NOT EXISTS warehouse_stocks (
		id SERIAL PRIMARY KEY,
		warehouse_id INT NOT NULL REFERENCES warehouses (id) ON DELETE CASCADE,
		stock_material_id INT NOT NULL REFERENCES stock_materials (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

CREATE TABLE
	IF NOT EXISTS supplier_materials (
		id SERIAL PRIMARY KEY,
		stock_material_id INT NOT NULL REFERENCES stock_materials (id) ON DELETE CASCADE,
		supplier_id INT NOT NULL REFERENCES suppliers (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);