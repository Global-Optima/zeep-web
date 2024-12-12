package storeWarehouses

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreWarehouseRepository interface {
	AddStock(storeId uint, input types.AddStockDTO) (uint, error)
	GetStockList(storeId uint, query types.GetStockQuery) ([]data.StoreWarehouseStock, error)
	GetStockById(storeId, stockId uint) (*data.StoreWarehouseStock, error)
	PartialUpdateStock(storeId, stockId uint, updateFields map[string]interface{}) error
	DeleteStockById(storeId, stockId uint) error
}

type storeWarehouseRepository struct {
	db *gorm.DB
}

func NewStoreWarehouseRepository(db *gorm.DB) StoreWarehouseRepository {
	return &storeWarehouseRepository{db: db}
}

func (r *storeWarehouseRepository) AddStock(storeID uint, input types.AddStockDTO) (uint, error) {
	// Fetch the StoreWarehouse for the given storeID
	var storeWarehouse data.StoreWarehouse
	err := r.db.
		Where("store_id = ?", storeID).
		First(&storeWarehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("store warehouse not found for store ID %d", storeID)
		}
		return 0, fmt.Errorf("failed to fetch store warehouse: %w", err)
	}

	// Check if a stock with the given IngredientID already exists for the StoreWarehouse
	var existingStock data.StoreWarehouseStock
	err = r.db.
		Where("store_warehouse_id = ? AND ingredient_id = ?", storeWarehouse.ID, input.IngredientID).
		First(&existingStock).Error
	if err == nil {
		return 0, fmt.Errorf("stock with ingredient ID %d already exists for store warehouse ID %d", input.IngredientID, storeWarehouse.ID)
	}

	// Fetch the Ingredient for the given IngredientID from the DTO
	var ingredient data.Ingredient
	err = r.db.
		Where("id = ?", input.IngredientID).
		First(&ingredient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("ingredient not found for ingredient ID %d", input.IngredientID)
		}
		return 0, fmt.Errorf("failed to fetch ingredient: %w", err)
	}

	storeWarehouseStock := data.StoreWarehouseStock{
		StoreWarehouseID: storeWarehouse.ID,
		IngredientID:     ingredient.ID,
		//TODO add units
		LowStockThreshold: input.LowStockThreshold,
		Quantity:          input.CurrentStock,
	}

	err = r.db.Create(&storeWarehouseStock).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create store warehouse stock: %w", err)
	}

	return storeWarehouseStock.ID, nil
}

func (r *storeWarehouseRepository) GetStockList(storeID uint, query types.GetStockQuery) ([]data.StoreWarehouseStock, error) {
	var StoreWarehouseStockList []data.StoreWarehouseStock
	var dbQuery *gorm.DB
	var totalRecords int64

	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	dbQuery = r.db.
		Model(&data.StoreWarehouseStock{}).
		Joins("JOIN store_warehouses ON store_warehouse_stocks.store_warehouse_id = store_warehouses.id").
		Joins("JOIN ingredients ON store_warehouse_stocks.ingredient_id = ingredients.id").
		Where("store_warehouses.store_id = ?", storeID)

	if query.LowStockOnly != nil && *query.LowStockOnly {
		dbQuery = dbQuery.Where("store_warehouse_stocks.quantity < store_warehouse_stocks.low_stock_threshold")
	}

	if query.SearchTerm != nil && *query.SearchTerm != "" {
		dbQuery = dbQuery.Where("ingredients.name ILIKE ?", "%"+*query.SearchTerm+"%")
	}

	dbQuery.Preload("Ingredient").Preload("StoreWarehouse")

	//Query with pagination
	err := dbQuery.Scopes(query.Pagination.PaginateGorm()).Find(&StoreWarehouseStockList).Count(&totalRecords).Error
	query.Pagination.SetTotal(totalRecords)

	if err != nil {
		return nil, err
	}

	return StoreWarehouseStockList, nil
}

func (r *storeWarehouseRepository) GetStockById(storeId, stockId uint) (*data.StoreWarehouseStock, error) {
	var StoreWarehouseStock data.StoreWarehouseStock

	if storeId == 0 {
		return nil, fmt.Errorf("store warhouse stock id cannot be 0")
	}

	if stockId == 0 {
		return nil, fmt.Errorf("store warhouse stock id cannot be 0")
	}

	dbQuery := r.db.Preload("Ingredient").Preload("StoreWarehouse")

	if err := dbQuery.
		Joins("JOIN store_warehouses ON store_warehouses.id = store_warehouse_stocks.store_warehouse_id").
		Where("store_warehouses.store_id = ?", storeId).
		Where("store_warehouse_stocks.id = ?", stockId).
		First(&StoreWarehouseStock).Error; err != nil {
		return nil, err
	}

	return &StoreWarehouseStock, nil
}

func (r *storeWarehouseRepository) PartialUpdateStock(storeId, stockId uint, updateFields map[string]interface{}) error {
	if storeId == 0 {
		return fmt.Errorf("storeId cannot be 0")
	}

	if stockId == 0 {
		return fmt.Errorf("stockId cannot be 0")
	}

	if len(updateFields) == 0 {
		return fmt.Errorf("no fields provided for update")
	}

	res := r.db.Model(&data.StoreWarehouseStock{}).
		Joins("JOIN store_warehouses ON store_warehouses.id = store_warehouse_stocks.store_warehouse_id").
		Where("store_warehouse_stocks.id = ?", stockId).
		Updates(updateFields)

	if res.Error != nil {
		return fmt.Errorf("failed to update store warehouse stock: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("update attempt had no changes for stockId=%d with storeId=%d", stockId, storeId)
	}

	return nil
}

func (r *storeWarehouseRepository) DeleteStockById(storeId, stockId uint) error {
	return r.db.Model(&data.StoreWarehouseStock{}).
		Joins("StoreWarehouse", r.db.Where(&data.StoreWarehouse{StoreID: storeId})).
		Where("store_warehouse_stocks.id = ?", stockId).
		Delete(&data.StoreWarehouseStock{}).Error
}
