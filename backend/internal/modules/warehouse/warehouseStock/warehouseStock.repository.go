package warehouseStock

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"gorm.io/gorm"
)

type WarehouseStockRepository interface {
	LogIncomingInventory(deliveries []data.SupplierWarehouseDelivery) error
	RecordDeliveriesAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error
	LogAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error
	TransferStock(sourceWarehouseID, targetWarehouseID uint, items []data.StockRequestIngredient) error

	GetDeliveryByID(deliveryID uint, delivery *data.SupplierWarehouseDelivery) error
	GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]data.SupplierWarehouseDelivery, error)

	ConvertInventoryItemsToStockRequest(items []types.ExistingInventoryItem) ([]data.StockRequestIngredient, error)
	SupplierMaterialExists(supplierID, stockMaterialID uint) (bool, error)
	CreateSupplierMaterial(association *data.SupplierMaterial) error

	AddToWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error
	DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error
	GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error)
	GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID uint) (*data.AggregatedWarehouseStock, []data.SupplierWarehouseDelivery, error)

	UpdateWarehouseStock(stock *data.WarehouseStock) error
	GetWarehouseStockByID(warehouseID, stockMaterialID uint) (*data.WarehouseStock, error)
	UpdateStockQuantity(stockID uint, quantity float64) error
	UpdateExpirationDate(stockMaterialID, warehouseID uint, newExpirationDate time.Time) error
}

type warehouseStockRepository struct {
	db *gorm.DB
}

func NewWarehouseStockRepository(db *gorm.DB) WarehouseStockRepository {
	return &warehouseStockRepository{db: db}
}

func (r *warehouseStockRepository) LogIncomingInventory(deliveries []data.SupplierWarehouseDelivery) error {
	return r.db.Create(&deliveries).Error
}

func (r *warehouseStockRepository) LogAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&deliveries).Error; err != nil {
			return fmt.Errorf("failed to log deliveries: %w", err)
		}

		for _, delivery := range deliveries {
			var stock data.WarehouseStock
			err := tx.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, delivery.StockMaterialID).
				First(&stock).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					newStock := data.WarehouseStock{
						WarehouseID:     warehouseID,
						StockMaterialID: delivery.StockMaterialID,
						Quantity:        delivery.Quantity,
					}
					if err := tx.Create(&newStock).Error; err != nil {
						return fmt.Errorf("failed to create warehouse stock: %w", err)
					}
				} else {
					return fmt.Errorf("failed to query warehouse stock: %w", err)
				}
			} else {
				err = tx.Model(&data.WarehouseStock{}).
					Where("id = ?", stock.ID).
					Update("quantity", gorm.Expr("quantity + ?", delivery.Quantity)).Error
				if err != nil {
					return fmt.Errorf("failed to update warehouse stock quantity: %w", err)
				}
			}
		}

		return nil
	})
}

func (r *warehouseStockRepository) RecordDeliveriesAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.recordDeliveries(tx, deliveries); err != nil {
			return fmt.Errorf("failed to record deliveries: %w", err)
		}

		for _, delivery := range deliveries {
			if err := r.updateOrInsertStock(tx, warehouseID, delivery); err != nil {
				return fmt.Errorf("failed to update or insert stock for stock_material_id %d: %w", delivery.StockMaterialID, err)
			}
		}

		return nil
	})
}

func (r *warehouseStockRepository) recordDeliveries(tx *gorm.DB, deliveries []data.SupplierWarehouseDelivery) error {
	if err := tx.Create(&deliveries).Error; err != nil {
		return fmt.Errorf("failed to log deliveries: %w", err)
	}
	return nil
}

func (r *warehouseStockRepository) updateOrInsertStock(tx *gorm.DB, warehouseID uint, delivery data.SupplierWarehouseDelivery) error {
	stock := data.WarehouseStock{
		WarehouseID:     warehouseID,
		StockMaterialID: delivery.StockMaterialID,
	}

	if err := tx.FirstOrCreate(&stock, "warehouse_id = ? AND stock_material_id = ?", warehouseID, delivery.StockMaterialID).Error; err != nil {
		return fmt.Errorf("failed to find or create warehouse stock: %w", err)
	}

	if stock.ID != 0 {
		if err := tx.Model(&data.WarehouseStock{}).
			Where("id = ?", stock.ID).
			Update("quantity", gorm.Expr("quantity + ?", delivery.Quantity)).Error; err != nil {
			return fmt.Errorf("failed to update warehouse stock quantity: %w", err)
		}
	}

	return nil
}

func (r *warehouseStockRepository) TransferStock(sourceWarehouseID, targetWarehouseID uint, items []data.StockRequestIngredient) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	for _, item := range items {
		err := tx.Model(&data.WarehouseStock{}).
			Where("warehouse_id = ? AND stock_material_id = ?", sourceWarehouseID, item.IngredientID).
			Update("quantity", gorm.Expr("quantity - ?", item.Quantity)).Error
		if err != nil {
			return fmt.Errorf("failed to deduct stock from source: %w", err)
		}

		err = tx.Model(&data.WarehouseStock{}).
			Where("warehouse_id = ? AND stock_material_id = ?", targetWarehouseID, item.IngredientID).
			Update("quantity", gorm.Expr("quantity + ?", item.Quantity)).Error
		if err != nil {
			return fmt.Errorf("failed to add stock to target: %w", err)
		}
	}

	return tx.Commit().Error
}

func (r *warehouseStockRepository) GetDeliveryByID(deliveryID uint, delivery *data.SupplierWarehouseDelivery) error {
	return r.db.First(delivery, "id = ?", deliveryID).Error
}

func (r *warehouseStockRepository) GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]data.SupplierWarehouseDelivery, error) {
	var deliveries []data.SupplierWarehouseDelivery
	query := r.db.Model(&data.SupplierWarehouseDelivery{})

	if warehouseID != nil {
		query = query.Where("warehouse_id = ?", *warehouseID)
	}
	if startDate != nil {
		query = query.Where("delivery_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("delivery_date <= ?", *endDate)
	}

	err := query.Find(&deliveries).Error
	return deliveries, err
}

func (r *warehouseStockRepository) ConvertInventoryItemsToStockRequest(items []types.ExistingInventoryItem) ([]data.StockRequestIngredient, error) {
	converted := make([]data.StockRequestIngredient, len(items))

	for i, item := range items {
		var stockMaterial data.StockMaterial
		if err := r.db.Preload("Ingredient").First(&stockMaterial, "id = ?", item.StockMaterialID).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve stock material for StockMaterialID %d: %w", item.StockMaterialID, err)
		}

		deliveredDate := time.Now()
		expirationDate := deliveredDate.AddDate(0, 0, stockMaterial.ExpirationPeriodInDays)

		converted[i] = data.StockRequestIngredient{
			IngredientID:   stockMaterial.IngredientID,
			Quantity:       item.Quantity,
			DeliveredDate:  deliveredDate,
			ExpirationDate: expirationDate,
		}
	}

	return converted, nil
}

func (r *warehouseStockRepository) SupplierMaterialExists(supplierID, stockMaterialID uint) (bool, error) {
	var count int64
	err := r.db.Model(&data.SupplierMaterial{}).
		Where("supplier_id = ? AND stock_material_id = ?", supplierID, stockMaterialID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *warehouseStockRepository) CreateSupplierMaterial(association *data.SupplierMaterial) error {
	return r.db.Create(association).Error
}

func (r *warehouseStockRepository) AddToWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error {
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

func (r *warehouseStockRepository) DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error {
	stock := &data.WarehouseStock{}
	if err := r.db.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		return fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	if stock.Quantity < quantityInPackages {
		return fmt.Errorf("insufficient stock for StockMaterialID %d in WarehouseID %d", stockMaterialID, warehouseID)
	}

	return r.db.Model(stock).Update("quantity", gorm.Expr("quantity - ?", quantityInPackages)).Error
}

func (r *warehouseStockRepository) GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error) {
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

func (r *warehouseStockRepository) GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID uint) (*data.AggregatedWarehouseStock, []data.SupplierWarehouseDelivery, error) {
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

// helper functions
func (r *warehouseStockRepository) findEarliestExpirationDate(deliveries []data.SupplierWarehouseDelivery) *time.Time {
	var earliest *time.Time
	for _, delivery := range deliveries {
		if earliest == nil || delivery.ExpirationDate.Before(*earliest) {
			earliest = &delivery.ExpirationDate
		}
	}
	return earliest
}

func (r *warehouseStockRepository) aggregateWarehouseStocks(
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

func (r *warehouseStockRepository) aggregateWarehouseStock(
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

func (r *warehouseStockRepository) getSupplierDeliveriesForStock(stockMaterialID, warehouseID uint) ([]data.SupplierWarehouseDelivery, error) {
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

func (r *warehouseStockRepository) getSupplierDeliveries(filter *types.GetWarehouseStockFilterQuery) ([]data.SupplierWarehouseDelivery, error) {
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

func (r *warehouseStockRepository) getWarehouseStock(stockMaterialID, warehouseID uint) (*data.WarehouseStock, error) {
	var stock data.WarehouseStock
	err := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Ingredient").
		Preload("StockMaterial.Ingredient.Unit").
		Preload("StockMaterial.Ingredient.IngredientCategory").
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Package").
		Preload("StockMaterial.Package.Unit").
		Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
		First(&stock).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch single warehouse stock: %w", err)
	}
	return &stock, nil
}

func (r *warehouseStockRepository) getWarehouseStocksWithPagination(filter *types.GetWarehouseStockFilterQuery) ([]data.WarehouseStock, int64, error) {
	var warehouseStocks []data.WarehouseStock
	var totalCount int64

	query := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Ingredient").
		Preload("StockMaterial.Ingredient.Unit").
		Preload("StockMaterial.Ingredient.IngredientCategory").
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Package").
		Preload("StockMaterial.Package.Unit").
		Joins("JOIN supplier_warehouse_deliveries ON supplier_warehouse_deliveries.stock_material_id = warehouse_stocks.stock_material_id").
		Joins("JOIN stock_materials ON warehouse_stocks.stock_material_id = stock_materials.id")

	// Apply filters
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.StockMaterialID != nil {
		query = query.Where("stock_material_id = ?", *filter.StockMaterialID)
	}

	if filter.IngredientID != nil {
		query = query.Where("stock_materials.ingredient_id = ?", *filter.IngredientID)
	}

	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("warehouse_stocks.quantity < stock_materials.safety_stock")
	}

	if filter.IsExpiring != nil && *filter.IsExpiring {
		days := 1095 // Default to 1095 days if not specified, as said
		if filter.ExpirationDays != nil {
			days = *filter.ExpirationDays
		}

		expirationThreshold := time.Now().AddDate(0, 0, days)
		query = query.Where("supplier_warehouse_deliveries.expiration_date <= ?", expirationThreshold)
	}

	if filter.CategoryID != nil {
		query = query.Where("stock_materials.category_id = ?", *filter.CategoryID)
	}

	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Where("stock_materials.name ILIKE ? OR stock_materials.description ILIKE ? OR stock_materials.barcode ILIKE ?", search, search, search)
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

func (r *warehouseStockRepository) UpdateWarehouseStock(stock *data.WarehouseStock) error {
	return r.db.Save(stock).Error
}

func (r *warehouseStockRepository) GetWarehouseStockByID(warehouseID, stockMaterialID uint) (*data.WarehouseStock, error) {
	var stock data.WarehouseStock
	if err := r.db.First(&stock).Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *warehouseStockRepository) UpdateStockQuantity(stockID uint, quantity float64) error {
	return r.db.Model(&data.WarehouseStock{}).
		Where("id = ?", stockID).
		Update("quantity", quantity).Error
}

func (r *warehouseStockRepository) UpdateExpirationDate(stockMaterialID, warehouseID uint, newExpirationDate time.Time) error {
	return r.db.Model(&data.SupplierWarehouseDelivery{}).
		Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
		Update("expiration_date", newExpirationDate).Error
}
