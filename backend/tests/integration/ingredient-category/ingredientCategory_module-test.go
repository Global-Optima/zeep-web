package ingredientCategories_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestIngredientCategoryEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create an Ingredient Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should create a new ingredient category",
				Method:      http.MethodPost,
				URL:         "/ingredient-categories",
				Body: map[string]interface{}{
					"name":        "Spices",
					"description": "Used in food preparation",
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Ingredient Categories", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of ingredient categories",
				Method:       http.MethodGet,
				URL:          "/ingredient-categories",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch an Ingredient Category by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return an ingredient category by ID",
				Method:       http.MethodGet,
				URL:          "/ingredient-categories/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent category",
				Method:       http.MethodGet,
				URL:          "/ingredient-categories/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update an Ingredient Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should update an ingredient category",
				Method:      http.MethodPut,
				URL:         "/ingredient-categories/1",
				Body: map[string]interface{}{
					"name":        "Herbs",
					"description": "Fresh and dried herbs",
				},
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete an Ingredient Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should delete an ingredient category",
				Method:       http.MethodDelete,
				URL:          "/ingredient-categories/1",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
