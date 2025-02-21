package functional

import (
	"fmt"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/tests"
	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}

func TestWarehouseService_GetWarehouseByID_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetWarehouseModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name        string
		id          uint
		expectError bool
	}{
		{
			name:        "Get existing warehouse",
			id:          1,
			expectError: false,
		},
		{
			name:        "Get non-existing warehouse",
			id:          0,
			expectError: true,
		},
		{
			name:        "Get non-existing warehouse",
			id:          999,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			warehouse, err := module.Service.GetWarehouseByID(tc.id)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, warehouse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, warehouse)
				assert.Equal(t, tc.id, warehouse.ID)
			}
		})
	}
}

// i will leave until getWarehouse and getAllWarehouse functions will be decided
func TestWarehouseService_GetWarehouses_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetWarehouseModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name          string
		filter        *types.WarehouseFilter
		expectedCount int
	}{
		{
			name:          "Get all warehouses",
			filter:        nil,
			expectedCount: 1,
		},
		{
			name: "Filter by search term",
			filter: &types.WarehouseFilter{
				Search: stringPtr("Warehouse"),
			},
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			warehouses, err := module.Service.GetWarehouses(tc.filter)
			if tc.filter == nil {
				assert.ErrorContains(t, err, fmt.Errorf("filter is nil").Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, warehouses, tc.expectedCount)
			}
		})
	}
}

func TestWarehouseService_UpdateWarehouse_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetWarehouseModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name        string
		id          uint
		update      types.UpdateWarehouseDTO
		expectError bool
	}{
		{
			name: "Update existing Warehouse",
			id:   1,
			update: types.UpdateWarehouseDTO{
				Name:              stringPtr("Primary Warehouse"),
				RegionID:          uintPtr(1),
				FacilityAddressID: uintPtr(1),
			},
			expectError: false,
		},
		{
			name: "Update non-existing Warehouse",
			id:   0,
			update: types.UpdateWarehouseDTO{
				Name:              stringPtr("Primary Warehouse"),
				RegionID:          uintPtr(1),
				FacilityAddressID: uintPtr(1),
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			warehouse, err := module.Service.UpdateWarehouse(tc.id, tc.update)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, warehouse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, warehouse)
				assert.Equal(t, *tc.update.Name, warehouse.Name)
				assert.Equal(t, tc.update.RegionID, uintPtr(warehouse.Region.ID))
				// assert.Equal(t, tc.update.FacilityAddressID, warehouse.FacilityAddress) // need to think how to check
			}
		})
	}
}

func TestWarehouseService_DeleteWarehouse_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetWarehouseModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name        string
		id          uint
		expectedErr error
	}{
		{
			name:        "Delete existing ID 1",
			id:          1,
			expectedErr: nil,
		},
		{
			name:        "Delete non-existing ID 0",
			id:          0,
			expectedErr: fmt.Errorf("warehouse with ID %d not found", 0),
		},
		{
			name:        "Delete non-existing ID 999",
			id:          999,
			expectedErr: fmt.Errorf("warehouse with ID %d not found", 999),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultErr := module.Service.DeleteWarehouse(tc.id)
			if tc.expectedErr == nil {
				assert.NoError(t, resultErr, "expected no error, but got one")

				_, getErr := module.Service.GetWarehouseByID(tc.id)
				assert.Error(t, getErr, "warehouse should not exist after deletion")
			} else {
				assert.EqualError(t, resultErr, tc.expectedErr.Error(), "expected error message does not match")
			}
		})
	}
}

// func TestWarehouseService_DeleteWarehouse(t *testing.T) {
// 	container := tests.NewTestContainer()
// 	db := container.GetDB()
// 	module := tests.GetWarehouseModule()

// 	if err := tests.TruncateAllTables(db); err != nil {
// 		t.Fatalf("Failed to truncate all tables: %v", err)
// 	}
// 	if err := tests.LoadTestData(db); err != nil {
// 		t.Fatalf("Failed to load test data: %v", err)
// 	}

// 	warehouseWithStores := types.CreateWarehouseDTO{
// 		Name:     "Test Warehouse With Stores",
// 		RegionID: 1,
// 		FacilityAddress: types.FacilityAddressDTO{
// 			Address: "123 Test St",
// 		},
// 	}
// 	createdWarehouse, err := module.Service.CreateWarehouse(warehouseWithStores)
// 	assert.NoError(t, err)

// 	err = module.Service.AssignStoreToWarehouse(types.AssignStoreToWarehouseRequest{
// 		WarehouseID: createdWarehouse.ID,
// 		StoreID:     1,
// 	})
// 	assert.NoError(t, err)

// 	testCases := []struct {
// 		name        string
// 		id          uint
// 		setupFunc   func()
// 		expectedErr error
// 		verifyFunc  func(t *testing.T, id uint)
// 	}{
// 		{
// 			name:        "Delete existing warehouse",
// 			id:          1,
// 			expectedErr: nil,
// 			verifyFunc: func(t *testing.T, id uint) {
// 				_, err := module.Service.GetWarehouseByID(id)
// 				assert.Error(t, err, "Warehouse should not exist after deletion")
// 			},
// 		},
// 		{
// 			name:        "Delete warehouse with assigned stores",
// 			id:          createdWarehouse.ID,
// 			expectedErr: nil,
// 			verifyFunc: func(t *testing.T, id uint) {
// 				_, err := module.Service.GetWarehouseByID(id)
// 				assert.Error(t, err, "Warehouse should not exist after deletion")

// 				stores, err := module.Service.GetAllStoresByWarehouse(id, &utils.Pagination{
// 					Page:     1,
// 					PageSize: 10,
// 				})
// 				assert.Error(t, err)
// 				assert.Empty(t, stores, "All store assignments should be removed")
// 			},
// 		},
// 		{
// 			name:        "Delete non-existing warehouse",
// 			id:          0,
// 			expectedErr: fmt.Errorf("warehouse with ID %d not found", 0),
// 		},
// 		{
// 			name: "Delete already deleted warehouse",
// 			id:   1,
// 			setupFunc: func() {
// 				_ = module.Service.DeleteWarehouse(1)
// 			},
// 			expectedErr: fmt.Errorf("warehouse with ID %d not found", 1),
// 		},
// 		{
// 			name:        "Delete warehouse with invalid ID",
// 			id:          999999,
// 			expectedErr: fmt.Errorf("warehouse with ID %d not found", 999999),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			if tc.setupFunc != nil {
// 				tc.setupFunc()
// 			}

// 			resultErr := module.Service.DeleteWarehouse(tc.id)

// 			if tc.expectedErr == nil {
// 				assert.NoError(t, resultErr, "expected no error, but got one")

// 				if tc.verifyFunc != nil {
// 					tc.verifyFunc(t, tc.id)
// 				}
// 			} else {
// 				assert.EqualError(t, resultErr, tc.expectedErr.Error(), "expected error message does not match")
// 			}
// 		})
// 	}
// }
