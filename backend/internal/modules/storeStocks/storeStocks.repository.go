package storeStocks

import (
	"database/sql"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"

	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StoreStockRepository interface {
	GetAvailableIngredientsToAdd(storeID uint, filter *ingredientTypes.IngredientFilter) ([]data.Ingredient, error)
	AddStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error)
	AddOrUpdateStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error)
	GetStockList(storeID uint, query *types.GetStockFilterQuery) ([]data.StoreStock, error)
	GetStockListByIDs(storeID uint, IDs []uint) ([]data.StoreStock, error)
	GetStockById(storeID, stockID uint) (*data.StoreStock, error)
	GetAllStockList(storeID uint) ([]data.StoreStock, error)
	UpdateStock(storeID, stockID uint, dto *types.UpdateStoreStockDTO) error
	DeleteStockById(storeID, stockID uint) error
	WithTransaction(txFunc func(txRepo storeStockRepository) error) error
	CloneWithTransaction(tx *gorm.DB) storeStockRepository
	DeductStockByProductSizeTechCart(storeID, productSizeID uint) ([]data.StoreStock, error)
	DeductStockByAdditiveTechCart(storeID, storeAdditiveID uint) ([]data.StoreStock, error)

	FindEarliestExpirationForIngredient(ingredientID, storeID uint) (*time.Time, error)
	GetStockByStoreAndIngredient(storeID, ingredientID uint, stock *data.StoreStock) error
}

type storeStockRepository struct {
	db *gorm.DB
}

func NewStoreStockRepository(db *gorm.DB) StoreStockRepository {
	return &storeStockRepository{db: db}
}

func (r *storeStockRepository) WithTransaction(txFunc func(txRepo storeStockRepository) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	txRepo := r.CloneWithTransaction(tx)

	if err := txFunc(txRepo); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *storeStockRepository) CloneWithTransaction(tx *gorm.DB) storeStockRepository {
	return storeStockRepository{
		db: tx,
	}
}

func (r *storeStockRepository) GetAvailableIngredientsToAdd(storeID uint, filter *ingredientTypes.IngredientFilter) ([]data.Ingredient, error) {
	var ingredients []data.Ingredient
	query := r.db.Model(&data.Ingredient{}).
		Preload("Unit").
		Preload("IngredientCategory")

	if filter.ProductSizeID != nil {
		query = query.Joins("JOIN product_size_ingredients psi ON psi.ingredient_id = ingredients.id").
			Where("psi.product_size_id = ?", *filter.ProductSizeID)
	}

	if filter.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.MinCalories != nil {
		query = query.Where("calories >= ?", *filter.MinCalories)
	}
	if filter.MaxCalories != nil {
		query = query.Where("calories <= ?", *filter.MaxCalories)
	}

	query = query.Where("ingredients.id NOT IN (?)", r.db.Model(&data.StoreStock{}).
		Select("ingredient_id").
		Where("store_id = ?", storeID))

	query, err := utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.Ingredient{})
	if err != nil {
		return nil, err
	}

	if err := query.Find(&ingredients).Error; err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (r *storeStockRepository) AddStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error) {
	var existingStock data.StoreStock
	err := r.db.
		Where("store_id = ? AND ingredient_id = ?", storeID, dto.IngredientID).
		First(&existingStock).Error
	if err == nil {
		return 0, fmt.Errorf("%w: stock with ingredient ID %d already exists for store ID %d",
			moduleErrors.ErrAlreadyExists, dto.IngredientID, storeID)
	}

	var ingredient data.Ingredient
	err = r.db.
		Where("id = ?", dto.IngredientID).
		First(&ingredient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("ingredient not found for ingredient ID %d", dto.IngredientID)
		}
		return 0, fmt.Errorf("failed to fetch ingredient: %w", err)
	}

	storeStock := types.AddToStock(*dto, storeID)
	err = r.db.Create(&storeStock).Error
	if err != nil {
		return 0, fmt.Errorf("failed to create store stock: %w", err)
	}

	return storeStock.ID, nil
}

func (r *storeStockRepository) AddOrUpdateStock(storeID uint, dto *types.AddStoreStockDTO) (uint, error) {
	var existingStock data.StoreStock
	err := r.db.
		Where("store_id = ? AND ingredient_id = ?", storeID, dto.IngredientID).
		First(&existingStock).Error

	if err == nil {
		existingStock.Quantity += dto.Quantity
		existingStock.LowStockThreshold = dto.LowStockThreshold
		err = r.db.Save(&existingStock).Error
		if err != nil {
			return 0, fmt.Errorf("failed to update store stock: %w", err)
		}
		return existingStock.ID, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		storeStock := data.StoreStock{
			StoreID:           storeID,
			IngredientID:      dto.IngredientID,
			Quantity:          dto.Quantity,
			LowStockThreshold: dto.LowStockThreshold,
		}

		err = r.db.Create(&storeStock).Error
		if err != nil {
			return 0, fmt.Errorf("failed to create store stock: %w", err)
		}
		return storeStock.ID, nil
	}

	return 0, fmt.Errorf("failed to add or update stock: %w", err)
}

func (r *storeStockRepository) GetStockList(storeID uint, filter *types.GetStockFilterQuery) ([]data.StoreStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeStockList []data.StoreStock

	query := r.db.Model(&data.StoreStock{}).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Joins("JOIN ingredients ON ingredient_id = ingredients.id").
		Where("store_id = ?", storeID)

	if filter.LowStockOnly != nil && *filter.LowStockOnly {
		query = query.Where("store_stocks.quantity < store_stocks.low_stock_threshold")
	}

	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("ingredients.name ILIKE ?", "%"+*filter.Search+"%")
	}

	var err error
	query, err = utils.ApplySortedPaginationForModel(query, filter.Pagination, filter.Sort, &data.StoreStock{})
	if err != nil {
		return nil, fmt.Errorf("failed to apply pagination: %w", err)
	}

	err = query.Find(&storeStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeStockList, nil
}

func (r *storeStockRepository) GetAllStockList(storeID uint) ([]data.StoreStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeWarehouseStockList []data.StoreStock

	query := r.db.Model(&data.StoreStock{}).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory").
		Where("store_id = ?", storeID)

	err := query.Find(&storeWarehouseStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeWarehouseStockList, nil
}

func (r *storeStockRepository) GetStockListByIDs(storeID uint, stockIds []uint) ([]data.StoreStock, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("storeId cannot be 0")
	}

	var storeWarehouseStockList []data.StoreStock

	query := r.db.Model(&data.StoreStock{}).
		Where("store_id = ? AND id IN (?)", storeID, stockIds).
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory")

	err := query.Find(&storeWarehouseStockList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stock list: %w", err)
	}

	return storeWarehouseStockList, nil
}

func (r *storeStockRepository) GetStockById(storeId, stockId uint) (*data.StoreStock, error) {
	var StoreWarehouseStock data.StoreStock

	if storeId == 0 {
		return nil, fmt.Errorf("store stock id cannot be 0")
	}

	if stockId == 0 {
		return nil, fmt.Errorf("store stock id cannot be 0")
	}

	dbQuery := r.db.Model(&data.StoreStock{}).
		Preload("Store").
		Preload("Store.Warehouse").
		Preload("Ingredient").
		Preload("Ingredient.Unit").
		Preload("Ingredient.IngredientCategory")

	if err := dbQuery.
		Where("store_id = ?", storeId).
		Where("store_stocks.id = ?", stockId).
		First(&StoreWarehouseStock).Error; err != nil {
		return nil, err
	}

	return &StoreWarehouseStock, nil
}

func (r *storeStockRepository) UpdateStock(storeId, stockId uint, dto *types.UpdateStoreStockDTO) error {

	if storeId == 0 {
		return fmt.Errorf("storeId cannot be 0")
	}

	if stockId == 0 {
		return fmt.Errorf("stockId cannot be 0")
	}

	updateFields := map[string]interface{}{}

	if dto.Quantity != nil {
		updateFields["quantity"] = *dto.Quantity
	}
	if dto.LowStockThreshold != nil {
		updateFields["low_stock_threshold"] = *dto.LowStockThreshold
	}

	var existingStock data.StoreStock

	res := r.db.Model(&data.StoreStock{}).
		Where("store_id = ?", storeId).
		Where("id = ?", stockId).
		First(&existingStock)

	if res.Error != nil {
		return utils.WrapError("failed to update store stock", res.Error)
	}

	updRes := r.db.Model(&data.StoreStock{}).
		Where(&data.StoreStock{BaseEntity: data.BaseEntity{ID: stockId}}).
		Updates(updateFields)

	if updRes.Error != nil {
		return utils.WrapError("failed to update store stock", updRes.Error)
	}

	if updRes.RowsAffected == 0 {
		return fmt.Errorf("update attempt had no changes for stockId=%d with storeId=%d", stockId, storeId)
	}

	return nil
}

func (r *storeStockRepository) DeleteStockById(storeId, stockId uint) error {
	res := r.db.
		Where("id = ? AND store_id = ?", stockId, storeId).
		Delete(&data.StoreStock{})

	if res.Error != nil {
		return fmt.Errorf("failed to delete store warehouse stock: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("stock not found")
	}

	return nil
}

func (r *storeStockRepository) DeductStockByProductSizeTechCart(storeID, storeProductSizeID uint) ([]data.StoreStock, error) {
	var store data.Store
	if err := r.getStoreWithWarehouse(storeID, &store); err != nil {
		return nil, err
	}

	var storeProductSize data.StoreProductSize

	err := r.db.Model(&data.StoreProductSize{}).
		Where("id = ?", storeProductSizeID).First(&storeProductSize).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("store product size not found for store ID %d", storeProductSizeID)
		}
		return nil, fmt.Errorf("failed to fetch store product size: %w", err)
	}

	productSizeIngredients, err := r.getProductSizeIngredients(storeProductSize.ProductSizeID)
	if err != nil {
		return nil, err
	}

	var updatedStocks []data.StoreStock

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, ingredient := range productSizeIngredients {
			updatedStock, err := r.deductProductSizeIngredientStock(tx, store.ID, ingredient)
			if err != nil {
				return err
			}
			updatedStocks = append(updatedStocks, *updatedStock)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedStocks, nil
}

func (r *storeStockRepository) DeductStockByAdditiveTechCart(storeID, storeAdditiveID uint) ([]data.StoreStock, error) {
	var store data.Store
	if err := r.getStoreWithWarehouse(storeID, &store); err != nil {
		return nil, err
	}

	var storeAdditive data.StoreAdditive
	err := r.db.Model(&data.StoreAdditive{}).
		Where("id = ?", storeAdditiveID).First(&storeAdditive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("store additive not found for store ID %d", storeAdditiveID)
		}
		return nil, fmt.Errorf("failed to fetch store additive: %w", err)
	}

	additiveIngredients, err := r.getAdditiveIngredients(storeAdditive.AdditiveID)
	if err != nil {
		return nil, err
	}

	var updatedStocks []data.StoreStock

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, ingredient := range additiveIngredients {
			updatedStock, err := r.deductAdditiveIngredientStock(tx, store.ID, ingredient)
			if err != nil {
				return err
			}
			updatedStocks = append(updatedStocks, *updatedStock)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedStocks, nil
}

func (r *storeStockRepository) getStoreWithWarehouse(storeID uint, storeWarehouse *data.Store) error {
	err := r.db.Preload("Warehouse").Where("id = ?", storeID).First(storeWarehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("store warehouse not found for store ID %d", storeID)
		}
		return fmt.Errorf("failed to fetch store warehouse: %w", err)
	}
	return nil
}

func (r *storeStockRepository) getProductSizeIngredients(productSizeID uint) ([]data.ProductSizeIngredient, error) {
	var productSizeIngredients []data.ProductSizeIngredient
	err := r.db.Preload("Ingredient").
		Preload("Ingredient.Unit").
		Where("product_size_id = ?", productSizeID).
		Find(&productSizeIngredients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product size ingredients: %w", err)
	}
	return productSizeIngredients, nil
}

func (r *storeStockRepository) getAdditiveIngredients(additiveID uint) ([]data.AdditiveIngredient, error) {
	var additiveIngredients []data.AdditiveIngredient
	err := r.db.Preload("Ingredient").
		Preload("Ingredient.Unit").
		Where("additive_id = ?", additiveID).
		Find(&additiveIngredients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch additive ingredients: %w", err)
	}
	return additiveIngredients, nil
}

func (r *storeStockRepository) deductProductSizeIngredientStock(tx *gorm.DB, storeID uint, ingredient data.ProductSizeIngredient) (*data.StoreStock, error) {
	var existingStock data.StoreStock
	err := tx.Preload("Ingredient").
		Preload("Ingredient.Unit").
		Where("store_id = ? AND ingredient_id = ?", storeID, ingredient.IngredientID).
		First(&existingStock).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock not found for ingredient ID %d", ingredient.IngredientID)
		}
		return nil, fmt.Errorf("failed to fetch store warehouse stock: %w", err)
	}

	var deductedQuantity float64
	if existingStock.Ingredient.UnitID == ingredient.Ingredient.UnitID {
		deductedQuantity = ingredient.Quantity
	} else {
		deductedQuantity = ingredient.Quantity * existingStock.Ingredient.Unit.ConversionFactor
	}

	if existingStock.Quantity < deductedQuantity {
		return nil, fmt.Errorf("insufficient stock for ingredient ID %d", ingredient.IngredientID)
	}

	newQuantity := existingStock.Quantity - deductedQuantity
	err = tx.Model(&data.StoreStock{}).
		Where("id = ?", existingStock.ID).
		Update("quantity", newQuantity).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update store warehouse stock for ingredient ID %d: %w", ingredient.IngredientID, err)
	}

	existingStock.Quantity = newQuantity
	return &existingStock, nil
}

func (r *storeStockRepository) deductAdditiveIngredientStock(tx *gorm.DB, storeID uint, ingredient data.AdditiveIngredient) (*data.StoreStock, error) {
	var existingStock data.StoreStock
	err := tx.Preload("Ingredient").
		Preload("Ingredient.Unit").
		Where("store_id = ? AND ingredient_id = ?", storeID, ingredient.IngredientID).
		First(&existingStock).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("stock not found for ingredient ID %d", ingredient.IngredientID)
		}
		return nil, fmt.Errorf("failed to fetch store warehouse stock: %w", err)
	}

	var deductedQuantity float64
	if existingStock.Ingredient.UnitID == ingredient.Ingredient.UnitID {
		deductedQuantity = ingredient.Quantity
	} else {
		deductedQuantity = ingredient.Quantity * existingStock.Ingredient.Unit.ConversionFactor
	}

	if existingStock.Quantity < deductedQuantity {
		return nil, fmt.Errorf("insufficient stock for ingredient ID %d", ingredient.IngredientID)
	}

	newQuantity := existingStock.Quantity - deductedQuantity
	err = tx.Model(&data.StoreStock{}).
		Where("id = ?", existingStock.ID).
		Update("quantity", newQuantity).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update store warehouse stock for ingredient ID %d: %w", ingredient.IngredientID, err)
	}

	existingStock.Quantity = newQuantity
	return &existingStock, nil
}

func (r *storeStockRepository) FindEarliestExpirationForIngredient(ingredientID, storeID uint) (*time.Time, error) {
	var earliestExpirationDate sql.NullTime

	err := r.db.Model(&data.StockRequestIngredient{}).
		Joins("JOIN stock_requests ON stock_requests.id = stock_request_ingredients.stock_request_id").
		Joins("JOIN stock_materials ON stock_materials.id = stock_request_ingredients.stock_material_id").
		Select("MIN(stock_request_ingredients.expiration_date) AS earliest_expiration_date").
		Where("stock_requests.store_id = ? AND stock_materials.ingredient_id = ?", storeID, ingredientID).
		Scan(&earliestExpirationDate).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch earliest expiration date for ingredient ID %d: %w", ingredientID, err)
	}

	if !earliestExpirationDate.Valid {
		return nil, nil
	}

	utcTime := earliestExpirationDate.Time.UTC()
	return &utcTime, nil
}

func (r *storeStockRepository) GetStockByStoreAndIngredient(
	storeID, ingredientID uint,
	stock *data.StoreStock,
) error {
	err := r.db.
		Where("store_id = ? AND ingredient_id = ?", storeID, ingredientID).
		First(stock).Error
	if err != nil {
		return fmt.Errorf("failed to fetch store stock: %w", err)
	}
	return nil
}
