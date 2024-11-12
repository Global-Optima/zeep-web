package product

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func SetupMockData(db *gorm.DB) {
	// Insert facility addresses
	err := db.Exec(`
        INSERT INTO facility_addresses (id, address, longitude, latitude) VALUES
        (1, '123 Main St', 37.7749, -122.4194)
    `).Error
	if err != nil {
		log.Fatalf("Error inserting facility addresses: %v", err)
	}

	// Insert categories
	err = db.Exec(`
        INSERT INTO categories (id, name, description) VALUES
        (1, 'Coffee', 'Coffee beverages'),
        (2, 'Tea', 'Various teas')
    `).Error
	if err != nil {
		log.Fatalf("Error inserting categories: %v", err)
	}

	// Insert ingredients
	err = db.Exec(`
        INSERT INTO ingredients (id, name, calories, fat, carbs, proteins, expires_at) VALUES
        (1, 'Coffee Beans', 10, 0, 0, 1, '2025-01-01'),
        (2, 'Milk', 60, 3.5, 5, 3, '2025-01-01'),
        (3, 'Black Tea Leaves', 0, 0, 1, 1, '2025-01-01'),
        (4, 'Green Tea Leaves', 0, 0, 1, 1, '2025-01-01'),
        (5, 'Chocolate Syrup', 70, 1, 10, 0.5, '2025-01-01'),
        (6, 'Tashkent Spice Mix', 5, 0.2, 1.5, 0.3, '2025-01-01'),
		(7, 'Coconut Milk', 60, 3.5, 5, 3, '2025-01-01')
    `).Error
	if err != nil {
		log.Fatalf("Error inserting ingredients: %v", err)
	}

	// Insert stores
	err = db.Exec(`
        INSERT INTO stores (id, name, facility_address_id) VALUES
        (1, 'Main Store', 1)
    `).Error
	if err != nil {
		log.Fatalf("Error inserting stores: %v", err)
	}

	// Insert products
	err = db.Exec(`
        INSERT INTO products (id, name, description, image_url, video_url, category_id) VALUES
        (1, 'Latte', 'A creamy coffee drink', 'https://example.com/latte.jpg', 'https://example.com/latte.mp4', 1),
        (2, 'Espresso', 'Strong coffee', 'https://example.com/espresso.jpg', NULL, 1),
        (3, 'Black Tea', 'Classic black tea', 'https://example.com/blacktea.jpg', NULL, 2),
        (4, 'Green Tea', 'Fresh green tea', 'https://example.com/greentea.jpg', NULL, 2),
        (5, 'Mocha', 'Chocolate-flavored coffee', 'https://example.com/mocha.jpg', NULL, 1),
        (6, 'Tashkentskiy Tea', 'Flavour of warm Tashkent', 'https://example.com/tashkentskiy.jpg', NULL, 2),
		(7, 'Cappuccino', 'A rich coffee with a layer of foam', 'https://example.com/cappuccino.jpg', NULL, 1)
    `).Error
	if err != nil {
		log.Fatalf("Error inserting products: %v", err)
	}

	// Insert store products
	err = db.Exec(`
        INSERT INTO store_products (store_id, product_id, is_available) VALUES
        (1, 1, true),  -- Latte available
        (1, 2, true),  -- Espresso available
        (1, 3, true),  -- Black Tea available but out of stock
        (1, 4, false), -- Green Tea unavailable
        (1, 5, true),   -- Mocha available
        (1, 6, true), -- Tashkentskiy available
		(1, 7, true)  -- Cappuccino available in Main Store
    `).Error
	if err != nil {
		log.Fatalf("Error inserting store products: %v", err)
	}

	// Insert city and store warehouses
	err = db.Exec(`
        INSERT INTO city_warehouses (id, name, facility_address_id) VALUES
        (1, 'Main City Warehouse', 1);
        
        INSERT INTO store_warehouses (id, store_id, city_warehouse_id) VALUES
        (1, 1, 1)
    `).Error
	if err != nil {
		log.Fatalf("Error inserting city and store warehouses: %v", err)
	}

	// Insert store warehouse stocks
	err = db.Exec(`
        INSERT INTO store_warehouse_stocks (ingredient_id, store_warehouse_id, quantity, status) VALUES
        (3, 1, 0, 'Out of Stock'),  -- Black Tea out of stock
        (2, 1, 50, 'In Stock'),     -- Coffee Beans in stock
        (3, 1, 0, 'In Stock'),     -- Black Tea Leaves in stock
        (4, 1, 30, 'In Stock'),     -- Green Tea Leaves in stock
        (5, 1, 20, 'In Stock'),     -- Chocolate Syrup in stock
        (6, 1, 40, 'In Stock'),      -- Tashkent Spice Mix in stock
		(7, 1, 20, 'In Stock') 		-- Coconut Milk is in stock 
    `).Error
	if err != nil {
		log.Fatalf("Error inserting store warehouse stocks: %v", err)
	}

	// Insert product sizes
	err = db.Exec(`
        INSERT INTO product_sizes (id, name, size, measure, base_price, is_default, product_id) VALUES
        (1, 'Small', 200, 'ml', 5.0, true, 1),  -- Latte small
        (2, 'Large', 300, 'ml', 6.5, false, 1),  -- Latte large
        (3, 'Small', 300, 'ml', 4.0, true, 2),  -- Espresso small
        (4, 'Standart', 300, 'ml', 3.5, true, 3),  -- Black Tea standard
        (5, 'Standart', 300, 'ml', 5.5, true, 5),  -- Mocha standard
        (6, 'Large', 700, 'ml', 5.0, true, 6),  -- Tashkent standard
		(7, 'Small', 250, 'ml', 5.0, true, 7),  -- Small Cappuccino, default size
    	(8, 'Large', 350, 'ml', 6.5, false, 7) -- Large Cappuccino
    `).Error
	if err != nil {
		log.Fatalf("Error inserting product sizes: %v", err)
	}

	// Insert additive categories and additives
	err = db.Exec(`
        INSERT INTO additive_categories (id, name, description) VALUES
        (1, 'Syrups', 'Sweet syrups');

        INSERT INTO additives (id, name, description, base_price, additive_category_id, image_url) VALUES
        (1, 'Vanilla Syrup', 'Sweet vanilla flavor', 1.5, 1, 'https://example.com/vanilla.jpg'),
        (2, 'Caramel Syrup', 'Rich caramel flavor', 1.75, 1, 'https://example.com/caramel.jpg'),
		(3, 'Hazelnut Syrup', 'Nutty hazelnut flavor', 1.5, 1, 'https://example.com/hazelnut.jpg'),
    	(4, 'Cinnamon', 'Spicy cinnamon flavor', 1.0, 1, 'https://example.com/cinnamon.jpg')
    `).Error
	if err != nil {
		log.Fatalf("Error inserting additives and additive categories: %v", err)
	}

	// Insert store-specific additive prices
	err = db.Exec(`
        INSERT INTO store_additives (store_id, additive_id, price) VALUES
        (1, 1, 1.5),  -- Vanilla Syrup available in Main Store at 1.5
        (1, 2, 1.75),  -- Caramel Syrup available in Main Store at 1.75
		(1, 3, 1.5),  -- Hazelnut Syrup available at 1.5 in Main Store
    	(1, 4, 1.0)   -- Cinnamon available at 1.0 in Main Store
    `).Error
	if err != nil {
		log.Fatalf("Error inserting store-specific additives: %v", err)
	}

	// Insert default product additives
	err = db.Exec(`
        INSERT INTO default_product_additives (product_id, additive_id) VALUES
        (1, 1), -- Vanilla Syrup for Latte
        (1, 2),  -- Caramel Syrup for Latte
		(7, 3),  -- Hazelnut Syrup for Small Cappuccino
    	(7, 4)   -- Cinnamon for Small Cappuccino
    `).Error
	if err != nil {
		log.Fatalf("Error inserting default product additives: %v", err)
	}

	// Insert product additives based on default sizes
	err = db.Exec(`
        INSERT INTO product_additives (product_size_id, additive_id) VALUES
        (1, 1),  -- Vanilla Syrup for Latte Small size
        (1, 2),  -- Caramel Syrup for Latte Small size
        (3, 1),  -- Vanilla Syrup for Espresso Small size
        (5, 2),  -- Caramel Syrup for Mocha Standard size
        (6, 1),   -- Vanilla Syrup for Tashkentskiy Large size
		(7, 3), -- Hazelnut Syrup for Small Cappuccino
		(7, 4) -- Cinnamon for Small Cappuccino
    `).Error
	if err != nil {
		log.Fatalf("Error inserting product additives: %v", err)
	}

	// Insert item ingredients
	err = db.Exec(`
		INSERT INTO item_ingredients (id, ingredient_id, name, weight, label) VALUES
		(1, 1, 'Coffee Beans', 30.0, 'g'),  -- Item for Coffee Beans
		(2, 2, 'Milk', 200.0, 'ml'),         -- Item for Milk
		(3, 7, 'Coconut Milk', 200.0, 'ml') -- Item for Coconut Milk
	`).Error
	if err != nil {
		log.Fatalf("Error inserting item ingredients: %v", err)
	}

	err = db.Exec(`
		INSERT INTO product_ingredients (product_id, item_ingredient_id) VALUES
		(1, 1),  -- Coffee Beans linked to Latte
    	(1, 2),  -- Milk linked to Latte
		(7, 1),  -- Coffee Beans for Cappuccino
		(7, 3)   -- Milk for Cappuccino
	`).Error
	if err != nil {
		log.Fatalf("Error linking ingredients to Cappuccino: %v", err)
	}
}

func TruncateTables(db *gorm.DB) error {
	tables := []string{
		"store_warehouse_stocks",    // dependent on other tables
		"store_product_sizes",       // dependent on products and sizes
		"store_additives",           // dependent on additives
		"product_additives",         // dependent on additives and products
		"default_product_additives", // dependent on additives and products
		"item_ingredients",          // dependent on ingredients
		"product_sizes",             // dependent on products
		"product_ingredients",       // dependent on products and ingredients
		"store_products",            // dependent on products and stores
		"products",                  // dependent on categories
		"categories",                // base table for products
		"additives",                 // base table for product_additives and default_product_additives
		"additive_categories",
		"ingredients", // base table for product_ingredients
		"stores",      // base table for store_products
		// Additional tables (if any) to ensure complete truncation:
		// TODO: add other tables for other modules ordered in right way
	}

	for _, table := range tables {
		truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)
		if err := db.Exec(truncateQuery).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	log.Println("TRUNCATED TABLES SUCCESSFULLY!")
	return nil
}
