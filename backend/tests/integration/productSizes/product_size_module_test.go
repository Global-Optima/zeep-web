package productSizes_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestProductEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Product Size", func(t *testing.T) {
		testCases := []utils.TestCase{
			// {
			// 	Description: "Admin should create a new product size",
			// 	Method:      http.MethodPost,
			// 	URL:         "/api/test/products/sizes",
			// 	Body: map[string]interface{}{
			// 		"productId":   1,
			// 		"name":        "S",
			// 		"size":        300,
			// 		"unitId":      1,
			// 		"basePrice":   4.99,
			// 		"ingredients": []map[string]interface{}{{"ingredientId": 2, "quantity": 5.0}},
			// 		"additives":   []map[string]interface{}{{"additiveId": 1, "isDefault": false}},
			// 		"machineId":   "TEST0000111122223333000010",
			// 	},
			// 	AuthRole:     data.RoleAdmin,
			// 	ExpectedCode: http.StatusCreated,
			// },
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
					"ingredients": []map[string]interface{}{{"ingredientId": 2, "quantity": 5.0}},
					"additives":   []map[string]interface{}{{"additiveId": 1, "isDefault": false}},
				},
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description: "Admin should create a new product size with provisions",
				Method:      http.MethodPost,
				URL:         "/api/test/products/sizes",
				Body: map[string]interface{}{
					"productId":   1,
					"name":        "L",
					"size":        350,
					"unitId":      1,
					"basePrice":   6.49,
					"ingredients": []map[string]interface{}{{"ingredientId": 2, "quantity": 3.0}},
					"additives":   []map[string]interface{}{{"additiveId": 1, "isDefault": false}},
					"provisions":  []map[string]interface{}{{"provisionId": 1, "volume": 0.5}, {"provisionId": 2, "volume": 1.1}},
					"machineId":   "TEST0000111122224444555566",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Should fail with invalid provision ID",
				Method:      http.MethodPost,
				URL:         "/api/test/products/sizes",
				Body: map[string]interface{}{
					"productId":   1,
					"name":        "L",
					"size":        400,
					"unitId":      1,
					"basePrice":   6.99,
					"ingredients": []map[string]interface{}{{"ingredientId": 2, "quantity": 3.0}},
					"provisions":  []map[string]interface{}{{"provisionId": 0, "volume": 1.0}},
					"machineId":   "TEST0000111122227777888899",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Should fail with zero provision volume",
				Method:      http.MethodPost,
				URL:         "/api/test/products/sizes",
				Body: map[string]interface{}{
					"productId":   1,
					"name":        "XL",
					"size":        450,
					"unitId":      1,
					"basePrice":   7.49,
					"ingredients": []map[string]interface{}{{"ingredientId": 3, "quantity": 2.0}},
					"provisions":  []map[string]interface{}{{"provisionId": 1, "volume": 0}},
					"machineId":   "TEST0000111122221010101010",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Admin should create a product size with no provisions",
				Method:      http.MethodPost,
				URL:         "/api/test/products/sizes",
				Body: map[string]interface{}{
					"productId":   1,
					"name":        "S",
					"size":        300,
					"unitId":      1,
					"basePrice":   4.99,
					"ingredients": []map[string]interface{}{{"ingredientId": 2, "quantity": 5.0}},
					"provisions":  []map[string]interface{}{},
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
			{
				Description: "Admin should update product size and assign provisions",
				Method:      http.MethodPut,
				URL:         "/api/test/products/sizes/1",
				Body: map[string]interface{}{
					"name":       "XL Updated",
					"size":       550,
					"basePrice":  8.00,
					"provisions": []map[string]interface{}{{"provisionId": 1, "volume": 2.0}, {"provisionId": 2, "volume": 1.2}},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Update should fail with invalid provision ID",
				Method:      http.MethodPut,
				URL:         "/api/test/products/sizes/1",
				Body: map[string]interface{}{
					"provisions": []map[string]interface{}{{"provisionId": 0, "volume": 1.5}},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Update should fail with zero volume",
				Method:      http.MethodPut,
				URL:         "/api/test/products/sizes/1",
				Body: map[string]interface{}{
					"provisions": []map[string]interface{}{{"provisionId": 1, "volume": 0}},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Product Size", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should not delete a product size in use",
				Method:       http.MethodDelete,
				URL:          "/api/test/products/sizes/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusConflict,
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
