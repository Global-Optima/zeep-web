package additives_test

import (
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
	"github.com/stretchr/testify/suite"
)

type AdditivesIntegrationTestSuite struct {
	suite.Suite
	Env *utils.TestEnvironment
}

func (suite *AdditivesIntegrationTestSuite) SetupSuite() {
	suite.Env = utils.NewTestEnvironment(suite.T())
}

func (suite *AdditivesIntegrationTestSuite) TearDownSuite() {
	suite.Env.Close()
}

func (suite *AdditivesIntegrationTestSuite) TestGetAdditivesByStoreAndProduct() {
	testCases := []utils.TestCase{
		// {
		// 	Description:  "Valid request with storeId=1 and productId=1",
		// 	Method:       "GET",
		// 	URL:          "/api/test/additives?storeId=1&productId=1",
		// 	ExpectedCode: 200,
		// 	ExpectedBody: []map[string]interface{}{
		// 		{
		// 			"id":   1,
		// 			"name": "Sweeteners",
		// 			"additives": []map[string]interface{}{
		// 				{
		// 					"id":          1,
		// 					"name":        "Sugar",
		// 					"description": "Sweet sugar",
		// 					"price":       0.60, // Adjust based on your mock data
		// 					"imageUrl":    "https://example.com/sugar.jpg",
		// 				},
		// 				{
		// 					"id":          3,
		// 					"name":        "Honey",
		// 					"description": "Natural sweetener",
		// 					"price":       1.00,
		// 					"imageUrl":    "https://example.com/honey.jpg",
		// 				},
		// 			},
		// 		},
		// 		{
		// 			"id":   2,
		// 			"name": "Dairy",
		// 			"additives": []map[string]interface{}{
		// 				{
		// 					"id":          2,
		// 					"name":        "Cream",
		// 					"description": "Adds richness",
		// 					"price":       1.50,
		// 					"imageUrl":    "https://example.com/cream.jpg",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	Description:  "Valid request with storeId=1 and productId=2",
		// 	Method:       "GET",
		// 	URL:          "/api/test/additives?storeId=1&productId=2",
		// 	ExpectedCode: 200,
		// 	ExpectedBody: []map[string]interface{}{
		// 		{
		// 			"id":   1,
		// 			"name": "Sweeteners",
		// 			"additives": []map[string]interface{}{
		// 				{
		// 					"id":          1,
		// 					"name":        "Sugar",
		// 					"description": "Sweet sugar",
		// 					"price":       0.60,
		// 					"imageUrl":    "https://example.com/sugar.jpg",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	Description:  "Valid request with storeId=2 and productId=3",
		// 	Method:       "GET",
		// 	URL:          "/api/test/additives?storeId=2&productId=3",
		// 	ExpectedCode: 200,
		// 	ExpectedBody: []map[string]interface{}{
		// 		{
		// 			"id":   1,
		// 			"name": "Sweeteners",
		// 			"additives": []map[string]interface{}{
		// 				{
		// 					"id":          3,
		// 					"name":        "Honey",
		// 					"description": "Natural sweetener",
		// 					"price":       1.20,
		// 					"imageUrl":    "https://example.com/honey.jpg",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	Description:  "Valid request with no additives available for storeId=1 and productId=3",
		// 	Method:       "GET",
		// 	URL:          "/api/test/additives?storeId=1&productId=3",
		// 	ExpectedCode: 200,
		// 	ExpectedBody: []map[string]interface{}{}, // Expecting an empty array when no additives are available
		// },
		// {
		// 	Description:  "Missing storeId parameter",
		// 	Method:       "GET",
		// 	URL:          "/api/test/additives?productId=1",
		// 	ExpectedCode: 400,
		// 	ExpectedBody: map[string]interface{}{
		// 		"error": "Invalid store ID",
		// 	},
		// },
		{
			Description:  "Missing productId parameter",
			Method:       "GET",
			URL:          "/api/test/additives?storeId=1",
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"error": "Invalid product ID",
			},
		},
		{
			Description:  "Invalid storeId format",
			Method:       "GET",
			URL:          "/api/test/additives?storeId=abc&productId=1",
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"error": "Invalid store ID",
			},
		},
		{
			Description:  "Invalid productId format",
			Method:       "GET",
			URL:          "/api/test/additives?storeId=1&productId=abc",
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"error": "Invalid product ID",
			},
		},
		{
			Description:  "Non-existent storeId and productId",
			Method:       "GET",
			URL:          "/api/test/additives?storeId=999&productId=999",
			ExpectedCode: 200,
			ExpectedBody: nil, // Expecting an empty array as no matching data exists
		},
	}

	suite.Env.RunTests(suite.T(), testCases)
}

func TestAdditivesIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AdditivesIntegrationTestSuite))
}
