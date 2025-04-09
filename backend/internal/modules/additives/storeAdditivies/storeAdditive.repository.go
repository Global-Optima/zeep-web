package storeAdditives

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies/types"
	additiveTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StoreAdditiveRepository interface {
	GetMissingStoreAdditiveIDsForProductSizes(storeID uint, productSizeIDs []uint) ([]uint, error)
	FilterMissingStoreAdditiveIDs(storeID uint, additivesIDs []uint) ([]uint, error)
	CreateStoreAdditive(storeID uint, dto *types.CreateStoreAdditiveDTO) (uint, error)
	CreateStoreAdditives(storeAdditives []data.StoreAdditive) ([]uint, error)
	GetStoreAdditiveWithDetailsByID(storeAdditiveID uint, filter *contexts.StoreContextFilter) (*data.StoreAdditive, error)
	GetStoreAdditiveByID(storeAdditiveID uint) (*data.StoreAdditive, error)
	GetSufficientStoreAdditiveByID(storeID, storeAdditiveID uint, frozenStock map[uint]float64) (*data.StoreAdditive, error)
	GetAvailableAdditivesToAdd(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.Additive, error)
	GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.StoreAdditive, error)
	GetStoreAdditivesByIDs(storeID uint, IDs []uint) ([]data.StoreAdditive, error)
	GetStoreAdditiveCategories(storeID, storeProductSizeID uint, filter *types.StoreAdditiveCategoriesFilter) ([]data.AdditiveCategory, error)
	UpdateStoreAdditive(storeID, storeAdditiveID uint, input *data.StoreAdditive) error
	DeleteStoreAdditive(storeID, storeAdditiveID uint) error
	GetStoreAdditiveWithProductSizeAdditive(
		storeID uint,
		storeProductSizeID uint,
		storeAdditiveID uint,
	) (*data.StoreAdditive, *data.ProductSizeAdditive, error)

	CloneWithTransaction(tx *gorm.DB) StoreAdditiveRepository
}

type storeAdditiveRepository struct {
	db *gorm.DB
}

func NewStoreAdditiveRepository(db *gorm.DB) StoreAdditiveRepository {
	return &storeAdditiveRepository{db: db}
}

func (r *storeAdditiveRepository) CloneWithTransaction(tx *gorm.DB) StoreAdditiveRepository {
	return &storeAdditiveRepository{db: tx}
}

func (r *storeAdditiveRepository) GetMissingStoreAdditiveIDsForProductSizes(
	storeID uint,
	productSizeIDs []uint,
) ([]uint, error) {
	if len(productSizeIDs) == 0 {
		return nil, nil
	}

	var allPSAdditiveIDs []uint
	if err := r.db.Model(&data.ProductSizeAdditive{}).
		Distinct().
		Where("product_size_id IN ?", productSizeIDs).
		Pluck("additive_id", &allPSAdditiveIDs).Error; err != nil {
		return nil, err
	}

	if len(allPSAdditiveIDs) == 0 {
		return nil, nil
	}

	var existingStoreAdditiveIDs []uint
	if err := r.db.Model(&data.StoreAdditive{}).
		Distinct().
		Where("store_id = ?", storeID).
		Where("additive_id IN ?", allPSAdditiveIDs).
		Pluck("additive_id", &existingStoreAdditiveIDs).Error; err != nil {
		return nil, err
	}

	storeAdditiveMap := make(map[uint]struct{}, len(existingStoreAdditiveIDs))
	for _, id := range existingStoreAdditiveIDs {
		storeAdditiveMap[id] = struct{}{}
	}

	var missingAdditiveIDs []uint
	for _, id := range allPSAdditiveIDs {
		if _, found := storeAdditiveMap[id]; !found {
			missingAdditiveIDs = append(missingAdditiveIDs, id)
		}
	}

	return missingAdditiveIDs, nil
}

func (r *storeAdditiveRepository) FilterMissingStoreAdditiveIDs(storeID uint, additivesIDs []uint) ([]uint, error) {
	if len(additivesIDs) == 0 {
		return []uint{}, nil
	}

	var existingAdditivesIDs []uint
	if err := r.db.
		Model(&data.StoreAdditive{}).
		Where("store_id = ? AND additive_id IN (?)", storeID, additivesIDs).
		Pluck("additive_id", &existingAdditivesIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch existing ingredient IDs: %w", err)
	}

	existingMap := make(map[uint]struct{}, len(existingAdditivesIDs))
	for _, id := range existingAdditivesIDs {
		existingMap[id] = struct{}{}
	}

	var missingIngredientIDs []uint
	for _, id := range additivesIDs {
		if _, found := existingMap[id]; !found {
			missingIngredientIDs = append(missingIngredientIDs, id)
		}
	}

	return missingIngredientIDs, nil
}

func (r *storeAdditiveRepository) CreateStoreAdditive(storeID uint, dto *types.CreateStoreAdditiveDTO) (uint, error) {
	var existingStock data.StoreStock
	err := r.db.
		Where("store_id = ? AND additive_id = ?", storeID, dto.AdditiveID).
		First(&existingStock).Error
	if err == nil {
		return 0, fmt.Errorf("%w: store additive with additive ID %d already exists for store ID %d",
			moduleErrors.ErrAlreadyExists, dto.AdditiveID, storeID)
	}

	storeAdditive := types.CreateToStoreAdditive(dto, storeID)

	if err := r.db.Create(storeAdditive).Error; err != nil {
		return 0, err
	}

	return storeAdditive.ID, nil
}

func (r *storeAdditiveRepository) CreateStoreAdditives(storeAdditives []data.StoreAdditive) ([]uint, error) {
	if len(storeAdditives) == 0 {
		return nil, nil
	}

	if err := r.db.Create(&storeAdditives).Error; err != nil {
		return nil, err
	}

	var ids []uint
	for _, sa := range storeAdditives {
		ids = append(ids, sa.ID)
	}

	return ids, nil
}

func (r *storeAdditiveRepository) GetAvailableAdditivesToAdd(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.Additive, error) {
	var additives []data.Additive

	query := r.db.
		Preload("Category").
		Preload("Unit").
		Joins("JOIN additive_categories ON additives.additive_category_id = additive_categories.id")

	var err error

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("additives.name ILIKE ? OR additives.description ILIKE ?", searchTerm, searchTerm)
	}

	if filter.MinPrice != nil {
		query = query.Where("additives.base_price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("additives.base_price <= ?", *filter.MaxPrice)
	}

	if filter.CategoryID != nil {
		query = query.Where("additive_categories.id = ?", *filter.CategoryID)
	}

	if filter.ProductSizeID != nil {
		query = query.Where("product_additives.product_size_id = ?", *filter.ProductSizeID)
	}

	query = query.Where("additives.id NOT IN (?)", r.db.Model(&data.StoreAdditive{}).Select("additive_id").Where("store_id = ?", storeID))

	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Additive{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	if err := query.Find(&additives).Error; err != nil {
		return nil, err
	}

	return additives, nil
}

func (r *storeAdditiveRepository) GetStoreAdditiveCategories(
	storeID, storeProductSizeID uint,
	filter *types.StoreAdditiveCategoriesFilter,
) ([]data.AdditiveCategory, error) {
	var categories []data.AdditiveCategory

	// First, get the relevant product_size_id to limit the join scope
	var productSizeID uint
	err := r.db.Model(&data.StoreProductSize{}).
		Select("product_size_id").
		Where("id = ? AND EXISTS (SELECT 1 FROM store_products WHERE store_products.id = store_product_sizes.store_product_id AND store_products.store_id = ?)",
			storeProductSizeID, storeID).
		Limit(1).
		Pluck("product_size_id", &productSizeID).Error
	if err != nil {
		return nil, err
	}
	if productSizeID == 0 {
		return nil, moduleErrors.ErrNotFound
	}

	// Build the main query for categories
	query := r.db.Model(&data.AdditiveCategory{}).
		Select("DISTINCT additive_categories.*").
		Joins("INNER JOIN additives ON additives.additive_category_id = additive_categories.id").
		Joins("INNER JOIN store_additives ON store_additives.additive_id = additives.id AND store_additives.store_id = ?", storeID).
		Joins("INNER JOIN product_size_additives ON product_size_additives.additive_id = additives.id AND product_size_additives.product_size_id = ?", productSizeID).
		Where("additive_categories.deleted_at IS NULL").
		Where("store_additives.deleted_at IS NULL")

	// Apply filters
	if filter.IsMultipleSelect != nil {
		query = query.Where("additive_categories.is_multiple_select = ?", *filter.IsMultipleSelect)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("(additive_categories.name ILIKE ? OR additive_categories.description ILIKE ?)", searchTerm, searchTerm)
	}

	// Fetch categories first
	err = query.Find(&categories).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, moduleErrors.ErrNotFound
		}
		return nil, err
	}

	if len(categories) == 0 {
		return categories, nil
	}

	// Extract category IDs
	var categoryIDs []uint
	for _, cat := range categories {
		categoryIDs = append(categoryIDs, cat.ID)
	}

	// Get additives for these categories with a separate query
	var additives []data.Additive
	err = r.db.Model(&data.Additive{}).
		Joins("INNER JOIN store_additives ON store_additives.additive_id = additives.id AND store_additives.store_id = ? AND store_additives.deleted_at IS NULL", storeID).
		Preload("Unit").
		Preload("Category").
		Preload("StoreAdditives", "store_id = ? AND deleted_at IS NULL", storeID).
		Preload("ProductSizeAdditives", "product_size_id = ?", productSizeID).
		Where("additives.additive_category_id IN ?", categoryIDs).
		Where("additives.deleted_at IS NULL").
		Find(&additives).Error
	if err != nil {
		return nil, err
	}

	// Organize additives by category
	additivesByCategoryID := make(map[uint][]data.Additive)
	for _, additive := range additives {
		additivesByCategoryID[additive.AdditiveCategoryID] = append(additivesByCategoryID[additive.AdditiveCategoryID], additive)
	}

	// Assign additives to their categories
	for i := range categories {
		if additiveList, ok := additivesByCategoryID[categories[i].ID]; ok {
			categories[i].Additives = additiveList
		}
	}

	return categories, nil
}

func (r *storeAdditiveRepository) GetStoreAdditiveWithDetailsByID(storeAdditiveID uint, filter *contexts.StoreContextFilter) (*data.StoreAdditive, error) {
	var storeAdditive *data.StoreAdditive

	query := r.db.Model(&data.StoreAdditive{}).
		Where(&data.StoreAdditive{BaseEntity: data.BaseEntity{ID: storeAdditiveID}}).
		Preload("Additive").
		Preload("Additive.Category").
		Preload("Additive.Unit").
		Preload("Additive.Ingredients.Ingredient.Unit").
		Preload("Additive.Ingredients.Ingredient.IngredientCategory")

	if filter != nil {
		if filter.StoreID != nil {
			query.Where(&data.StoreAdditive{StoreID: *filter.StoreID})
		}

		if filter.FranchiseeID != nil {
			query.Joins("JOIN stores ON stores.id = store_additives.store_id").
				Where(&data.StoreAdditive{
					Store: data.Store{
						FranchiseeID: filter.FranchiseeID,
					},
				})
		}
	}

	err := query.First(&storeAdditive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreAdditiveNotFound
		}
		return nil, err
	}

	return storeAdditive, nil
}

func (r *storeAdditiveRepository) GetStoreAdditiveByID(storeAdditiveID uint) (*data.StoreAdditive, error) {
	var storeAdditive *data.StoreAdditive

	err := r.db.Model(&data.StoreAdditive{}).
		Where(&data.StoreAdditive{BaseEntity: data.BaseEntity{ID: storeAdditiveID}}).
		First(&storeAdditive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreAdditiveNotFound
		}
		return nil, err
	}

	return storeAdditive, nil
}

func (r *storeAdditiveRepository) GetStoreAdditiveWithProductSizeAdditive(
	storeID uint,
	storeProductSizeId uint,
	storeAdditiveId uint,
) (*data.StoreAdditive, *data.ProductSizeAdditive, error) {
	var sa data.StoreAdditive
	err := r.db.
		Model(&data.StoreAdditive{}).
		Preload("Additive").
		Where("id = ? AND store_id = ?", storeAdditiveId, storeID).
		First(&sa).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("storeAdditive not found (storeID=%d, storeAdditiveID=%d)", storeID, storeAdditiveId)
		}
		return nil, nil, fmt.Errorf("failed to load StoreAdditive: %w", err)
	}

	var sps data.StoreProductSize
	err = r.db.
		Model(&data.StoreProductSize{}).
		Preload("ProductSize").
		Where("id = ?", storeProductSizeId).
		First(&sps).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &sa, nil, fmt.Errorf("storeProductSize not found (id=%d)", storeProductSizeId)
		}
		return &sa, nil, fmt.Errorf("failed to load StoreProductSize: %w", err)
	}

	productSizeId := sps.ProductSizeID

	var psa data.ProductSizeAdditive
	err = r.db.
		Model(&data.ProductSizeAdditive{}).
		Where("product_size_id = ? AND additive_id = ?", productSizeId, sa.AdditiveID).
		Preload("Additive").
		First(&psa).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &sa, nil, fmt.Errorf(
				"productSizeAdditive not found (productSizeId=%d, additiveId=%d)",
				productSizeId, sa.AdditiveID,
			)
		}
		return &sa, nil, fmt.Errorf("failed to load ProductSizeAdditive: %w", err)
	}

	return &sa, &psa, nil
}

func (r *storeAdditiveRepository) GetSufficientStoreAdditiveByID(
	storeID, storeAdditiveID uint, frozenMap map[uint]float64,
) (*data.StoreAdditive, error) {
	var storeAdditive data.StoreAdditive

	err := r.db.
		Where("store_id = ?", storeID).
		Where("id = ?", storeAdditiveID).
		Preload("Additive").
		Preload("Additive.Category").
		Preload("Additive.Unit").
		Preload("Additive.Ingredients").
		Preload("Additive.Ingredients.Ingredient.Unit").
		Preload("Additive.Ingredients.Ingredient.IngredientCategory").
		First(&storeAdditive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStoreAdditiveNotFound
		}
		return nil, fmt.Errorf("%w: failed to get storeAdditive (id=%d): %w",
			types.ErrFailedToFetchStoreAdditive, storeAdditiveID, err)
	}

	for _, ingrUsage := range storeAdditive.Additive.Ingredients {
		requiredAmount := ingrUsage.Quantity

		var stock data.StoreStock
		err := r.db.
			Where("store_id = ? AND ingredient_id = ?", storeID, ingrUsage.IngredientID).
			First(&stock).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, types.ErrStoreStockNotFound
			}
			return nil, fmt.Errorf("%w: failed to get stock for ingredient %d: %w",
				types.ErrFailedToFetchStoreStock, ingrUsage.IngredientID, err)
		}

		if stock.Quantity < frozenMap[ingrUsage.IngredientID] {
			return nil, fmt.Errorf("%w: insufficient stock for ingredient %q (ID=%d): already pending %.2f, need %.2f, have left %.2f",
				types.ErrInsufficientStock, ingrUsage.Ingredient.Name, ingrUsage.IngredientID,
				frozenMap[ingrUsage.IngredientID], requiredAmount, stock.Quantity)
		}

		effectiveAvailable := stock.Quantity - frozenMap[ingrUsage.IngredientID]
		if effectiveAvailable < requiredAmount {
			return nil, fmt.Errorf("%w: insufficient effective stock for ingredient %q (ID=%d): need %.2f, have %.2f",
				types.ErrInsufficientStock, ingrUsage.Ingredient.Name, ingrUsage.IngredientID,
				requiredAmount, effectiveAvailable)
		}
	}

	return &storeAdditive, nil
}

func (r *storeAdditiveRepository) GetStoreAdditives(storeID uint, filter *additiveTypes.AdditiveFilterQuery) ([]data.StoreAdditive, error) {
	var storeAdditives []data.StoreAdditive

	query := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ?", storeID).
		Joins("JOIN additives ON additives.id = store_additives.additive_id").
		Preload("Additive.Category").
		Preload("Additive.Unit")

	var err error

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("additives.name ILIKE ? OR additives.description ILIKE ?", searchTerm, searchTerm)
	}

	if filter.MinPrice != nil {
		query = query.Where("store_additives.price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("store_additives.price <= ?", *filter.MaxPrice)
	}

	if filter.CategoryID != nil {
		query = query.Where("additive_categories.id = ?", *filter.CategoryID)
	}
	if filter.ProductSizeID != nil {
		query = query.Where("product_additives.product_size_id = ?", *filter.ProductSizeID)
	}

	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreAdditive{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&storeAdditives).Error; err != nil {
		return nil, err
	}

	return storeAdditives, nil
}

func (r *storeAdditiveRepository) GetStoreAdditivesByIDs(storeID uint, IDs []uint) ([]data.StoreAdditive, error) {
	var storeAdditives []data.StoreAdditive

	query := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ?", storeID).
		Preload("Additive.Category").
		Preload("Additive.Unit").
		Where("id IN (?)", IDs)

	if err := query.Find(&storeAdditives).Error; err != nil {
		return nil, err
	}

	return storeAdditives, nil
}

func (r *storeAdditiveRepository) UpdateStoreAdditive(storeID, storeAdditiveID uint, input *data.StoreAdditive) error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}

	res := r.db.Model(&data.StoreAdditive{}).
		Where("store_id = ? AND id = ?", storeID, storeAdditiveID).
		Updates(input)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return types.ErrStoreAdditiveNotFound
	}

	return nil
}

func (r *storeAdditiveRepository) DeleteStoreAdditive(storeID, storeAdditiveID uint) error {
	var exists bool
	err := r.db.
		Table("product_size_additives AS psa").
		Select("1").
		Joins("JOIN store_additives sa ON psa.additive_id = sa.additive_id").
		Joins("JOIN store_product_sizes sps ON psa.product_size_id = sps.product_size_id").
		Joins("JOIN store_products sp ON sps.store_product_id = sp.id").
		Where("sa.id = ? AND sa.store_id = ?", storeAdditiveID, storeID).
		Where("sp.store_id = ?", storeID).
		Where("psa.deleted_at IS NULL").
		Where("sa.deleted_at IS NULL").
		Where("sps.deleted_at IS NULL").
		Where("sp.deleted_at IS NULL").
		Limit(1).
		Scan(&exists).Error
	if err != nil {
		return errors.Wrap(err, "failed to check store additive usage")
	}

	if exists {
		return types.ErrStoreAdditiveInUse
	}

	// Proceed with deletion if not in use
	result := r.db.Where("store_id = ? AND id = ?", storeID, storeAdditiveID).Delete(&data.StoreAdditive{})
	if result.Error != nil {
		return errors.Wrap(result.Error, "failed to delete store additive")
	}

	if result.RowsAffected == 0 {
		return types.ErrStoreAdditiveNotFound
	}

	return nil
}
