package warehouse

import (
	"fmt"
	"strings"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	AssignStoreToWarehouse(storeID, warehouseID uint) error
	ReassignStoreToWarehouse(storeID, newWarehouseID uint) error
	GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]data.Store, error)

	CreateWarehouse(warehouse *data.Warehouse, facilityAddress *data.FacilityAddress) error
	GetWarehouseByID(id uint) (*data.Warehouse, error)
	GetAllWarehouses(pagination *utils.Pagination) ([]data.Warehouse, error)
	UpdateWarehouse(warehouse *data.Warehouse) error
	DeleteWarehouse(id uint) error

	AddToWarehouseStock(warehouseID, stockMaterialID uint, quantity float64) error
	DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantity float64) error
	GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error)
	ResetWarehouseStock(warehouseID uint, stocks []data.WarehouseStock) error
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

// stores
func (r *warehouseRepository) AssignStoreToWarehouse(storeID, warehouseID uint) error {
	storeWarehouse := data.StoreWarehouse{
		StoreID:     storeID,
		WarehouseID: warehouseID,
	}
	return r.db.Create(&storeWarehouse).Error
}

func (r *warehouseRepository) ReassignStoreToWarehouse(storeID, newWarehouseID uint) error {
	return r.db.Model(&data.StoreWarehouse{}).
		Where("store_id = ?", storeID).
		Update("warehouse_id", newWarehouseID).Error
}

func (r *warehouseRepository) GetAllStoresByWarehouse(warehouseID uint, pagination *utils.Pagination) ([]data.Store, error) {
	var stores []data.Store

	query := r.db.Joins("JOIN store_warehouses ON stores.id = store_warehouses.store_id").
		Where("store_warehouses.warehouse_id = ?", warehouseID)

	if _, err := utils.ApplyPagination(query, pagination, &data.Store{}); err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}

// warehouses CRUD
func (r *warehouseRepository) CreateWarehouse(warehouse *data.Warehouse, facilityAddress *data.FacilityAddress) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(facilityAddress).Error; err != nil {
			return err
		}
		warehouse.FacilityAddressID = facilityAddress.ID
		return tx.Create(warehouse).Error
	})
}

func (r *warehouseRepository) GetWarehouseByID(id uint) (*data.Warehouse, error) {
	var warehouse data.Warehouse
	if err := r.db.Preload("FacilityAddress").First(&warehouse, id).Error; err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) GetAllWarehouses(pagination *utils.Pagination) ([]data.Warehouse, error) {
	var warehouses []data.Warehouse

	query := r.db.Preload("FacilityAddress").Model(&data.Warehouse{})
	if _, err := utils.ApplyPagination(query, pagination, &data.Warehouse{}); err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	if err := query.Find(&warehouses).Error; err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (r *warehouseRepository) UpdateWarehouse(warehouse *data.Warehouse) error {
	return r.db.Save(warehouse).Error
}

func (r *warehouseRepository) DeleteWarehouse(id uint) error {
	return r.db.Delete(&data.Warehouse{}, id).Error
}

// stocks
func (r *warehouseRepository) AddToWarehouseStock(warehouseID, stockMaterialID uint, quantity float64) error {
	stock := &data.WarehouseStock{}
	if err := r.db.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			stock = &data.WarehouseStock{
				WarehouseID:     warehouseID,
				StockMaterialID: stockMaterialID,
				Quantity:        quantity,
			}
			return r.db.Create(stock).Error
		}
		return fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	return r.db.Model(stock).Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
}

func (r *warehouseRepository) DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantity float64) error {
	stock := &data.WarehouseStock{}
	if err := r.db.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		return fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	if stock.Quantity < quantity {
		return fmt.Errorf("insufficient stock for StockMaterialID %d in WarehouseID %d", stockMaterialID, warehouseID)
	}

	return r.db.Model(stock).Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

func (r *warehouseRepository) GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error) {
	var stocks []data.AggregatedWarehouseStock

	query := r.db.Model(&data.WarehouseStock{}).
		Select(`
			warehouse_stocks.warehouse_id,
			warehouse_stocks.stock_material_id,
			SUM(supplier_warehouse_deliveries.quantity) AS total_quantity,
			MIN(supplier_warehouse_deliveries.expiration_date) AS earliest_expiration_date
		`).
		Joins("LEFT JOIN supplier_warehouse_deliveries ON warehouse_stocks.stock_material_id = supplier_warehouse_deliveries.stock_material_id").
		Joins("LEFT JOIN stock_materials ON warehouse_stocks.stock_material_id = stock_materials.id").
		Preload("StockMaterial.Unit").
		Group("warehouse_stocks.warehouse_id, warehouse_stocks.stock_material_id, stock_materials.id")

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_stocks.warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.StockMaterialID != nil {
		query = query.Where("warehouse_stocks.stock_material_id = ?", *filter.StockMaterialID)
	}

	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("SUM(supplier_warehouse_deliveries.quantity) < stock_materials.safety_stock")
	}

	if filter.Category != nil && *filter.Category != "" {
		query = query.Where("LOWER(stock_materials.category) LIKE ?", "%"+strings.ToLower(*filter.Category)+"%")
	}

	if filter.ExpirationDays != nil && *filter.ExpirationDays > 0 {
		expirationThreshold := time.Now().AddDate(0, 0, *filter.ExpirationDays)
		query = query.Where("MIN(supplier_warehouse_deliveries.expiration_date) <= ?", expirationThreshold)
	}

	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("LOWER(stock_materials.name) LIKE ?", "%"+strings.ToLower(*filter.Search)+"%")
	}

	// Apply pagination
	var err error
	query, err = utils.ApplyPagination(query, filter.Pagination, &data.AggregatedWarehouseStock{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	if err := query.Scan(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch stocks: %w", err)
	}

	return stocks, nil
}

func (r *warehouseRepository) ResetWarehouseStock(warehouseID uint, stocks []data.WarehouseStock) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("warehouse_id = ?", warehouseID).Delete(&data.WarehouseStock{}).Error; err != nil {
			return fmt.Errorf("failed to clear warehouse stock: %w", err)
		}

		if len(stocks) > 0 {
			if err := tx.Create(&stocks).Error; err != nil {
				return fmt.Errorf("failed to create new warehouse stock: %w", err)
			}
		}
		return nil
	})
}
