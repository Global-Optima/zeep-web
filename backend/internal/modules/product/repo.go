package product

import (
	"strings"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductDAO, error)
	SearchStoreProducts(storeID uint, searchQuery string, category string, offset, limit int) ([]types.ProductDAO, error)
	GetStoreProductDetails(storeID, productID uint) (*types.ProductDAO, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductDAO, error) {
	var products []types.ProductDAO

	query := r.db.Table("products").
		Select(`
        products.id, 
        products.name, 
        products.description, 
        products.image_url, 
        c.name as category_name, 
        COALESCE(store_product_sizes.price, product_sizes.base_price) as base_price, 
        store_products.is_available, 
        (CASE WHEN COALESCE(store_warehouse_stocks.quantity, 0) = 0 THEN true ELSE false END) as is_out_of_stock
    `).
		Joins(`
        JOIN store_products 
        ON store_products.product_id = products.id 
        AND store_products.store_id = ?
    `, storeID).
		Joins(`
        LEFT JOIN categories c 
        ON c.id = products.category_id
    `).
		Joins(`
        LEFT JOIN product_sizes 
        ON product_sizes.product_id = products.id 
        AND product_sizes.is_default = true
    `).
		Joins(`
        LEFT JOIN store_product_sizes 
        ON store_product_sizes.product_size_id = product_sizes.id 
        AND store_product_sizes.store_id = ?
    `, storeID).
		Joins(`
        LEFT JOIN store_warehouse_stocks 
        ON store_warehouse_stocks.ingredient_id = products.id 
        AND store_warehouse_stocks.store_warehouse_id = store_products.store_id
    `).
		Where("store_products.is_available = ?", true)

	if category != "" {
		query = query.Where("LOWER(c.name) = ?", strings.ToLower(category))
	}

	query = query.Offset(offset).Limit(limit)

	if err := query.Scan(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) SearchStoreProducts(storeID uint, searchQuery, category string, offset, limit int) ([]types.ProductDAO, error) {
	var products []types.ProductDAO

	query := r.db.Table("products").
		Select(`
        products.id, 
        products.name, 
        products.description, 
        products.image_url, 
        c.name as category_name, 
        COALESCE(store_product_sizes.price, product_sizes.base_price) as base_price, 
        store_products.is_available, 
        (CASE WHEN COALESCE(store_warehouse_stocks.quantity, 0) = 0 THEN true ELSE false END) as is_out_of_stock
    `).
		Joins(`
        JOIN store_products 
        ON store_products.product_id = products.id 
        AND store_products.store_id = ?
    `, storeID).
		Joins(`
        LEFT JOIN categories c 
        ON c.id = products.category_id
    `).
		Joins(`
        LEFT JOIN product_sizes 
        ON product_sizes.product_id = products.id 
        AND product_sizes.is_default = true
    `).
		Joins(`
        LEFT JOIN store_product_sizes 
        ON store_product_sizes.product_size_id = product_sizes.id 
        AND store_product_sizes.store_id = ?
    `, storeID).
		Joins(`
        LEFT JOIN store_warehouse_stocks 
        ON store_warehouse_stocks.ingredient_id = products.id 
        AND store_warehouse_stocks.store_warehouse_id = store_products.store_id
    `).
		Where("store_products.is_available = ?", true).
		Where("LOWER(products.name) LIKE ?", "%"+strings.ToLower(searchQuery)+"%")

	if category != "" {
		query = query.Where("LOWER(c.name) = ?", strings.ToLower(category))
	}

	query = query.Offset(offset).Limit(limit)

	if err := query.Scan(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) GetStoreProductDetails(storeID, productID uint) (*types.ProductDAO, error) {
	var product types.ProductDAO

	query := r.db.Table("products").
		Select(`
        products.id, 
        products.name, 
        products.description, 
        products.image_url, 
        products.video_url, 
        c.id as category_id, 
        c.name as category_name, 
        COALESCE(store_product_sizes.price, product_sizes.base_price) as base_price, 
        store_products.is_available, 
        (CASE WHEN COALESCE(store_warehouse_stocks.quantity, 0) = 0 THEN true ELSE false END) as is_out_of_stock
    `).
		Joins(`
        JOIN store_products 
        ON store_products.product_id = products.id 
        AND store_products.store_id = ?
    `, storeID).
		Joins(`
        LEFT JOIN categories c 
        ON c.id = products.category_id
    `).
		Joins(`
        LEFT JOIN product_sizes 
        ON product_sizes.product_id = products.id 
        AND product_sizes.is_default = true
    `).
		Joins(`
        LEFT JOIN store_product_sizes 
        ON store_product_sizes.product_size_id = product_sizes.id 
        AND store_product_sizes.store_id = ?
    `, storeID).
		Joins(`
        LEFT JOIN store_warehouse_stocks 
        ON store_warehouse_stocks.ingredient_id = products.id 
        AND store_warehouse_stocks.store_warehouse_id = store_products.store_id
    `).
		Where("products.id = ?", productID)

	if err := query.First(&product).Error; err != nil {
		return nil, err
	}

	if sizes, err := r.loadSizes(productID, storeID); err != nil {
		return nil, err
	} else {
		product.Sizes = sizes
	}

	if additives, err := r.loadAdditives(productID, storeID); err != nil {
		return nil, err
	} else {
		product.Additives = additives
	}

	if defaultAdditives, err := r.loadDefaultAdditives(productID); err != nil {
		return nil, err
	} else {
		product.DefaultAdditives = defaultAdditives
	}

	if nutrition, err := r.loadNutrition(productID); err != nil {
		return nil, err
	} else {
		product.Nutrition = nutrition
	}

	return &product, nil
}

func (r *productRepository) loadSizes(productID, storeID uint) ([]types.SizeDAO, error) {
	var sizes []types.SizeDAO

	err := r.db.Table("product_sizes").
		Select(`
            product_sizes.id, 
            product_sizes.name, 
            product_sizes.size, 
            product_sizes.measure, 
            COALESCE(store_product_sizes.price, product_sizes.base_price) as price, 
            product_sizes.is_default
        `).
		Joins(`
            LEFT JOIN store_product_sizes 
            ON store_product_sizes.product_size_id = product_sizes.id 
            AND store_product_sizes.store_id = ?
        `, storeID).
		Where("product_sizes.product_id = ?", productID).
		Or("store_product_sizes.store_id IS NULL").
		Scan(&sizes).Error

	if err != nil {
		return nil, err
	}

	return sizes, nil
}

func (r *productRepository) loadAdditives(productID, storeID uint) ([]types.AdditiveDAO, error) {
	var additives []types.AdditiveDAO

	err := r.db.Table("additives").
		Select(`
            additives.id, 
            additives.name, 
            additives.description, 
            additive_categories.name as category_name, 
            COALESCE(store_additives.price, additives.base_price) as price
        `).
		Joins(`
            JOIN product_additives 
            ON product_additives.additive_id = additives.id 
            AND product_additives.product_size_id = ?
        `, productID).
		Joins(`
            LEFT JOIN store_additives 
            ON store_additives.additive_id = additives.id 
            AND store_additives.store_id = ?
        `, storeID).
		Joins(`
            LEFT JOIN additive_categories 
            ON additive_categories.id = additives.additive_category_id
        `).
		Joins(`
            LEFT JOIN store_warehouse_stocks 
            ON store_warehouse_stocks.ingredient_id = additives.id 
            AND store_warehouse_stocks.store_warehouse_id = ?
        `, storeID).
		Where("store_additives.additive_id IS NOT NULL OR store_warehouse_stocks.quantity > 0").
		Scan(&additives).Error

	if err != nil {
		return nil, err
	}

	return additives, nil
}

func (r *productRepository) loadNutrition(productID uint) (types.NutritionDAO, error) {
	var nutrition types.NutritionDAO

	err := r.db.Table("ingredients").
		Select(`
            SUM(ingredients.calories) as calories, 
            SUM(ingredients.fat) as fat, 
            SUM(ingredients.carbs) as carbohydrates, 
            SUM(ingredients.proteins) as proteins
        `).
		Joins(`
            JOIN product_ingredients 
            ON product_ingredients.item_ingredient_id = ingredients.id
        `).
		Where("product_ingredients.product_id = ?", productID).
		Group("product_ingredients.product_id").
		Scan(&nutrition).Error

	if err != nil {
		return types.NutritionDAO{}, err
	}

	return nutrition, nil
}

func (r *productRepository) loadDefaultAdditives(productID uint) ([]types.AdditiveDAO, error) {
	var defaultAdditives []types.AdditiveDAO

	err := r.db.Table("additives").
		Select(`
            additives.id, 
            additives.name, 
            additives.description, 
            additive_categories.name as category_name, 
            additives.base_price as price
        `).
		Joins(`
            JOIN default_product_additives 
            ON default_product_additives.additive_id = additives.id
        `).
		Joins(`
            LEFT JOIN additive_categories 
            ON additive_categories.id = additives.additive_category_id
        `).
		Where("default_product_additives.product_id = ?", productID).
		Scan(&defaultAdditives).Error
	if err != nil {
		return nil, err
	}

	return defaultAdditives, nil
}
