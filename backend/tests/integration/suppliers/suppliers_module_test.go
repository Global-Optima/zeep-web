package supplier_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestSupplierEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Supplier", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new supplier",
				Method:      http.MethodPost,
				URL:         "/api/test/suppliers",
				Body: map[string]interface{}{
					"name":         "Supplier A",
					"contactEmail": "supplier@example.com",
					"contactPhone": "+1234567890",
					"city":         "New York",
					"address":      "123 Supplier St",
				},
				AuthRole:     data.RoleAdmin, // Only Admin can create
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a new supplier",
				Method:      http.MethodPost,
				URL:         "/api/test/suppliers",
				Body: map[string]interface{}{
					"name":         "Supplier A",
					"contactEmail": "supplier@example.com",
					"contactPhone": "1234567890",
					"city":         "New York",
					"address":      "123 Supplier St",
				},
				AuthRole:     data.RoleBarista, // Not allowed to create
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Suppliers", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all suppliers",
				Method:       http.MethodGet,
				URL:          "/api/test/suppliers",
				AuthRole:     data.RoleAdmin, // Admin access
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager shouldn't fetch all suppliers",
				Method:       http.MethodGet,
				URL:          "/api/test/suppliers",
				AuthRole:     data.RoleStoreManager, // Allowed to read suppliers
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Supplier by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a supplier by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/suppliers/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Unauthorized user should not fetch a supplier",
				Method:       http.MethodGet,
				URL:          "/api/test/suppliers/1",
				AuthRole:     data.RoleBarista, // For example, a Barista may not have read permission
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description:  "Should return 404 for non-existent supplier",
				Method:       http.MethodGet,
				URL:          "/api/test/suppliers/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Supplier", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update a supplier",
				Method:      http.MethodPut,
				URL:         "/api/test/suppliers/1",
				Body: map[string]interface{}{
					"name": "Updated Supplier A",
					"city": "Los Angeles",
				},
				AuthRole:     data.RoleAdmin, // Only Admin can update
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Store Manager should NOT be able to update a supplier",
				Method:      http.MethodPut,
				URL:         "/api/test/suppliers/1",
				Body: map[string]interface{}{
					"name": "Updated Supplier A",
					"city": "Los Angeles",
				},
				AuthRole:     data.RoleStoreManager, // Not permitted to update
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Supplier", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should delete a supplier",
				Method:       http.MethodDelete,
				URL:          "/api/test/suppliers/1",
				AuthRole:     data.RoleAdmin, // Only Admin can delete
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete a supplier",
				Method:       http.MethodDelete,
				URL:          "/api/test/suppliers/1",
				AuthRole:     data.RoleBarista, // Not permitted to delete
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
