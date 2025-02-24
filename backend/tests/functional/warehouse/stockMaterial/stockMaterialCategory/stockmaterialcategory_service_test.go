package functional

import (
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory/types"
	"github.com/Global-Optima/zeep-web/backend/tests"
	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

// func uintPtr(u uint) *uint {
// 	return &u
// }

// func floatPtr(f float64) *float64 {
// 	return &f
// }

func TestStockMaterialCategoryService_GetByID_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetStockMaterialCategoryModule()

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
			name:        "Get existing category",
			id:          1,
			expectError: false,
		},
		{
			name:        "Get non-existing category",
			id:          0,
			expectError: true,
		},
		{
			name:        "Get non-existing category with high ID",
			id:          999,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			category, err := module.Service.GetByID(tc.id)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, category)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, category)
				assert.Equal(t, tc.id, category.ID)
				assert.Equal(t, "Raw Materials", category.Name)
				assert.Equal(t, "Materials used in production", category.Description)
			}
		})
	}
}

func TestStockMaterialCategoryService_GetAll_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetStockMaterialCategoryModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name          string
		filter        types.StockMaterialCategoryFilter
		expectedCount int
	}{
		{
			name:          "Get all categories without filter",
			filter:        types.StockMaterialCategoryFilter{},
			expectedCount: 1,
		},
		{
			name: "Filter by search term - matching",
			filter: types.StockMaterialCategoryFilter{
				Search: stringPtr("Raw"),
			},
			expectedCount: 1,
		},
		{
			name: "Filter by search term - non-matching",
			filter: types.StockMaterialCategoryFilter{
				Search: stringPtr("NonExisting"),
			},
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			categories, err := module.Service.GetAll(tc.filter)
			assert.NoError(t, err)
			assert.Len(t, categories, tc.expectedCount)
			if tc.expectedCount > 0 {
				assert.Equal(t, "Raw Materials", categories[0].Name)
			}
		})
	}
}

func TestStockMaterialCategoryService_Create_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetStockMaterialCategoryModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name        string
		input       types.CreateStockMaterialCategoryDTO
		expectError bool
	}{
		// ('Raw Materials', 'Materials used in production');
		{
			name: "Create valid category",
			input: types.CreateStockMaterialCategoryDTO{
				Name:        "Test Category",
				Description: "Test Description",
			},
			expectError: false,
		},
		{
			name: "Create with empty name",
			input: types.CreateStockMaterialCategoryDTO{
				Name:        "",
				Description: "Test Description",
			},
			expectError: true,
		},
		{
			name: "Create with duplicate name",
			input: types.CreateStockMaterialCategoryDTO{
				Name:        "Raw Materials", // Already exists in preloaded data
				Description: "Test Description",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := module.Service.Create(tc.input)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), id)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uint(0), id)

				// Verify the created category
				category, err := module.Service.GetByID(id)
				assert.NoError(t, err)
				assert.Equal(t, tc.input.Name, category.Name)
				assert.Equal(t, tc.input.Description, category.Description)
				assert.NotEmpty(t, category.CreatedAt)
				assert.NotEmpty(t, category.UpdatedAt)
			}
		})
	}
}

func TestStockMaterialCategoryService_Update_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetStockMaterialCategoryModule()

	if err := tests.TruncateAllTables(db); err != nil {
		t.Fatalf("Failed to truncate all tables: %v", err)
	}
	if err := tests.LoadTestData(db); err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	testCases := []struct {
		name        string
		id          uint
		input       types.UpdateStockMaterialCategoryDTO
		expectError bool
	}{
		{
			name: "Update existing category",
			id:   1,
			input: types.UpdateStockMaterialCategoryDTO{
				Name:        stringPtr("Updated Category"),
				Description: stringPtr("Updated Description"),
			},
			expectError: false,
		},
		{
			name: "Update non-existing category",
			id:   999,
			input: types.UpdateStockMaterialCategoryDTO{
				Name: stringPtr("Updated Name"),
			},
			expectError: true,
		},
		{
			name: "Update with partial fields",
			id:   1,
			input: types.UpdateStockMaterialCategoryDTO{
				Description: stringPtr("Only Description Updated"),
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := module.Service.Update(tc.id, tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the update
				category, err := module.Service.GetByID(tc.id)
				assert.NoError(t, err)
				if tc.input.Name != nil {
					assert.Equal(t, *tc.input.Name, category.Name)
				}
				if tc.input.Description != nil {
					assert.Equal(t, *tc.input.Description, category.Description)
				}
			}
		})
	}
}

func TestStockMaterialCategoryService_Delete_WithPreloadedData(t *testing.T) {
	container := tests.NewTestContainer()
	db := container.GetDB()
	module := tests.GetStockMaterialCategoryModule()

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
			name:        "Delete existing category",
			id:          1,
			expectError: false,
		},
		{
			name:        "Delete non-existing category",
			id:          999,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := module.Service.Delete(tc.id)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify deletion
				_, err := module.Service.GetByID(tc.id)
				assert.Error(t, err)
			}
		})
	}
}
