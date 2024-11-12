CREATE TABLE "facility_addresses" (
  "id" BIGSERIAL PRIMARY KEY,
  "address" VARCHAR(255) NOT NULL,
  "longitude" DECIMAL(9,6),
  "latitude" DECIMAL(9,6)
);

CREATE TABLE "categories" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(100) NOT NULL,
  "description" TEXT
);

CREATE TABLE "additive_categories" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(100) NOT NULL,
  "description" TEXT
);

CREATE TABLE "stores" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "facility_address_id" BIGINT REFERENCES "facility_addresses"("id") ON DELETE SET NULL,
  "is_franchise" BOOLEAN DEFAULT FALSE,
  "admin_id" BIGINT,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "city_warehouses" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "facility_address_id" BIGINT REFERENCES "facility_addresses"("id") ON DELETE SET NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "store_warehouses" (
  "id" BIGSERIAL PRIMARY KEY,
  "store_id" BIGINT REFERENCES "stores"("id") ON DELETE CASCADE,
  "city_warehouse_id" BIGINT REFERENCES "city_warehouses"("id") ON DELETE SET NULL
);

CREATE TABLE "ingredients" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "calories" DECIMAL(5,2) CHECK (calories >= 0),
  "fat" DECIMAL(5,2) CHECK (fat >= 0),
  "carbs" DECIMAL(5,2) CHECK (carbs >= 0),
  "proteins" DECIMAL(5,2) CHECK (proteins >= 0),
  "expires_at" TIMESTAMP
);

CREATE TABLE "products" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(100) NOT NULL,
  "description" TEXT,
  "image_url" VARCHAR(2048),
  "video_url" VARCHAR(2048),
  "category_id" BIGINT REFERENCES "categories"("id") ON DELETE SET NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "store_warehouse_stocks" (
  "id" BIGSERIAL PRIMARY KEY,
  "ingredient_id" BIGINT REFERENCES "ingredients"("id") ON DELETE CASCADE,
  "store_warehouse_id" BIGINT REFERENCES "store_warehouses"("id") ON DELETE CASCADE,
  "quantity" INT CHECK (quantity >= 0) DEFAULT 0,
  "status" VARCHAR(50) NOT NULL
);

CREATE TABLE "stock_requests" (
  "id" BIGSERIAL PRIMARY KEY,
  "city_warehouse_id" BIGINT REFERENCES "city_warehouses"("id") ON DELETE SET NULL,
  "store_warehouse_id" BIGINT REFERENCES "store_warehouses"("id") ON DELETE CASCADE,
  "status" VARCHAR(50) NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "request_ingredients" (
  "id" BIGSERIAL PRIMARY KEY,
  "stock_request_id" BIGINT REFERENCES "stock_requests"("id") ON DELETE CASCADE,
  "ingredient_id" BIGINT REFERENCES "ingredients"("id") ON DELETE CASCADE,
  "quantity" INT CHECK (quantity > 0) DEFAULT 1
);

CREATE TABLE "item_ingredients" (
  "id" BIGSERIAL PRIMARY KEY,
  "ingredient_id" BIGINT REFERENCES "ingredients"("id") ON DELETE CASCADE,
  "name" VARCHAR(255) NOT NULL,
  "weight" DECIMAL(5,2) CHECK (weight > 0),
  "label" VARCHAR(20)
);

CREATE TABLE "product_ingredients" (
  "id" BIGSERIAL PRIMARY KEY,
  "item_ingredient_id" BIGINT REFERENCES "item_ingredients"("id") ON DELETE CASCADE,
  "product_id" BIGINT REFERENCES "products"("id") ON DELETE CASCADE
);

CREATE TABLE "employee_roles" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(50) NOT NULL
);

CREATE TABLE "employees" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "phone" VARCHAR(15) UNIQUE,
  "email" VARCHAR(255) UNIQUE,
  "role_id" BIGINT REFERENCES "employee_roles"("id") ON DELETE SET NULL,
  "store_id" BIGINT REFERENCES "stores"("id") ON DELETE SET NULL,
  "is_active" BOOLEAN DEFAULT TRUE
);

CREATE TABLE "employee_workdays" (
  "id" BIGSERIAL PRIMARY KEY,
  "day" VARCHAR(15) NOT NULL,
  "start_at" TIME NOT NULL,
  "end_at" TIME NOT NULL,
  "employee_id" BIGINT REFERENCES "employees"("id") ON DELETE CASCADE
);

CREATE TABLE "employee_audit" (
  "id" BIGSERIAL PRIMARY KEY,
  "start_work_at" TIMESTAMP,
  "end_work_at" TIMESTAMP,
  "employee_id" BIGINT REFERENCES "employees"("id") ON DELETE CASCADE
);

CREATE TABLE "customers" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "phone" VARCHAR(15) UNIQUE,
  "is_verified" BOOLEAN DEFAULT FALSE,
  "is_banned" BOOLEAN DEFAULT FALSE
);

CREATE TABLE "customer_addresses" (
  "id" BIGSERIAL PRIMARY KEY,
  "customer_id" BIGINT REFERENCES "customers"("id") ON DELETE CASCADE,
  "address" VARCHAR(255) NOT NULL,
  "longitude" VARCHAR(20),
  "latitude" VARCHAR(20)
);

CREATE TABLE "referrals" (
  "id" BIGSERIAL PRIMARY KEY,
  "customer_id" BIGINT REFERENCES "customers"("id") ON DELETE CASCADE,
  "refereal_id" BIGINT REFERENCES "customers"("id") ON DELETE CASCADE,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "verification_codes" (
  "id" BIGSERIAL PRIMARY KEY,
  "customer_id" BIGINT REFERENCES "customers"("id") ON DELETE CASCADE,
  "code" VARCHAR(6) NOT NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP,
  "expires_at" TIMESTAMP NOT NULL
);

CREATE TABLE "bonuses" (
  "id" BIGSERIAL PRIMARY KEY,
  "bonuses" DECIMAL(10,2) CHECK (bonuses >= 0),
  "customer_id" BIGINT REFERENCES "customers"("id") ON DELETE CASCADE,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "expires_at" TIMESTAMP
);

CREATE TABLE "orders" (
  "id" BIGSERIAL PRIMARY KEY,
  "customer_name" VARCHAR(255) NOT NULL,
  "customer_id" BIGINT REFERENCES "customers"("id") ON DELETE SET NULL,
  "store_id" BIGINT REFERENCES "stores"("id") ON DELETE SET NULL,
  "price" DECIMAL(10,2) CHECK (price >= 0),
  "date" TIMESTAMP,
  "status" VARCHAR(50),
  "employee_id" BIGINT REFERENCES "employees"("id") ON DELETE SET NULL,
  "delivery_address_id" BIGINT REFERENCES "customer_addresses"("id") ON DELETE SET NULL,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "order_products" (
  "id" BIGSERIAL PRIMARY KEY,
  "order_id" BIGINT REFERENCES "orders"("id") ON DELETE CASCADE,
  "product_id" BIGINT REFERENCES "products"("id") ON DELETE CASCADE,
  "quantity" INT CHECK (quantity > 0) DEFAULT 1
);

CREATE TABLE "product_sizes" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(10) NOT NULL,
  "measure" VARCHAR(255),
  "base_price" DECIMAL(10,2) DEFAULT 0,
  "size" INT,
  "is_default" BOOLEAN DEFAULT FALSE,
  "product_id" BIGINT REFERENCES "products"("id") ON DELETE SET NULL,
  "discount_id" BIGINT
);

CREATE TABLE "recipe_steps" (
  "id" BIGSERIAL PRIMARY KEY,
  "product_id" BIGINT REFERENCES "products"("id") ON DELETE CASCADE,
  "step" INT NOT NULL,
  "name" VARCHAR(100),
  "description" TEXT,
  "image_url" VARCHAR(2048)
);

CREATE TABLE "store_products" (
  "id" BIGSERIAL PRIMARY KEY,
  "product_id" BIGINT REFERENCES "products"("id") ON DELETE CASCADE,
  "store_id" BIGINT REFERENCES "stores"("id") ON DELETE CASCADE,
  "is_available" BOOLEAN DEFAULT TRUE
);

CREATE TABLE "store_product_sizes" (
  "id" BIGSERIAL PRIMARY KEY,
  "product_size_id" BIGINT REFERENCES "product_sizes"("id") ON DELETE CASCADE,
  "store_id" BIGINT REFERENCES "stores"("id") ON DELETE CASCADE,
  "price" DECIMAL(10,2) DEFAULT 0
);

CREATE TABLE "additives" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "base_price" DECIMAL(10,2) DEFAULT 0,
  "size" VARCHAR(200),
  "additive_category_id" BIGINT REFERENCES "additive_categories"("id") ON DELETE SET NULL,
  "image_url" VARCHAR(2048),
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "product_additives" (
  "id" BIGSERIAL PRIMARY KEY,
  "product_size_id" BIGINT REFERENCES "product_sizes"("id") ON DELETE CASCADE,
  "additive_id" BIGINT REFERENCES "additives"("id") ON DELETE CASCADE
);

CREATE TABLE "default_product_additives" (
  "id" BIGSERIAL PRIMARY KEY,
  "product_id" BIGINT REFERENCES "products"("id") ON DELETE CASCADE,
  "additive_id" BIGINT REFERENCES "additives"("id") ON DELETE CASCADE
);

CREATE TABLE "store_additives" (
  "id" BIGSERIAL PRIMARY KEY,
  "additive_id" BIGINT REFERENCES "additives"("id") ON DELETE CASCADE,
  "store_id" BIGINT REFERENCES "stores"("id") ON DELETE CASCADE,
  "price" DECIMAL(10,2) DEFAULT 0
);

CREATE TABLE "order_product_additives" (
  "id" BIGSERIAL PRIMARY KEY,
  "additive_id" BIGINT REFERENCES "additives"("id") ON DELETE CASCADE,
  "order_products_id" BIGINT REFERENCES "order_products"("id") ON DELETE CASCADE
);

-- CREATE TABLE "additives_ingredients" (
--     "id" BIGSERIAL PRIMARY KEY,
--     "additive_id" BIGINT REFERENCES "additives"("id") ON DELETE CASCADE,
--     "item_ingredient_id" BIGINT REFERENCES "item_ingredients"("id") ON DELETE CASCADE
-- ); -- based on ERD (still ? count additives as ingredients : they are complex so they have ingredients)


