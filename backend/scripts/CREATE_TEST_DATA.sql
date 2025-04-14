-- Set timezone for consistency
SET timezone TO 'UTC';

-- 1. Insert into facility_addresses
INSERT INTO facility_addresses (address, longitude, latitude)
VALUES
  ('123 Main St, Test City', 12.345678, 98.765432),
  ('456 Secondary St, Test City', 12.346000, 98.765900);

-- 2. Insert into units
INSERT INTO units (name, conversion_factor)
VALUES
  ('Piece', 1),
  ('Kilogram', 1),
  ('Unused', 1);

-- 3. Insert into product_categories
INSERT INTO product_categories (name, description)
VALUES
  ('Coffee', 'Brewed coffee products'),
  ('Tea', 'Varieties of tea');

-- 4. Insert into additive_categories
INSERT INTO additive_categories (name, description, is_multiple_select)
VALUES
  ('Milk Additives', 'Milk based additives', TRUE),
  ('Syrups', 'Flavored syrups', TRUE),
  ('Unused Additives', 'Unused additives description', TRUE);

-- 5. Insert into ingredient_categories
INSERT INTO ingredient_categories (name, description)
VALUES
  ('Dairy', 'Milk and related products'),
  ('Sweeteners', 'Sugar and sweeteners'),
  ('Unused Category', 'Not used category');

-- 6. Insert into stock_material_categories
INSERT INTO stock_material_categories (name, description)
VALUES
  ('Raw Materials', 'Materials used in production'),
  ('Unused Materials', 'Not used materials');

-- 7. Insert into products
INSERT INTO products (
  name,
  description,
  image_key,
  video_key,
  category_id,
  created_at,
  updated_at
)
VALUES
  ('Espresso', 'Strong coffee shot', 'http://example.com/espresso.png', NULL, 1, NOW(), NOW()),
  ('Green Tea', 'Light and healthy', 'http://example.com/greentea.png', NULL, 2, NOW(), NOW()),
  ('Americano', 'Light coffee shot', 'http://example.com/americano.png', NULL, 1, NOW(), NOW());

-- 8. Insert into recipe_steps (for Espresso)
INSERT INTO recipe_steps (product_id, step, name, description, image_key)
VALUES
  (1, 1, 'Grind Beans', 'Grind fresh coffee beans', 'http://example.com/grind.png'),
  (1, 2, 'Tamp', 'Tamp the coffee', 'http://example.com/tamp.png');

-- 9. Insert into product_sizes
INSERT INTO product_sizes (name, unit_id, base_price, size, product_id, machine_id)
VALUES
  ('S', 1, 2.50, 250, 1, 'TEST0000111122223333000001'),
  ('M', 1, 3.50, 350, 1, 'TEST0000111122223333000002'),
  ('S', 2, 4.50, 245, 3, 'TEST0000111122223333000003');

-- 10. Insert into additives
INSERT INTO additives (name, description, base_price, size, unit_id, additive_category_id, image_key, machine_id)
VALUES
  ('Extra Milk', 'Additional milk', 0.50, 100, 1, 1, 'http://example.com/milk.png', 'TEST0000111122223333000009901'),
  ('Vanilla Syrup', 'Sweet vanilla flavor', 0.75, 50, 1, 2, 'http://example.com/vanilla.png', 'TEST0000111122223333000009902');

-- 11. Insert into franchisees
INSERT INTO franchisees (name, description)
VALUES
  ('Franchisee A', 'Main franchise');

-- 12. Insert into regions
INSERT INTO regions (name)
VALUES ('Test Region');

-- 13. Insert into warehouses
INSERT INTO warehouses (facility_address_id, region_id, name)
VALUES (1, 1, 'Main Warehouse');

-- 14. Insert into stores
INSERT INTO stores (name, facility_address_id, franchisee_id, warehouse_id, is_active, contact_phone, contact_email, store_hours)
VALUES
  ('Test Store', 1, 1, 1, true, '+1234567890', 'store@test.com', '8-20'),
  ('Second Store', 2, 1, 1, true, '+1987654321', 'store2@test.com', '9-21');

-- 15. Insert into store_additives
INSERT INTO store_additives (additive_id, store_id, store_price)
VALUES
  (1, 1, 0.55),
  (2, 1, 0.80);

-- 16. Insert into store_products
INSERT INTO store_products (product_id, store_id, is_available)
VALUES
  (1, 1, true),
  (2, 1, true);

-- 17. Insert into store_product_sizes
INSERT INTO store_product_sizes (product_size_id, store_product_id, store_price)
VALUES
  (1, 1, 2.75),
  (2, 1, 3.75);

-- 18. Insert into product_size_additives
INSERT INTO product_size_additives (product_size_id, additive_id, is_default)
VALUES
    (1, 1, true),
    (1, 2, false);

-- 19. Insert into ingredients
INSERT INTO ingredients (name, calories, fat, carbs, proteins, expiration_in_days, unit_id, category_id, is_allergen)
VALUES
  ('Coffee Beans', 5, 0, 0, 0, 365, 1, 1, true),
  ('Milk', 42, 1, 5, 3, 7, 1, 1, false),
  ('Unused ingredient', 5, 0, 0, 0, 365, 1, 1, false);

-- 20. Insert into product_size_ingredients
INSERT INTO product_size_ingredients (ingredient_id, product_size_id, quantity)
VALUES
  (1, 1, 20),
  (2, 1, 50);

-- 21. Insert into additive_ingredients
INSERT INTO additive_ingredients (ingredient_id, additive_id, quantity)
VALUES
  (2, 1, 30);

-- 22. Insert into store_stocks
INSERT INTO store_stocks (store_id, ingredient_id, low_stock_threshold, quantity)
VALUES
  (1, 1, 10, 1000),
  (1, 2, 5, 500);

-- 23. Insert into customers
INSERT INTO customers (first_name, last_name, password, phone, is_verified, is_banned)
VALUES
  ('John', 'Doe', 'hashedpass', '+1111111111', true, false),
  ('Alice', 'Smith', 'hashedpass', '+2222222222', true, false);


-- 24. Insert multiple employees for different roles (using the given hashed password)
INSERT INTO employees (first_name, last_name, phone, email, hashed_password, is_active)
VALUES
  ('Alice', 'Smith', '+2222222222', 'alice@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- STORE_MANAGER
  ('Bob', 'Johnson', '+3333333333', 'bob@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- BARISTA
  ('Charlie', 'Williams', '+4444444444', 'charlie@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- STORE_MANAGER
  ('David', 'Brown', '+5555555555', 'david@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- BARISTA
  ('Emma', 'Davis', '+6666666666', 'emma@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- WAREHOUSE_MANAGER
  ('Frank', 'Wilson', '+7777777777', 'frank@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- WAREHOUSE_EMPLOYEE
  ('Grace', 'Moore', '+8888888888', 'grace@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- REGION_WAREHOUSE_MANAGER
  ('Henry', 'Taylor', '+9999999999', 'henry@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- FRANCHISEE_MANAGER
  ('Ivy', 'Anderson', '+1010101010', 'ivy@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true),  -- FRANCHISEE_OWNER
  ('Jack', 'Thomas', '+1111111111', 'jack@test.com', '$2a$10$GEmb44LusyHrWXXaz5BKce5N8CvBvz3lPK7CuNS.S86.Quec12Xgy', true);  -- ADMIN

-- 25. Insert store employees
INSERT INTO store_employees (employee_id, store_id, role)
VALUES
  (1, 1, 'STORE_MANAGER'),
  (2, 1, 'BARISTA'),
  (3, 2, 'STORE_MANAGER'),
  (4, 2, 'BARISTA');

-- 26. Insert warehouse employees
INSERT INTO warehouse_employees (employee_id, warehouse_id, role)
VALUES
  (5, 1, 'WAREHOUSE_MANAGER'),
  (6, 1, 'WAREHOUSE_EMPLOYEE');

-- 27. Insert region employees
INSERT INTO region_employees (employee_id, region_id, role)
VALUES
  (7, 1, 'REGION_WAREHOUSE_MANAGER');

-- 28. Insert franchisee employees
INSERT INTO franchisee_employees (employee_id, franchisee_id, role)
VALUES
  (8, 1, 'FRANCHISEE_MANAGER'),
  (9, 1, 'FRANCHISEE_OWNER');

-- 29. Insert admin employees
INSERT INTO admin_employees (employee_id, role)
VALUES
  (10, 'ADMIN'),
  (1, 'OWNER');

-- 30. Insert into employee_work_tracks
INSERT INTO employee_work_tracks (start_work_at, end_work_at, employee_id)
VALUES (NOW(), NOW() + INTERVAL '8 hours', 1);

-- 31. Insert into employee_workdays
INSERT INTO employee_workdays (day, start_at, end_at, employee_id)
VALUES ('Monday', '08:00', '16:00', 1);

-- 32. Insert into employee_notifications
INSERT INTO employee_notifications (event_type, priority, details)
VALUES ('OrderCreated', 'HIGH', '{"info": "Test notification"}');

-- 33. Insert into employee_notification_recipients
INSERT INTO employee_notification_recipients (notification_id, employee_id, is_read)
VALUES (1, 1, false);

-- 34. Insert into referrals
INSERT INTO referrals (customer_id, referee_id)
VALUES (1, 1);

-- 35. Insert into verification_codes
INSERT INTO verification_codes (customer_id, code, expires_at)
VALUES (1, 'ABC123', NOW() + INTERVAL '1 day');

-- 36. Insert into customer_addresses
INSERT INTO customer_addresses (customer_id, address, longitude, latitude)
VALUES (1, '789 Customer Rd, Test City', '12.345678', '98.765432');

-- 37. Insert into bonuses
INSERT INTO bonuses (bonuses, customer_id, expires_at)
VALUES (10.00, 1, NOW() + INTERVAL '30 days');

-- 38. Insert into orders
INSERT INTO orders (customer_id, display_number, customer_name, store_employee_id, store_id, delivery_address_id, status, total, created_at)
VALUES
  (1, 1, 'John Doe', 1, 1, 1, 'PENDING', 5.00, now() - interval '2 days'),
  (1, 2, 'John Doe', 1, 1, 1, 'COMPLETED', 7.50, now() - interval '1 day');

-- Insert into orders for Second Store (updated)
INSERT INTO orders (customer_id, display_number, customer_name, store_employee_id, store_id, delivery_address_id, status, total, created_at)
VALUES
  (2, 1, 'Alice Smith', 1, 2, 1, 'PENDING', 10.00, now() - interval '3 days');

-- 39. Insert into suborders

-- For order 1 (John Doe, PENDING)
INSERT INTO suborders (order_id, store_product_size_id, price, status, created_at)
VALUES
  (1, 1, 5.00, 'PENDING', now() - interval '2 days');

-- For order 2 (John Doe, COMPLETED)
INSERT INTO suborders (order_id, store_product_size_id, price, status, created_at)
VALUES
  (2, 1, 7.50, 'COMPLETED', now() - interval '1 day');

-- For order 3 (Alice Smith, PENDING)
INSERT INTO suborders (order_id, store_product_size_id, price, status, created_at)
VALUES
  (3, 1, 10.00, 'PENDING', now() - interval '3 days');

-- 40. Insert into suborder_additives
-- For order 1 suborder
INSERT INTO suborder_additives (suborder_id, store_additive_id, price, created_at)
VALUES
  (1, 1, 0.50, now() - interval '2 days');

-- For order 2 suborder
INSERT INTO suborder_additives (suborder_id, store_additive_id, price, created_at)
VALUES
  (2, 1, 0.75, now() - interval '1 day');

-- 41. Insert into stock_requests
-- INSERT INTO stock_requests (store_id, warehouse_id, status, details, store_comment, warehouse_comment)
-- VALUES (1, 1, 'CREATED', '{}', 'Store comment', 'Warehouse comment');

-- 42. Insert into stock_materials
INSERT INTO stock_materials (name, description, ingredient_id, safety_stock, unit_id, size, category_id, barcode, expiration_period_in_days, is_active)
VALUES
    ('Coffee Beans Jacobs', 'Material for coffee', 1, 50, 1, 100, 1, 'CB001', 365, true),
    ('Coffee Beans Nescafe', 'Material for coffee', 1, 50, 1, 300, 1, 'CB002', 365, true),
    ('Unused stock material', 'not used', 1, 50, 1, 300, 1, 'CB003', 365, true);

-- 43. Insert into stock_request_ingredients
-- INSERT INTO stock_request_ingredients (stock_request_id, stock_material_id, quantity, delivered_date, expiration_date)
-- VALUES (1, 1, 20, NOW(), NOW() + INTERVAL '1 year');

-- 44. Insert into suppliers
INSERT INTO suppliers (name, contact_email, contact_phone, city, address)
VALUES ('Supplier A', 'supplierA@test.com', '+3333333333', 'Test City', '101 Supplier Rd');

-- 45. Insert into supplier_warehouse_deliveries
INSERT INTO supplier_warehouse_deliveries (supplier_id, warehouse_id, delivery_date)
VALUES (1, 1, NOW());

-- 46. Insert into supplier_warehouse_delivery_materials
INSERT INTO supplier_warehouse_delivery_materials (delivery_id, stock_material_id, barcode, quantity, expiration_date)
VALUES (1, 1, 'CB001', 100, NOW() + INTERVAL '1 year');

-- 47. Insert into warehouse_stocks
INSERT INTO warehouse_stocks (warehouse_id, stock_material_id, quantity)
VALUES
  (1, 1, 200),
  (1, 2, 500);

-- 48. Insert into supplier_materials
INSERT INTO supplier_materials (stock_material_id, supplier_id)
VALUES (1, 1);

-- 49. Insert into supplier_prices
INSERT INTO supplier_prices (supplier_material_id, base_price)
VALUES (1, 1.50);

-- 50. Insert into provisions
INSERT INTO
    provisions (
    name,
    absolute_volume,
    net_cost,
    unit_id,
    preparation_in_minutes,
    default_expiration_in_minutes,
    limit_per_day
)
VALUES
    ('Test Заготовка черного чая', 1.8, 700, 1, 35, 850, 3),
    ('Test Фруктовый микс', 3150, 1200, 2, 5, 600, 10),
    ('Test Unused provision', 1.8, 700, 1, 35, 850, 3);

-- 51. Insert into provision_ingredients
INSERT INTO
    provision_ingredients (
    provision_id,
    ingredient_id,
    quantity
)
VALUES
    (1, 1, 0.4),
    (1, 2, 5),
    (2, 1, 0.2),
    (2, 2, 0.5),
    (2, 3, 0.07);

-- 52. Insert into product_size_provisions
INSERT INTO
    product_size_provisions (
    product_size_id,
    provision_id,
    volume
)
VALUES
    (1, 1, 10);

INSERT INTO store_provisions(
    store_id,
    provision_id,
    expiration_in_minutes,
    status,
    volume,
    initial_volume
)
VALUES
    (1, 1, 100, 'PREPARING', 15, 15),
    (1, 2, 100, 'PREPARING', 20, 20),
    (1, 1, 100, 'COMPLETED', 35, 40),
    (1, 2, 100, 'COMPLETED', 30, 30);