INSERT INTO facility_addresses (address, longitude, latitude) VALUES
  ('123 Main St', 76.123456, 43.123456),
  ('456 Elm St', 76.654321, 43.654321),
  ('789 Oak St', 76.789012, 43.789012),
  ('101 Maple St', 76.345678, 43.345678);

INSERT INTO categories (name, description) VALUES
  ('Beverages', 'Drinks and refreshments'),
  ('Desserts', 'Sweet treats'),
  ('Snacks', 'Light snacks and appetizers');

INSERT INTO additive_categories (name, description) VALUES
  ('Syrups', 'Flavoring syrups'),
  ('Toppings', 'Various toppings for drinks and snacks');

INSERT INTO stores (name, facility_address_id, is_franchise, admin_id, created_at, updated_at) VALUES
  ('URBO Coffee Central', 1, TRUE, 101, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('URBO Coffee South', 2, FALSE, 102, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('URBO Coffee West', 3, TRUE, 103, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO city_warehouses (name, facility_address_id, created_at, updated_at) VALUES
  ('Central Warehouse', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Southern Warehouse', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO store_warehouses (store_id, city_warehouse_id) VALUES
  (1, 1),
  (2, 2),
  (3, 1);

INSERT INTO ingredients (name, calories, fat, carbs, proteins, expires_at) VALUES
  ('Sugar', 387, 0, 100, 0, '2025-01-01'),
  ('Milk', 42, 1, 5, 3.5, '2024-11-30'),
  ('Coffee Beans', 0, 0, 0, 0, '2025-06-01'),
  ('Flour', 364, 1, 76, 10, '2025-03-01'),
  ('Butter', 717, 81, 0, 0.85, '2024-12-15'),
  ('Eggs', 155, 11, 1.1, 13, '2024-11-20'),
  ('Vanilla Extract', 288, 0.1, 12.7, 0.1, '2025-05-10'),
  ('Baking Powder', 53, 0, 28.1, 0, '2025-04-05'),
  ('Salt', 0, 0, 0, 0, '2025-07-01'),
  ('Cocoa Powder', 228, 13.7, 57.9, 19.6, '2025-02-20');

INSERT INTO products (name, description, image_url, video_url, category_id, created_at, updated_at) VALUES
  ('Latte', 'A creamy coffee drink', 'https://example.com/latte.jpg', 'https://example.com/latte.mp4', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Cheesecake', 'Classic cheesecake dessert', 'https://example.com/cheesecake.jpg', 'https://example.com/cheesecake.mp4', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Espresso', 'Strong coffee shot', 'https://example.com/espresso.jpg', 'https://example.com/espresso.mp4', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Cappuccino', 'Espresso with steamed milk foam', 'https://example.com/cappuccino.jpg', 'https://example.com/cappuccino.mp4', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Muffin', 'Blueberry muffin', 'https://example.com/muffin.jpg', 'https://example.com/muffin.mp4', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Bagel', 'Freshly baked bagel', 'https://example.com/bagel.jpg', 'https://example.com/bagel.mp4', 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Mocha', 'Chocolate flavored coffee', 'https://example.com/mocha.jpg', 'https://example.com/mocha.mp4', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Croissant', 'Buttery flaky pastry', 'https://example.com/croissant.jpg', 'https://example.com/croissant.mp4', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Americano', 'Espresso with hot water', 'https://example.com/americano.jpg', 'https://example.com/americano.mp4', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Brownie', 'Chocolate brownie', 'https://example.com/brownie.jpg', 'https://example.com/brownie.mp4', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO store_warehouse_stocks (ingredient_id, store_warehouse_id, quantity, status) VALUES
  (1, 1, 50, 'Available'),
  (2, 2, 30, 'Low Stock'),
  (3, 3, 100, 'Available'),
  (4, 1, 20, 'Available'),
  (5, 2, 10, 'Low Stock'),
  (6, 3, 60, 'Available'),
  (7, 1, 15, 'Available'),
  (8, 2, 25, 'Available'),
  (9, 3, 40, 'Available'),
  (10, 1, 35, 'Available');

INSERT INTO stock_requests (city_warehouse_id, store_warehouse_id, status, created_at, updated_at) VALUES
  (1, 1, 'Pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (2, 2, 'Approved', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (1, 3, 'Pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO request_ingredients (stock_request_id, ingredient_id, quantity) VALUES
  (1, 1, 20),
  (2, 2, 15),
  (3, 3, 25);

INSERT INTO item_ingredients (ingredient_id, name, weight, label) VALUES
  (1, 'Sweetener', 15.0, 'g'),
  (2, 'Dairy', 200.0, 'ml'),
  (3, 'Coffee', 50.0, 'g'),
  (4, 'Flour', 120.0, 'g'),
  (5, 'Butter', 30.0, 'g');

INSERT INTO product_ingredients (item_ingredient_id, product_id) VALUES
  (1, 1),  -- Sweetener for Latte
  (2, 1),  -- Dairy for Latte
  (3, 1),  -- Coffee for Latte
  (4, 5),  -- Flour for Muffin
  (5, 8);  -- Butter for Croissant

INSERT INTO employee_roles (name) VALUES
  ('Barista'),
  ('Manager'),
  ('Cashier');

INSERT INTO employees (name, phone, email, role_id, store_id, is_active) VALUES
  ('John Doe', '1234567890', 'johndoe@example.com', 1, 1, TRUE),
  ('Jane Smith', '0987654321', 'janesmith@example.com', 2, 2, TRUE),
  ('Jim Brown', '5678901234', 'jimbrown@example.com', 3, 3, TRUE);

INSERT INTO employee_workdays (day, start_at, end_at, employee_id) VALUES
  ('Monday', '09:00:00', '17:00:00', 1),
  ('Tuesday', '09:00:00', '17:00:00', 2),
  ('Wednesday', '09:00:00', '17:00:00', 3);

INSERT INTO customers (name, password, phone, is_verified, is_banned) VALUES
  ('Alice Brown', 'password123', '1231231234', TRUE, FALSE),
  ('Bob Green', 'password456', '4564564567', TRUE, FALSE),
  ('Eve White', 'password789', '7897897890', TRUE, FALSE);

INSERT INTO customer_addresses (customer_id, address, longitude, latitude) VALUES
  (1, '789 Pine St', '76.234567', '43.234567'),
  (2, '101 Maple St', '76.345678', '43.345678'),
  (3, '202 Oak St', '76.456789', '43.456789');

INSERT INTO verification_codes (customer_id, code, created_at, updated_at, expires_at) VALUES
  (1, 'AB1234', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '2024-12-01 23:59:59'),
  (2, 'CD5678', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '2024-12-01 23:59:59');

INSERT INTO orders (customer_name, customer_id, store_id, price, date, status, employee_id, delivery_address_id, created_at, updated_at) VALUES
  ('Alice Brown', 1, 1, 150.00, CURRENT_TIMESTAMP, 'Completed', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Bob Green', 2, 2, 75.00, CURRENT_TIMESTAMP, 'Processing', 2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Eve White', 3, 3, 50.00, CURRENT_TIMESTAMP, 'Pending', 3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO order_products (order_id, product_id, quantity) VALUES
  (1, 1, 2),  -- Alice Brown orders 2 Lattes
  (2, 5, 1),  -- Bob Green orders 1 Muffin
  (3, 8, 1);  -- Eve White orders 1 Croissant

INSERT INTO product_sizes (name, measure, base_price, size, is_default, product_id) VALUES
  ('Small', '200ml', 5.00, 200, TRUE, 1),
  ('Large', '400ml', 7.50, 400, FALSE, 1),
  ('Single', '150g', 3.50, 150, TRUE, 5),
  ('Double', '300g', 6.50, 300, FALSE, 5);

INSERT INTO recipe_steps (product_id, step, name, description, image_url) VALUES
  (1, 1, 'Grind Coffee', 'Grind the coffee beans', 'https://example.com/grind.jpg'),
  (1, 2, 'Steam Milk', 'Steam the milk for latte', 'https://example.com/steam.jpg'),
  (5, 1, 'Mix Ingredients', 'Mix flour, sugar, and milk for muffin', 'https://example.com/mix.jpg');

INSERT INTO store_products (product_id, store_id, is_available) VALUES
  (1, 1, TRUE),
  (2, 2, TRUE),
  (3, 3, FALSE),
  (5, 1, TRUE),
  (8, 2, TRUE);

INSERT INTO store_product_sizes (product_size_id, store_id, price) VALUES
  (1, 1, 6.00),
  (2, 1, 8.00),
  (3, 2, 4.00),
  (4, 3, 7.00);

INSERT INTO additives (name, description, base_price, size, additive_category_id, image_url, created_at, updated_at) VALUES
  ('Vanilla Syrup', 'Sweet vanilla syrup', 1.50, '30ml', 1, 'https://example.com/vanilla.jpg', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Caramel Syrup', 'Rich caramel syrup', 1.75, '30ml', 1, 'https://example.com/caramel.jpg', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Whipped Cream', 'Cream topping for drinks', 0.75, '50g', 2, 'https://example.com/cream.jpg', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('Cinnamon', 'Spicy cinnamon powder', 0.50, '5g', 2, 'https://example.com/cinnamon.jpg', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO default_product_additives (product_id, additive_id) VALUES
  (1, 1),  -- Vanilla Syrup for Latte
  (2, 2);  -- Caramel Syrup for Cheesecake

INSERT INTO order_product_additives (additive_id, order_products_id) VALUES
  (1, 1),
  (2, 2);

INSERT INTO referrals (customer_id, refereal_id, created_at) VALUES
  (1, 2, CURRENT_TIMESTAMP),  -- Alice referred Bob
  (2, 3, CURRENT_TIMESTAMP);  -- Bob referred Eve

INSERT INTO bonuses (bonuses, customer_id, created_at, expires_at) VALUES
  (50.00, 1, CURRENT_TIMESTAMP, '2025-01-01'),
  (30.00, 2, CURRENT_TIMESTAMP, '2025-06-01');

INSERT INTO employee_audit (start_work_at, end_work_at, employee_id) VALUES
  ('2024-11-01 09:00:00', '2024-11-01 17:00:00', 1),
  ('2024-11-02 09:00:00', '2024-11-02 17:00:00', 2);

INSERT INTO order_products (order_id, product_id, quantity) VALUES
  (1, 5, 3);  -- Alice ordered 3 muffins


-- Edge Cases:
-- INSERT INTO products (name, description, image_url, category_id, created_at, updated_at)
-- VALUES ('No-Size Product', 'Product without defined size', 'https://example.com/no-size.jpg', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- INSERT INTO products (name, description, image_url, category_id, created_at, updated_at)
-- VALUES ('No-Additive Product', 'Product with no additives', 'https://example.com/no-additive.jpg', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- INSERT INTO store_warehouse_stocks (ingredient_id, store_warehouse_id, quantity, status)
-- VALUES (3, 1, 0, 'Out of Stock');  -- Coffee Beans set to 0 quantity

-- INSERT INTO ingredients (name, calories, fat, carbs, proteins, expires_at)
-- VALUES ('Expired Ingredient', 50, 0, 12, 2, '2023-01-01');

-- INSERT INTO default_product_additives (product_id, additive_id) VALUES
--   (1, 2),  -- Adding Caramel Syrup to Latte (in addition to Vanilla Syrup)
--   (1, 3);  -- Adding Whipped Cream to Latte

-- UPDATE employees SET is_active = FALSE WHERE id = 1;  -- Set John Doe to inactive
