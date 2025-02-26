package productSizes_test

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestProductEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Product Size", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new product size",
				Method:      http.MethodPost,
				URL:         "/api/test/products/sizes",
				Body: map[string]interface{}{
					"productId":   1,
					"name":        "S",
					"size":        300,
					"unitId":      1,
					"basePrice":   4.99,
					"isDefault":   false,
					"ingredients": []map[string]interface{}{{"ingredientId": 2, "quantity": 5.0}},
					"additives":   []map[string]interface{}{{"additiveId": 1, "isDefault": false}},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a product size",
				Method:      http.MethodPost,
				URL:         "/api/test/products/sizes",
				Body: map[string]interface{}{
					"productId":   1,
					"name":        "S",
					"size":        300,
					"unitId":      1,
					"basePrice":   4.99,
					"isDefault":   false,
					"ingredients": []map[string]interface{}{{"ingredientId": 2, "quantity": 5.0}},
					"additives":   []map[string]interface{}{{"additiveId": 1, "isDefault": false}},
				},
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch Product Sizes by Product ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all product sizes for a product",
				Method:       http.MethodGet,
				URL:          "/api/test/products/1/sizes",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch all product sizes for a product",
				Method:       http.MethodGet,
				URL:          "/api/test/products/1/sizes",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should be able to fetch all product sizes for a product",
				Method:       http.MethodGet,
				URL:          "/api/test/products/1/sizes",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Product Size by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a product size by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/products/sizes/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch a product size by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/products/sizes/1",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should be able to fetch a product size by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/products/sizes/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for a non-existent product size",
				Method:       http.MethodGet,
				URL:          "/api/test/products/sizes/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Product Size", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update a product size",
				Method:      http.MethodPut,
				URL:         "/api/test/products/sizes/1",
				Body: map[string]interface{}{
					"name":      "Large",
					"size":      500,
					"basePrice": 5.99,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Store Manager should NOT be able to update a product size",
				Method:      http.MethodPut,
				URL:         "/api/test/products/sizes/1",
				Body: map[string]interface{}{
					"name":      "Large",
					"size":      500,
					"basePrice": 5.99,
				},
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Product Size", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should delete a product size",
				Method:       http.MethodDelete,
				URL:          "/api/test/products/sizes/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete a product size",
				Method:       http.MethodDelete,
				URL:          "/api/test/products/sizes/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
