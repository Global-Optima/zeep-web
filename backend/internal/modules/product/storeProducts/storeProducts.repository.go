package storeProducts

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreProductRepository interface {
	GetStoreProductById(storeID uint, storeProductID uint) (*data.StoreProduct, error)
	GetStoreProductsByProductIDs(storeID uint, productIDs []uint) ([]data.StoreProduct, error)
	GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]data.StoreProduct, error)
	CreateStoreProduct(product *data.StoreProduct) (uint, error)
	CreateMultipleStoreProducts(storeProducts []data.StoreProduct) ([]uint, error)
	UpdateStoreProductByID(storeID, storeProductID uint, updateModels *types.StoreProductModels) error
	DeleteStoreProductWithSizes(storeID, storeProductID uint) error

	GetStoreProductSizeById(storeID, storeProductSizeID uint) (*data.StoreProductSize, error)
	UpdateProductSize(storeID, productSizeID uint, size *data.StoreProductSize) error
	DeleteStoreProductSize(storeID, productSizeID uint) error
	CloneWithTransaction(tx *gorm.DB) StoreProductRepository
}

type storeProductRepository struct {
	db *gorm.DB
}

func NewStoreProductRepository(db *gorm.DB) StoreProductRepository {
	return &storeProductRepository{db: db}
}

func (r *storeProductRepository) CloneWithTransaction(tx *gorm.DB) StoreProductRepository {
	return &storeProductRepository{
		db: tx,
	}
}

func (r *storeProductRepository) GetStoreProductById(storeID uint, storeProductID uint) (*data.StoreProduct, error) {
	var storeProduct data.StoreProduct
	err := r.db.Model(&data.StoreProduct{}).
		Where("store_id = ? AND id = ?", storeID, storeProductID).
		Preload("Product.ProductSizes").
		Preload("StoreProductSizes.ProductSize.Unit").
		Preload("Product.Category").
		Preload("StoreProductSizes.ProductSize.Additives.Additive.Category").
		Preload("StoreProductSizes.ProductSize.Additives.Additive.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory").
		First(&storeProduct).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreProductNotFound
		}
		return nil, err
	}

	return &storeProduct, nil
}

func (r *storeProductRepository) GetStoreProductsByProductIDs(storeID uint, productIDs []uint) ([]data.StoreProduct, error) {
	var storeProducts []data.StoreProduct
	query := r.db.Model(&data.StoreProduct{}).
		Where("store_id = ? AND product_id IN (?)", storeID, productIDs).
		Joins("JOIN products ON store_products.product_id = products.id").
		Preload("Product.ProductSizes").
		Preload("Product.Category").
		Preload("StoreProductSizes.ProductSize.Unit").
		Preload("StoreProductSizes.ProductSize.Additives.Additive.Category").
		Preload("StoreProductSizes.ProductSize.Additives.Additive.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory")

	if err := query.Find(&storeProducts).Error; err != nil {
		return nil, err
	}

	return storeProducts, nil
}

func (r *storeProductRepository) GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]data.StoreProduct, error) {
	var storeProducts []data.StoreProduct
	query := r.db.Model(&data.StoreProduct{}).
		Where("store_id = ?", storeID).
		Joins("JOIN products ON store_products.product_id = products.id").
		Preload("Product.ProductSizes").
		Preload("Product.Category").
		Preload("StoreProductSizes.ProductSize.Unit").
		Preload("StoreProductSizes.ProductSize.Additives.Additive.Category").
		Preload("StoreProductSizes.ProductSize.Additives.Additive.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory")

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
		if filter.MinPrice != nil {
			query = query.Where("store_product_sizes.price >= ", *filter.MinPrice)
		}
		if filter.MaxPrice != nil {
			query = query.Where("store_product_sizes.price >= ", *filter.MaxPrice)
		}
	}

	if filter == nil {
		return nil, fmt.Errorf("filter is not binded")
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreProduct{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&storeProducts).Error; err != nil {
		return nil, err
	}

	return storeProducts, nil
}

func (r *storeProductRepository) CreateStoreProduct(product *data.StoreProduct) (uint, error) {
	err := r.db.Create(product).Error
	if err != nil {
		return 0, err
	}

	return product.ID, nil
}

func (r *storeProductRepository) CreateMultipleStoreProducts(storeProducts []data.StoreProduct) ([]uint, error) {
	if err := r.db.Create(&storeProducts).Error; err != nil {
		return nil, err
	}

	var ids []uint
	for _, storeProduct := range storeProducts {
		ids = append(ids, storeProduct.ID)
	}

	return ids, nil
}

func (r *storeProductRepository) UpdateStoreProductByID(storeID, storeProductID uint, updateModels *types.StoreProductModels) error {
	if updateModels == nil {
		return types.ErrNoUpdateContext
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&data.StoreProduct{}).
			Where("store_id = ? AND id = ?", storeID, storeProductID).
			Update("is_available", &updateModels.StoreProduct.IsAvailable).Error

		if err != nil {
			return err
		}

		if len(updateModels.StoreProductSizes) > 0 {
			if err := tx.Where("store_product_id = ?", storeProductID).
				Delete(&data.StoreProductSize{}).Error; err != nil {
				return err
			}

			for i := range updateModels.StoreProductSizes {
				updateModels.StoreProductSizes[i].StoreProductID = storeProductID
			}

			if err := tx.Create(&updateModels.StoreProductSizes).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (r *storeProductRepository) DeleteStoreProductWithSizes(storeID, storeProductID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var storeProduct data.StoreProduct

		if err := tx.Where("id = ? AND store_id = ?", storeProductID, storeID).
			First(&storeProduct).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: no storeProduct %d in store %d", types.ErrStoreProductNotFound, storeProductID, storeID)
			}
			return err
		}

		if err := tx.Where("store_product_id = ?", storeProductID).
			Delete(&data.StoreProductSize{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&data.StoreProduct{}, storeProductID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *storeProductRepository) GetStoreProductSizeById(storeID, storeProductSizeID uint) (*data.StoreProductSize, error) {
	var storeProductSize data.StoreProductSize

	err := r.db.Model(&data.StoreProductSize{}).
		Joins("JOIN store_products sp ON store_product_sizes.store_product_id = sp.id").
		Where("sp.store_id = ? AND store_product_sizes.id = ?", storeID, storeProductSizeID).
		Preload("ProductSize.Unit").
		Preload("ProductSize.Additives.Additive.Category").
		Preload("ProductSize.Additives.Additive.Unit").
		Preload("ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory").
		First(&storeProductSize).Error

	if err != nil {
		return &storeProductSize, err
	}
	return &storeProductSize, nil
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
