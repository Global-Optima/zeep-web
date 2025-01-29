package storeWarehouses

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreWarehouseRepository interface {
	AddStock(storeId uint, dto *types.AddStoreStockDTO) (uint, error)
	AddOrUpdateStock(storeId uint, dto *types.AddStoreStockDTO) (uint, error)
	GetStockList(storeId uint, query *types.GetStockFilterQuery) ([]data.StoreWarehouseStock, error)
	GetStockListByIDs(storeId uint, IDs []uint) ([]data.StoreWarehouseStock, error)
	GetStockById(storeId, stockId uint) (*data.StoreWarehouseStock, error)
	GetStockListForNotifications(storeID uint) ([]data.StoreWarehouseStock, error)
	UpdateStock(storeId, stockId uint, dto *types.UpdateStoreStockDTO) error
	DeleteStockById(storeId, stockId uint) error
	WithTransaction(txFunc func(txRepo storeWarehouseRepository) error) error
	CloneWithTransaction(tx *gorm.DB) storeWarehouseRepository
}

type storeWarehouseRepository struct {
	db *gorm.DB
}

func NewStoreWarehouseRepository(db *gorm.DB) StoreWarehouseRepository {
	return &storeWarehouseRepository{db: db}
}

func (r *storeWarehouseRepository) WithTransaction(txFunc func(txRepo storeWarehouseRepository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Clone the repository with the transaction instance
	txRepo := r.CloneWithTransaction(tx)

	// Execute the transaction logic
	if err := txFunc(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *storeWarehouseRepository) CloneWithTransaction(tx *gorm.DB) storeWarehouseRepository {
	return storeWarehouseRepository{
		db: tx,
	}
}

func (r *storeWarehouseRepository) AddOrUpdateStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error) {
	// Fetch the StoreWarehouse for the given storeID
	var storeWarehouse data.StoreWarehouse
	err := r.db.
		Where("store	_id = ?", storeID).
		First(&storeWarehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("store warehouse not found for store ID %d", storeID)
		}
		return 0, fmt.Errorf("failed to fetch store warehouse: %w", err)
	}

	// Fetch the existing stock (if any) for the given ingredient
	var existingStock data.StoreWarehouseStock
	err = r.db.
		Where("store_warehouse_id = ? AND ingredient_id = ?", storeWarehouse.ID, dto.IngredientID).
		First(&existingStock).Error

	if err == nil {
		// Update existing stock
		existingStock.Quantity += dto.Quantity
		existingStock.LowStockThreshold = dto.LowStockThreshold
		err = r.db.Save(&existingStock).Error
		if err != nil {
			return 0, fmt.Errorf("failed to update store warehouse stock: %w", err)
		}
		return existingStock.ID, nil
	}

	// If no existing stock, create a new stock entry
	if errors.Is(err, gorm.ErrRecordNotFound) {
		storeWarehouseStock := data.StoreWarehouseStock{
			StoreWarehouseID:  storeWarehouse.ID,
			IngredientID:      dto.IngredientID,
			Quantity:          dto.Quantity,
			LowStockThreshold: dto.LowStockThreshold,
		}

		err = r.db.Create(&storeWarehouseStock).Error
		if err != nil {
			return 0, fmt.Errorf("failed to create store warehouse stock: %w", err)
		}
		return storeWarehouseStock.ID, nil
	}

	// Handle unexpected errors
	return 0, fmt.Errorf("failed to add or update stock: %w", err)
}

func (r *storeWarehouseRepository) AddStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error) {
	// Fetch the StoreWarehouse for the given storeID
	storeWarehouse, err := r.getStoreWarehouseByStoreId(storeID)
	if err != nil {
		return 0, err
	}

	// Check if a stock with the given IngredientID already exists for the StoreWarehouse
	var existingStock data.StoreWarehouseStock
	err = r.db.
		Where("store_warehouse_id = ? AND ingredient_id = ?", storeWarehouse.ID, dto.IngredientID).
		First(&existingStock).Error
	if err == nil {
		return 0, fmt.Errorf("%w: stock with ingredient ID %d already exists for store warehouse ID %d",
			types.ErrStockAlreadyExists, dto.IngredientID, storeWarehouse.ID)
	}

	// Fetch the Ingredient for the given IngredientID from the DTO
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

	storeWarehouseStock := types.AddToStock(*dto, storeWarehouse.ID)

	err = r.db.Create(&storeWarehouseStock).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create store warehouse stock: %w", err)
	}

	return storeWarehouseStock.ID, nil
}

func (r *storeWarehouseRepository) getStoreWarehouseByStoreId(storeID uint) (*data.StoreWarehouse, error) {
	var storeWarehouse data.StoreWarehouse
	err := r.db.
		Where("store_id = ?", storeID).
		First(&storeWarehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("store warehouse not found for store ID %d", storeID)
		}
		return nil, fmt.Errorf("failed to fetch store warehouse: %w", err)
	}
	return &storeWarehouse, nil
}

func (r *storeWarehouseRepository) GetStockList(storeID uint, filter *types.GetStockFilterQuery) ([]data.StoreWarehouseStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeWarehouseStockList []data.StoreWarehouseStock

	query := r.db.Model(&data.StoreWarehouseStock{}).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Preload("StoreWarehouse").
		Joins("JOIN store_warehouses ON store_warehouse_stocks.store_warehouse_id = store_warehouses.id").
		Joins("JOIN ingredients ON store_warehouse_stocks.ingredient_id = ingredients.id").
		Where("store_warehouses.store_id = ?", storeID)

	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("store_warehouse_stocks.quantity < store_warehouse_stocks.low_stock_threshold")
	}

	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("ingredients.name ILIKE ?", "%"+*filter.Search+"%")
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreWarehouseStock{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&storeWarehouseStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeWarehouseStockList, nil
}

func (r *storeWarehouseRepository) GetStockListForNotifications(storeID uint) ([]data.StoreWarehouseStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeWarehouseStockList []data.StoreWarehouseStock

	query := r.db.Model(&data.StoreWarehouseStock{}).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Preload("StoreWarehouse").
		Joins("JOIN store_warehouses ON store_warehouse_stocks.store_warehouse_id = store_warehouses.id").
		Joins("JOIN ingredients ON store_warehouse_stocks.ingredient_id = ingredients.id").
		Where("store_warehouses.store_id = ?", storeID)

	err := query.Find(&storeWarehouseStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeWarehouseStockList, nil
}

func (r *storeWarehouseRepository) GetStockListByIDs(storeID uint, stockIds []uint) ([]data.StoreWarehouseStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeWarehouseStockList []data.StoreWarehouseStock

	query := r.db.Model(&data.StoreWarehouseStock{}).
		Preload("Ingredient.Unit").
		Joins("JOIN store_warehouses ON store_warehouse_stocks.store_warehouse_id = store_warehouses.id").
		Joins("JOIN ingredients ON store_warehouse_stocks.ingredient_id = ingredients.id").
		Where("store_warehouses.store_id = ? AND store_warehouse_stocks.id IN (?)", storeID, stockIds).
		Preload("Ingredient").
		Preload("StoreWarehouse")

	err := query.Find(&storeWarehouseStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeWarehouseStockList, nil
}

func (r *storeWarehouseRepository) GetStockById(storeId, stockId uint) (*data.StoreWarehouseStock, error) {
	var StoreWarehouseStock data.StoreWarehouseStock

	if storeId == 0 {
		return nil, fmt.Errorf("store warhouse stock id cannot be 0")
	}

	if stockId == 0 {
		return nil, fmt.Errorf("store warhouse stock id cannot be 0")
	}

	dbQuery := r.db.Model(&data.StoreWarehouseStock{}).
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Preload("StoreWarehouse").
		Preload("StoreWarehouse.Store")

	if err := dbQuery.
		Joins("JOIN store_warehouses ON store_warehouses.id = store_warehouse_stocks.store_warehouse_id").
		Where("store_warehouses.store_id = ?", storeId).
		Where("store_warehouse_stocks.id = ?", stockId).
		First(&StoreWarehouseStock).Error; err != nil {
		return nil, err
	}

	return &StoreWarehouseStock, nil
}

func (r *storeWarehouseRepository) UpdateStock(storeId, stockId uint, dto *types.UpdateStoreStockDTO) error {

	if storeId == 0 {
		return fmt.Errorf("storeId cannot be 0")
	}

	if stockId == 0 {
		return fmt.Errorf("stockId cannot be 0")
	}

	updateFields := map[string]interface{}{}

	if dto.Quantity != nil {
		updateFields["quantity"] = *dto.Quantity
	}
	if dto.LowStockThreshold != nil {
		updateFields["low_stock_threshold"] = *dto.LowStockThreshold
	}

	var existingStock data.StoreWarehouseStock

	res := r.db.Model(&data.StoreWarehouseStock{}).
		Joins("JOIN store_warehouses ON store_warehouses.id = store_warehouse_stocks.store_warehouse_id").
		Where("store_warehouses.store_id = ?", storeId).
		Where("store_warehouse_stocks.id = ?", stockId).
		First(&existingStock)

	if res.Error != nil {
		return utils.WrapError("failed to update store warehouse stock", res.Error)
	}

	updRes := r.db.Model(&data.StoreWarehouseStock{}).
		Where(&data.StoreWarehouseStock{BaseEntity: data.BaseEntity{ID: stockId}}).
		Updates(updateFields)

	if updRes.Error != nil {
		return utils.WrapError("failed to update store warehouse stock", updRes.Error)
	}

	if updRes.RowsAffected == 0 {
		return fmt.Errorf("update attempt had no changes for stockId=%d with storeId=%d", stockId, storeId)
	}

	return nil
}

func (r *storeWarehouseRepository) DeleteStockById(storeId, stockId uint) error {
	res := r.db.
		Where("id = ? AND store_warehouse_id IN (SELECT id FROM store_warehouses WHERE store_id = ?)", stockId, storeId).
		Delete(&data.StoreWarehouseStock{})

	if res.Error != nil {
		return fmt.Errorf("failed to delete store warehouse stock: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("stock not found")
	}

	return nil
}
