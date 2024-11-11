package product

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetStoreProducts(storeID, categoryID uint, searchQuery string, limit, offset int) ([]data.Product, error)
	GetStoreProductDetails(storeID, productID uint) (*data.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetStoreProducts(storeID, categoryID uint, searchQuery string, limit, offset int) ([]data.Product, error) {
	var products []data.Product

	query := r.db.
		Where("category_id = ?", categoryID).
		Where("EXISTS (?)",
			r.db.Model(&data.StoreProduct{}).
				Select("1").
				Where("store_products.product_id = products.id").
				Where("store_products.store_id = ? AND store_products.is_available = TRUE", storeID),
		).Preload("ProductSizes", "is_default = TRUE")

	if searchQuery != "" {
		searchPattern := "%" + searchQuery + "%"
		query = query.Where("(name ILIKE ? OR description ILIKE ?)", searchPattern, searchPattern)
	}

	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	if products == nil {
		products = []data.Product{}
	}

	return products, nil
}

func (r *productRepository) GetStoreProductDetails(storeID, productID uint) (*data.Product, error) {
	var product data.Product

	err := r.db.Model(&data.Product{}).
		Where("products.id = ?", productID).
		Where("EXISTS (?)",
			r.db.Model(&data.StoreProduct{}).
				Select("1").
				Where("store_products.product_id = products.id").
				Where("store_products.store_id = ? AND store_products.is_available = TRUE", storeID),
		).
		Preload("ProductSizes").
		Preload("DefaultAdditives.Additive").
		Preload("RecipeSteps").
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}
