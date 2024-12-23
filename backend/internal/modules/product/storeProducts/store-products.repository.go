package storeProducts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreProductRepository interface {
	GetStoreProductById(storeID uint, storeProductID uint) (*types.StoreProductDetailsDTO, error)
	GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]data.StoreProduct, error)
	CreateStoreProduct(storeProduct *data.StoreProduct) (uint, error)
	UpdateStoreProduct(storeID, storeProductID uint, product *data.StoreProduct) error
	DeleteStoreProduct(storeID, storeProductID uint) error

	GetStoreProductSizeById(storeID, storeProductSizeID uint) (*data.StoreProductSize, error)
	//GetStoreProductSizes(storeID, storeProductSizeID uint, filter *types.StoreProductSizesFilterDTO) ([]data.StoreProductSize, error)
	CreateStoreProductSize(storeProductSize *data.StoreProductSize) (uint, error)
	UpdateProductSize(storeID, productSizeID uint, size *data.StoreProductSize) error
	DeleteStoreProductSize(storeID, productSizeID uint) error
}

type storeProductRepository struct {
	db *gorm.DB
}

func NewStoreProductRepository(db *gorm.DB) StoreProductRepository {
	return &storeProductRepository{db: db}
}

// GetStoreProductById retrieves a specific store product by ID and store ID
func (r *storeProductRepository) GetStoreProductById(storeID uint, storeProductID uint) (*types.StoreProductDetailsDTO, error) {
	var storeProduct data.StoreProduct
	err := r.db.
		Preload("Product").
		Preload("Store").
		Preload("Store.ProductSizes", "store_id = ?", storeID).
		Preload("Store.ProductSizes.ProductSize").
		Where("store_id = ? AND id = ?", storeID, storeProductID).
		Find(&storeProduct).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	var defaultAdditives []data.DefaultProductAdditive
	err = r.db.
		Preload("Additive").
		Where("product_id = ?", storeProduct.Product.ID).
		Find(&defaultAdditives).
		Error

	if err != nil {
		return nil, err
	}

	dto := types.MapToStoreProductDetailsDTO(&storeProduct, defaultAdditives)
	return &dto, nil
}

// GetStoreProducts retrieves all store products for a given store with optional filters
func (r *storeProductRepository) GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]data.StoreProduct, error) {
	var storeProducts []data.StoreProduct
	query := r.db.Model(&data.StoreProduct{}).Where("store_id = ?", storeID)

	// Apply filters if provided
	if filter != nil {
		if filter.Name != nil {
			query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
		}
		if filter.IsAvailable != nil {
			query = query.Where("is_available = ?", *filter.IsAvailable)
		}
	}

	if err := query.
		Preload("Store").
		Preload("Store.ProductSizes").
		Preload("Product").
		Preload("Store.ProductSizes", "store_id = ?", storeID).
		Preload("Store.ProductSizes.ProductSize", "is_default = TRUE").
		Preload("Store.ProductSizes.ProductSize.ProductIngredients").
		Preload("Store.ProductSizes.ProductSize.ProductIngredients.Ingredient").
		Where("store_id = ?", storeID).
		Find(&storeProducts).Error; err != nil {
		return nil, err
	}

	// Map results to DTOs

	return storeProducts, nil
}

// CreateStoreProduct creates a new store product
func (r *storeProductRepository) CreateStoreProduct(storeProduct *data.StoreProduct) (uint, error) {
	if err := r.db.Create(storeProduct).Error; err != nil {
		return 0, err
	}
	return storeProduct.ID, nil
}

// UpdateStoreProduct updates an existing store product
func (r *storeProductRepository) UpdateStoreProduct(storeID, productID uint, storeProduct *data.StoreProduct) error {
	result := r.db.Model(&data.StoreProduct{}).
		Where("store_id = ? AND id = ?", storeID, productID).
		Updates(storeProduct)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteStoreProduct deletes a store product
func (r *storeProductRepository) DeleteStoreProduct(storeID, storeProductID uint) error {
	result := r.db.
		Where("store_id = ? AND id = ?", storeID, storeProductID).
		Delete(&data.StoreProduct{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetStoreProductSizeById retrieves a store product size by ID
func (r *storeProductRepository) GetStoreProductSizeById(storeID, storeProductSizeID uint) (*data.StoreProductSize, error) {
	var storeProductSize data.StoreProductSize
	if err := r.db.
		Where("store_id = ? AND id = ?", storeID, storeProductSizeID).
		Find(&storeProductSize).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &storeProductSize, nil
}

// CreateStoreProductSize creates a new store product size
func (r *storeProductRepository) CreateStoreProductSize(storeProductSize *data.StoreProductSize) (uint, error) {
	if err := r.db.Create(storeProductSize).Error; err != nil {
		return 0, err
	}
	return storeProductSize.ID, nil
}

// UpdateProductSize updates an existing store product size
func (r *storeProductRepository) UpdateProductSize(storeID, productSizeID uint, size *data.StoreProductSize) error {
	result := r.db.Model(&data.StoreProductSize{}).
		Where("store_id = ? AND id = ?", storeID, productSizeID).
		Updates(size)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteStoreProductSize deletes a store product size
func (r *storeProductRepository) DeleteStoreProductSize(storeID, productSizeID uint) error {
	result := r.db.
		Where("store_id = ? AND id = ?", storeID, productSizeID).
		Delete(&data.StoreProductSize{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
