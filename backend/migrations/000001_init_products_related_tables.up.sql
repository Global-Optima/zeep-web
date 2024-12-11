-- FacilityAddress Table
CREATE TABLE
	IF NOT EXISTS facility_addresses (
		id SERIAL PRIMARY KEY,
		address VARCHAR(255) NOT NULL,
		longitude DECIMAL(9, 6),
		latitude DECIMAL(9, 6),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- ProductCategory Table
CREATE TABLE
	IF NOT EXISTS product_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- AdditiveCategory Table
CREATE TABLE
	IF NOT EXISTS additive_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		is_multiple_select BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Product Table
CREATE TABLE
	IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
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
		name VARCHAR(255) NOT NULL,
		description TEXT,
		base_price DECIMAL(10, 2) DEFAULT 0,
		size VARCHAR(200),
		additive_category_id INT REFERENCES additive_categories (id) ON DELETE SET NULL,
		image_url VARCHAR(2048),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Store Table
CREATE TABLE
	IF NOT EXISTS stores (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		facility_address_id INT REFERENCES facility_addresses (id) ON UPDATE CASCADE ON DELETE SET NULL,
		is_franchise BOOLEAN DEFAULT FALSE,
		status VARCHAR(20) DEFAULT 'ACTIVE',
		contact_phone VARCHAR(20),
		contact_email VARCHAR(255),
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

-- StoreProductSize Table
CREATE TABLE
	IF NOT EXISTS store_product_sizes (
		id SERIAL PRIMARY KEY,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		price DECIMAL(10, 2) DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

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

-- ProductAdditive Table
CREATE TABLE
	IF NOT EXISTS product_additives (
		id SERIAL PRIMARY KEY,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- DefaultProductAdditive Table
CREATE TABLE
	IF NOT EXISTS default_product_additives (
		id SERIAL PRIMARY KEY,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Ingredient Table
CREATE TABLE
	IF NOT EXISTS ingredients (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		calories DECIMAL(5, 2) CHECK (calories >= 0),
		fat DECIMAL(5, 2) CHECK (fat >= 0),
		carbs DECIMAL(5, 2) CHECK (carbs >= 0),
		proteins DECIMAL(5, 2) CHECK (proteins >= 0),
		expires_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- ProductIngredient Table
CREATE TABLE
	IF NOT EXISTS product_ingredients (
		id SERIAL PRIMARY KEY,
		item_ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- ItemIngredients Table
CREATE TABLE IF NOT EXISTS item_ingredients (
    id SERIAL PRIMARY KEY,
    ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    item_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ
);

-- CityWarehouses Table
CREATE TABLE
	IF NOT EXISTS city_warehouses (
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
		city_warehouse_id INT NOT NULL REFERENCES city_warehouses (id) ON DELETE CASCADE,
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
		low_stock_threshold DECIMAL(10, 2) NOT NULL CHECK (quantity > 0)
        quantity DECIMAL(10, 2) NOT NULL CHECK (quantity >= 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StockRequests Table
CREATE TABLE
	IF NOT EXISTS stock_requests (
		id SERIAL PRIMARY KEY,
		city_warehouse_id INT NOT NULL REFERENCES city_warehouses (id) ON DELETE CASCADE,
		status VARCHAR(50) NOT NULL,
		request_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- StockRequestIngredients Table
CREATE TABLE
	IF NOT EXISTS stock_request_ingredients (
		id SERIAL PRIMARY KEY,
		stock_request_id INT NOT NULL REFERENCES stock_requests (id) ON DELETE CASCADE,
		ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		quantity DECIMAL(10, 2) NOT NULL CHECK (quantity > 0),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);

-- Customer Table
CREATE TABLE
	IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		phone VARCHAR(15) UNIQUE,
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
		name VARCHAR(255) NOT NULL,
		phone VARCHAR(15) UNIQUE,
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
		warehouse_id INT NOT NULL REFERENCES city_warehouses (id) ON DELETE CASCADE,
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

-- Suppliers Table
CREATE TABLE
	IF NOT EXISTS suppliers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		contact_email VARCHAR(255),
		contact_phone VARCHAR(20),
		address VARCHAR(255),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMPTZ
	);