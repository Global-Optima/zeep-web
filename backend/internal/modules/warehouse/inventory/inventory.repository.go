package inventory

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory/types"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	LogIncomingInventory(deliveries []data.Delivery) error
	UpdateStockLevels(warehouseID uint, items []data.StockRequestIngredient) error
	TransferStock(sourceWarehouseID, targetWarehouseID uint, items []data.StockRequestIngredient) error
	GetInventoryLevels(warehouseID uint) ([]data.StoreWarehouseStock, error)
	PickupStock(storeWarehouseID uint, items []data.StockRequestIngredient) error
	GetExpiringItems(warehouseID uint, thresholdDays int) ([]data.Delivery, error)
	ExtendExpiration(deliveryID uint, newExpirationDate time.Time) error
	GetDeliveryByID(deliveryID uint, delivery *data.Delivery) error
	ConvertInventoryItemsToStockRequest(items []types.InventoryItem) ([]data.StockRequestIngredient, error)
	GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]data.Delivery, error)
	CreateAuditLog(log data.AuditLog) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) LogIncomingInventory(deliveries []data.Delivery) error {
	return r.db.Create(&deliveries).Error
}

func (r *inventoryRepository) UpdateStockLevels(warehouseID uint, items []data.StockRequestIngredient) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	for _, item := range items {
		err := tx.Model(&data.StoreWarehouseStock{}).
			Where("store_warehouse_id = ? AND ingredient_id = ?", warehouseID, item.IngredientID).
			Update("quantity", gorm.Expr("quantity + ?", item.Quantity)).Error
		if err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (r *inventoryRepository) TransferStock(sourceWarehouseID, targetWarehouseID uint, items []data.StockRequestIngredient) error {
	tx := r.db.Begin()
	defer tx.Rollback()

	for _, item := range items {
		var pkg data.Package
		if err := r.db.Where("sku_id = ?", item.IngredientID).First(&pkg).Error; err != nil {
			return fmt.Errorf("failed to retrieve package for SKU %d: %w", item.IngredientID, err)
		}

		var unit data.Unit
		if err := r.db.First(&unit, pkg.PackageUnitID).Error; err != nil {
			return fmt.Errorf("failed to retrieve unit for SKU %d: %w", item.IngredientID, err)
		}

		measurementQuantity := item.Quantity * pkg.PackageSize * unit.ConversionFactor

		if err := tx.Model(&data.StoreWarehouseStock{}).
			Where("store_warehouse_id = ? AND ingredient_id = ?", sourceWarehouseID, item.IngredientID).
			Update("quantity", gorm.Expr("quantity - ?", item.Quantity)).Error; err != nil {
			return err
		}

		if err := tx.Model(&data.StoreWarehouseStock{}).
			Where("store_warehouse_id = ? AND ingredient_id = ?", targetWarehouseID, item.IngredientID).
			Update("quantity", gorm.Expr("quantity + ?", measurementQuantity)).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (r *inventoryRepository) GetInventoryLevels(warehouseID uint) ([]data.StoreWarehouseStock, error) {
	var stocks []data.StoreWarehouseStock
	err := r.db.Where("store_warehouse_id = ?", warehouseID).Find(&stocks).Error
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
			return err
		}
	}

	return tx.Commit().Error
}

func (r *inventoryRepository) GetExpiringItems(warehouseID uint, thresholdDays int) ([]data.Delivery, error) {
	var deliveries []data.Delivery
	err := r.db.Where("target = ? AND expiration_date <= ?", warehouseID, time.Now().AddDate(0, 0, thresholdDays)).
		Find(&deliveries).Error
	return deliveries, err
}

func (r *inventoryRepository) ExtendExpiration(deliveryID uint, newExpirationDate time.Time) error {
	return r.db.Model(&data.Delivery{}).
		Where("id = ?", deliveryID).
		Update("expiration_date", newExpirationDate).Error
}

func (r *inventoryRepository) GetDeliveryByID(deliveryID uint, delivery *data.Delivery) error {
	return r.db.First(delivery, "id = ?", deliveryID).Error
}

func (r *inventoryRepository) ConvertInventoryItemsToStockRequest(items []types.InventoryItem) ([]data.StockRequestIngredient, error) {
	converted := make([]data.StockRequestIngredient, len(items))

	for i, item := range items {
		var pkg data.Package
		if err := r.db.Where("sku_id = ?", item.SKU_ID).First(&pkg).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve package for SKU_ID %d: %w", item.SKU_ID, err)
		}

		var unit data.Unit
		if err := r.db.First(&unit, pkg.PackageUnitID).Error; err != nil {
			return nil, fmt.Errorf("failed to retrieve unit for SKU_ID %d: %w", item.SKU_ID, err)
		}

		measurementQuantity := item.Quantity * pkg.PackageSize * unit.ConversionFactor

		converted[i] = data.StockRequestIngredient{
			IngredientID: item.SKU_ID,
			Quantity:     measurementQuantity,
		}
	}

	return converted, nil
}

func (r *inventoryRepository) GetDeliveries(warehouseID *uint, startDate, endDate *time.Time) ([]data.Delivery, error) {
	var deliveries []data.Delivery
	query := r.db.Model(&data.Delivery{})

	if warehouseID != nil {
		query = query.Where("target = ?", *warehouseID)
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

func (r *inventoryRepository) CreateAuditLog(log data.AuditLog) error {
	return r.db.Create(&log).Error
}
