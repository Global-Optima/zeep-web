package storeStocks

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreStockRepository interface {
	GetAvailableIngredientsToAdd(storeID uint, filter *ingredientTypes.IngredientFilter) ([]data.Ingredient, error)
	AddStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error)
	AddMultipleStocks(stocks []data.StoreStock) ([]uint, error)
	AddOrSaveStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error)
	GetStockList(storeID uint, query *types.GetStockFilterQuery) ([]data.StoreStock, error)
	GetStockListByIDs(storeID uint, IDs []uint) ([]data.StoreStock, error)
	GetRawStockByID(storeID, stockID uint) (*data.StoreStock, error)
	GetStockById(stockID uint, filter *contexts.StoreContextFilter) (*data.StoreStock, error)
	GetAllStockList(storeID uint) ([]data.StoreStock, error)
	SaveStock(storeStock *data.StoreStock) error
	DeleteStockById(storeID, stockID uint) error
	WithTransaction(txFunc func(txRepo StoreStockRepository) error) error
	CloneWithTransaction(tx *gorm.DB) StoreStockRepository

	FindEarliestExpirationForIngredient(ingredientID, storeID uint) (*time.Time, error)
	FilterMissingIngredientsIDs(storeID uint, ingredientsIDs []uint) ([]uint, error)
	CheckStoreStockUsage(stockID uint) (bool, error)
}

type storeStockRepository struct {
	db *gorm.DB
}

func NewStoreStockRepository(db *gorm.DB) StoreStockRepository {
	return &storeStockRepository{db: db}
}

func (r *storeStockRepository) WithTransaction(txFunc func(txRepo StoreStockRepository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	txRepo := r.CloneWithTransaction(tx)

	if err := txFunc(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *storeStockRepository) CloneWithTransaction(tx *gorm.DB) StoreStockRepository {
	return &storeStockRepository{
		db: tx,
	}
}

func (r *storeStockRepository) GetAvailableIngredientsToAdd(storeID uint, filter *ingredientTypes.IngredientFilter) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	query := r.db.Model(&data.Ingredient{}).
		Preload("Unit").
		Preload("IngredientCategory")

	if filter.ProductSizeID != nil {
		query = query.Joins("JOIN product_size_ingredients psi ON psi.ingredient_id = ingredients.id").
			Where("psi.product_size_id = ?", *filter.ProductSizeID)
	}

	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.MinCalories != nil {
		query = query.Where("calories >= ?", *filter.MinCalories)
	}
	if filter.MaxCalories != nil {
		query = query.Where("calories <= ?", *filter.MaxCalories)
	}

	query = query.Where("ingredients.id NOT IN (?)", r.db.Model(&data.StoreStock{}).
		Select("ingredient_id").
		Where("store_id = ?", storeID))

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Ingredient{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *storeStockRepository) AddStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error) {
	var existingStock data.StoreStock
	err := r.db.
		Where("store_id = ? AND ingredient_id = ?", storeID, dto.IngredientID).
		First(&existingStock).Error
	if err == nil {
		return 0, fmt.Errorf("%w: stock with ingredient ID %d already exists for store ID %d",
			moduleErrors.ErrAlreadyExists, dto.IngredientID, storeID)
	}

	var ingredient data.Ingredient
	err = r.db.
		Where("id = ?", dto.IngredientID).
		First(&ingredient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("ingredient not found for ingredient ID %d", dto.IngredientID)
		}
		return 0, fmt.Errorf("failed to fetch ingredient: %w", err)
	}

	storeStock := types.AddToStock(*dto, storeID)
	err = r.db.Create(&storeStock).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create store stock: %w", err)
	}

	return storeStock.ID, nil
}

func (r *storeStockRepository) AddMultipleStocks(stocks []data.StoreStock) ([]uint, error) {
	stockIDs := make([]uint, len(stocks))

	err := r.db.Create(&stocks).Error
	if err != nil {
		return nil, err
	}

	for i, stock := range stocks {
		stockIDs[i] = stock.ID
	}
	return stockIDs, nil
}

func (r *storeStockRepository) AddOrSaveStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error) {
	var existingStock data.StoreStock
	err := r.db.
		Where("store_id = ? AND ingredient_id = ?", storeID, dto.IngredientID).
		First(&existingStock).Error

	if err == nil {
		existingStock.Quantity += dto.Quantity
		existingStock.LowStockThreshold = dto.LowStockThreshold
		err = r.db.Save(&existingStock).Error
		if err != nil {
			return 0, fmt.Errorf("failed to update store stock: %w", err)
		}
		return existingStock.ID, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		storeStock := data.StoreStock{
			StoreID:           storeID,
			IngredientID:      dto.IngredientID,
			Quantity:          dto.Quantity,
			LowStockThreshold: dto.LowStockThreshold,
		}

		err = r.db.Create(&storeStock).Error
		if err != nil {
			return 0, fmt.Errorf("failed to create store stock: %w", err)
		}
		return storeStock.ID, nil
	}

	return 0, fmt.Errorf("failed to add or update stock: %w", err)
}

func (r *storeStockRepository) GetStockList(storeID uint, filter *types.GetStockFilterQuery) ([]data.StoreStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeStockList []data.StoreStock

	query := r.db.Model(&data.StoreStock{}).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Joins("JOIN ingredients ON ingredient_id = ingredients.id").
		Where("store_id = ?", storeID)

	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("store_stocks.quantity < store_stocks.low_stock_threshold")
	}

	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("ingredients.name ILIKE ?", "%"+*filter.Search+"%")
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreStock{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&storeStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeStockList, nil
}

func (r *storeStockRepository) GetAllStockList(storeID uint) ([]data.StoreStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeWarehouseStockList []data.StoreStock

	query := r.db.Model(&data.StoreStock{}).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Where("store_id = ?", storeID)

	err := query.Find(&storeWarehouseStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeWarehouseStockList, nil
}

func (r *storeStockRepository) GetStockListByIDs(storeID uint, stockIds []uint) ([]data.StoreStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeWarehouseStockList []data.StoreStock

	query := r.db.Model(&data.StoreStock{}).
		Where("store_id = ? AND id IN (?)", storeID, stockIds).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory")

	err := query.Find(&storeWarehouseStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeWarehouseStockList, nil
}

func (r *storeStockRepository) GetRawStockByID(storeID, stockID uint) (*data.StoreStock, error) {
	var StoreStock data.StoreStock

	err := r.db.Model(&data.StoreStock{}).
		Where("store_id = ? AND id = ?", storeID, stockID).
		First(&StoreStock).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStockNotFound
		}
		return nil, fmt.Errorf("failed to get store stock: %w", err)
	}

	return &StoreStock, nil
}

func (r *storeStockRepository) GetStockById(stockId uint, filter *contexts.StoreContextFilter) (*data.StoreStock, error) {
	var StoreStock data.StoreStock

	if stockId == 0 {
		return nil, fmt.Errorf("store stock id cannot be 0")
	}

	query := r.db.Model(&data.StoreStock{}).
		Where(&data.StoreStock{BaseEntity: data.BaseEntity{ID: stockId}}).
		Preload("Store").
		Preload("Store.Warehouse").
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory")

	if filter != nil {
		if filter.StoreID != nil {
			query.Where(&data.StoreStock{StoreID: *filter.StoreID})
		}

		if filter.FranchiseeID != nil {
			query.Joins("JOIN stores ON stores.id = store_stocks.store_id").
				Where(&data.StoreStock{
					Store: data.Store{
						FranchiseeID: filter.FranchiseeID,
					},
				})
		}
	}

	err := query.First(&StoreStock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrStockNotFound
		}
		return nil, err
	}

	return &StoreStock, nil
}

func (r *storeStockRepository) SaveStock(storeStock *data.StoreStock) error {
	if storeStock == nil || storeStock.ID == 0 {
		return fmt.Errorf("not enough data to save store stock")
	}

	updRes := r.db.Save(storeStock)

	if updRes.Error != nil {
		return utils.WrapError("failed to update store stock", updRes.Error)
	}

	if updRes.RowsAffected == 0 {
		return fmt.Errorf("update attempt had no changes for stockId=%d with storeId=%d", storeStock.ID, storeStock.StoreID)
	}

	return nil
}

func (r *storeStockRepository) DeleteStockById(storeId, stockId uint) error {
	res := r.db.
		Unscoped().
		Where("id = ? AND store_id = ?", stockId, storeId).
		Delete(&data.StoreStock{})

	if res.Error != nil {
		return fmt.Errorf("failed to delete store stock: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("stock not found")
	}

	return nil
}

func (r *storeStockRepository) CheckStoreStockUsage(stockID uint) (bool, error) {
	var isInUse bool

	err := r.db.
		Table("store_stocks AS ss").
		Select("1").
		Joins("JOIN ingredients AS i ON ss.ingredient_id = i.id").
		Joins("LEFT JOIN product_size_ingredients AS psi ON psi.ingredient_id = i.id").
		Joins("LEFT JOIN store_product_sizes AS sps ON psi.product_size_id = sps.product_size_id").
		Joins("LEFT JOIN store_products AS sp ON sps.store_product_id = sp.id").
		Joins("LEFT JOIN additive_ingredients AS ai ON ai.ingredient_id = i.id").
		Joins("LEFT JOIN store_additives AS sa ON ai.additive_id = sa.additive_id").
		Where("ss.id = ?", stockID).
		Where("ss.deleted_at IS NULL").
		Where("i.deleted_at IS NULL").
		Where("(psi.id IS NOT NULL OR ai.id IS NOT NULL)"). // Ensures ingredient is used
		Limit(1).
		Scan(&isInUse).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	return isInUse, nil
}

func (r *storeStockRepository) FindEarliestExpirationForIngredient(ingredientID, storeID uint) (*time.Time, error) {
	var earliestExpirationDate sql.NullTime

	err := r.db.Model(&data.StockRequestIngredient{}).
		Joins("JOIN stock_requests ON stock_requests.id = stock_request_ingredients.stock_request_id").
		Joins("JOIN stock_materials ON stock_materials.id = stock_request_ingredients.stock_material_id").
		Select("MIN(stock_request_ingredients.expiration_date) AS earliest_expiration_date").
		Where("stock_requests.store_id = ? AND stock_materials.ingredient_id = ?", storeID, ingredientID).
		Scan(&earliestExpirationDate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch earliest expiration date for ingredient ID %d: %w", ingredientID, err)
	}

	if !earliestExpirationDate.Valid {
		return nil, nil
	}

	utcTime := earliestExpirationDate.Time.UTC()
	return &utcTime, nil
}

func (r *storeStockRepository) FilterMissingIngredientsIDs(storeID uint, ingredientsIDs []uint) ([]uint, error) {
	if len(ingredientsIDs) == 0 {
		return []uint{}, nil
	}

	var existingIngredientIDs []uint
	if err := r.db.
		Model(&data.StoreStock{}).
		Where("store_id = ? AND ingredient_id IN (?)", storeID, ingredientsIDs).
		Pluck("ingredient_id", &existingIngredientIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch existing ingredient IDs: %w", err)
	}

	existingMap := make(map[uint]struct{}, len(existingIngredientIDs))
	for _, id := range existingIngredientIDs {
		existingMap[id] = struct{}{}
	}

	var missingIngredientIDs []uint
	for _, id := range ingredientsIDs {
		if _, found := existingMap[id]; !found {
			missingIngredientIDs = append(missingIngredientIDs, id)
		}
	}

	return missingIngredientIDs, nil
}
