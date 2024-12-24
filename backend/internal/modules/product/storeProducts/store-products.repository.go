package storeProducts

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreProductRepository interface {
	GetStoreProductById(storeID uint, storeProductID uint) (*types.StoreProductDetailsDTO, error)
	GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]data.StoreProduct, error)
	CreateStoreProduct(storeProduct *data.StoreProduct) (uint, error)
	CreateMultipleStoreProducts(storeProducts []data.StoreProduct) ([]uint, error)
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

	dto := types.MapToStoreProductDetailsDTO(&storeProduct)
	return &dto, nil
}

func (r *storeProductRepository) GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]data.StoreProduct, error) {
	var storeProducts []data.StoreProduct
	query := r.db.Model(&data.StoreProduct{}).Where("store_id = ?", storeID).
		Preload("Store").
		Preload("Store.ProductSizes").
		Preload("Product").
		Preload("Store.ProductSizes", "store_id = ?", storeID).
		Preload("Store.ProductSizes.ProductSize", "is_default = TRUE")

	if filter != nil {
		if filter.Search != nil {
			query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
		}
		if filter.IsAvailable != nil {
			query = query.Where("is_available = ?", *filter.IsAvailable)
		}
		if filter.CategoryID != nil {
			query = query.Where("products.category_id = ?", *filter.CategoryID)
		}
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreProduct{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&storeProducts).Error; err != nil {
		return nil, err
	}

	return storeProducts, nil
}

func (r *storeProductRepository) CreateStoreProduct(storeProduct *data.StoreProduct) (uint, error) {
	if err := r.db.Create(storeProduct).Error; err != nil {
		return 0, err
	}
	return storeProduct.ID, nil
}

func (r *storeProductRepository) CreateMultipleStoreProducts(storeProducts []data.StoreProduct) ([]uint, error) {
	if err := r.db.Create(storeProducts).Error; err != nil {
		return nil, err
	}

	var ids []uint
	for _, storeProduct := range storeProducts {
		ids = append(ids, storeProduct.ID)
	}

	return ids, nil
}

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

func (r *storeProductRepository) CreateStoreProductSize(storeProductSize *data.StoreProductSize) (uint, error) {
	if err := r.db.Create(storeProductSize).Error; err != nil {
		return 0, err
	}
	return storeProductSize.ID, nil
}

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
