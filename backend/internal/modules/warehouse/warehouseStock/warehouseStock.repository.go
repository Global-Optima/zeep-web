package warehouseStock

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseStockRepository interface {
	RecordDeliveriesAndUpdateStock(delivery data.SupplierWarehouseDelivery, materials []data.SupplierWarehouseDeliveryMaterial, warehouseID uint) error

	GetDeliveryByID(deliveryID uint, delivery *data.SupplierWarehouseDelivery) error
	GetDeliveries(filter types.WarehouseDeliveryFilter) ([]data.SupplierWarehouseDelivery, error)

	ConvertInventoryItemsToStockRequest(items []types.ReceiveWarehouseStockMaterial) ([]data.StockRequestIngredient, error)

	AddToWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error
	DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) (*data.WarehouseStock, error)
	GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error)
	GetWarehouseStockMaterialDetails(stockMaterialID uint, filter *contexts.WarehouseContextFilter) (*data.AggregatedWarehouseStock, error)
	AddWarehouseStocks(warehouseID uint, stocks []data.WarehouseStock) error

	UpdateWarehouseStock(stock *data.WarehouseStock) error
	GetWarehouseStockByID(warehouseID, stockMaterialID uint) (*data.WarehouseStock, error)
	GetWarehouseStocksForNotifications(warehouseID uint) ([]data.WarehouseStock, error)
	UpdateStockQuantity(stockID uint, warehouseID uint, quantity float64) (*data.WarehouseStock, error)
	UpdateExpirationDate(stockMaterialID, warehouseID uint, newExpirationDate time.Time) error

	GetAvailableToAddStockMaterials(storeID uint, filter *types.AvailableStockMaterialFilter) ([]data.StockMaterial, error)

	FindEarliestExpirationDateForStock(stockMaterialID uint, filter *contexts.WarehouseContextFilter) (*time.Time, error)
}

type warehouseStockRepository struct {
	db *gorm.DB
}

func NewWarehouseStockRepository(db *gorm.DB) WarehouseStockRepository {
	return &warehouseStockRepository{db: db}
}

func (r *warehouseStockRepository) RecordDeliveriesAndUpdateStock(delivery data.SupplierWarehouseDelivery, materials []data.SupplierWarehouseDeliveryMaterial, warehouseID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&delivery).Error; err != nil {
			return fmt.Errorf("failed to create delivery: %w", err)
		}

		for i := range materials {
			materials[i].DeliveryID = delivery.ID
		}

		err := r.recordDeliveryMaterials(tx, materials)
		if err != nil {
			return fmt.Errorf("failed to record delivery materials: %w", err)
		}

		for _, material := range materials {
			if err := r.updateOrInsertStock(tx, warehouseID, material); err != nil {
				return fmt.Errorf("failed to update or insert stock for stock_material_id %d: %w", material.StockMaterialID, err)
			}
		}

		return nil
	})
}

func (r *warehouseStockRepository) recordDeliveryMaterials(tx *gorm.DB, materials []data.SupplierWarehouseDeliveryMaterial) error {
	if err := tx.Create(&materials).Error; err != nil {
		return fmt.Errorf("failed to log delivery materials: %w", err)
	}
	return nil
}

func (r *warehouseStockRepository) updateOrInsertStock(tx *gorm.DB, warehouseID uint, material data.SupplierWarehouseDeliveryMaterial) error {
	stock := data.WarehouseStock{
		WarehouseID:     warehouseID,
		StockMaterialID: material.StockMaterialID,
	}

	if err := tx.FirstOrCreate(&stock, "warehouse_id = ? AND stock_material_id = ?", warehouseID, material.StockMaterialID).Error; err != nil {
		return fmt.Errorf("failed to find or create warehouse stock: %w", err)
	}

	if err := tx.Model(&data.WarehouseStock{}).
		Where("id = ?", stock.ID).
		Update("quantity", gorm.Expr("quantity + ?", material.Quantity)).Error; err != nil {
		return fmt.Errorf("failed to update warehouse stock quantity: %w", err)
	}

	return nil
}

func (r *warehouseStockRepository) GetDeliveryByID(deliveryID uint, delivery *data.SupplierWarehouseDelivery) error {
	return r.db.Preload("Materials").
		Preload("Supplier").
		Preload("Warehouse").
		Preload("Materials.StockMaterial").
		Preload("Materials.StockMaterial.Unit").
		Preload("Materials.StockMaterial.Ingredient").
		Preload("Materials.StockMaterial.Ingredient.Unit").
		Preload("Materials.StockMaterial.Ingredient.IngredientCategory").
		Preload("Materials.StockMaterial.StockMaterialCategory").
		First(delivery, "id = ?", deliveryID).Error
}

func (r *warehouseStockRepository) GetDeliveries(filter types.WarehouseDeliveryFilter) ([]data.SupplierWarehouseDelivery, error) {
	var deliveries []data.SupplierWarehouseDelivery
	query := r.db.Model(&data.SupplierWarehouseDelivery{}).
		Preload("Materials").
		Preload("Supplier").
		Preload("Warehouse").
		Preload("Materials.StockMaterial").
		Preload("Materials.StockMaterial.Unit").
		Preload("Materials.StockMaterial.Ingredient").
		Preload("Materials.StockMaterial.Ingredient.Unit").
		Preload("Materials.StockMaterial.Ingredient.IngredientCategory").
		Preload("Materials.StockMaterial.StockMaterialCategory").
		Joins("JOIN suppliers ON suppliers.id = supplier_warehouse_deliveries.supplier_id")

	if filter.WarehouseID != nil {
		query = query.Where("supplier_warehouse_deliveries.warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.StartDate != nil {
		query = query.Where("delivery_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("delivery_date <= ?", *filter.EndDate)
	}

	search := "%"
	if filter.Search != nil {
		search = "%" + *filter.Search + "%"
	}

	query = query.
		Joins("JOIN supplier_warehouse_delivery_materials ON supplier_warehouse_delivery_materials.delivery_id = supplier_warehouse_deliveries.id").
		Joins("JOIN stock_materials ON supplier_warehouse_delivery_materials.stock_material_id = stock_materials.id").
		Where(`
			(
				stock_materials.name ILIKE ? OR
				stock_materials.description ILIKE ? OR
				stock_materials.barcode ILIKE ? OR
				suppliers.name ILIKE ?
			)
		`, search, search, search, search).
		Group("supplier_warehouse_deliveries.id")

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.SupplierWarehouseDelivery{})
	if err != nil {
		return nil, err
	}

	err = query.Find(&deliveries).Error
	return deliveries, err
}

func (r *warehouseStockRepository) ConvertInventoryItemsToStockRequest(items []types.ReceiveWarehouseStockMaterial) ([]data.StockRequestIngredient, error) {
	converted := make([]data.StockRequestIngredient, len(items))

	for i, item := range items {
		var stockMaterial data.StockMaterial
		if err := r.db.Preload("Ingredient").First(&stockMaterial, "id = ?", item.StockMaterialID).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve stock material for StockMaterialID %d: %w", item.StockMaterialID, err)
		}

		deliveredDate := time.Now()
		expirationDate := deliveredDate.AddDate(0, 0, stockMaterial.ExpirationPeriodInDays)

		converted[i] = data.StockRequestIngredient{
			Quantity:       item.Quantity,
			DeliveredDate:  deliveredDate,
			ExpirationDate: expirationDate,
		}
	}

	return converted, nil
}

func (r *warehouseStockRepository) AddToWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error {
	stock := &data.WarehouseStock{}
	if err := r.db.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			stock = &data.WarehouseStock{
				WarehouseID:     warehouseID,
				StockMaterialID: stockMaterialID,
				Quantity:        quantityInPackages,
			}
			return r.db.Create(stock).Error
		}
		return types.ErrAddWarehouseStockMaterial
	}

	return r.db.Model(stock).Update("quantity", gorm.Expr("quantity + ?", quantityInPackages)).Error
}

func (r *warehouseStockRepository) DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) (*data.WarehouseStock, error) {
	stock := &data.WarehouseStock{}
	if err := r.db.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	if stock.Quantity < quantityInPackages {
		return nil, fmt.Errorf("insufficient stock for StockMaterialID %d in WarehouseID %d: available %.2f, requested %.2f", stockMaterialID, warehouseID, stock.Quantity, quantityInPackages)
	}

	if err := r.db.Model(stock).
		Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
		Update("quantity", gorm.Expr("quantity - ?", quantityInPackages)).Error; err != nil {
		return nil, fmt.Errorf("failed to deduct stock for StockMaterialID %d in WarehouseID %d: %w", stockMaterialID, warehouseID, err)
	}

	if err := r.db.Preload("StockMaterial").Preload("Warehouse").Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).First(stock).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated stock: %w", err)
	}

	return stock, nil
}

func (r *warehouseStockRepository) GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error) {
	warehouseStocks, err := r.getWarehouseStocksWithPagination(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stocks: %w", err)
	}

	materials, err := r.getDeliveryMaterials(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delivery materials: %w", err)
	}

	materialMap := make(map[uint][]data.SupplierWarehouseDeliveryMaterial)
	for _, material := range materials {
		materialMap[material.StockMaterialID] = append(materialMap[material.StockMaterialID], material)
	}

	aggregatedStocks := r.aggregateWarehouseStocks(warehouseStocks, materialMap)

	return aggregatedStocks, nil
}

func (r *warehouseStockRepository) GetWarehouseStockMaterialDetails(stockMaterialID uint, filter *contexts.WarehouseContextFilter) (*data.AggregatedWarehouseStock, error) {
	warehouseStock, err := r.getWarehouseStock(stockMaterialID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	earliestExpirationDate, err := r.FindEarliestExpirationDateForStock(stockMaterialID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch earliest expiration date: %w", err)
	}

	return &data.AggregatedWarehouseStock{
		WarehouseID:            warehouseStock.WarehouseID,
		StockMaterialID:        warehouseStock.StockMaterialID,
		StockMaterial:          warehouseStock.StockMaterial,
		TotalQuantity:          warehouseStock.Quantity,
		EarliestExpirationDate: earliestExpirationDate,
	}, nil
}

// helper functions
func (r *warehouseStockRepository) findEarliestMaterialExpirationDate(materials []data.SupplierWarehouseDeliveryMaterial) *time.Time {
	if len(materials) == 0 {
		return nil
	}

	var earliest *time.Time
	for _, material := range materials {
		if earliest == nil || material.ExpirationDate.Before(*earliest) {
			earliest = &material.ExpirationDate
		}
	}

	if earliest == nil {
		return nil
	}

	utcTime := earliest.UTC()
	return &utcTime
}

func (r *warehouseStockRepository) FindEarliestExpirationDateForStock(stockMaterialID uint, filter *contexts.WarehouseContextFilter) (*time.Time, error) {
	var earliestExpirationDate sql.NullTime
	query := r.db.Model(&data.SupplierWarehouseDeliveryMaterial{}).
		Where(&data.SupplierWarehouseDeliveryMaterial{StockMaterialID: stockMaterialID}).
		Select("MIN(supplier_warehouse_delivery_materials.expiration_date) AS earliest_expiration_date")

	if filter != nil {
		if filter.WarehouseID != nil {
			query.Where(&data.SupplierWarehouseDeliveryMaterial{
				Delivery: data.SupplierWarehouseDelivery{
					WarehouseID: *filter.WarehouseID,
				},
			})
		}

		if filter.RegionID != nil {
			query.Joins("JOIN supplier_warehouse_deliveries ON supplier_warehouse_deliveries.id = supplier_warehouse_delivery_materials.delivery_id").
				Where(&data.SupplierWarehouseDeliveryMaterial{
					Delivery: data.SupplierWarehouseDelivery{
						Warehouse: data.Warehouse{RegionID: *filter.RegionID},
					},
				})
		}
	}

	err := query.Scan(&earliestExpirationDate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch earliest expiration date for stock material ID %d: %w", stockMaterialID, err)
	}

	if !earliestExpirationDate.Valid {
		return nil, nil
	}

	UTCTime := earliestExpirationDate.Time.UTC()
	return &UTCTime, nil
}

func (r *warehouseStockRepository) aggregateWarehouseStocks(
	warehouseStocks []data.WarehouseStock,
	materialMap map[uint][]data.SupplierWarehouseDeliveryMaterial,
) []data.AggregatedWarehouseStock {
	var aggregatedStocks []data.AggregatedWarehouseStock

	for _, stock := range warehouseStocks {
		materials := materialMap[stock.StockMaterialID]

		earliestExpirationDate := r.findEarliestMaterialExpirationDate(materials)

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

func (r *warehouseStockRepository) getDeliveryMaterials(filter *types.GetWarehouseStockFilterQuery) ([]data.SupplierWarehouseDeliveryMaterial, error) {
	var materials []data.SupplierWarehouseDeliveryMaterial

	query := r.db.Model(&data.SupplierWarehouseDeliveryMaterial{}).
		Preload("StockMaterial")

	if filter.WarehouseID != nil {
		query = query.Joins("JOIN supplier_warehouse_deliveries ON supplier_warehouse_deliveries.id = supplier_warehouse_delivery_materials.delivery_id").
			Where("supplier_warehouse_deliveries.warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.StockMaterialID != nil {
		query = query.Where("stock_material_id = ?", *filter.StockMaterialID)
	}

	if err := query.Find(&materials).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to fetch delivery materials: %w", err)
	}

	return materials, nil
}

func (r *warehouseStockRepository) getWarehouseStock(stockMaterialID uint, filter *contexts.WarehouseContextFilter) (*data.WarehouseStock, error) {
	var stock data.WarehouseStock
	query := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Ingredient").
		Preload("StockMaterial.Ingredient.Unit").
		Preload("StockMaterial.Ingredient.IngredientCategory").
		Preload("StockMaterial.StockMaterialCategory").
		Where(&data.WarehouseStock{StockMaterialID: stockMaterialID})

	if filter != nil {
		if filter.WarehouseID != nil {
			query.Where(&data.WarehouseStock{
				WarehouseID: *filter.WarehouseID,
			})
		}

		if filter.RegionID != nil {
			query.Joins("JOIN warehouses ON warehouses.id = warehouse_stocks.warehouse_id")
			query.Where(&data.WarehouseStock{Warehouse: data.Warehouse{RegionID: *filter.RegionID}})
		}
	}

	err := query.First(&stock).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch single warehouse stock: %w", err)
	}
	return &stock, nil
}

func (r *warehouseStockRepository) getWarehouseStocksWithPagination(filter *types.GetWarehouseStockFilterQuery) ([]data.WarehouseStock, error) {
	var warehouseStocks []data.WarehouseStock

	// Base query
	query := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Ingredient").
		Preload("StockMaterial.Ingredient.Unit").
		Preload("StockMaterial.Ingredient.IngredientCategory").
		Preload("StockMaterial.StockMaterialCategory").
		Joins("JOIN stock_materials ON warehouse_stocks.stock_material_id = stock_materials.id")

	// Filter by warehouse
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_stocks.warehouse_id = ?", *filter.WarehouseID)
	}

	// Filter by stock material
	if filter.StockMaterialID != nil {
		query = query.Where("warehouse_stocks.stock_material_id = ?", *filter.StockMaterialID)
	}

	// Filter by ingredient
	if filter.IngredientID != nil {
		query = query.Where("stock_materials.ingredient_id = ?", *filter.IngredientID)
	}

	// Filter by low stock
	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("warehouse_stocks.quantity <= stock_materials.safety_stock")
	}

	// Filter by expiration-related conditions
	if filter.IsExpiring != nil && *filter.IsExpiring {
		query = query.Joins(`
			LEFT JOIN supplier_warehouse_delivery_materials
			ON supplier_warehouse_delivery_materials.stock_material_id = warehouse_stocks.stock_material_id
		`).
			Where(`
				supplier_warehouse_delivery_materials.expiration_date IS NOT NULL AND
				supplier_warehouse_delivery_materials.expiration_date <= NOW() + INTERVAL '7 days'
			`).
			Group("warehouse_stocks.id")
	}

	// Filter by expiration days
	if filter.ExpirationDays != nil {
		expirationThreshold := time.Now().AddDate(0, 0, *filter.ExpirationDays)
		query = query.Joins(`
			LEFT JOIN supplier_warehouse_delivery_materials
			ON supplier_warehouse_delivery_materials.stock_material_id = warehouse_stocks.stock_material_id
		`).
			Where(`
				supplier_warehouse_delivery_materials.expiration_date IS NOT NULL AND
				supplier_warehouse_delivery_materials.expiration_date <= ?
			`, expirationThreshold).
			Group("warehouse_stocks.id")
	}

	// Filter by category
	if filter.CategoryID != nil {
		query = query.Where("stock_materials.category_id = ?", *filter.CategoryID)
	}

	// Filter by search query with proper grouping of OR conditions
	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Where(`
			(
				stock_materials.name ILIKE ? OR
				stock_materials.description ILIKE ? OR
				stock_materials.barcode ILIKE ?
			)
		`, search, search, search)
	}

	// Apply pagination and sorting
	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.WarehouseStock{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	// Execute query
	if err := query.Find(&warehouseStocks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stocks: %w", err)
	}

	return warehouseStocks, nil
}

func (r *warehouseStockRepository) UpdateExpirationDate(stockMaterialID, warehouseID uint, newExpirationDate time.Time) error {
	var deliveryMaterialIDs []uint
	err := r.db.Model(&data.SupplierWarehouseDeliveryMaterial{}).
		Joins("JOIN supplier_warehouse_deliveries ON supplier_warehouse_deliveries.id = supplier_warehouse_delivery_materials.delivery_id").
		Where("supplier_warehouse_deliveries.warehouse_id = ? AND supplier_warehouse_delivery_materials.stock_material_id = ?", warehouseID, stockMaterialID).
		Pluck("supplier_warehouse_delivery_materials.id", &deliveryMaterialIDs).Error
	if err != nil {
		return fmt.Errorf("failed to identify delivery materials to update: %w", err)
	}

	if len(deliveryMaterialIDs) == 0 {
		return fmt.Errorf("no deliveries found for stock material ID %d in warehouse ID %d", stockMaterialID, warehouseID)
	}

	err = r.db.Model(&data.SupplierWarehouseDeliveryMaterial{}).
		Where("id IN ?", deliveryMaterialIDs).
		Update("expiration_date", newExpirationDate).Error
	if err != nil {
		return fmt.Errorf("failed to update expiration date: %w", err)
	}

	return nil
}

func (r *warehouseStockRepository) UpdateStockQuantity(stockMaterialID, warehouseID uint, quantity float64) (*data.WarehouseStock, error) {
	var updatedStock data.WarehouseStock

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var warehouseStock data.WarehouseStock
		if err := tx.Model(&data.WarehouseStock{}).
			Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
			First(&warehouseStock).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("stock material ID %d not found in warehouse ID %d", stockMaterialID, warehouseID)
			}
			return fmt.Errorf("failed to check warehouse stock existence: %w", err)
		}

		if err := tx.Model(&data.WarehouseStock{}).
			Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
			Update("quantity", quantity).Error; err != nil {
			return fmt.Errorf("failed to update quantity for stock material ID %d in warehouse ID %d: %w", stockMaterialID, warehouseID, err)
		}

		if err := tx.Model(&data.WarehouseStock{}).Preload("StockMaterial").
			Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
			First(&updatedStock).Error; err != nil {
			return fmt.Errorf("failed to fetch updated stock: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &updatedStock, nil
}

func (r *warehouseStockRepository) AddWarehouseStocks(warehouseID uint, stocks []data.WarehouseStock) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, stock := range stocks {
			var existingStock data.WarehouseStock
			err := tx.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stock.StockMaterialID).
				First(&existingStock).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("failed to check existing warehouse stock: %w", err)
			}

			if existingStock.ID > 0 {
				existingStock.Quantity += stock.Quantity
				if err := tx.Save(&existingStock).Error; err != nil {
					return fmt.Errorf("failed to update existing warehouse stock: %w", err)
				}
			} else {
				if err := tx.Create(&stock).Error; err != nil {
					return fmt.Errorf("failed to insert new warehouse stock: %w", err)
				}
			}
		}
		return nil
	})
}

func (r *warehouseStockRepository) GetWarehouseStockByID(warehouseID, stockMaterialID uint) (*data.WarehouseStock, error) {
	var stock data.WarehouseStock
	if err := r.db.Model(&stock).
		Preload("Warehouse").
		Preload("StockMaterial").
		First(&stock).Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *warehouseStockRepository) GetWarehouseStocksForNotifications(warehouseID uint) ([]data.WarehouseStock, error) {
	var stocks []data.WarehouseStock
	err := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial").
		Preload("Warehouse").
		Where("warehouse_id = ?", warehouseID).
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *warehouseStockRepository) UpdateWarehouseStock(stock *data.WarehouseStock) error {
	return r.db.Save(stock).Error
}

func (r *warehouseStockRepository) GetAvailableToAddStockMaterials(storeID uint, filter *types.AvailableStockMaterialFilter) ([]data.StockMaterial, error) {
	var stockMaterials []data.StockMaterial
	var warehouseID uint

	if err := r.db.
		Model(&data.Store{}).
		Select("warehouse_id").
		Where("id = ?", storeID).
		Limit(1).
		Scan(&warehouseID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse for store: %w", err)
	}

	if warehouseID == 0 {
		return nil, fmt.Errorf("no warehouse found for the given store ID: %d", storeID)
	}

	query := r.db.Model(&data.StockMaterial{}).
		Select("DISTINCT stock_materials.*").
		Joins("JOIN warehouse_stocks ON warehouse_stocks.stock_material_id = stock_materials.id").
		Where("warehouse_stocks.warehouse_id = ?", warehouseID).
		Where("stock_materials.is_active = ?", true).
		Preload("Unit").
		Preload("StockMaterialCategory").
		Preload("Ingredient").
		Preload("Ingredient.IngredientCategory").
		Preload("Ingredient.Unit")

	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Where(`
			stock_materials.name ILIKE ? OR
			stock_materials.description ILIKE ? OR
			stock_materials.barcode ILIKE ?`,
			search, search, search)
	}

	if filter.IngredientID != nil {
		query = query.Where("stock_materials.ingredient_id = ?", *filter.IngredientID)
	}

	if filter.StockMaterialID != nil {
		query = query.Where("stock_materials.id = ?", *filter.StockMaterialID)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StockMaterial{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	if err := query.Find(&stockMaterials).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch available stock materials: %w", err)
	}

	return stockMaterials, nil
}
