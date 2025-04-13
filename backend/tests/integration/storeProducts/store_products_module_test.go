package storeProducts_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestStoreProductEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Get Store Product", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a valid store product",
				Method:       http.MethodGet,
				URL:          "/api/test/store-products/1?storeId=1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent store product",
				Method:       http.MethodGet,
				URL:          "/api/test/store-products/9999?storeId=1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Description:  "Unauthorized user should NOT fetch store product",
				Method:       http.MethodGet,
				URL:          "/api/test/store-products/1?storeId=1",
				AuthRole:     "",
				ExpectedCode: http.StatusUnauthorized,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Get Available Products to Add", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin fetches available products without filter",
				Method:       http.MethodGet,
				URL:          "/api/test/store-products/available-to-add?storeId=1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Filter available products by category",
				Method:       http.MethodGet,
				URL:          "/api/test/store-products/available-to-add?categoryId=2&storeId=1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Search for a specific product",
				Method:       http.MethodGet,
				URL:          "/api/test/store-products/available-to-add?search=Latte&storeId=1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Get Recommended Products for order", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin fetches recommended products",
				Method:      http.MethodGet,
				URL:         "/api/test/store-products/recommended?storeId=1",
				Body: map[string]interface{}{
					"storeProductIds": []uint{1, 2},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should fail if no products chosen",
				Method:       http.MethodGet,
				URL:          "/api/test/store-products/recommended?storeId=1",
				Body:         map[string]interface{}{},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Create Store Product", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new store product",
				Method:      http.MethodPost,
				URL:         "/api/test/store-products?storeId=1",
				Body: map[string]interface{}{
					"productId":   3,
					"isAvailable": true,
					"productSizes": []map[string]interface{}{
						{"productSizeID": 3, "storePrice": 10.99},
					},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Should fail with missing required fields",
				Method:      http.MethodPost,
				URL:         "/api/test/store-products?storeId=1",
				Body: map[string]interface{}{
					"isAvailable": true, // Missing productId
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Unauthorized user should NOT create store product",
				Method:      http.MethodPost,
				URL:         "/api/test/store-products?storeId=1",
				Body: map[string]interface{}{
					"productId":   1,
					"isAvailable": true,
				},
				AuthRole:     "",
				ExpectedCode: http.StatusUnauthorized,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update Store Product", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update store product availability",
				Method:      http.MethodPut,
				URL:         "/api/test/store-products/1?storeId=1",
				Body: map[string]interface{}{
					"isAvailable": false,
					"productSizes": []map[string]interface{}{
						{"productSizeID": 1, "storePrice": 10.99},
						{"productSizeID": 2, "storePrice": 12.99},
					},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Should return 400 for invalid product ID",
				Method:      http.MethodPut,
				URL:         "/api/test/store-products/-1?storeId=1",
				Body: map[string]interface{}{
					"isAvailable": false,
					"productSizes": []map[string]interface{}{
						{"productSizeID": 1, "storePrice": 10.99},
						{"productSizeID": 2, "storePrice": 12.99},
					},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Unauthorized user should NOT update store product",
				Method:      http.MethodPut,
				URL:         "/api/test/store-products/1?storeId=1",
				Body: map[string]interface{}{
					"isAvailable": false,
					"productSizes": []map[string]interface{}{
						{"productSizeID": 1, "storePrice": 10.99},
						{"productSizeID": 2, "storePrice": 12.99},
					},
				},
				AuthRole:     "",
				ExpectedCode: http.StatusUnauthorized,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete Store Product", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should not delete store product in use",
				Method:       http.MethodDelete,
				URL:          "/api/test/store-products/1?storeId=1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusConflict,
			},
			{
				Description:  "Should return 404 for non-existent product",
				Method:       http.MethodDelete,
				URL:          "/api/test/store-products/9999?storeId=1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Description:  "Unauthorized user should NOT delete store product",
				Method:       http.MethodDelete,
				URL:          "/api/test/store-products/1?storeId=1",
				AuthRole:     "",
				ExpectedCode: http.StatusUnauthorized,
			},
		}
		env.RunTests(t, testCases)
	})
}
