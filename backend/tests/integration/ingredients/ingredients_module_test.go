package ingredients_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestIngredientEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create an Ingredient", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should create a new ingredient",
				Method:      http.MethodPost,
				URL:         "/ingredients",
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
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Ingredients", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of ingredients",
				Method:       http.MethodGet,
				URL:          "/ingredients",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch an Ingredient by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return an ingredient by ID",
				Method:       http.MethodGet,
				URL:          "/ingredients/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent ingredient",
				Method:       http.MethodGet,
				URL:          "/ingredients/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update an Ingredient", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should update an ingredient",
				Method:      http.MethodPut,
				URL:         "/ingredients/1",
				Body: map[string]interface{}{
					"name":     "Brown Sugar",
					"calories": 380,
				},
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete an Ingredient", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should delete an ingredient",
				Method:       http.MethodDelete,
				URL:          "/ingredients/1",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
