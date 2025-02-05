package supplier_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestSupplierEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Supplier", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should create a new supplier",
				Method:      http.MethodPost,
				URL:         "/suppliers",
				Body: map[string]interface{}{
					"name":         "Supplier A",
					"contactEmail": "supplier@example.com",
					"contactPhone": "1234567890",
					"city":         "New York",
					"address":      "123 Supplier St",
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Suppliers", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of suppliers",
				Method:       http.MethodGet,
				URL:          "/suppliers",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Supplier by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a supplier by ID",
				Method:       http.MethodGet,
				URL:          "/suppliers/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent supplier",
				Method:       http.MethodGet,
				URL:          "/suppliers/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Supplier", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should update a supplier",
				Method:      http.MethodPut,
				URL:         "/suppliers/1",
				Body: map[string]interface{}{
					"name": "Updated Supplier A",
					"city": "Los Angeles",
				},
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Supplier", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should delete a supplier",
				Method:       http.MethodDelete,
				URL:          "/suppliers/1",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
