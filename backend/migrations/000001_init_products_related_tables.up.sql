-- FacilityAddress Table
CREATE TABLE
	facility_addresses (
		id SERIAL PRIMARY KEY,
		address VARCHAR(255) NOT NULL,
		longitude DECIMAL(9, 6),
		latitude DECIMAL(9, 6),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- ProductCategory Table
CREATE TABLE
	product_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- AdditiveCategory Table
CREATE TABLE
	additive_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Product Table
CREATE TABLE
	products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		image_url VARCHAR(2048),
		video_url VARCHAR(2048),
		category_id INT REFERENCES product_categories (id) ON UPDATE CASCADE ON DELETE SET NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

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
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- ProductSize Table
CREATE TABLE
	product_sizes (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		measure VARCHAR(50),
		base_price DECIMAL(10, 2) NOT NULL,
		size INT NOT NULL,
		is_default BOOLEAN DEFAULT FALSE,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		discount_id INT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Additive Table
CREATE TABLE
	additives (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		base_price DECIMAL(10, 2) DEFAULT 0,
		size VARCHAR(200),
		additive_category_id INT REFERENCES additive_categories (id) ON DELETE SET NULL,
		image_url VARCHAR(2048),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Store Table
CREATE TABLE
	stores (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		facility_address_id INT REFERENCES facility_addresses (id) ON UPDATE CASCADE ON DELETE SET NULL,
		is_franchise BOOLEAN DEFAULT FALSE,
		admin_id INT,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- StoreAdditive Table
CREATE TABLE
	store_additives (
		id SERIAL PRIMARY KEY,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		price DECIMAL(10, 2) DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- StoreProductSize Table
CREATE TABLE
	store_product_sizes (
		id SERIAL PRIMARY KEY,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		price DECIMAL(10, 2) DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- StoreProduct Table
CREATE TABLE
	store_products (
		id SERIAL PRIMARY KEY,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		store_id INT NOT NULL REFERENCES stores (id) ON DELETE CASCADE,
		is_available BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- ProductAdditive Table
CREATE TABLE
	product_additives (
		id SERIAL PRIMARY KEY,
		product_size_id INT NOT NULL REFERENCES product_sizes (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- DefaultProductAdditive Table
CREATE TABLE
	default_product_additives (
		id SERIAL PRIMARY KEY,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		additive_id INT NOT NULL REFERENCES additives (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Ingredient Table
CREATE TABLE
	ingredients (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		calories DECIMAL(5, 2) CHECK (calories >= 0),
		fat DECIMAL(5, 2) CHECK (fat >= 0),
		carbs DECIMAL(5, 2) CHECK (carbs >= 0),
		proteins DECIMAL(5, 2) CHECK (proteins >= 0),
		expires_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- ProductIngredient Table
CREATE TABLE
	product_ingredients (
		id SERIAL PRIMARY KEY,
		item_ingredient_id INT NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
		product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Customer Table
CREATE TABLE
	customers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		phone VARCHAR(15) UNIQUE,
		is_verified BOOLEAN DEFAULT FALSE,
		is_banned BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- EmployeeRole Table
CREATE TABLE
	employee_roles (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Employee Table
CREATE TABLE
	employees (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		phone VARCHAR(15) UNIQUE,
		email VARCHAR(255) UNIQUE,
		role_id INT REFERENCES employee_roles (id) ON UPDATE CASCADE ON DELETE SET NULL,
		store_id INT REFERENCES stores (id) ON UPDATE CASCADE ON DELETE SET NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- EmployeeAudit Table
CREATE TABLE
	employee_audits (
		id SERIAL PRIMARY KEY,
		start_work_at TIMESTAMPTZ,
		end_work_at TIMESTAMPTZ,
		employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- EmployeeWorkday Table
CREATE TABLE
	employee_workdays (
		id SERIAL PRIMARY KEY,
		day VARCHAR(15) NOT NULL,
		start_at TIME NOT NULL,
		end_at TIME NOT NULL,
		employee_id INT NOT NULL REFERENCES employees (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Referral Table
CREATE TABLE
	referrals (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		referee_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- VerificationCode Table
CREATE TABLE
	verification_codes (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
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
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);

-- Bonus Table
CREATE TABLE
	bonuses (
		id SERIAL PRIMARY KEY,
		bonuses DECIMAL(10, 2) CHECK (bonuses >= 0),
		customer_id INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
		expires_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);