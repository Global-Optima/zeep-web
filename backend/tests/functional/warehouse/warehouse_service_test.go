package functional

import (
	"fmt"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/types"
	"github.com/Global-Optima/zeep-web/backend/tests"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var container = tests.NewTestContainer()

func ResetTestData(t *testing.T) *gorm.DB {
	db := container.GetDB()
	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}
	return db
}

func TestWarehouseService_GetWarehouseByID_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetWarehouseModule()

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

func TestWarehouseService_GetWarehouses_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetWarehouseModule()

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
				Search:     tests.StringPtr("Warehouse"),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
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
	_ = ResetTestData(t)
	module := tests.GetWarehouseModule()

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
				Name:     tests.StringPtr("Primary Warehouse"),
				RegionID: tests.UintPtr(1),
				FacilityAddress: &types.FacilityAddressDTO{
					Address: "New Address",
				},
			},
			expectError: false,
		},
		{
			name: "Update non-existing Warehouse",
			id:   0,
			update: types.UpdateWarehouseDTO{
				Name:     tests.StringPtr("Primary Warehouse"),
				RegionID: tests.UintPtr(1),
				FacilityAddress: &types.FacilityAddressDTO{
					Address: "New Address",
				},
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
				assert.Equal(t, tc.update.RegionID, tests.UintPtr(warehouse.Region.ID))
				// assert.Equal(t, tc.update.FacilityAddressID, warehouse.FacilityAddress) // need to think how to check
			}
		})
	}
}

func TestWarehouseService_DeleteWarehouse_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetWarehouseModule()

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
