package provisions

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestProvisionEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Provision", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin creates a valid provision with ingredients",
				Method:      http.MethodPost,
				URL:         "/api/test/provisions",
				Body: map[string]interface{}{
					"name":                       "Test Provision",
					"absoluteVolume":             0.5,
					"netCost":                    200.0,
					"unitId":                     1,
					"preparationInMinutes":       15,
					"defaultExpirationInMinutes": 60,
					"limitPerDay":                3,
					"ingredients": []map[string]interface{}{
						{"ingredientId": 1, "quantity": 0.2},
						{"ingredientId": 2, "quantity": 0.3},
					},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Should fail with invalid ingredient ID",
				Method:      http.MethodPost,
				URL:         "/api/test/provisions",
				Body: map[string]interface{}{
					"name":                       "Invalid Ingredient Provision",
					"absoluteVolume":             1.0,
					"netCost":                    100,
					"unitId":                     1,
					"preparationInMinutes":       20,
					"defaultExpirationInMinutes": 60,
					"limitPerDay":                3,
					"ingredients": []map[string]interface{}{
						{"ingredientId": 9999, "quantity": 1.0},
					},
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Get Provisions", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin fetches all provisions",
				Method:       http.MethodGet,
				URL:          "/api/test/provisions",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Get Provision by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin fetches provision by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/provisions/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent ID",
				Method:       http.MethodGet,
				URL:          "/api/test/provisions/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Provision", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin updates an existing provision",
				Method:      http.MethodPut,
				URL:         "/api/test/provisions/1",
				Body: map[string]interface{}{
					"name":        "Updated Provision Name",
					"netCost":     300.0,
					"limitPerDay": 5,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Should fail updating with invalid unit ID",
				Method:      http.MethodPut,
				URL:         "/api/test/provisions/1",
				Body: map[string]interface{}{
					"unitId": 9999,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Provision", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should fail to delete provision in use",
				Method:       http.MethodDelete,
				URL:          "/api/test/provisions/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusConflict,
			},
			{
				Description:  "Admin deletes an unused provision",
				Method:       http.MethodDelete,
				URL:          "/api/test/provisions/3",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
