package functional

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware/contexts"
	"strings"
	"testing"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks/types"
	"github.com/Global-Optima/zeep-web/backend/tests"
	"github.com/stretchr/testify/assert"
)

var container = tests.NewTestContainer()

func stringPtr(s string) *string {
	return &s
}

func setupStoreStockTest(t *testing.T) storeStocks.StoreStockService {
	db := container.GetDB()
	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("TruncateAllTables error: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("LoadTestData error: %v", err)
	}

	return tests.GetStoreStocksModule().Service
}

func createTestStock(t *testing.T, service storeStocks.StoreStockService, storeID uint, dto types.AddStoreStockDTO) uint {
	id, err := service.AddStock(storeID, &dto)
	assert.NoError(t, err, "AddStock should succeed")
	assert.NotZero(t, id, "Returned stock ID should not be zero")
	return id
}

func deleteTestStock(t *testing.T, service storeStocks.StoreStockService, storeID, stockID uint) {
	err := service.DeleteStockById(storeID, stockID)
	assert.NoError(t, err, "DeleteStockById should succeed")
}

func clearMockStock(t *testing.T, service storeStocks.StoreStockService) {
	deleteTestStock(t, service, 1, 1)
	deleteTestStock(t, service, 1, 2)
}

func TestStoreStockService_AddStock(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	deleteTestStock(t, service, storeID, 1)

	t.Run("Success - New Stock", func(t *testing.T) {
		dto := types.AddStoreStockDTO{
			IngredientID:      1,
			Quantity:          100,
			LowStockThreshold: 50,
		}
		id, err := service.AddStock(storeID, &dto)
		assert.NoError(t, err)
		stockDto, err := service.GetStockById(id, &contexts.StoreContextFilter{StoreID: &storeID})
		assert.NoError(t, err)
		assert.Equal(t, 100.0, stockDto.Quantity, "Quantity should be 100")
		assert.Equal(t, false, stockDto.LowStockAlert, "LowStockAlert should be false when quantity is above threshold")
	})

	t.Run("Failure - Duplicate Stock", func(t *testing.T) {
		dto := types.AddStoreStockDTO{
			IngredientID:      1,
			Quantity:          100,
			LowStockThreshold: 50,
		}

		_, err := service.AddStock(1, &dto)
		assert.Error(t, err, "Should not allow duplicate stock for same ingredient")
		assert.True(t, strings.Contains(err.Error(), "already exists"),
			"Error should indicate stock already exists")
	})

	t.Run("Failure - Ingredient Not Found", func(t *testing.T) {
		dto := types.AddStoreStockDTO{
			IngredientID:      9999,
			Quantity:          100,
			LowStockThreshold: 50,
		}
		_, err := service.AddStock(storeID, &dto)
		assert.Error(t, err, "Should error if ingredient does not exist")
	})
}

func TestStoreStockService_AddMultipleStock(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	deleteTestStock(t, service, storeID, 1)

	t.Run("Add Multiple - New and Update", func(t *testing.T) {

		dto1 := types.AddStoreStockDTO{
			IngredientID:      1,
			Quantity:          100,
			LowStockThreshold: 50,
		}
		id1, err := service.AddStock(storeID, &dto1)
		assert.NoError(t, err)

		multiDto := types.AddMultipleStoreStockDTO{
			IngredientStocks: []types.AddStoreStockDTO{
				{IngredientID: 1, Quantity: 50, LowStockThreshold: 50},
				{IngredientID: 2, Quantity: 200, LowStockThreshold: 100},
			},
		}
		ids, err := service.AddMultipleStock(storeID, &multiDto)
		assert.NoError(t, err)
		assert.Len(t, ids, 2, "Should return two stock IDs")

		stock1, err := service.GetStockById(id1, nil)
		assert.NoError(t, err)
		assert.Equal(t, 150.0, stock1.Quantity, "Quantity should be updated to 150")
	})
}

func TestStoreStockService_GetStockListAndByIDs(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	clearMockStock(t, service)

	t.Run("GetStockList with Filtering", func(t *testing.T) {

		dto1 := types.AddStoreStockDTO{IngredientID: 1, Quantity: 40, LowStockThreshold: 50}
		dto2 := types.AddStoreStockDTO{IngredientID: 2, Quantity: 20, LowStockThreshold: 30}
		_ = createTestStock(t, service, storeID, dto1)
		_ = createTestStock(t, service, storeID, dto2)

		low := true
		filter := types.GetStockFilterQuery{
			BaseFilter:   tests.BaseFilterWithPagination(1, 10),
			LowStockOnly: &low,
		}
		list, err := service.GetStockList(storeID, &filter)
		assert.NoError(t, err)

		for _, stock := range list {
			assert.True(t, stock.LowStockAlert, "Each returned stock should be flagged as low stock")
			deleteTestStock(t, service, storeID, stock.ID)
		}
	})

	t.Run("GetStockListByIDs", func(t *testing.T) {

		id1 := createTestStock(t, service, 1, types.AddStoreStockDTO{IngredientID: 1, Quantity: 100, LowStockThreshold: 50})
		id2 := createTestStock(t, service, 1, types.AddStoreStockDTO{IngredientID: 2, Quantity: 200, LowStockThreshold: 100})
		ids := []uint{id1, id2}
		list, err := service.GetStockListByIDs(1, ids)
		assert.NoError(t, err)
		assert.Len(t, list, 2, "Should retrieve two stocks by IDs")
	})
}

func TestStoreStockService_GetStockById(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	clearMockStock(t, service)

	t.Run("Existing Stock", func(t *testing.T) {
		id := createTestStock(t, service, storeID, types.AddStoreStockDTO{IngredientID: 1, Quantity: 100, LowStockThreshold: 50})
		stock, err := service.GetStockById(id, &contexts.StoreContextFilter{StoreID: &storeID})
		assert.NoError(t, err)
		assert.Equal(t, id, stock.ID, "Retrieved stock ID should match")
	})

	t.Run("Non-existent Stock", func(t *testing.T) {
		_, err := service.GetStockById(9999, &contexts.StoreContextFilter{StoreID: &storeID})
		assert.Error(t, err, "Retrieving a non-existent stock should error")
	})
}

func TestStoreStockService_UpdateStockById(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	clearMockStock(t, service)

	t.Run("Successful Update", func(t *testing.T) {
		id := createTestStock(t, service, storeID, types.AddStoreStockDTO{IngredientID: 1, Quantity: 100, LowStockThreshold: 50})
		updateDTO := types.UpdateStoreStockDTO{
			Quantity:          floatPtr(80),
			LowStockThreshold: floatPtr(40),
		}
		err := service.UpdateStockById(storeID, id, &updateDTO)
		assert.NoError(t, err)
		stock, err := service.GetStockById(id, &contexts.StoreContextFilter{StoreID: &storeID})
		assert.NoError(t, err)
		assert.Equal(t, 80.0, stock.Quantity, "Quantity should be updated to 80")
		assert.Equal(t, 40.0, stock.LowStockThreshold, "Threshold should be updated to 40")
	})

	t.Run("Invalid Update - Non-existent Stock", func(t *testing.T) {
		updateDTO := types.UpdateStoreStockDTO{
			Quantity:          floatPtr(80),
			LowStockThreshold: floatPtr(40),
		}
		err := service.UpdateStockById(1, 9999, &updateDTO)
		assert.Error(t, err, "Updating a non-existent stock should error")
	})
}

func TestStoreStockService_DeleteStockById(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	clearMockStock(t, service)

	t.Run("Successful Deletion", func(t *testing.T) {
		id := createTestStock(t, service, storeID, types.AddStoreStockDTO{IngredientID: 1, Quantity: 100, LowStockThreshold: 50})
		err := service.DeleteStockById(storeID, id)
		assert.NoError(t, err)
		_, err = service.GetStockById(id, &contexts.StoreContextFilter{StoreID: &storeID})
		assert.Error(t, err, "Deleted stock should no longer be retrievable")
	})

	t.Run("Deletion of Non-existent Stock", func(t *testing.T) {
		err := service.DeleteStockById(storeID, 9999)
		assert.Error(t, err, "Deleting a non-existent stock should error")
	})
}

func TestStoreStockService_GetAvailableIngredientsToAdd(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	clearMockStock(t, service)

	t.Run("Exclude Existing Ingredients", func(t *testing.T) {

		_ = createTestStock(t, service, storeID, types.AddStoreStockDTO{IngredientID: 1, Quantity: 100, LowStockThreshold: 50})

		filter := &ingredientTypes.IngredientFilter{
			BaseFilter: tests.BaseFilterWithPagination(1, 10),
			Name:       stringPtr(""),
		}
		ingList, err := service.GetAvailableIngredientsToAdd(storeID, filter)
		assert.NoError(t, err)
		for _, ing := range ingList {
			assert.NotEqual(t, uint(1), ing.ID, "Ingredient already in stock should be excluded")
		}
	})
}

func TestStoreStockService_CheckStockNotifications(t *testing.T) {
	service := setupStoreStockTest(t)
	var storeID uint = 1
	clearMockStock(t, service)

	t.Run("Low Stock Notification", func(t *testing.T) {
		id := createTestStock(t, service, storeID, types.AddStoreStockDTO{IngredientID: 1, Quantity: 30, LowStockThreshold: 50})
		stock, err := service.GetStockById(id, &contexts.StoreContextFilter{StoreID: &storeID})
		assert.NoError(t, err)

		err = service.CheckStockNotifications(storeID, data.StoreStock{
			BaseEntity: data.BaseEntity{
				ID:        stock.ID,
				UpdatedAt: time.Now(),
			},
			Quantity:          stock.Quantity,
			LowStockThreshold: stock.LowStockThreshold,
			Ingredient:        data.Ingredient{BaseEntity: data.BaseEntity{ID: stock.Ingredient.ID}, Name: stock.Ingredient.Name, ExpirationInDays: 10},
			Store:             data.Store{BaseEntity: data.BaseEntity{ID: 1}, Name: "Test Store"},
		})

		assert.NoError(t, err, "CheckStockNotifications should not return an error even if notifications fail")
	})
}

func floatPtr(f float64) *float64 {
	return &f
}
