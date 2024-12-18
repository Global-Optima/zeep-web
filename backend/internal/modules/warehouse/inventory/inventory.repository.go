package inventory

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory/types"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	LogIncomingInventory(deliveries []data.SupplierWarehouseDelivery) error
	RecordDeliveriesAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error
	LogAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error

	TransferStock(sourceWarehouseID, targetWarehouseID uint, items []data.StockRequestIngredient) error
	GetInventoryLevels(warehouseID uint) ([]data.WarehouseStock, error)
	PickupStock(storeWarehouseID uint, items []data.StockRequestIngredient) error

	GetDeliveryByID(deliveryID uint, delivery *data.SupplierWarehouseDelivery) error
	GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]data.SupplierWarehouseDelivery, error)
	GetExpiringItems(warehouseID uint, thresholdDays int) ([]data.SupplierWarehouseDelivery, error)
	ExtendExpiration(deliveryID uint, newExpirationDate time.Time) error

	ConvertInventoryItemsToStockRequest(items []types.InventoryItem) ([]data.StockRequestIngredient, error)
	ResolveIngredientID(stockMaterialID uint) (uint, error)
	SupplierMaterialExists(supplierID, stockMaterialID uint) (bool, error)
	CreateSupplierMaterial(association *data.SupplierMaterial) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) LogIncomingInventory(deliveries []data.SupplierWarehouseDelivery) error {
	return r.db.Create(&deliveries).Error
}

func (r *inventoryRepository) LogAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error {
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

func (r *inventoryRepository) RecordDeliveriesAndUpdateStock(deliveries []data.SupplierWarehouseDelivery, warehouseID uint) error {
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

func (r *inventoryRepository) recordDeliveries(tx *gorm.DB, deliveries []data.SupplierWarehouseDelivery) error {
	if err := tx.Create(&deliveries).Error; err != nil {
		return fmt.Errorf("failed to log deliveries: %w", err)
	}
	return nil
}

func (r *inventoryRepository) updateOrInsertStock(tx *gorm.DB, warehouseID uint, delivery data.SupplierWarehouseDelivery) error {
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

func (r *inventoryRepository) TransferStock(sourceWarehouseID, targetWarehouseID uint, items []data.StockRequestIngredient) error {
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

func (r *inventoryRepository) GetInventoryLevels(warehouseID uint) ([]data.WarehouseStock, error) {
	var stocks []data.WarehouseStock
	err := r.db.Preload("StockMaterial").Preload("Warehouse").
		Where("warehouse_id = ?", warehouseID).
		Find(&stocks).Error
	return stocks, err
}

func (r *inventoryRepository) PickupStock(storeWarehouseID uint, items []data.StockRequestIngredient) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	for _, item := range items {
		err := tx.Model(&data.StoreWarehouseStock{}).
			Where("store_warehouse_id = ? AND ingredient_id = ?", storeWarehouseID, item.IngredientID).
			Update("quantity", gorm.Expr("quantity - ?", item.Quantity)).Error
		if err != nil {
			return fmt.Errorf("failed to deduct stock from store warehouse: %w", err)
		}
	}

	return tx.Commit().Error
}

func (r *inventoryRepository) GetExpiringItems(warehouseID uint, thresholdDays int) ([]data.SupplierWarehouseDelivery, error) {
	var deliveries []data.SupplierWarehouseDelivery
	err := r.db.Preload("StockMaterial").Where("warehouse_id = ? AND expiration_date <= ?", warehouseID, time.Now().AddDate(0, 0, thresholdDays)).
		Find(&deliveries).Error
	return deliveries, err
}

func (r *inventoryRepository) ExtendExpiration(deliveryID uint, newExpirationDate time.Time) error {
	return r.db.Model(&data.SupplierWarehouseDelivery{}).
		Where("id = ?", deliveryID).
		Update("expiration_date", newExpirationDate).Error
}

func (r *inventoryRepository) GetDeliveryByID(deliveryID uint, delivery *data.SupplierWarehouseDelivery) error {
	return r.db.First(delivery, "id = ?", deliveryID).Error
}

func (r *inventoryRepository) GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]data.SupplierWarehouseDelivery, error) {
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

func (r *inventoryRepository) ConvertInventoryItemsToStockRequest(items []types.InventoryItem) ([]data.StockRequestIngredient, error) {
	converted := make([]data.StockRequestIngredient, len(items))

	for i, item := range items {

		var stockMaterial data.StockMaterial
		if err := r.db.First(&stockMaterial, item.StockMaterialID).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve stock material for StockMaterialID %d: %w", item.StockMaterialID, err)
		}

		var mapping data.IngredientStockMaterialMapping
		if err := r.db.Where("stock_material_id = ?", stockMaterial.ID).First(&mapping).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve ingredient mapping for StockMaterialID %d: %w", stockMaterial.ID, err)
		}

		deliveredDate := time.Now()
		expirationDate := deliveredDate.AddDate(0, 0, stockMaterial.ExpirationPeriodInDays)

		converted[i] = data.StockRequestIngredient{
			IngredientID:   mapping.IngredientID,
			Quantity:       item.Quantity,
			DeliveredDate:  deliveredDate,
			ExpirationDate: expirationDate,
		}
	}

	return converted, nil
}

func (r *inventoryRepository) ResolveIngredientID(stockMaterialID uint) (uint, error) {
	var mapping data.IngredientStockMaterialMapping
	err := r.db.Where("stock_material_id = ?", stockMaterialID).First(&mapping).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no mapping found for stock material ID %d", stockMaterialID)
		}
		return 0, fmt.Errorf("failed to resolve ingredient ID for stock material ID %d: %w", stockMaterialID, err)
	}

	return mapping.IngredientID, nil
}

func (r *inventoryRepository) SupplierMaterialExists(supplierID, stockMaterialID uint) (bool, error) {
	var count int64
	err := r.db.Model(&data.SupplierMaterial{}).
		Where("supplier_id = ? AND stock_material_id = ?", supplierID, stockMaterialID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *inventoryRepository) CreateSupplierMaterial(association *data.SupplierMaterial) error {
	return r.db.Create(association).Error
}
