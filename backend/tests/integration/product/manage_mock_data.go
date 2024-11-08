package product

import (
	"fmt"

	"gorm.io/gorm"
)

func SetupMockData(db *gorm.DB) error {
	// Insert stores
	if err := db.Exec(`INSERT INTO stores (id, name) VALUES (1, 'Downtown Store');`).Error; err != nil {
		return err
	}

	// Insert categories
	if err := db.Exec(`
		INSERT INTO categories (id, name) VALUES 
		(1, 'Coffee'), 
		(2, 'Tea');`).Error; err != nil {
		return err
	}

	// Insert products with a mix of those having default additives, optional additives, or none
	if err := db.Exec(`
		INSERT INTO products (id, name, description, image_url, video_url, category_id) VALUES 
		(1, 'Latte', 'A creamy coffee drink', 'https://example.com/latte.jpg', 'https://example.com/latte.mp4', 1),
		(2, 'Espresso', 'Strong coffee', 'https://example.com/espresso.jpg', NULL, 1),
		(3, 'Black Tea', 'Classic black tea', 'https://example.com/black_tea.jpg', NULL, 2),
		(4, 'Green Tea', 'Fresh green tea', 'https://example.com/green_tea.jpg', NULL, 2),
		(5, 'Mocha', 'Chocolate-flavored coffee', 'https://example.com/mocha.jpg', NULL, 1);`).Error; err != nil {
		return err
	}

	// Insert store products with some unavailable items
	if err := db.Exec(`
		INSERT INTO store_products (store_id, product_id, is_available) VALUES 
		(1, 1, true),  -- Latte is available
		(1, 2, true),  -- Espresso is available
		(1, 3, true),  -- Black Tea is available
		(1, 4, false), -- Green Tea is unavailable
		(1, 5, true);  -- Mocha is available
	`).Error; err != nil {
		return err
	}

	// Insert store warehouse stocks to handle stock status
	if err := db.Exec(`
		INSERT INTO store_warehouse_stocks (store_warehouse_id, ingredient_id, quantity) VALUES 
		(1, 1, 50),  -- Latte has stock
		(1, 2, 30),  -- Espresso has stock
		(1, 4, 0);   -- Green Tea is out of stock
	`).Error; err != nil {
		return err
	}

	// Insert product sizes with different configurations for each product
	if err := db.Exec(`
		INSERT INTO product_sizes (id, product_id, name, size, measure, base_price, is_default) VALUES 
		(1, 1, 'Small', 200, 'ml', 5.0, true),    -- Latte default size
		(2, 1, 'Large', 300, 'ml', 6.5, false),   -- Latte optional size
		(3, 2, 'Standard', 100, 'ml', 4.0, true); -- Espresso default size only
	`).Error; err != nil {
		return err
	}

	// Insert additives with different configurations
	if err := db.Exec(`
		INSERT INTO additives (id, name, description, base_price, additive_category_id) VALUES 
		(1, 'Vanilla Syrup', 'Sweet vanilla flavor', 1.5, 1),
		(2, 'Caramel Syrup', 'Rich caramel flavor', 1.75, 1),
		(3, 'Honey', 'Natural sweetener', 2.0, 2);  -- Additional additive for variety
	`).Error; err != nil {
		return err
	}

	// Link additives to products (optional additives)
	if err := db.Exec(`
		INSERT INTO product_additives (product_size_id, additive_id) VALUES 
		(1, 1), -- Vanilla Syrup for Latte Small
		(2, 2); -- Caramel Syrup for Latte Large
	`).Error; err != nil {
		return err
	}

	// Set default additives (automatically included with the product)
	if err := db.Exec(`
		INSERT INTO default_product_additives (product_id, additive_id) VALUES 
		(1, 1), -- Latte has Vanilla Syrup as default additive
		(1, 2); -- Latte also has Caramel Syrup as default additive
	`).Error; err != nil {
		return err
	}

	// Insert nutrition details with a mix of full and partial data
	if err := db.Exec(`
		INSERT INTO ingredients (id, calories, fat, carbs, proteins) VALUES 
		(1, 100, 5, 15, 3), -- Full nutrition for Latte
		(2, 80, 3, 10, 1), -- Full nutrition for Espresso
		(3, NULL, NULL, NULL, NULL); -- Missing nutrition data for Green Tea
	`).Error; err != nil {
		return err
	}

	// Link nutrition data to products
	if err := db.Exec(`
		INSERT INTO product_ingredients (product_id, item_ingredient_id) VALUES 
		(1, 1), -- Latte with full nutrition
		(2, 2); -- Espresso with full nutrition
	`).Error; err != nil {
		return err
	}

	return nil
}

func TruncateTables(db *gorm.DB) error {
	// Define the order of truncation to respect foreign key dependencies.
	tables := []string{
		"store_warehouse_stocks",    // dependent on other tables
		"store_product_sizes",       // dependent on products and sizes
		"store_additives",           // dependent on additives
		"product_additives",         // dependent on additives and products
		"default_product_additives", // dependent on additives and products
		"product_sizes",             // dependent on products
		"product_ingredients",       // dependent on products and ingredients
		"store_products",            // dependent on products and stores
		"products",                  // dependent on categories
		"categories",                // base table for products
		"additives",                 // base table for product_additives and default_product_additives
		"ingredients",               // base table for product_ingredients
		"stores",                    // base table for store_products

		// TODO: add other tables for other modules ordered in right way
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}
