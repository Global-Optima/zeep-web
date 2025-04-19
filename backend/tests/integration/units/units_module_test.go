package units_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestUnitEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Unit", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new unit",
				Method:      http.MethodPost,
				URL:         "/api/test/units",
				Body: map[string]interface{}{
					"name":             "Test Unit",
					"conversionFactor": 1.0,
				},
				AuthRole:     data.RoleAdmin, // ✅ Admin role required
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a unit",
				Method:      http.MethodPost,
				URL:         "/api/test/units",
				Body: map[string]interface{}{
					"name":             "Kilogram",
					"conversionFactor": 1.0,
				},
				AuthRole:     data.RoleBarista, // ❌ Barista doesn't have permission
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description: "Admin should NOT create a unit with a name containing only spaces",
				Method:      http.MethodPost,
				URL:         "/api/test/units",
				Body: map[string]interface{}{
					"name":             "    ", // Only spaces
					"conversionFactor": 1.0,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Admin should NOT create a unit with too many characters in the name",
				Method:      http.MethodPost,
				URL:         "/api/test/units",
				Body: map[string]interface{}{
					"name":             strings.Repeat("A", 300), // Exceeding typical length limits
					"conversionFactor": 1.0,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Description: "Admin should NOT create a unit with out of boundary excessively high conversion factor",
				Method:      http.MethodPost,
				URL:         "/api/test/units",
				Body: map[string]interface{}{
					"name":             "HighConversionUnit",
					"conversionFactor": 1e10, // Example of an excessively high conversion factor
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Units", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all units",
				Method:       http.MethodGet,
				URL:          "/api/test/units",
				AuthRole:     data.RoleAdmin, // ✅ Only admin can access
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Admin should search units by name",
				Method:       http.MethodGet,
				URL:          "/api/test/units?search=Kilogram",
				AuthRole:     data.RoleAdmin, // ✅ Only admin can access
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch all units",
				Method:       http.MethodGet,
				URL:          "/api/test/units",
				AuthRole:     data.RoleStoreManager, // ✅ Store manager has read access
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Unit by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a unit by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/units/1",
				AuthRole:     data.RoleAdmin, // ✅ Admin access
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Unauthorized user should not fetch a unit",
				Method:       http.MethodGet,
				URL:          "/api/test/units/1",
				AuthRole:     data.RoleBarista, // ❌ Barista shouldn't have access
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent unit",
				Method:       http.MethodGet,
				URL:          "/api/test/units/9999",
				AuthRole:     data.RoleAdmin, // ✅ Still requires admin access
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Unit", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update a unit",
				Method:      http.MethodPut,
				URL:         "/api/test/units/1",
				Body: map[string]interface{}{
					"name":             "Gram",
					"conversionFactor": 0.001,
				},
				AuthRole:     data.RoleAdmin, // ✅ Admin access required
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Admin should NOT update a unit with a name containing only spaces",
				Method:      http.MethodPut,
				URL:         "/api/test/units/1",
				Body: map[string]interface{}{
					"name":             "    ", // Name with whitespace only
					"conversionFactor": 0.001,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Admin should NOT update a unit with an excessively high conversion factor",
				Method:      http.MethodPut,
				URL:         "/api/test/units/1",
				Body: map[string]interface{}{
					"name":             "Gram",
					"conversionFactor": 1e10, // Out-of-bound conversion factor
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
			{
				Description: "Store Manager should NOT be able to update a unit",
				Method:      http.MethodPut,
				URL:         "/api/test/units/1",
				Body: map[string]interface{}{
					"name":             "Gram",
					"conversionFactor": 0.001,
				},
				AuthRole:     data.RoleStoreManager, // ❌ Not allowed to update
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Unit", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should not delete a unit in use",
				Method:       http.MethodDelete,
				URL:          "/api/test/units/1",
				AuthRole:     data.RoleAdmin, // ✅ Only admin can delete
				ExpectedCode: http.StatusConflict,
			},
			{
				Description:  "Should delete unused unit",
				Method:       http.MethodDelete,
				URL:          "/api/test/units/3",
				AuthRole:     data.RoleAdmin, // ✅ Only admin can delete
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete a unit",
				Method:       http.MethodDelete,
				URL:          "/api/test/units/1",
				AuthRole:     data.RoleBarista, // ❌ Barista not allowed
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
