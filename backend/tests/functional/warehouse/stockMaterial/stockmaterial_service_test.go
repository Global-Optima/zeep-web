package functional

import (
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/types"
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

func TestStockMaterialService_GetStockMaterialByID_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetStockMaterialModule()

	testCases := []struct {
		name        string
		id          uint
		expectError bool
	}{
		{
			name:        "Get existing stock material",
			id:          1,
			expectError: false,
		},
		{
			name:        "Get non-existing stock material",
			id:          0,
			expectError: true,
		},
		{
			name:        "Get non-existing stock material with high ID",
			id:          999,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			material, err := module.Service.GetStockMaterialByID(tc.id)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, material)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, material)
				assert.Equal(t, tc.id, material.ID)
				assert.Equal(t, "Coffee Beans Jacobs", material.Name)
			}
		})
	}
}

// failed test
func TestStockMaterialService_GetAllStockMaterials_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetStockMaterialModule()

	testCases := []struct {
		name          string
		filter        *types.StockMaterialFilter
		expectedCount int
		expectError   bool
	}{
		{
			name: "Get all stock materials without filter",
			filter: &types.StockMaterialFilter{
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "Filter by search term",
			filter: &types.StockMaterialFilter{
				Search:     tests.StringPtr("Jacobs"),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "Filter by non-existing term",
			filter: &types.StockMaterialFilter{
				Search:     tests.StringPtr("NonExisting"),
				BaseFilter: tests.BaseFilterWithPagination(1, 10),
			},
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			materials, err := module.Service.GetAllStockMaterials(tc.filter)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, materials, tc.expectedCount)
				if tc.expectedCount > 0 {
					assert.Equal(t, "Coffee Beans Jacobs", materials[0].Name)
				}
			}
		})
	}
}

// test failed
func TestStockMaterialService_CreateStockMaterial_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetStockMaterialModule()

	testCases := []struct {
		name        string
		input       *types.CreateStockMaterialDTO
		expectError bool
	}{
		{
			name: "Create valid stock material",
			input: &types.CreateStockMaterialDTO{
				Name:                   "New Test Material",
				Description:            "Test material description",
				SafetyStock:            50.5,
				UnitID:                 1, // From preloaded data
				CategoryID:             1, // From preloaded data
				IngredientID:           1, // From preloaded data
				Barcode:                "TEST123",
				ExpirationPeriodInDays: 30,
				Size:                   100.5,
			},
			expectError: false,
		},
		{
			name: "Create with missing required field (name)",
			input: &types.CreateStockMaterialDTO{
				Description:  "Test material description",
				SafetyStock:  50.5,
				UnitID:       1,
				CategoryID:   1,
				IngredientID: 1,
			},
			expectError: true,
		},
		{
			name: "Create with invalid safety stock (0)",
			input: &types.CreateStockMaterialDTO{
				Name:         "Test Material",
				Description:  "Test description",
				SafetyStock:  0, // Should be > 0
				UnitID:       1,
				CategoryID:   1,
				IngredientID: 1,
			},
			expectError: true,
		},
		{
			name: "Create with non-existing unit ID",
			input: &types.CreateStockMaterialDTO{
				Name:         "Test Material",
				Description:  "Test description",
				SafetyStock:  50.5,
				UnitID:       999,
				CategoryID:   1,
				IngredientID: 1,
			},
			expectError: true,
		},
		{
			name: "Create with non-existing category ID",
			input: &types.CreateStockMaterialDTO{
				Name:         "Test Material",
				Description:  "Test description",
				SafetyStock:  50.5,
				UnitID:       1,
				CategoryID:   999,
				IngredientID: 1,
			},
			expectError: true,
		},
		{
			name: "Create with non-existing ingredient ID",
			input: &types.CreateStockMaterialDTO{
				Name:         "Test Material",
				Description:  "Test description",
				SafetyStock:  50.5,
				UnitID:       1,
				CategoryID:   1,
				IngredientID: 999,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := module.Service.CreateStockMaterial(tc.input)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.input.Name, result.Name)
				assert.Equal(t, tc.input.Description, result.Description)
				assert.Equal(t, tc.input.SafetyStock, result.SafetyStock)
				assert.Equal(t, tc.input.UnitID, result.Unit.ID)
				assert.Equal(t, tc.input.CategoryID, result.Category.ID)
				assert.Equal(t, tc.input.IngredientID, result.Ingredient.ID)
				assert.Equal(t, tc.input.Barcode, result.Barcode)
				assert.Equal(t, tc.input.ExpirationPeriodInDays, result.ExpirationPeriodInDays)
				assert.Equal(t, tc.input.Size, result.Size)
				assert.True(t, result.IsActive) // Should be active by default
				assert.NotEmpty(t, result.CreatedAt)
				assert.NotEmpty(t, result.UpdatedAt)
			}
		})
	}
}

func TestStockMaterialService_UpdateStockMaterial_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetStockMaterialModule()

	testCases := []struct {
		name        string
		id          uint
		update      *types.UpdateStockMaterialDTO
		expectError bool
	}{
		{
			name: "Update existing stock material",
			id:   1,
			update: &types.UpdateStockMaterialDTO{
				Name:        tests.StringPtr("Updated Material"),
				Description: tests.StringPtr("Updated description"),
				SafetyStock: tests.FloatPtr(150.0),
			},
			expectError: false,
		},
		{
			name: "Update non-existing stock material",
			id:   999,
			update: &types.UpdateStockMaterialDTO{
				Name: tests.StringPtr("Updated Material"),
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := module.Service.UpdateStockMaterial(tc.id, tc.update)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the update
				updated, err := module.Service.GetStockMaterialByID(tc.id)
				assert.NoError(t, err)
				assert.Equal(t, *tc.update.Name, updated.Name)
				assert.Equal(t, *tc.update.Description, updated.Description)
				assert.Equal(t, *tc.update.SafetyStock, updated.SafetyStock)
			}
		})
	}
}

func TestStockMaterialService_DeleteStockMaterial_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetStockMaterialModule()

	testCases := []struct {
		name        string
		id          uint
		expectError bool
	}{
		{
			name:        "Delete existing stock material",
			id:          1,
			expectError: false,
		},
		{
			name:        "Delete non-existing stock material",
			id:          999,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := module.Service.DeleteStockMaterial(tc.id)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify deletion
				_, err := module.Service.GetStockMaterialByID(tc.id)
				assert.Error(t, err)
			}
		})
	}
}

func TestStockMaterialService_DeactivateStockMaterial_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetStockMaterialModule()

	testCases := []struct {
		name        string
		id          uint
		expectError bool
	}{
		{
			name:        "Deactivate existing stock material",
			id:          1,
			expectError: false,
		},
		{
			name:        "Deactivate non-existing stock material",
			id:          999,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := module.Service.DeactivateStockMaterial(tc.id)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify deactivation
				material, err := module.Service.GetStockMaterialByID(tc.id)
				assert.NoError(t, err)
				assert.False(t, material.IsActive)
			}
		})
	}
}

func TestStockMaterialService_BarcodeOperations_WithPreloadedData(t *testing.T) {
	_ = ResetTestData(t)
	module := tests.GetStockMaterialModule()

	t.Run("Get existing barcode", func(t *testing.T) {
		barcode, err := module.Service.GetStockMaterialBarcode(1)
		assert.NoError(t, err)
		assert.NotNil(t, barcode)
	})

	t.Run("Generate barcode PDF", func(t *testing.T) {
		pdf, err := module.Service.GenerateStockMaterialBarcodePDF(1)
		assert.NoError(t, err)
		assert.NotNil(t, pdf)
	})

	t.Run("Generate new barcode", func(t *testing.T) {
		response, err := module.Service.GenerateBarcode()
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Barcode)
	})

	t.Run("Retrieve by barcode", func(t *testing.T) {
		material, err := module.Service.RetrieveStockMaterialByBarcode("CB001")
		assert.NoError(t, err)
		assert.NotNil(t, material)
		assert.Equal(t, "Coffee Beans Material", material.Name)
	})

	t.Run("Retrieve by non-existing barcode", func(t *testing.T) {
		material, err := module.Service.RetrieveStockMaterialByBarcode("NONEXIST")
		assert.Error(t, err)
		assert.Nil(t, material)
	})
}
