package storeProducts

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts/types"
	productTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreProductRepository interface {
	GetStoreProductCategories(storeID uint) ([]data.ProductCategory, error)
	GetStoreProductById(storeProductID uint, filter *contexts.StoreContextFilter) (*data.StoreProduct, error)
	GetStoreProductsByStoreProductIDs(storeID uint, storeProductIDs []uint) ([]data.StoreProduct, error)
	GetStoreProducts(storeID uint, filter *types.StoreProductsFilterDTO) ([]data.StoreProduct, error)
	GetAvailableProductsToAdd(storeID uint, filter *productTypes.ProductsFilterDto) ([]data.Product, error)
	GetRecommendedStoreProducts(storeID uint, excludedStoreProductIDs []uint) ([]data.StoreProduct, error)
	CreateStoreProduct(product *data.StoreProduct) (uint, error)
	CreateMultipleStoreProducts(storeProducts []data.StoreProduct) ([]uint, error)
	UpdateStoreProductByID(storeID, storeProductID uint, updateModels *types.StoreProductModels) error
	DeleteStoreProductWithSizes(storeID, storeProductID uint) error

	GetStoreProductSizeById(storeID, storeProductSizeID uint) (*data.StoreProductSize, error)
	GetSufficientStoreProductSizeById(storeID, storeProductSizeID uint, frozenMap map[uint]float64) (*data.StoreProductSize, error)
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

func (r *storeProductRepository) GetStoreProductCategories(storeID uint) ([]data.ProductCategory, error) {
	var categories []data.ProductCategory

	err := r.db.Model(&data.ProductCategory{}).
		Joins("JOIN products ON products.category_id = product_categories.id").
		Joins("JOIN store_products ON store_products.product_id = products.id").
		Joins("JOIN store_product_sizes ON store_product_sizes.store_product_id = store_products.id").
		Where("store_products.store_id = ? AND store_products.is_available = ?", storeID, true).
		Group("product_categories.id").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *storeProductRepository) GetStoreProductById(storeProductID uint, filter *contexts.StoreContextFilter) (*data.StoreProduct, error) {
	var storeProduct data.StoreProduct
	query := r.db.Model(&data.StoreProduct{}).
		Where(&data.StoreProduct{BaseEntity: data.BaseEntity{ID: storeProductID}}).
		Preload("Product.ProductSizes").
		Preload("Product.Category").
		Preload("StoreProductSizes.ProductSize.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory")

	if filter != nil {
		if filter.StoreID != nil {
			query.Where(&data.StoreProduct{StoreID: *filter.StoreID})
		}

		if filter.FranchiseeID != nil {
			query.Joins("JOIN stores ON stores.id = store_products.store_id").
				Where(&data.StoreProduct{
					Store: data.Store{
						FranchiseeID: filter.FranchiseeID,
					},
				})
		}
	}

	err := query.First(&storeProduct).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, moduleErrors.ErrNotFound
		}
		return nil, err
	}

	return &storeProduct, nil
}

func (r *storeProductRepository) GetStoreProductsByStoreProductIDs(storeID uint, storeProductIDs []uint) ([]data.StoreProduct, error) {
	var storeProducts []data.StoreProduct
	query := r.db.Model(&data.StoreProduct{}).
		Where("store_products.store_id = ? AND store_products.id IN (?)", storeID, storeProductIDs).
		Joins("JOIN products ON store_products.product_id = products.id").
		Preload("Product.ProductSizes").
		Preload("Product.Category").
		Preload("StoreProductSizes.ProductSize.Unit").
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
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory")

	if filter != nil {
		if filter.Search != nil {
			query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
		}
		if filter.IsAvailable != nil {
			query = query.Where("is_available = ?", *filter.IsAvailable).
				Where("EXISTS (SELECT 1 FROM store_product_sizes sps WHERE sps.store_product_id = store_products.id)")
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

func (r *storeProductRepository) GetRecommendedStoreProducts(storeID uint, excludedStoreProductIDs []uint) ([]data.StoreProduct, error) {
	var storeProducts []data.StoreProduct

	query := r.db.Model(&data.StoreProduct{}).
		Where("store_id = ?", storeID).
		Where("id NOT IN (?)", excludedStoreProductIDs).
		Order("RANDOM()").
		Limit(5).
		Preload("Product.ProductSizes").
		Preload("Product.Category").
		Preload("StoreProductSizes.ProductSize.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("StoreProductSizes.ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory")

	if err := query.Find(&storeProducts).Error; err != nil {
		return nil, err
	}

	return storeProducts, nil
}

func (r *storeProductRepository) GetAvailableProductsToAdd(storeID uint, filter *productTypes.ProductsFilterDto) ([]data.Product, error) {
	var products []data.Product

	query := r.db.
		Model(&data.Product{}).
		Preload("Category").
		Preload("ProductSizes.Unit").
		Where("EXISTS (SELECT 1 FROM product_sizes WHERE product_sizes.product_id = products.id)")

	if filter.CategoryID != nil {
		query = query.Where("products.category_id = ?", *filter.CategoryID)
	}

	if filter.Search != nil {
		searchPattern := "%" + *filter.Search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", searchPattern, searchPattern)
	}

	query = query.Where("products.id NOT IN (?)", r.db.Model(&data.StoreProduct{}).Select("product_id").Where("store_id = ?", storeID))

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Product{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&products).Error
	if err != nil {
		return nil, err
	}

	if products == nil {
		products = []data.Product{}
	}

	return products, nil
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
		Preload("ProductSize.Product").
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

func (r *storeProductRepository) GetSufficientStoreProductSizeById(
	storeID, storeProductSizeID uint, frozenMap map[uint]float64,
) (*data.StoreProductSize, error) {
	var storeProductSize data.StoreProductSize

	err := r.db.Model(&data.StoreProductSize{}).
		Joins("JOIN store_products sp ON store_product_sizes.store_product_id = sp.id").
		Where("sp.store_id = ? AND store_product_sizes.id = ?", storeID, storeProductSizeID).
		Preload("StoreProduct.Store").
		Preload("StoreProduct.Product").
		Preload("ProductSize").
		Preload("ProductSize.Unit").
		Preload("ProductSize.Product").
		Preload("ProductSize.ProductSizeIngredients").
		Preload("ProductSize.ProductSizeIngredients.Ingredient.Unit").
		Preload("ProductSize.ProductSizeIngredients.Ingredient.IngredientCategory").
		Preload("ProductSize.Additives.Additive.Category").
		Preload("ProductSize.Additives.Additive.Unit").
		First(&storeProductSize).Error

	if err != nil {
		return nil, fmt.Errorf("failed to load StoreProductSize (ID=%d): %w", storeProductSizeID, err)
	}

	for _, usage := range storeProductSize.ProductSize.ProductSizeIngredients {
		needed := usage.Quantity

		var stock data.StoreStock
		err := r.db.Model(&data.StoreStock{}).
			Where("store_id = ?", storeID).
			Where("ingredient_id = ?", usage.IngredientID).
			First(&stock).Error
		if err != nil {
			return nil, fmt.Errorf("failed to find stock for ingredient %d: %w",
				usage.IngredientID, err)
		}

		if stock.Quantity < frozenMap[usage.IngredientID] {
			return nil, fmt.Errorf(
				"insufficient stock for ingredient %q (ID=%d): already pending %.2f, need %.2f, have left %.2f",
				usage.Ingredient.Name, usage.IngredientID, frozenMap[usage.IngredientID], needed, stock.Quantity,
			)
		}

		effectiveAvailable := stock.Quantity - frozenMap[usage.IngredientID]
		if effectiveAvailable < needed {
			return nil, fmt.Errorf(
				"insufficient stock for ingredient %q (ID=%d): need %.2f, have effective available %.2f",
				usage.Ingredient.Name, usage.IngredientID, needed, effectiveAvailable,
			)
		}
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
