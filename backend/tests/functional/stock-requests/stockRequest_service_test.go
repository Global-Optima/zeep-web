package functional

import (
	"strings"
	"testing"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests/types"
	"github.com/Global-Optima/zeep-web/backend/tests"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var container = tests.NewTestContainer()

func uintPtr(u uint) *uint {
	return &u
}

func stringPtr(s string) *string {
	return &s
}

func createTestStockRequest(t *testing.T, service stockRequests.StockRequestService, storeID uint) uint {
	dto := types.CreateStockRequestDTO{
		StockMaterials: []types.StockRequestStockMaterialDTO{
			{StockMaterialID: 1, Quantity: 10},
		},
	}
	id, _, err := service.CreateStockRequest(storeID, dto)
	assert.NoError(t, err, "CreateStockRequest should succeed")
	return id
}

func setupTest(t *testing.T) stockRequests.StockRequestService {
	db := container.GetDB()
	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("TruncateAllTables error: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("LoadTestData error: %v", err)
	}

	return tests.GetStockRequestsModule().Service
}

func getTestDB() *gorm.DB {
	return tests.NewTestContainer().GetDB()
}

func TestStockRequestService_CreateStockRequest(t *testing.T) {
	service := setupTest(t)

	t.Run("Normal Creation", func(t *testing.T) {
		dto := types.CreateStockRequestDTO{
			StockMaterials: []types.StockRequestStockMaterialDTO{
				{StockMaterialID: 1, Quantity: 10},
			},
		}
		id, storeName, err := service.CreateStockRequest(1, dto)
		assert.NoError(t, err, "Creation should succeed when no open cart exists")
		assert.NotZero(t, id, "New stock request ID should not be zero")
		assert.NotEmpty(t, storeName, "Store name should not be empty")
		req, err := service.GetStockRequestByID(id)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestCreated, req.Status, "New request should have status CREATED")
	})

	t.Run("Duplicate Open Cart", func(t *testing.T) {
		dto := types.CreateStockRequestDTO{
			StockMaterials: []types.StockRequestStockMaterialDTO{
				{StockMaterialID: 1, Quantity: 5},
			},
		}
		_, _, err := service.CreateStockRequest(1, dto)
		assert.Error(t, err, "A duplicate open cart should not be allowed")
		assert.True(t, strings.Contains(err.Error(), "existing"),
			"Error message should mention that an open cart already exists")
	})
}

func TestStockRequestService_GetStockRequests(t *testing.T) {
	service := setupTest(t)
	_ = createTestStockRequest(t, service, 1)

	t.Run("Filter by Store", func(t *testing.T) {
		filter := types.GetStockRequestsFilter{
			BaseFilter: tests.BaseFilterWithPagination(1, 10),
			StoreID:    uintPtr(1),
		}
		requests, err := service.GetStockRequests(filter)
		assert.NoError(t, err)

		assert.Len(t, requests, 1, "Expected one stock request for store 1")
	})

	t.Run("Filter by Warehouse", func(t *testing.T) {
		filter := types.GetStockRequestsFilter{
			BaseFilter:  tests.BaseFilterWithPagination(1, 10),
			StoreID:     uintPtr(1),
			WarehouseID: uintPtr(1),
		}
		_, err := service.GetStockRequests(filter)
		assert.NoError(t, err, "Filtering by warehouse should succeed")
	})
}

func TestStockRequestService_GetStockRequestByID(t *testing.T) {
	service := setupTest(t)

	t.Run("Fetch Existing Request", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)
		req, err := service.GetStockRequestByID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, req.RequestID, "Fetched request ID should match created ID")
	})
}

func TestStockRequestService_RejectStockRequestByStore(t *testing.T) {
	service := setupTest(t)
	db := getTestDB()

	t.Run("Valid Transition", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)

		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestInDelivery).Error
		assert.NoError(t, err)
		dto := types.RejectStockRequestStatusDTO{
			Comment: stringPtr("Store rejects the request"),
		}
		req, err := service.RejectStockRequestByStore(id, dto)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestRejectedByStore, req.Status, "Status should change to REJECTED_BY_STORE")
	})

	t.Run("Invalid Transition", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)
		dto := types.RejectStockRequestStatusDTO{
			Comment: stringPtr("Store rejects the request"),
		}
		_, err := service.RejectStockRequestByStore(id, dto)
		assert.Error(t, err, "Store rejection from CREATED status should error")
		assert.True(t, strings.Contains(err.Error(), "invalid status transition"),
			"Error should mention invalid status transition")
	})
}

func TestStockRequestService_RejectStockRequestByWarehouse(t *testing.T) {
	service := setupTest(t)
	db := getTestDB()

	t.Run("Valid Transition", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)
		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestProcessed).Error
		assert.NoError(t, err)
		dto := types.RejectStockRequestStatusDTO{
			Comment: stringPtr("Warehouse rejects the request"),
		}
		req, err := service.RejectStockRequestByWarehouse(id, dto)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestRejectedByWarehouse, req.Status, "Status should change to REJECTED_BY_WAREHOUSE")
	})

	t.Run("Invalid Transition", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)
		dto := types.RejectStockRequestStatusDTO{
			Comment: stringPtr("Warehouse rejects the request"),
		}
		_, err := service.RejectStockRequestByWarehouse(id, dto)
		assert.Error(t, err, "Warehouse rejection from CREATED status should error")
	})
}

func TestStockRequestService_SetProcessedStatus(t *testing.T) {
	service := setupTest(t)
	db := getTestDB()

	t.Run("Successful Processing", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)

		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("created_at", time.Now().Add(-25*time.Hour)).Error
		assert.NoError(t, err)
		req, err := service.SetProcessedStatus(id)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestProcessed, req.Status, "Status should change to PROCESSED")
	})

	// t.Run("Rate Limit Violation", func(t *testing.T) {
	// 	// Create first request (valid)
	// 	id1 := createTestStockRequest(t, service, 1)

	// 	// Ensure it's finalized so it's counted in the 24-hour limit
	// 	err := db.Model(&data.StockRequest{}).Where("id = ?", id1).
	// 		Update("created_at", time.Now().Add(-1*time.Hour)).
	// 		Update("status", data.StockRequestCompleted).Error
	// 	assert.NoError(t, err)

	// 	// Create second request (valid)
	// 	id2 := createTestStockRequest(t, service, 1)
	// 	err = db.Model(&data.StockRequest{}).Where("id = ?", id2).
	// 		Update("created_at", time.Now().Add(-1*time.Hour)).
	// 		Update("status", data.StockRequestCompleted).Error
	// 	assert.NoError(t, err)

	// 	_, _, err = service.CreateStockRequest(1, types.CreateStockRequestDTO{
	// 		StockMaterials: []types.StockRequestStockMaterialDTO{
	// 			{StockMaterialID: 1, Quantity: 10},
	// 		},
	// 	})
	// 	assert.Error(t, err, "Processing should fail due to rate limit")
	// 	assert.True(t, strings.Contains(err.Error(), "one request allowed per day"),
	// 		"Error should indicate rate limit violation")
	// })
}

func TestStockRequestService_SetInDeliveryStatus(t *testing.T) {
	service := setupTest(t)
	db := getTestDB()

	id := createTestStockRequest(t, service, 1)

	t.Run("Successful Transition", func(t *testing.T) {
		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestProcessed).Error
		assert.NoError(t, err)
		req, err := service.SetInDeliveryStatus(id)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestInDelivery, req.Status, "Status should change to IN_DELIVERY")
	})

	t.Run("Insufficient Warehouse Stock", func(t *testing.T) {
		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestProcessed).Error
		assert.NoError(t, err)

		err = db.Exec("UPDATE warehouse_stocks SET quantity = ? WHERE warehouse_id = ? AND stock_material_id = ?", 0, 1, 1).Error
		assert.NoError(t, err)
		_, err = service.SetInDeliveryStatus(id)
		assert.Error(t, err, "Transition to IN_DELIVERY should fail due to insufficient stock")
	})
}

func TestStockRequestService_SetCompletedStatus(t *testing.T) {
	service := setupTest(t)
	db := getTestDB()

	t.Run("Successful Completion", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)

		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestInDelivery).Error
		assert.NoError(t, err)
		req, err := service.SetCompletedStatus(id)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestCompleted, req.Status, "Status should change to COMPLETED")
	})
}

func TestStockRequestService_AcceptStockRequestWithChange(t *testing.T) {
	service := setupTest(t)
	db := getTestDB()
	id := createTestStockRequest(t, service, 1)

	t.Run("Existing Ingredient - Quantity Change", func(t *testing.T) {
		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestInDelivery).Error
		assert.NoError(t, err)

		dto := types.AcceptWithChangeRequestStatusDTO{
			Comment: stringPtr("Adjusting quantity"),
			Items: []types.StockRequestStockMaterialDTO{
				{StockMaterialID: 1, Quantity: 5},
			},
		}
		req, err := service.AcceptStockRequestWithChange(id, dto)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestAcceptedWithChange, req.Status, "Status should change to ACCEPTED_WITH_CHANGE")
	})

	t.Run("New Ingredient Addition", func(t *testing.T) {
		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestInDelivery).Error
		assert.NoError(t, err)
		dto := types.AcceptWithChangeRequestStatusDTO{
			Comment: stringPtr("Adding new ingredient"),
			Items: []types.StockRequestStockMaterialDTO{
				{StockMaterialID: 2, Quantity: 5},
			},
		}
		req, err := service.AcceptStockRequestWithChange(id, dto)
		assert.NoError(t, err)
		assert.Equal(t, data.StockRequestAcceptedWithChange, req.Status, "Status should change to ACCEPTED_WITH_CHANGE")
	})
}

func TestStockRequestService_UpdateStockRequest(t *testing.T) {
	service := setupTest(t)
	id := createTestStockRequest(t, service, 1)

	t.Run("Valid Update", func(t *testing.T) {
		items := []types.StockRequestStockMaterialDTO{
			{StockMaterialID: 1, Quantity: 20},
		}
		_, err := service.UpdateStockRequest(id, items)
		assert.NoError(t, err)
		req, err := service.GetStockRequestByID(id)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(req.StockMaterials), 1, "There should be at least one ingredient")
		assert.Equal(t, 20.0, req.StockMaterials[0].Quantity, "Quantity should be updated to 20")
	})

	t.Run("Invalid Stock Material", func(t *testing.T) {
		items := []types.StockRequestStockMaterialDTO{
			{StockMaterialID: 9999, Quantity: 20},
		}
		_, err := service.UpdateStockRequest(id, items)
		assert.Error(t, err, "Updating with an invalid stock material should error")
	})
}

func TestStockRequestService_DeleteStockRequest(t *testing.T) {
	service := setupTest(t)
	db := getTestDB()

	t.Run("Valid Deletion", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)
		req, err := service.DeleteStockRequest(id)
		assert.NoError(t, err)
		_, err = service.GetStockRequestByID(req.ID)
		assert.Error(t, err, "Deleted stock request should not be retrievable")
	})

	t.Run("Invalid Deletion - Wrong Status", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)

		err := db.Model(&data.StockRequest{}).Where("id = ?", id).
			Update("status", data.StockRequestProcessed).Error
		assert.NoError(t, err)
		_, err = service.DeleteStockRequest(id)
		assert.Error(t, err, "Deletion should fail when status is not CREATED")
	})

	t.Run("Non-Existent Request", func(t *testing.T) {
		_, err := service.DeleteStockRequest(9999)
		assert.Error(t, err, "Deleting a non-existent stock request should error")
	})
}

func TestStockRequestService_GetLastCreatedStockRequest(t *testing.T) {
	service := setupTest(t)

	t.Run("Existing Open Cart", func(t *testing.T) {
		id := createTestStockRequest(t, service, 1)
		lastCart, err := service.GetLastCreatedStockRequest(1)
		assert.NoError(t, err)
		assert.NotNil(t, lastCart, "Should return an open cart for store 1")
		assert.Equal(t, id, lastCart.RequestID, "Last created stock request ID should match")
	})

	t.Run("No Open Cart", func(t *testing.T) {
		_, err := service.GetLastCreatedStockRequest(9999)
		assert.Error(t, err, "Expected error when no open cart exists for store 9999")
	})
}

func TestStockRequestService_AddStockMaterialToCart(t *testing.T) {
	service := setupTest(t)

	t.Run("Valid Addition and Update", func(t *testing.T) {
		cart, err := service.AddStockMaterialToCart(1, types.StockRequestStockMaterialDTO{StockMaterialID: 1, Quantity: 5})
		assert.NoError(t, err)

		_, err = service.AddStockMaterialToCart(1, types.StockRequestStockMaterialDTO{StockMaterialID: 1, Quantity: 3})
		assert.NoError(t, err)
		updatedCart, err := service.GetStockRequestByID(cart.ID)
		assert.NoError(t, err)
		var found bool
		for _, item := range updatedCart.StockMaterials {
			if item.StockMaterial.ID == 1 {
				found = true
				assert.Equal(t, 8.0, item.Quantity, "Quantity should be updated to 8")
			}
		}
		assert.True(t, found, "Stock material with ID 1 should be present in the cart")
	})

	t.Run("Invalid Stock Material", func(t *testing.T) {
		_, err := service.AddStockMaterialToCart(1, types.StockRequestStockMaterialDTO{StockMaterialID: 9999, Quantity: 5})
		assert.Error(t, err, "Adding a non-existent stock material should error")
	})
}
