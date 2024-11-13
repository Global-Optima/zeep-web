package stores_test

import (
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
	"github.com/stretchr/testify/suite"
)

type StoresIntegrationTestSuite struct {
	suite.Suite
	Env *utils.TestEnvironment
}

func (suite *StoresIntegrationTestSuite) SetupSuite() {
	suite.Env = utils.NewTestEnvironment(suite.T())
}

func (suite *StoresIntegrationTestSuite) TearDownSuite() {
	suite.Env.Close()
}

// Test cases for GetAllStores endpoint
func (suite *StoresIntegrationTestSuite) TestGetAllStores() {
	testCases := []utils.TestCase{
		{
			Description:  "Valid request for all stores",
			Method:       "GET",
			URL:          "/api/test/stores",
			ExpectedCode: 200,
			ExpectedBody: []map[string]interface{}{
				{
					"id":          1,
					"name":        "Downtown Coffee",
					"isFranchise": true,
					"facilityAddress": map[string]interface{}{
						"id":      1,
						"address": "123 Coffee St",
					},
				},
				{
					"id":          2,
					"name":        "Uptown Snacks",
					"isFranchise": false,
					"facilityAddress": map[string]interface{}{
						"id":      2,
						"address": "456 Snack Ave",
					},
				},
			},
		},
	}

	suite.Env.RunTests(suite.T(), testCases)
}

func (suite *StoresIntegrationTestSuite) TestGetStoreEmployees() {
	testCases := []utils.TestCase{
		{
			Description:  "Valid request with storeId=1",
			Method:       "GET",
			URL:          "/api/test/stores/1/employees",
			ExpectedCode: 200,
			ExpectedBody: []map[string]interface{}{
				{
					"id":       1,
					"name":     "Alice Smith",
					"phone":    "+1234567890",
					"email":    "alice@example.com",
					"isActive": true,
					"role": map[string]interface{}{
						"id":   1,
						"name": "Manager",
					},
				},
				{
					"id":       2,
					"name":     "Bob Johnson",
					"phone":    "+0987654321",
					"email":    "bob@example.com",
					"isActive": true,
					"role": map[string]interface{}{
						"id":   2,
						"name": "Barista",
					},
				},
			},
		},
		{
			Description:  "Non-existent storeId",
			Method:       "GET",
			URL:          "/api/test/stores/999/employees",
			ExpectedCode: 200,
			ExpectedBody: []map[string]interface{}{}, // Expecting an empty array
		},
		{
			Description:  "Invalid storeId format",
			Method:       "GET",
			URL:          "/api/test/stores/abc/employees",
			ExpectedCode: 400,
			ExpectedBody: map[string]interface{}{
				"error": "Invalid store ID",
			},
		},
	}

	suite.Env.RunTests(suite.T(), testCases)
}

func TestStoresIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(StoresIntegrationTestSuite))
}
