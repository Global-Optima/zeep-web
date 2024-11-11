-- Insert into FacilityAddress
INSERT INTO
  facility_addresses (address, longitude, latitude)
VALUES
  ('123 Main St, Springfield', -93.2923, 44.0481),
  ('456 Oak St, Shelbyville', -92.1145, 43.9501),
  ('789 Elm St, Capital City', -91.4045, 42.8912);

-- Insert into ProductCategory
INSERT INTO
  product_categories (name, description)
VALUES
  ('Beverages', 'Drinks of all kinds'),
  ('Snacks', 'Light food items and appetizers'),
  ('Desserts', 'Sweet items and confections');

-- Insert into AdditiveCategory
INSERT INTO
  additive_categories (name, description)
VALUES
  ('Flavorings', 'Extra flavors to enhance taste'),
  ('Sweeteners', 'Additional sweetness'),
  ('Toppings', 'Toppings for desserts and drinks');

-- Insert into Products
INSERT INTO
  products (
    name,
    description,
    image_url,
    video_url,
    category_id
  )
VALUES
  (
    'Latte',
    'A smooth blend of espresso and milk',
    'https://example.com/images/latte.jpg',
    'https://example.com/videos/latte.mp4',
    1
  ),
  (
    'Cappuccino',
    'Espresso with steamed milk foam',
    'https://example.com/images/cappuccino.jpg',
    'https://example.com/videos/cappuccino.mp4',
    1
  ),
  (
    'Brownie',
    'Rich chocolate dessert',
    'https://example.com/images/brownie.jpg',
    NULL,
    3
  );

-- Insert into RecipeStep
INSERT INTO
  recipe_steps (product_id, step, name, description, image_url)
VALUES
  (
    1,
    1,
    'Brew Espresso',
    'Brew the espresso shot',
    'https://example.com/images/espresso.jpg'
  ),
  (
    1,
    2,
    'Steam Milk',
    'Steam milk to desired texture',
    'https://example.com/images/steamed-milk.jpg'
  ),
  (
    2,
    1,
    'Prepare Foam',
    'Make a thick foam for topping',
    'https://example.com/images/foam.jpg'
  );

-- Insert into ProductSize
INSERT INTO
  product_sizes (
    name,
    measure,
    base_price,
    size,
    is_default,
    product_id
  )
VALUES
  ('Small', 'oz', 3.50, 8, true, 1),
  ('Medium', 'oz', 4.50, 12, false, 1),
  ('Large', 'oz', 5.50, 16, false, 1),
  ('Regular', 'oz', 4.00, 12, true, 2);

-- Insert into Additive
INSERT INTO
  additives (
    name,
    description,
    base_price,
    size,
    additive_category_id,
    image_url
  )
VALUES
  (
    'Vanilla Syrup',
    'Sweet vanilla flavoring',
    0.50,
    'ml',
    1,
    'https://example.com/images/vanilla.jpg'
  ),
  (
    'Caramel Syrup',
    'Sweet caramel flavor',
    0.75,
    'ml',
    1,
    'https://example.com/images/caramel.jpg'
  ),
  (
    'Whipped Cream',
    'Creamy topping for beverages',
    0.80,
    'ml',
    3,
    'https://example.com/images/whipped-cream.jpg'
  );

-- Insert into Store
INSERT INTO
  stores (name, facility_address_id, is_franchise, admin_id)
VALUES
  ('Downtown Coffee House', 1, false, NULL),
  ('Uptown Coffee Spot', 2, true, NULL);

-- Insert into StoreAdditive
INSERT INTO
  store_additives (additive_id, store_id, price)
VALUES
  (1, 1, 0.60),
  (2, 1, 0.75),
  (3, 2, 0.85);

-- Insert into StoreProductSize
INSERT INTO
  store_product_sizes (product_size_id, store_id, price)
VALUES
  (1, 1, 3.50),
  (2, 1, 4.00),
  (3, 2, 5.00);

-- Insert into StoreProduct
INSERT INTO
  store_products (product_id, store_id, is_available)
VALUES
  (1, 1, true),
  (2, 1, true),
  (3, 2, false);

-- Insert into ProductAdditive
INSERT INTO
  product_additives (product_size_id, additive_id)
VALUES
  (1, 1),
  (1, 2),
  (2, 3);

-- Insert into DefaultProductAdditive
INSERT INTO
  default_product_additives (product_id, additive_id)
VALUES
  (1, 1),
  (1, 2),
  (2, 3);

-- Insert into Ingredient
INSERT INTO
  ingredients (name, calories, fat, carbs, proteins, expires_at)
VALUES
  ('Sugar', 387, 0, 100, 0, '2024-12-31 00:00:00+00'),
  ('Milk', 42, 1, 5, 3, '2024-01-15 00:00:00+00');

-- Insert into ProductIngredient
INSERT INTO
  product_ingredients (item_ingredient_id, product_id)
VALUES
  (1, 1),
  (2, 1),
  (1, 2);

-- Insert into Customer
INSERT INTO
  customers (name, password, phone, is_verified, is_banned)
VALUES
  (
    'John Doe',
    'hashed_password_123',
    '123-456-7890',
    true,
    false
  ),
  (
    'Jane Smith',
    'hashed_password_456',
    '987-654-3210',
    false,
    false
  );

-- Insert into EmployeeRole
INSERT INTO
  employee_roles (name)
VALUES
  ('Manager'),
  ('Barista'),
  ('Cashier');

-- Insert into Employee
INSERT INTO
  employees (name, phone, email, role_id, store_id, is_active)
VALUES
  (
    'Alice Johnson',
    '111-222-3333',
    'alice@example.com',
    1,
    1,
    true
  ),
  (
    'Bob Brown',
    '444-555-6666',
    'bob@example.com',
    2,
    2,
    true
  );

-- Insert into EmployeeAudit
INSERT INTO
  employee_audits (start_work_at, end_work_at, employee_id)
VALUES
  (
    '2024-10-01 09:00:00+00',
    '2024-10-01 17:00:00+00',
    1
  ),
  (
    '2024-10-02 09:00:00+00',
    '2024-10-02 17:00:00+00',
    2
  );

-- Insert into EmployeeWorkday
INSERT INTO
  employee_workdays (day, start_at, end_at, employee_id)
VALUES
  ('Monday', '08:00:00', '16:00:00', 1),
  ('Tuesday', '08:00:00', '16:00:00', 2);

-- Insert into Referral
INSERT INTO
  referrals (customer_id, referee_id)
VALUES
  (1, 2);

-- Insert into VerificationCode
INSERT INTO
  verification_codes (customer_id, code, expires_at)
VALUES
  (1, '123456', '2024-12-31 23:59:59+00'),
  (2, '654321', '2024-12-31 23:59:59+00');

-- Insert into CustomerAddress
INSERT INTO
  customer_addresses (customer_id, address, longitude, latitude)
VALUES
  (1, '123 Maple St', '-93.2923', '44.0481'),
  (2, '456 Oak St', '-92.1145', '43.9501');

-- Insert into Bonus
INSERT INTO
  bonuses (bonuses, customer_id, expires_at)
VALUES
  (100.00, 1, '2024-12-31 23:59:59+00'),
  (50.00, 2, '2024-06-30 23:59:59+00');