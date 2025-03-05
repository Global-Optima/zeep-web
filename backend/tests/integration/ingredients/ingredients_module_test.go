package ingredients_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestIngredientEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create an Ingredient", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new ingredient",
				Method:      http.MethodPost,
				URL:         "/api/test/ingredients",
				Body: map[string]interface{}{
					"name":             "Sugar",
					"calories":         400,
					"fat":              0,
					"carbs":            100,
					"proteins":         0,
					"categoryId":       1,
					"unitId":           1,
					"expirationInDays": 365,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a new ingredient",
				Method:      http.MethodPost,
				URL:         "/api/test/ingredients",
				Body: map[string]interface{}{
					"name":             "Sugar",
					"calories":         400,
					"fat":              0,
					"carbs":            100,
					"proteins":         0,
					"categoryId":       1,
					"unitId":           1,
					"expirationInDays": 365,
				},
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Ingredients", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all ingredients",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredients",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch all ingredients",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredients",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should be able to fetch ingredients",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredients",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch an Ingredient by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch an ingredient by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredients/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch an ingredient by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredients/1",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should be able to fetch an ingredient by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredients/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 500 for non-existent ingredient", // Note: Your service returns 404 instead of 500 here if not found.
				Method:       http.MethodGet,
				URL:          "/api/test/ingredients/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update an Ingredient", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update an ingredient",
				Method:      http.MethodPut,
				URL:         "/api/test/ingredients/1",
				Body: map[string]interface{}{
					"name":     "Brown Sugar",
					"calories": 380,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Store Manager should NOT be able to update an ingredient",
				Method:      http.MethodPut,
				URL:         "/api/test/ingredients/1",
				Body: map[string]interface{}{
					"name":     "Brown Sugar",
					"calories": 380,
				},
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete an Ingredient", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should not delete used ingredient",
				Method:       http.MethodDelete,
				URL:          "/api/test/ingredients/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusConflict,
			},
			{
				Description:  "Admin should delete unused ingredient",
				Method:       http.MethodDelete,
				URL:          "/api/test/ingredients/3",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete an ingredient",
				Method:       http.MethodDelete,
				URL:          "/api/test/ingredients/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
