package stockRequests

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const DefaultLowStockThreshold = 100

type StockRequestRepository interface {
	CreateStockRequest(stockRequest *data.StockRequest) error
	AddIngredientsToStockRequest(ingredients []data.StockRequestIngredient) error
	DeleteStockRequestIngredient(ingredientID uint) error
	ReplaceStockRequestIngredients(request data.StockRequest, ingredients []data.StockRequestIngredient) error
	UpdateStockRequestIngredientDates(stockRequestIngredientID uint, dates *types.UpdateIngredientDates) error

	GetStockRequests(filter types.GetStockRequestsFilter) ([]data.StockRequest, error)
	GetStockRequestByID(requestID uint) (*data.StockRequest, error)
	UpdateStockRequestStatus(stockRequest *data.StockRequest) error
	AddStoreComment(requestID uint, comment string) error
	AddWarehouseComment(requestID uint, comment string) error
	AddDetails(stockRequestID uint, details []types.StockRequestDetails) error

	DeductWarehouseStock(stockMaterialID, warehouseID uint, quantityInPackages float64) (*data.WarehouseStock, error)
	UpsertToStoreStock(storeID, stockMaterialID uint, quantityInPackages float64) error
	ReturnWarehouseStock(stockMaterialID, warehouseID uint, quantity float64) (*data.WarehouseStock, error)
	GetWarehouseStockQuantity(warehouseID, stockMaterialID uint) (float64, error)
	GetStoreWarehouse(storeID uint) (*data.Store, error)

	GetLastStockRequestDate(storeID uint) (*time.Time, error)

	DeleteStockRequest(requestID uint) error
	GetOpenCartByStoreID(storeID uint) (*data.StockRequest, error)
	UpdateStockRequestIngredientQuantity(ingredientID uint, quantity float64) error

	CountStockRequestsInLast24Hours(storeID uint) (int64, error)
	CountFinalizedStockRequestsInLast24Hours(storeID uint) (int64, error)
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

func (r *stockRequestRepository) UpdateStockRequestIngredientDates(stockRequestIngredientID uint, dates *types.UpdateIngredientDates) error {
	return r.db.Model(&data.StockRequestIngredient{}).
		Where("id = ?", stockRequestIngredientID).
		Updates(&data.StockRequestIngredient{
			DeliveredDate:  dates.DeliveredDate,
			ExpirationDate: dates.ExpirationDate,
		}).Error
}

func (r *stockRequestRepository) GetStockRequests(filter types.GetStockRequestsFilter) ([]data.StockRequest, error) {
	var requests []data.StockRequest
	query := r.db.Model(&data.StockRequest{}).
		Preload("Warehouse.FacilityAddress").
		Preload("Store.FacilityAddress").
		Preload("Ingredients.StockMaterial").
		Preload("Ingredients.StockMaterial.StockMaterialCategory").
		Preload("Ingredients.StockMaterial.Ingredient").
		Preload("Ingredients.StockMaterial.Ingredient.Unit").
		Preload("Ingredients.StockMaterial.Ingredient.IngredientCategory").
		Preload("Ingredients.StockMaterial.Unit")

	// Apply basic filters
	if filter.StoreID != nil {
		query = query.Where("store_id = ?", *filter.StoreID)
	}
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if len(filter.Statuses) > 0 {
		query = query.Where("status IN (?)", filter.Statuses)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Incorporate search logic via a subquery so we don't perform extra round trips.
	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		subQuery := r.db.Model(&data.StockRequest{}).
			Joins("JOIN stock_request_ingredients sri ON sri.stock_request_id = stock_requests.id").
			Joins("JOIN stock_materials sm ON sri.stock_material_id = sm.id").
			Where(`
				sm.name ILIKE ? OR
				sm.description ILIKE ? OR
				sm.barcode ILIKE ?
			`, search, search, search).
			Group("stock_requests.id").
			Select("stock_requests.id")
		query = query.Where("stock_requests.id IN (?)", subQuery)
	}

	// Order results and apply pagination.
	query = query.Order("stock_requests.created_at DESC")
	var err error
	query, err = utils.ApplyPagination(query, filter.Pagination, &data.StockRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&requests).Error
	return requests, err
}

func (r *stockRequestRepository) GetStockRequestByID(requestID uint) (*data.StockRequest, error) {
	var stockRequest data.StockRequest

	err := r.db.Model(&data.StockRequest{}).
		Preload("Warehouse").
		Preload("Store").
		Preload("Warehouse.FacilityAddress").
		Preload("Store.FacilityAddress").
		Preload("Ingredients.StockMaterial").
		Preload("Ingredients.StockMaterial.StockMaterialCategory").
		Preload("Ingredients.StockMaterial.Ingredient").
		Preload("Ingredients.StockMaterial.Ingredient.Unit").
		Preload("Ingredients.StockMaterial.Ingredient.IngredientCategory").
		Preload("Ingredients.StockMaterial.Unit").
		First(&stockRequest, requestID).
		Error
	if err != nil {
		return nil, err
	}

	return &stockRequest, nil
}

func (r *stockRequestRepository) UpdateStockRequestStatus(stockRequest *data.StockRequest) error {
	return r.db.Model(&data.StockRequest{}).Where("id = ?", stockRequest.ID).Update("status", stockRequest.Status).Error
}

func (r *stockRequestRepository) DeductWarehouseStock(stockMaterialID, warehouseID uint, quantity float64) (*data.WarehouseStock, error) {
	var updatedStock data.WarehouseStock

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&data.WarehouseStock{}).
			Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
			Update("quantity", gorm.Expr("quantity - ?", quantity)).Error; err != nil {
			return fmt.Errorf("failed to update stock quantity: %w", err)
		}
		if err := tx.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
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

func (r *stockRequestRepository) UpsertToStoreStock(storeID, stockMaterialID uint, quantityInPackages float64) error {
	var sm data.StockMaterial
	if err := r.db.Preload("Ingredient.Unit").First(&sm, stockMaterialID).Error; err != nil {
		return fmt.Errorf("failed to fetch stock material details for ID %d: %w", stockMaterialID, err)
	}

	q := sm.Size * quantityInPackages
	if sm.Ingredient.UnitID != sm.UnitID {
		q = sm.Ingredient.Unit.ConversionFactor * q
	}

	var storeStock data.StoreStock
	err := r.db.Where("store_id = ? AND ingredient_id = ?", storeID, sm.IngredientID).First(&storeStock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newStock := data.StoreStock{
				StoreID:           storeID,
				IngredientID:      sm.IngredientID,
				Quantity:          q,
				LowStockThreshold: DefaultLowStockThreshold,
			}
			return r.db.Create(&newStock).Error
		}
		return err
	}

	storeStock.Quantity += q
	return r.db.Save(&storeStock).Error
}

func (r *stockRequestRepository) GetLastStockRequestDate(storeID uint) (*time.Time, error) {
	var count int64
	var lastRequest data.StockRequest

	err := r.db.Model(&data.StockRequest{}).Where("store_id = ?", storeID).Count(&count).Error
	if err != nil {
		return nil, err
	}

	if count <= 1 {
		return nil, nil
	}

	err = r.db.Select("created_at").Where("store_id = ?", storeID).
		Order("created_at DESC").First(&lastRequest).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &lastRequest.CreatedAt, err
}

func (r *stockRequestRepository) GetWarehouseStockQuantity(warehouseID, stockMaterialID uint) (float64, error) {
	var stock data.WarehouseStock
	err := r.db.Select("quantity").
		Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
		First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return stock.Quantity, nil
}

func (r *stockRequestRepository) GetStoreWarehouse(storeID uint) (*data.Store, error) {
	var storeWarehouse data.Store
	err := r.db.Model(&data.Store{}).Preload("Warehouse").Where("id = ?", storeID).First(&storeWarehouse).Error
	if err != nil {
		return nil, err
	}
	return &storeWarehouse, nil
}

func (r *stockRequestRepository) ReplaceStockRequestIngredients(request data.StockRequest, ingredients []data.StockRequestIngredient) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var warehouseStocks []data.WarehouseStock
		err := tx.Model(&data.WarehouseStock{}).Where("warehouse_id = ?", request.WarehouseID).Find(&warehouseStocks).Error
		if err != nil {
			return fmt.Errorf("failed to fetch the warehouse stock: %w", err)
		}

		warehouseStockMap := make(map[uint]*data.WarehouseStock)
		for _, stock := range warehouseStocks {
			warehouseStockMap[stock.StockMaterialID] = &stock
		}

		// Ensure all requested materials exist in the warehouse
		for _, ingredient := range ingredients {
			if warehouseStockMap[ingredient.StockMaterialID] == nil {
				return fmt.Errorf("material ID %d is not present in warehouse ID %d", ingredient.StockMaterialID, request.WarehouseID)
			}
		}

		if err := tx.Where("stock_request_id = ?", request.ID).Delete(&data.StockRequestIngredient{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing ingredients for stock request ID %d: %w", request.ID, err)
		}

		if err := tx.Create(&ingredients).Error; err != nil {
			return fmt.Errorf("failed to add new ingredients for stock request ID %d: %w", request.ID, err)
		}

		return nil
	})
}

func (r *stockRequestRepository) AddStoreComment(requestID uint, comment string) error {
	return r.db.Model(&data.StockRequest{}).
		Where("id = ?", requestID).
		Update("store_comment", comment).Error
}

func (r *stockRequestRepository) AddWarehouseComment(requestID uint, comment string) error {
	return r.db.Model(&data.StockRequest{}).
		Where("id = ?", requestID).
		Update("warehouse_comment", comment).Error
}

func (r *stockRequestRepository) DeleteStockRequest(requestID uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("stock_request_id = ?", requestID).Delete(&data.StockRequestIngredient{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to delete stock request ingredients: %w", err)
		}

		if err := tx.Where("id = ?", requestID).Delete(&data.StockRequest{}).Error; err != nil {
			return fmt.Errorf("failed to delete stock request: %w", err)
		}

		return nil
	})
	return err
}

func (r *stockRequestRepository) GetOpenCartByStoreID(storeID uint) (*data.StockRequest, error) {
	var request data.StockRequest
	err := r.db.Model(&data.StockRequest{}).
		Preload("Store").
		Preload("Warehouse").
		Preload("Warehouse.FacilityAddress").
		Preload("Store.FacilityAddress").
		Preload("Ingredients.StockMaterial").
		Preload("Ingredients.StockMaterial.StockMaterialCategory").
		Preload("Ingredients.StockMaterial.Ingredient").
		Preload("Ingredients.StockMaterial.Ingredient.Unit").
		Preload("Ingredients.StockMaterial.Ingredient.IngredientCategory").
		Preload("Ingredients.StockMaterial.Unit").
		Where("store_id = ? AND status = ?", storeID, data.StockRequestCreated).
		First(&request).Error
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *stockRequestRepository) UpdateStockRequestIngredientQuantity(ingredientID uint, quantity float64) error {
	return r.db.Model(&data.StockRequestIngredient{}).
		Where("id = ?", ingredientID).
		Update("quantity", quantity).
		Error
}

func (r *stockRequestRepository) AddDetails(stockRequestID uint, newDetails []types.StockRequestDetails) error {
	updatedDetailsJSON, err := json.Marshal(newDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal updated details: %w", err)
	}

	err = r.db.Model(&data.StockRequest{}).
		Where("id = ?", stockRequestID).
		Update("details", datatypes.JSON(updatedDetailsJSON)).Error
	if err != nil {
		return fmt.Errorf("failed to update details for stock request ID %d: %w", stockRequestID, err)
	}

	return nil
}

func (r *stockRequestRepository) ReturnWarehouseStock(stockMaterialID, warehouseID uint, quantity float64) (*data.WarehouseStock, error) {
	var updatedStock data.WarehouseStock
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&data.WarehouseStock{}).
			Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
			Update("quantity", gorm.Expr("quantity + ?", quantity)).Error; err != nil {
			return fmt.Errorf("failed to return stock quantity: %w", err)
		}
		if err := tx.Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
			First(&updatedStock).Error; err != nil {
			return fmt.Errorf("failed to fetch updated warehouse stock: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updatedStock, nil
}

func (r *stockRequestRepository) CountStockRequestsInLast24Hours(storeID uint) (int64, error) {
	var count int64
	threshold := time.Now().Add(-24 * time.Hour)
	err := r.db.Model(&data.StockRequest{}).
		Where("store_id = ? AND created_at >= ?", storeID, threshold).
		Count(&count).Error
	return count, err
}

func (r *stockRequestRepository) CountFinalizedStockRequestsInLast24Hours(storeID uint) (int64, error) {
	var count int64
	threshold := time.Now().Add(-24 * time.Hour)
	finalizedStatuses := []data.StockRequestStatus{
		data.StockRequestProcessed,
		data.StockRequestInDelivery,
		data.StockRequestCompleted,
		data.StockRequestRejectedByStore,
		data.StockRequestRejectedByWarehouse,
		data.StockRequestAcceptedWithChange,
	}
	err := r.db.Model(&data.StockRequest{}).
		Where("store_id = ? AND created_at >= ? AND status IN (?)", storeID, threshold, finalizedStatuses).
		Count(&count).Error
	return count, err
}
