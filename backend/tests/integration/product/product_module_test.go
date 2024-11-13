package product_test

import (
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
	"github.com/stretchr/testify/suite"
)

type ProductIntegrationTestSuite struct {
	suite.Suite
	Env *utils.TestEnvironment
}

func (suite *ProductIntegrationTestSuite) SetupSuite() {
	suite.Env = utils.NewTestEnvironment(suite.T())
}

func (suite *ProductIntegrationTestSuite) TearDownSuite() {
	suite.Env.Close()
}
func (suite *ProductIntegrationTestSuite) TestGetStoreProducts() {
	// Test cases for GetStoreProducts endpointfunc (suite *ProductIntegrationTestSuite) TestGetStoreProducts() {
	testCases := []utils.TestCase{
		// {
		// 	Description:  "Valid request with storeId=1 and categoryId=1",
		// 	Method:       "GET",
		// 	URL:          "/api/test/products?storeId=1&categoryId=1",
		// 	ExpectedCode: 200,
		// 	ExpectedBody: []map[string]interface{}{
		// 		{
		// 			"id":          1,
		// 			"name":        "Coffee",
		// 			"description": "Hot brewed coffee",
		// 			"imageUrl":    "https://example.com/coffee.jpg",
		// 			"basePrice":   3.25,
		// 		},
		// 		{
		// 			"id":          2,
		// 			"name":        "Latte",
		// 			"description": "Coffee with steamed milk",
		// 			"imageUrl":    "https://example.com/latte.jpg",
		// 			"basePrice":   4.25,
		// 		},
		// 	},
		// },
		// {
		// 	Description:  "Request with search query",
		// 	Method:       "GET",
		// 	URL:          "/api/test/products?storeId=1&categoryId=1&search=Latte",
		// 	ExpectedCode: 200,
		// 	ExpectedBody: []map[string]interface{}{
		// 		{
		// 			"id":          2,
		// 			"name":        "Latte",
		// 			"description": "Coffee with steamed milk",
		// 			"imageUrl":    "https://example.com/latte.jpg",
		// 			"basePrice":   4.25,
		// 		},
		// 	},
		// },
		{
			Description:  "Non-existent storeId or categoryId",
			Method:       "GET",
			URL:          "/api/test/products?storeId=999&categoryId=999",
			ExpectedCode: 200,
			ExpectedBody: []map[string]interface{}{}, // Empty array expected
		},
		{
			Description:  "Invalid storeId format",
			Method:       "GET",
			URL:          "/api/test/products?storeId=abc&categoryId=1",
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"error": "Invalid store Id",
			},
		},
		{
			Description:  "Invalid categoryId format",
			Method:       "GET",
			URL:          "/api/test/products?storeId=1&categoryId=abc",
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"error": "Invalid category Id",
			},
		},
	}

	suite.Env.RunTests(suite.T(), testCases)
}

func (suite *ProductIntegrationTestSuite) TestGetStoreProductDetails() {
	testCases := []utils.TestCase{
		// {
		// 	Description:  "Valid request with storeId=1 and productId=1",
		// 	Method:       "GET",
		// 	URL:          "/api/test/products/1?storeId=1",
		// 	ExpectedCode: 200,
		// 	ExpectedBody: map[string]interface{}{
		// 		"id":          1,
		// 		"name":        "Coffee",
		// 		"description": "Hot brewed coffee",
		// 		"imageUrl":    "https://example.com/coffee.jpg",
		// 		"sizes": []map[string]interface{}{
		// 			{
		// 				"id":        1,
		// 				"name":      "Small",
		// 				"basePrice": 3.25,
		// 				"measure":   "oz",
		// 			},
		// 		},
		// 		"defaultAdditives": []map[string]interface{}{
		// 			{
		// 				"id":          1,
		// 				"name":        "Sugar",
		// 				"description": "Adds sweetness",
		// 				"imageUrl":    "https://example.com/sugar.jpg",
		// 			},
		// 		},
		// 		"recipeSteps": []map[string]interface{}{
		// 			{
		// 				"id":          1,
		// 				"description": "Add coffee grounds",
		// 				"imageUrl":    "https://example.com/step1.jpg",
		// 				"step":        1,
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	Description:  "Non-existent storeId or productId",
		// 	Method:       "GET",
		// 	URL:          "/api/test/products/999?storeId=999",
		// 	ExpectedCode: 404,
		// 	ExpectedBody: map[string]interface{}{
		// 		"error": "Product not found",
		// 	},
		// },
		// {
		// 	Description:  "Invalid storeId format",
		// 	Method:       "GET",
		// 	URL:          "/api/test/products/1?storeId=abc",
		// 	ExpectedCode: 400,
		// 	ExpectedBody: map[string]interface{}{
		// 		"error": "Invalid store ID",
		// 	},
		// },
		// {
		// 	Description:  "Invalid productId format",
		// 	Method:       "GET",
		// 	URL:          "/api/test/products/abc?storeId=1",
		// 	ExpectedCode: 400,
		// 	ExpectedBody: map[string]interface{}{
		// 		"error": "Invalid product ID",
		// 	},
		// },
	}

	suite.Env.RunTests(suite.T(), testCases)
}

func TestProductIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ProductIntegrationTestSuite))
}
