package storeWarehouses

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses/types"
	"gorm.io/gorm"
)

type StoreWarehouseRepository interface {
	//AddIngredient(input types.AddIngredientDTO) error
	GetStoreWarehouseStockList(query types.GetStoreWarehouseStockQuery) ([]data.StoreWarehouseStock, error)
	GetStoreWarehouseStockById(storeId, ingredientId uint) (*data.StoreWarehouseStock, error)
	PartialUpdateStoreWarehouseStock(id uint, updateFields map[string]interface{}) error
}

type storeWarehouseRepository struct {
	db *gorm.DB
}

func NewStoreWarehouseRepository(db *gorm.DB) StoreWarehouseRepository {
	return &storeWarehouseRepository{db: db}
}

/*func (r *storeWarehouseRepository) AddIngredient(input types.AddIngredientDTO) error {
	var storeWarehouseStock data.StoreWarehouseStock

	return r.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(input).Error; err != nil {
				return err
			}

			storeWarehouse := &data.StoreWarehouse{}

			if err := tx.Create(&storeWarehouseStock).Error; err != nil {
				if err := {}
			}
		})
}*/

func (r *storeWarehouseRepository) GetStoreWarehouseStockList(query types.GetStoreWarehouseStockQuery) ([]data.StoreWarehouseStock, error) {
	var StoreWarehouseStockList []data.StoreWarehouseStock
	var dbQuery *gorm.DB

	if query.StoreID != 0 {
		dbQuery = r.db.
			Model(&data.StoreWarehouseStock{}).
			Joins("JOIN store_warehouses ON store_warehouse_stock.store_warehouse_id = store_warehouses.id").
			Joins("JOIN ingredients ON ingredients.id = store_warehouse_stock.ingredient_id = ingredients.id").
			Where("store_warehouse.store_id = ?", query.StoreID)
	} else {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	if query.LowStockOnly != nil && *query.LowStockOnly {
		dbQuery = dbQuery.Where("store_warehouse_stock.quantity < store_warehouse_stock.low_stock_threshold")
	}

	if query.SearchTerm != nil && *query.SearchTerm != "" {
		dbQuery = dbQuery.Where("ingredients.name ILIKE ?", "%"+*query.SearchTerm+"%")
	}

	dbQuery.Limit(query.Limit)
	dbQuery.Offset(query.Offset)

	dbQuery.Preload("Ingredient").Preload("StoreWarehouse")

	err := dbQuery.Find(&StoreWarehouseStockList).Error
	if err != nil {
		return nil, err
	}

	return StoreWarehouseStockList, nil
}

func (r *storeWarehouseRepository) GetStoreWarehouseStockById(storeId, ingredientId uint) (*data.StoreWarehouseStock, error) {
	var StoreWarehouseStock data.StoreWarehouseStock

	if storeId == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}
	if ingredientId == 0 {
		return nil, fmt.Errorf("ingredientId cannot be 0")
	}

	dbQuery := r.db.Preload("Ingredient").Preload("StoreWarehouse")

	if err := dbQuery.First(&StoreWarehouseStock, "store_id = ? AND ingredient.id = ?", storeId, ingredientId).Error; err != nil {
		return nil, err
	}

	return &StoreWarehouseStock, nil
}

func (r *storeWarehouseRepository) PartialUpdateStoreWarehouseStock(id uint, updateFields map[string]interface{}) error {

	return nil
}
