package product

import (
	data "github.com/Global-Optima/zeep-web/backend/internal/data/products"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductDTO, error)
	SearchStoreProducts(storeID uint, searchQuery string, category string, offset, limit int) ([]types.ProductDTO, error)
	GetStoreProductDetails(storeID, productID uint) (*types.ProductDTO, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetStoreProducts(storeID uint, category string, offset, limit int) ([]types.ProductDTO, error) {
	var products []types.ProductDTO

	query := r.db.Model(&data.Product{}).
		Select("products.id AS product_id, products.name AS product_name, products.description, products.image_url, "+
			"store_products.is_available, "+
			"(CASE WHEN store_warehouse_stocks.quantity = 0 THEN true ELSE false END) AS is_out_of_stock, "+
			"COALESCE(store_product_sizes.price, product_sizes.base_price) AS price").
		Joins("JOIN store_products ON store_products.product_id = products.id AND store_products.store_id = ?", storeID).
		Joins("JOIN product_sizes ON product_sizes.product_id = products.id").
		Joins("LEFT JOIN store_product_sizes ON store_product_sizes.product_size_id = product_sizes.id AND store_product_sizes.store_id = ?", storeID).
		Joins("LEFT JOIN categories ON categories.id = products.category_id").
		Joins("LEFT JOIN store_warehouse_stocks ON store_warehouse_stocks.store_warehouse_id = store_products.store_id AND store_warehouse_stocks.ingredient_id = products.id").
		Where("store_products.is_available = ?", true).
		Offset(offset).Limit(limit)

	if category != "" {
		query = query.Where("categories.name = ?", category)
	}

	err := query.Scan(&products).Error
	return products, err
}

func (r *productRepository) SearchStoreProducts(storeID uint, searchQuery string, category string, offset, limit int) ([]types.ProductDTO, error) {
	var products []types.ProductDTO

	query := r.db.Model(&data.Product{}).
		Select("products.id AS product_id, products.name AS product_name, products.description, products.image_url, "+
			"store_products.is_available, "+
			"(CASE WHEN store_warehouse_stocks.quantity = 0 THEN true ELSE false END) AS is_out_of_stock, "+
			"COALESCE(store_product_sizes.price, product_sizes.base_price) AS price").
		Joins("JOIN store_products ON store_products.product_id = products.id AND store_products.store_id = ?", storeID).
		Joins("JOIN product_sizes ON product_sizes.product_id = products.id").
		Joins("LEFT JOIN store_product_sizes ON store_product_sizes.product_size_id = product_sizes.id AND store_product_sizes.store_id = ?", storeID).
		Joins("LEFT JOIN categories ON categories.id = products.category_id").
		Joins("LEFT JOIN store_warehouse_stocks ON store_warehouse_stocks.store_warehouse_id = store_products.store_id AND store_warehouse_stocks.ingredient_id = products.id").
		Where("store_products.is_available = ?", true).
		Where("LOWER(products.name) LIKE ?", "%"+searchQuery+"%").
		Offset(offset).Limit(limit)

	if category != "" {
		query = query.Where("categories.name = ?", category)
	}

	err := query.Scan(&products).Error
	return products, err
}

func (r *productRepository) GetStoreProductDetails(storeID, productID uint) (*types.ProductDTO, error) {
	var product types.ProductDTO

	err := r.db.Model(&data.Product{}).
		Select("products.id AS product_id, products.name AS product_name, products.description, products.image_url, products.video_url, "+
			"store_products.is_available, "+
			"(CASE WHEN store_warehouse_stocks.quantity = 0 THEN true ELSE false END) AS is_out_of_stock, "+
			"COALESCE(store_product_sizes.price, product_sizes.base_price) AS price").
		Joins("JOIN store_products ON store_products.product_id = products.id AND store_products.store_id = ?", storeID).
		Joins("JOIN product_sizes ON product_sizes.product_id = products.id").
		Joins("LEFT JOIN store_product_sizes ON store_product_sizes.product_size_id = product_sizes.id AND store_product_sizes.store_id = ?", storeID).
		Joins("LEFT JOIN store_warehouse_stocks ON store_warehouse_stocks.store_warehouse_id = store_products.store_id AND store_warehouse_stocks.ingredient_id = products.id").
		Where("products.id = ?", productID).
		First(&product).Error

	return &product, err
}
