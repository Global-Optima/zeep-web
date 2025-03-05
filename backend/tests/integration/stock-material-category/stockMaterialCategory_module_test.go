package stockMaterialCategory_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestStockMaterialCategoryEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Stock Material Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new stock material category",
				Method:      http.MethodPost,
				URL:         "/api/test/stock-material-categories",
				Body: map[string]interface{}{
					"name":        "Metal",
					"description": "Used in construction",
				},
				AuthRole:     data.RoleAdmin, // Only Admin can create
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a new stock material category",
				Method:      http.MethodPost,
				URL:         "/api/test/stock-material-categories",
				Body: map[string]interface{}{
					"name":        "Metal",
					"description": "Used in construction",
				},
				AuthRole:     data.RoleBarista, // Unauthorized role
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Stock Material Categories", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all stock material categories",
				Method:       http.MethodGet,
				URL:          "/api/test/stock-material-categories",
				AuthRole:     data.RoleAdmin, // Admin access
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch all stock material categories",
				Method:       http.MethodGet,
				URL:          "/api/test/stock-material-categories",
				AuthRole:     data.RoleStoreManager, // Allowed read access
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Stock Material Category by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a stock material category by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/stock-material-categories/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 500 for non-existent category",
				Method:       http.MethodGet,
				URL:          "/api/test/stock-material-categories/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Stock Material Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update a stock material category",
				Method:      http.MethodPut,
				URL:         "/api/test/stock-material-categories/1",
				Body: map[string]interface{}{
					"name":        "Plastic",
					"description": "Synthetic material",
				},
				AuthRole:     data.RoleAdmin, // Only Admin can update
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Store Manager should NOT be able to update a stock material category",
				Method:      http.MethodPut,
				URL:         "/api/test/stock-material-categories/1",
				Body: map[string]interface{}{
					"name":        "Plastic",
					"description": "Synthetic material",
				},
				AuthRole:     data.RoleStoreManager, // Not permitted to update
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Stock Material Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should not delete used stock material category",
				Method:       http.MethodDelete,
				URL:          "/api/test/stock-material-categories/1",
				AuthRole:     data.RoleAdmin, // Only Admin can delete
				ExpectedCode: http.StatusConflict,
			},
			{
				Description:  "Admin should delete unsed stock material category",
				Method:       http.MethodDelete,
				URL:          "/api/test/stock-material-categories/2",
				AuthRole:     data.RoleAdmin, // Only Admin can delete
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete a stock material category",
				Method:       http.MethodDelete,
				URL:          "/api/test/stock-material-categories/1",
				AuthRole:     data.RoleBarista, // Unauthorized role
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
