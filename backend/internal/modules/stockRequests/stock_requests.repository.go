package stockRequests

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"gorm.io/gorm"
)

type StockRequestRepository interface {
	CreateStockRequest(stockRequest *data.StockRequest) error
	AddIngredientsToStockRequest(ingredients []data.StockRequestIngredient) error
	DeleteStockRequestIngredient(ingredientID uint) error
	ReplaceStockRequestIngredients(requestID uint, ingredients []data.StockRequestIngredient) error
	UpdateStockRequestIngredientDates(stockRequestIngredientID uint, dates *types.UpdateIngredientDates) error

	GetStockRequests(filter types.GetStockRequestsFilter) ([]data.StockRequest, error)
	GetStockRequestByID(requestID uint) (*data.StockRequest, error)
	UpdateStockRequestStatus(stockRequest *data.StockRequest) error
	AddStoreComment(requestID uint, comment string) error
	AddWarehouseComment(requestID uint, comment string) error

	DeductWarehouseStock(stockMaterialID, warehouseID uint, quantityInPackages float64) error
	AddToStoreWarehouseStock(storeWarehouseID, stockMaterialID uint, quantityInPackages float64) error
	GetWarehouseStockQuantity(warehouseID, stockMaterialID uint) (float64, error)
	GetStoreWarehouse(storeID uint) (*data.StoreWarehouse, error)
	GetLowStockIngredients(storeWarehouseID uint) ([]data.StoreWarehouseStock, error)

	GetLastStockRequestDate(storeWarehouseID uint) (*time.Time, error)

	GetAllStockMaterials(storeID uint, filters types.StockMaterialFilter) ([]types.StockMaterialDTO, error)
	GetStockMaterialsByIngredient(ingredientID uint, warehouseID *uint) ([]data.WarehouseStock, error)
	GetStockMaterialByID(stockMaterialID uint, stockMaterial *data.StockMaterial) error
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
		Updates(dates).Error
}

func (r *stockRequestRepository) GetStockRequests(filter types.GetStockRequestsFilter) ([]data.StockRequest, error) {
	var requests []data.StockRequest
	query := r.db.Model(&data.StockRequest{}).
		Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
			return db.Preload("StockMaterial", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "category_id", "barcode").
					Preload("StockMaterialCategory", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					}).
					Preload("Package", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "size", "unit_id", "stock_material_id").
							Preload("Unit", func(db *gorm.DB) *gorm.DB {
								return db.Select("id", "name")
							})
					})
			}).Preload("Ingredient", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "category_id", "unit_id").
					Preload("IngredientCategory", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					}).
					Preload("Unit", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					})
			})
		}).
		Preload("Store", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "facility_address_id").
				Preload("FacilityAddress", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "address")
				})
		}).
		Preload("Warehouse", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		})

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

	query = query.Order("created_at DESC")

	var err error
	query, err = utils.ApplyPagination(query, filter.Pagination, &data.StockRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&requests).Error
	return requests, err
}

func (r *stockRequestRepository) GetStockRequestByID(requestID uint) (*data.StockRequest, error) {
	var request data.StockRequest
	err := r.db.Model(&data.StockRequest{}).
		Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
			return db.Preload("StockMaterial", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "category_id", "barcode").
					Preload("StockMaterialCategory", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					}).
					Preload("Package", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "size", "unit_id", "stock_material_id").
							Preload("Unit", func(db *gorm.DB) *gorm.DB {
								return db.Select("id", "name")
							})
					})
			}).Preload("Ingredient", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "category_id", "unit_id").
					Preload("IngredientCategory", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					}).
					Preload("Unit", func(db *gorm.DB) *gorm.DB {
						return db.Select("id", "name")
					})
			})
		}).
		Preload("Store", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "facility_address_id").
				Preload("FacilityAddress", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "address")
				})
		}).
		Preload("Warehouse", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Where("id = ?", requestID).
		First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *stockRequestRepository) UpdateStockRequestStatus(stockRequest *data.StockRequest) error {
	return r.db.Model(&data.StockRequest{}).Where("id = ?", stockRequest.ID).Update("status", stockRequest.Status).Error
}

func (r *stockRequestRepository) DeductWarehouseStock(stockMaterialID, warehouseID uint, quantity float64) error {
	return r.db.Model(&data.WarehouseStock{}).
		Where("warehouse_id = ? AND stock_material_id = ?", warehouseID, stockMaterialID).
		Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

func (r *stockRequestRepository) AddToStoreWarehouseStock(storeWarehouseID, stockMaterialID uint, quantityInPackages float64) error {
	var stockMaterial data.StockMaterial
	if err := r.db.Preload("Package.Unit").Preload("Ingredient.Unit").
		Where("id = ?", stockMaterialID).First(&stockMaterial).Error; err != nil {
		return fmt.Errorf("failed to fetch stock material details for ID %d: %w", stockMaterialID, err)
	}

	quantityInUnits := utils.ConvertPackagesToUnits(stockMaterial, quantityInPackages)

	return r.db.Model(&data.StoreWarehouseStock{}).
		Where("store_warehouse_id = ? AND ingredient_id = ?", storeWarehouseID, stockMaterial.IngredientID).
		Update("quantity", gorm.Expr("quantity + ?", quantityInUnits)).Error
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

	var stockMaterials []data.StockMaterial
	query := r.db.Preload("Unit").Preload("Ingredient").Preload("StockMaterialCategory")

	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}

	if filters.Search != nil {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("name ILIKE ? OR barcode ILIKE ?", searchTerm, searchTerm)
	}

	if err := query.Find(&stockMaterials).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve stock materials: %w", err)
	}

	products := []types.StockMaterialDTO{}
	for _, stockMaterial := range stockMaterials {
		quantity, err := r.GetWarehouseStockQuantity(storeWarehouse.WarehouseID, stockMaterial.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve stock quantity for StockMaterialID %d: %w", stockMaterial.ID, err)
		}

		products = append(products, types.StockMaterialDTO{
			StockMaterialID: stockMaterial.ID,
			Name:            stockMaterial.Name,
			Category:        stockMaterial.StockMaterialCategory.Name,
			Unit:            stockMaterial.Unit.Name,
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
			return 0, nil
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

func (r *stockRequestRepository) GetStockMaterialsByIngredient(ingredientID uint, warehouseID *uint) ([]data.WarehouseStock, error) {
	var stocks []data.WarehouseStock

	query := r.db.Preload("StockMaterial.Unit").
		Preload("StockMaterial.StockMaterialCategory").
		Preload("Warehouse").
		Joins("JOIN stock_materials ON warehouse_stocks.stock_material_id = stock_materials.id").
		Where("stock_materials.ingredient_id = ?", ingredientID)

	if warehouseID != nil {
		query = query.Where("warehouse_stocks.warehouse_id = ?", *warehouseID)
	}

	err := query.Find(&stocks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock materials for ingredient ID %d: %w", ingredientID, err)
	}

	return stocks, nil
}

func (r *stockRequestRepository) GetStockMaterialByID(stockMaterialID uint, stockMaterial *data.StockMaterial) error {
	return r.db.Preload("Ingredient").Preload("StockMaterialCategory").First(stockMaterial, "id = ?", stockMaterialID).Error
}

func (r *stockRequestRepository) ReplaceStockRequestIngredients(requestID uint, ingredients []data.StockRequestIngredient) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("stock_request_id = ?", requestID).Delete(&data.StockRequestIngredient{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing ingredients for stock request ID %d: %w", requestID, err)
		}

		if err := tx.Create(&ingredients).Error; err != nil {
			return fmt.Errorf("failed to add new ingredients for stock request ID %d: %w", requestID, err)
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
