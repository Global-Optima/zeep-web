package categories_test

import (
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
	"github.com/stretchr/testify/suite"
)

type CategoriesIntegrationTestSuite struct {
	suite.Suite
	Env *utils.TestEnvironment
}

func (suite *CategoriesIntegrationTestSuite) SetupSuite() {
	suite.Env = utils.NewTestEnvironment(suite.T())
}

func (suite *CategoriesIntegrationTestSuite) TearDownSuite() {
	suite.Env.Close()
}
func (suite *CategoriesIntegrationTestSuite) TestGetAllCategories() {
	testCases := []utils.TestCase{
		{
			Description:  "Valid request for all categories",
			Method:       "GET",
			URL:          "/api/test/categories",
			ExpectedCode: 200,
			ExpectedBody: []map[string]interface{}{
				{
					"id":          1,
					"name":        "Beverages",
					"description": "Category for drinks",
				},
				{
					"id":          2,
					"name":        "Snacks",
					"description": "Category for snacks",
				},
			},
		},
		// {
		// 	Description:  "Simulate internal server error",
		// 	Method:       "GET",
		// 	URL:          "/api/test/categories",
		// 	ExpectedCode: 500,
		// 	ExpectedBody: map[string]interface{}{
		// 		"error": "Failed to retrieve categories",
		// 	},
		// },
	}

	suite.Env.RunTests(suite.T(), testCases)
}

func TestCategoriesIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CategoriesIntegrationTestSuite))
}
