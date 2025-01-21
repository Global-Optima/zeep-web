package warehouseStock

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type WarehouseStockRepository interface {
	RecordDeliveriesAndUpdateStock(delivery data.SupplierWarehouseDelivery, materials []data.SupplierWarehouseDeliveryMaterial, warehouseID uint) error
	TransferStock(sourceWarehouseID, targetWarehouseID uint, items []data.StockRequestIngredient) error

	GetDeliveryByID(deliveryID uint, delivery *data.SupplierWarehouseDelivery) error
	GetDeliveries(filter types.DeliveryFilter) ([]data.SupplierWarehouseDelivery, error)

	ConvertInventoryItemsToStockRequest(items []types.ReceiveWarehouseStockMaterial) ([]data.StockRequestIngredient, error)

	AddToWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error
	DeductFromWarehouseStock(warehouseID, stockMaterialID uint, quantityInPackages float64) error
	GetWarehouseStock(filter *types.GetWarehouseStockFilterQuery) ([]data.AggregatedWarehouseStock, error)
	GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID uint) (*data.AggregatedWarehouseStock, error)
	AddWarehouseStocks(warehouseID uint, stocks []data.WarehouseStock) error

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

	var unitQuantity float64
	if material.Package.UnitID == material.StockMaterial.UnitID {
		unitQuantity = material.Quantity * material.Package.Size
	} else {
		unitQuantity = material.Quantity * material.Package.Size * material.StockMaterial.Unit.ConversionFactor
	}

	if err := tx.FirstOrCreate(&stock, "warehouse_id = ? AND stock_material_id = ?", warehouseID, material.StockMaterialID).Error; err != nil {
		return fmt.Errorf("failed to find or create warehouse stock: %w", err)
	}

	if err := tx.Model(&data.WarehouseStock{}).
		Where("id = ?", stock.ID).
		Update("quantity", gorm.Expr("quantity + ?", unitQuantity)).Error; err != nil {
		return fmt.Errorf("failed to update warehouse stock quantity: %w", err)
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
	return r.db.Preload("Materials").
		Preload("Supplier").
		Preload("Warehouse").
		Preload("Materials.Package").
		Preload("Materials.Package.Unit").
		Preload("Materials.StockMaterial").
		Preload("Materials.StockMaterial.Unit").
		Preload("Materials.StockMaterial.Ingredient").
		Preload("Materials.StockMaterial.Ingredient.Unit").
		Preload("Materials.StockMaterial.Ingredient.IngredientCategory").
		Preload("Materials.StockMaterial.StockMaterialCategory").
		Preload("Materials.StockMaterial.Packages").
		Preload("Materials.StockMaterial.Packages.Unit").First(delivery, "id = ?", deliveryID).Error
}

func (r *warehouseStockRepository) GetDeliveries(filter types.DeliveryFilter) ([]data.SupplierWarehouseDelivery, error) {
	var deliveries []data.SupplierWarehouseDelivery
	query := r.db.Model(&data.SupplierWarehouseDelivery{}).
		Preload("Materials").
		Preload("Supplier").
		Preload("Warehouse").
		Preload("Materials.Package").
		Preload("Materials.Package.Unit").
		Preload("Materials.StockMaterial").
		Preload("Materials.StockMaterial.Unit").
		Preload("Materials.StockMaterial.Ingredient").
		Preload("Materials.StockMaterial.Ingredient.Unit").
		Preload("Materials.StockMaterial.Ingredient.IngredientCategory").
		Preload("Materials.StockMaterial.StockMaterialCategory").
		Preload("Materials.StockMaterial.Packages").
		Preload("Materials.StockMaterial.Packages.Unit").
		Joins("JOIN suppliers ON suppliers.id = supplier_warehouse_deliveries.supplier_id")

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.StartDate != nil {
		query = query.Where("delivery_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("delivery_date <= ?", *filter.EndDate)
	}

	if filter.SearchBySupplier != nil {
		search := "%" + *filter.SearchBySupplier + "%"
		query = query.Where("suppliers.name ILIKE ?", search)
	}

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
			IngredientID:   stockMaterial.IngredientID,
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

func (r *warehouseStockRepository) GetWarehouseStockMaterialDetails(stockMaterialID, warehouseID uint) (*data.AggregatedWarehouseStock, error) {
	warehouseStock, err := r.getWarehouseStock(stockMaterialID, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stock: %w", err)
	}

	earliestExpirationDate, err := r.findEarliestExpirationDateForStock(stockMaterialID, warehouseID)
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
	var earliest *time.Time
	for _, material := range materials {
		if earliest == nil || material.ExpirationDate.Before(*earliest) {
			earliest = &material.ExpirationDate
		}
	}
	UTCTime := utils.ToUTC(*earliest)
	return &UTCTime
}

func (r *warehouseStockRepository) findEarliestExpirationDateForStock(stockMaterialID, warehouseID uint) (*time.Time, error) {
	var earliestExpirationDate time.Time
	err := r.db.Model(&data.SupplierWarehouseDeliveryMaterial{}).
		Joins("JOIN supplier_warehouse_deliveries ON supplier_warehouse_deliveries.id = supplier_warehouse_delivery_materials.delivery_id").
		Where("supplier_warehouse_deliveries.warehouse_id = ? AND supplier_warehouse_delivery_materials.stock_material_id = ?", warehouseID, stockMaterialID).
		Select("MIN(supplier_warehouse_delivery_materials.expiration_date)").
		Scan(&earliestExpirationDate).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No expiration date found for this stock
		}
		return nil, fmt.Errorf("failed to fetch earliest expiration date for stock material ID %d: %w", stockMaterialID, err)
	}

	UTCTime := utils.ToUTC(earliestExpirationDate)
	return &UTCTime, nil
}

func (r *warehouseStockRepository) aggregateWarehouseStocks(
	warehouseStocks []data.WarehouseStock,
	materialMap map[uint][]data.SupplierWarehouseDeliveryMaterial,
) []data.AggregatedWarehouseStock {
	aggregatedStocks := []data.AggregatedWarehouseStock{}

	for _, stock := range warehouseStocks {
		materials := materialMap[stock.StockMaterialID]

		earliestExpirationDate := utils.ToUTC(*r.findEarliestMaterialExpirationDate(materials))

		aggregatedStocks = append(aggregatedStocks, data.AggregatedWarehouseStock{
			WarehouseID:            stock.WarehouseID,
			StockMaterialID:        stock.StockMaterialID,
			StockMaterial:          stock.StockMaterial,
			TotalQuantity:          stock.Quantity,
			EarliestExpirationDate: &earliestExpirationDate,
		})
	}

	return aggregatedStocks
}

func (r *warehouseStockRepository) getDeliveryMaterials(filter *types.GetWarehouseStockFilterQuery) ([]data.SupplierWarehouseDeliveryMaterial, error) {
	var materials []data.SupplierWarehouseDeliveryMaterial

	query := r.db.Model(&data.SupplierWarehouseDeliveryMaterial{}).
		Preload("StockMaterial").
		Preload("Package").
		Preload("Package.Unit")

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

func (r *warehouseStockRepository) getWarehouseStock(stockMaterialID, warehouseID uint) (*data.WarehouseStock, error) {
	var stock data.WarehouseStock
	err := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Ingredient").
		Preload("StockMaterial.Ingredient.Unit").
		Preload("StockMaterial.Ingredient.IngredientCategory").
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Packages").
		Preload("StockMaterial.Packages.Unit").
		Where("stock_material_id = ? AND warehouse_id = ?", stockMaterialID, warehouseID).
		First(&stock).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch single warehouse stock: %w", err)
	}
	return &stock, nil
}

func (r *warehouseStockRepository) getWarehouseStocksWithPagination(filter *types.GetWarehouseStockFilterQuery) ([]data.WarehouseStock, error) {
	var warehouseStocks []data.WarehouseStock

	query := r.db.Model(&data.WarehouseStock{}).
		Preload("StockMaterial").
		Preload("StockMaterial.Unit").
		Preload("StockMaterial.Ingredient").
		Preload("StockMaterial.Ingredient.Unit").
		Preload("StockMaterial.Ingredient.IngredientCategory").
		Preload("StockMaterial.StockMaterialCategory").
		Preload("StockMaterial.Packages").
		Preload("StockMaterial.Packages.Unit").
		Joins("LEFT JOIN supplier_warehouse_delivery_materials ON supplier_warehouse_delivery_materials.stock_material_id = warehouse_stocks.stock_material_id").
		Joins("JOIN stock_materials ON warehouse_stocks.stock_material_id = stock_materials.id")

	// Apply filters
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_stocks.warehouse_id = ?", *filter.WarehouseID)
	}

	if filter.StockMaterialID != nil {
		query = query.Where("stock_materials.stock_material_id = ?", *filter.StockMaterialID)
	}

	if filter.IngredientID != nil {
		query = query.Where("stock_materials.ingredient_id = ?", *filter.IngredientID)
	}

	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("warehouse_stocks.quantity <= stock_materials.safety_stock")
	}

	if filter.IsExpiring != nil && *filter.IsExpiring {
		query = query.Where("supplier_warehouse_delivery_materials.expiration_date <= stock_materials.expiration_period_in_days")
	}

	var days int
	if filter.ExpirationDays != nil {
		days = *filter.ExpirationDays
		expirationThreshold := time.Now().AddDate(0, 0, days)
		query = query.Where("supplier_warehouse_delivery_materials.expiration_date <= ?", expirationThreshold)
	}

	if filter.CategoryID != nil {
		query = query.Where("stock_materials.category_id = ?", *filter.CategoryID)
	}

	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Where("stock_materials.name ILIKE ? OR stock_materials.description ILIKE ? OR stock_materials.barcode ILIKE ?", search, search, search)
	}

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.WarehouseStock{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&warehouseStocks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch warehouse stocks: %w", err)
	}

	return warehouseStocks, nil
}

func (r *warehouseStockRepository) UpdateExpirationDate(stockMaterialID, warehouseID uint, newExpirationDate time.Time) error {
	return r.db.Model(&data.SupplierWarehouseDeliveryMaterial{}).
		Joins("JOIN supplier_warehouse_deliveries ON supplier_warehouse_deliveries.id = supplier_warehouse_delivery_materials.delivery_id").
		Where("supplier_warehouse_deliveries.warehouse_id = ? AND supplier_warehouse_delivery_materials.stock_material_id = ?", warehouseID, stockMaterialID).
		Update("expiration_date", newExpirationDate).Error
}

func (r *warehouseStockRepository) UpdateStockQuantity(stockID uint, quantity float64) error {
	return r.db.Model(&data.WarehouseStock{}).
		Where("id = ?", stockID).
		Update("quantity", quantity).Error
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
	if err := r.db.First(&stock).Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *warehouseStockRepository) UpdateWarehouseStock(stock *data.WarehouseStock) error {
	return r.db.Save(stock).Error
}
