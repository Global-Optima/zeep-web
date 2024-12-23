package stockRequests

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"gorm.io/gorm"
)

type StockRequestRepository interface {
	CreateStockRequest(stockRequest *data.StockRequest) error
	AddIngredientsToStockRequest(ingredients []data.StockRequestIngredient) error
	DeleteStockRequestIngredient(ingredientID uint) error
	UpdateStockRequestIngredientDates(dates *types.UpdateIngredientDates) error

	GetStockRequests(filter types.StockRequestFilter) ([]data.StockRequest, error)
	GetStockRequestByID(requestID uint) (*data.StockRequest, error)
	UpdateStockRequestStatus(stockRequest *data.StockRequest) error
	GetMappingByIngredientID(ingredientID uint, mapping *data.IngredientStockMaterialMapping) error
	DeductWarehouseStock(stockMaterialID, warehouseID uint, quantity float64) error
	AddToStoreWarehouseStock(storeWarehouseID, ingredientID uint, quantity float64) error
	GetLastStockRequestDate(storeWarehouseID uint) (*time.Time, error)
	GetLowStockIngredients(storeWarehouseID uint) ([]data.StoreWarehouseStock, error)
	GetAllStockMaterials(storeID uint, filters types.StockMaterialFilter) ([]types.StockMaterialDTO, error)
	GetWarehouseStockQuantity(warehouseID, stockMaterialID uint) (float64, error)
	GetStoreWarehouse(storeID uint) (*data.StoreWarehouse, error)
}

type stockRequestRepository struct {
	db *gorm.DB
}

func NewStockRequestRepository(db *gorm.DB) StockRequestRepository {
	return &stockRequestRepository{db: db}
}

func (r *stockRequestRepository) CreateStockRequest(stockRequest *data.StockRequest) error {
	return r.db.Create(stockRequest).Error
}

func (r *stockRequestRepository) AddIngredientsToStockRequest(ingredients []data.StockRequestIngredient) error {
	return r.db.Create(&ingredients).Error
}

func (r *stockRequestRepository) DeleteStockRequestIngredient(ingredientID uint) error {
	return r.db.Delete(&data.StockRequestIngredient{}, ingredientID).Error
}

func (r *stockRequestRepository) UpdateStockRequestIngredientDates(dates *types.UpdateIngredientDates) error {
	return r.db.Model(&data.StockRequestIngredient{}).
		Updates(dates).Error
}

func (r *stockRequestRepository) GetStockRequests(filter types.StockRequestFilter) ([]data.StockRequest, error) {
	var requests []data.StockRequest
	query := r.db.Preload("Ingredients.Ingredient.Unit").
		Preload("Ingredients.Ingredient.IngredientCategory").
		Preload("Store").
		Preload("Warehouse")

	if filter.StoreID != nil {
		query = query.Where("store_id = ?", *filter.StoreID)
	}
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	err := query.Find(&requests).Error
	return requests, err
}

func (r *stockRequestRepository) GetStockRequestByID(requestID uint) (*data.StockRequest, error) {
	var request data.StockRequest
	err := r.db.Preload("Ingredients.Ingredient.Unit").
		Preload("Ingredients.Ingredient.IngredientCategory").
		Preload("Store").
		Preload("Warehouse").
		Where("id = ?", requestID).
		First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *stockRequestRepository) UpdateStockRequestStatus(stockRequest *data.StockRequest) error {
	return r.db.Model(stockRequest).Update("status", stockRequest.Status).Error
}

func (r *stockRequestRepository) GetMappingByIngredientID(ingredientID uint, mapping *data.IngredientStockMaterialMapping) error {
	return r.db.Where("ingredient_id = ?", ingredientID).First(mapping).Error
}

func (r *stockRequestRepository) DeductWarehouseStock(stockMaterialID, warehouseID uint, quantity float64) error {
	return r.db.Model(&data.WarehouseStock{}).
		Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
		Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

func (r *stockRequestRepository) AddToStoreWarehouseStock(storeWarehouseID, ingredientID uint, quantity float64) error {
	return r.db.Model(&data.StoreWarehouseStock{}).
		Where("store_warehouse_id = ? AND ingredient_id = ?", storeWarehouseID, ingredientID).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
}

func (r *stockRequestRepository) GetLastStockRequestDate(storeID uint) (*time.Time, error) {
	var lastRequest data.StockRequest
	err := r.db.Select("created_at").Where("store_id = ?", storeID).
		Order("created_at DESC").First(&lastRequest).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return lastRequest.RequestDate, err
}

func (r *stockRequestRepository) GetLowStockIngredients(storeWarehouseID uint) ([]data.StoreWarehouseStock, error) {
	var stocks []data.StoreWarehouseStock
	err := r.db.Preload("Ingredient.Unit").
		Where("store_warehouse_id = ? AND quantity <= low_stock_threshold", storeWarehouseID).
		Find(&stocks).Error
	return stocks, err
}

func (r *stockRequestRepository) GetAllStockMaterials(storeID uint, filters types.StockMaterialFilter) ([]types.StockMaterialDTO, error) {
	storeWarehouse, err := r.GetStoreWarehouse(storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve store warehouse for store ID %d: %w", storeID, err)
	}

	var mappings []data.IngredientStockMaterialMapping
	query := r.db.Preload("Ingredient.Unit").Preload("StockMaterial")

	if filters.Category != nil {
		query = query.Where("stock_materials.category = ?", *filters.Category)
	}

	if filters.Search != nil {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("stock_materials.name ILIKE ? OR stock_materials.barcode ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Find(&mappings).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve ingredient mappings: %w", err)
	}

	products := []types.StockMaterialDTO{}
	for _, mapping := range mappings {
		quantity, err := r.GetWarehouseStockQuantity(storeWarehouse.WarehouseID, mapping.StockMaterialID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve stock quantity for stock material ID %d: %w", mapping.StockMaterialID, err)
		}

		products = append(products, types.StockMaterialDTO{
			StockMaterialID: mapping.StockMaterialID,
			Name:            mapping.StockMaterial.Name,
			Category:        mapping.StockMaterial.Category,
			Unit:            mapping.Ingredient.Unit.Name,
			AvailableQty:    quantity,
		})
	}

	return products, nil
}

func (r *stockRequestRepository) GetWarehouseStockQuantity(warehouseID, stockMaterialID uint) (float64, error) {
	var stock data.WarehouseStock
	err := r.db.Select("quantity").
		Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
		First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // No stock found, return 0 quantity
		}
		return 0, err
	}
	return stock.Quantity, nil
}

func (r *stockRequestRepository) GetStoreWarehouse(storeID uint) (*data.StoreWarehouse, error) {
	var storeWarehouse data.StoreWarehouse
	err := r.db.Where("store_id = ?", storeID).First(&storeWarehouse).Error
	if err != nil {
		return nil, err
	}
	return &storeWarehouse, nil
}
