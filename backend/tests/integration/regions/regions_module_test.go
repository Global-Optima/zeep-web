package regions_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestRegionEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Region", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new region",
				Method:      http.MethodPost,
				URL:         "/api/test/regions",
				Body: map[string]interface{}{
					"name": "North Region",
				},
				AuthRole:     data.RoleAdmin, // Only Admin allowed
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a region",
				Method:      http.MethodPost,
				URL:         "/api/test/regions",
				Body: map[string]interface{}{
					"name": "North Region",
				},
				AuthRole:     data.RoleBarista, // Unauthorized role
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Regions", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all regions",
				Method:       http.MethodGet,
				URL:          "/api/test/regions",
				AuthRole:     data.RoleAdmin, // Admin access
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Region Manager should NOT(why?) fetch all regions",
				Method:       http.MethodGet,
				URL:          "/api/test/regions",
				AuthRole:     data.RoleRegionWarehouseManager, // Region Manager access
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description:  "Barista should NOT be able to fetch regions",
				Method:       http.MethodGet,
				URL:          "/api/test/regions",
				AuthRole:     data.RoleBarista, // Unauthorized role
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Region by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a region by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/regions/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Region Manager should NOT(why?) fetch a region by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/regions/1",
				AuthRole:     data.RoleRegionWarehouseManager,
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description:  "Unauthorized user should NOT fetch a region",
				Method:       http.MethodGet,
				URL:          "/api/test/regions/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description:  "Should return 404 for non-existent region", // 404 not shown, using 500 instead
				Method:       http.MethodGet,
				URL:          "/api/test/regions/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Region", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update a region",
				Method:      http.MethodPut,
				URL:         "/api/test/regions/1",
				Body: map[string]interface{}{
					"name": "Updated Region",
				},
				AuthRole:     data.RoleAdmin, // Only Admin allowed
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Region Manager should NOT be able to update a region",
				Method:      http.MethodPut,
				URL:         "/api/test/regions/1",
				Body: map[string]interface{}{
					"name": "Updated Region",
				},
				AuthRole:     data.RoleRegionWarehouseManager, // Not allowed to update
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Region", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should delete a region",
				Method:       http.MethodDelete,
				URL:          "/api/test/regions/1",
				AuthRole:     data.RoleAdmin, // Only Admin allowed
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete a region",
				Method:       http.MethodDelete,
				URL:          "/api/test/regions/1",
				AuthRole:     data.RoleBarista, // Unauthorized role
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
