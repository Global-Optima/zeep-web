package warehouse

import (
	"fmt"
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

	AddToWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error
	DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error
	GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error)
	GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID uint) (*data.AggregatedWarehouseStock, []data.SupplierWarehouseDelivery, error)
	ResetWarehouseStock(warehouseID uint, stocks []data.WarehouseStock) error
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

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

func (r *warehouseRepository) AddToWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error {
	stock := &data.WarehouseStock{}
	if err := r.db.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			stock = &data.WarehouseStock{
				WarehouseID:     warehouseID,
				StockMaterialID: stockMaterialID,
				Quantity:        quantityInPackages,
			}
			return r.db.Create(stock).Error
		}
		return fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	return r.db.Model(stock).Update("quantity", gorm.Expr("quantity + ?", quantityInPackages)).Error
}

func (r *warehouseRepository) DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error {
	stock := &data.WarehouseStock{}
	if err := r.db.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		return fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	if stock.Quantity < quantityInPackages {
		return fmt.Errorf("insufficient stock for StockMaterialID %d in WarehouseID %d", stockMaterialID, warehouseID)
	}

	return r.db.Model(stock).Update("quantity", gorm.Expr("quantity - ?", quantityInPackages)).Error
}

func (r *warehouseRepository) GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error) {
	warehouseStocks, totalCount, err := r.getWarehouseStocksWithPagination(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stocks: %w", err)
	}

	supplierDeliveries, err := r.getSupplierDeliveries(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch supplier deliveries: %w", err)
	}

	deliveryMap := make(map[uint][]data.SupplierWarehouseDelivery)
	for _, delivery := range supplierDeliveries {
		deliveryMap[delivery.StockMaterialID] = append(deliveryMap[delivery.StockMaterialID], delivery)
	}

	aggregatedStocks := r.aggregateWarehouseStocks(warehouseStocks, deliveryMap)

	if filter.Pagination != nil {
		filter.Pagination.TotalCount = int(totalCount)
		if filter.Pagination.PageSize > 0 {
			filter.Pagination.TotalPages = (int(totalCount) + filter.Pagination.PageSize - 1) / filter.Pagination.PageSize
		}
	}

	return aggregatedStocks, nil
}

func (r *warehouseRepository) GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID uint) (*data.AggregatedWarehouseStock, []data.SupplierWarehouseDelivery, error) {
	warehouseStock, err := r.getWarehouseStock(stockMaterialID, warehouseID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	deliveries, err := r.getSupplierDeliveriesForStock(stockMaterialID, warehouseID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch supplier deliveries: %w", err)
	}

	aggregatedStock := r.aggregateWarehouseStock(*warehouseStock, deliveries)
	return aggregatedStock, deliveries, nil
}

func (r *warehouseRepository) ResetWarehouseStock(warehouseID uint, stocks []data.WarehouseStock) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Unscoped().Where("warehouse_id = ?", warehouseID).Delete(&data.WarehouseStock{}).Error; err != nil {
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

// helper functions
func (r *warehouseRepository) findEarliestExpirationDate(deliveries []data.SupplierWarehouseDelivery) *time.Time {
	var earliest *time.Time
	for _, delivery := range deliveries {
		if earliest == nil || delivery.ExpirationDate.Before(*earliest) {
			earliest = &delivery.ExpirationDate
		}
	}
	return earliest
}

func (r *warehouseRepository) aggregateWarehouseStocks(
	warehouseStocks []data.WarehouseStock,
	deliveryMap map[uint][]data.SupplierWarehouseDelivery,
) []data.AggregatedWarehouseStock {
	aggregatedStocks := []data.AggregatedWarehouseStock{}

	for _, stock := range warehouseStocks {
		deliveries := deliveryMap[stock.StockMaterialID]

		earliestExpirationDate := r.findEarliestExpirationDate(deliveries)

		aggregatedStocks = append(aggregatedStocks, data.AggregatedWarehouseStock{
			WarehouseID:            stock.WarehouseID,
			StockMaterialID:        stock.StockMaterialID,
			StockMaterial:          stock.StockMaterial,
			TotalQuantity:          stock.Quantity,
			EarliestExpirationDate: earliestExpirationDate,
		})
	}

	return aggregatedStocks
}

func (r *warehouseRepository) aggregateWarehouseStock(
	warehouseStock data.WarehouseStock,
	deliveries []data.SupplierWarehouseDelivery,
) *data.AggregatedWarehouseStock {
	earliestExpirationDate := r.findEarliestExpirationDate(deliveries)

	return &data.AggregatedWarehouseStock{
		WarehouseID:            warehouseStock.WarehouseID,
		StockMaterialID:        warehouseStock.StockMaterialID,
		StockMaterial:          warehouseStock.StockMaterial,
		TotalQuantity:          warehouseStock.Quantity,
		EarliestExpirationDate: earliestExpirationDate,
	}
}

func (r *warehouseRepository) getSupplierDeliveriesForStock(stockMaterialID, warehouseID uint) ([]data.SupplierWarehouseDelivery, error) {
	var deliveries []data.SupplierWarehouseDelivery
	err := r.db.Model(&data.SupplierWarehouseDelivery{}).
		Preload("Supplier").
		Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
		Find(&deliveries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries for stock material: %w", err)
	}
	return deliveries, nil
}

func (r *warehouseRepository) getSupplierDeliveries(filter *types.GetWarehouseStockFilterQuery) ([]data.SupplierWarehouseDelivery, error) {
	var deliveries []data.SupplierWarehouseDelivery

	query := r.db.Model(&data.SupplierWarehouseDelivery{})

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.StockMaterialID != nil {
		query = query.Where("stock_material_id = ?", *filter.StockMaterialID)
	}

	if err := query.Find(&deliveries).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch supplier deliveries: %w", err)
	}

	return deliveries, nil
}

func (r *warehouseRepository) getWarehouseStock(stockMaterialID, warehouseID uint) (*data.WarehouseStock, error) {
	var stock data.WarehouseStock
	err := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Package").
		Preload("StockMaterial.Package.Unit").
		Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
		First(&stock).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch single warehouse stock: %w", err)
	}
	return &stock, nil
}

func (r *warehouseRepository) getWarehouseStocksWithPagination(filter *types.GetWarehouseStockFilterQuery) ([]data.WarehouseStock, int64, error) {
	var warehouseStocks []data.WarehouseStock
	var totalCount int64

	query := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Package").
		Preload("StockMaterial.Package.Unit").
		Preload("StockMaterial")

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.StockMaterialID != nil {
		query = query.Where("stock_material_id = ?", *filter.StockMaterialID)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count warehouse stocks: %w", err)
	}

	if filter.Pagination != nil {
		offset := (filter.Pagination.Page - 1) * filter.Pagination.PageSize
		query = query.Offset(offset).Limit(filter.Pagination.PageSize)
	}

	if err := query.Find(&warehouseStocks).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch warehouse stocks: %w", err)
	}

	return warehouseStocks, totalCount, nil
}
