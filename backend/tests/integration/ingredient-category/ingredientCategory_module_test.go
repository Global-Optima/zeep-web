package ingredientCategories_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestIngredientCategoryEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create an Ingredient Category", func(t *testing.T) {
		longDescription := strings.Repeat("This is a very long description. ", 20)
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new ingredient category",
				Method:      http.MethodPost,
				URL:         "/api/test/ingredient-categories",
				Body: map[string]interface{}{
					"name":        "Spices",
					"description": "Used in food preparation",
				},
				AuthRole:     data.RoleAdmin, // Only Admin can create
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a new ingredient category",
				Method:      http.MethodPost,
				URL:         "/api/test/ingredient-categories",
				Body: map[string]interface{}{
					"name":        "Spices",
					"description": "Used in food preparation",
				},
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description: "Admin should create a category with special characters in the name",
				Method:      http.MethodPost,
				URL:         "/api/test/ingredient-categories",
				Body: map[string]interface{}{
					"name":        "Spices!@#$%^&*()_+",
					"description": "Category with special characters in the name",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Admin should create a category with a long description",
				Method:      http.MethodPost,
				URL:         "/api/test/ingredient-categories",
				Body: map[string]interface{}{
					"name":        "Long Desc Category",
					"description": longDescription,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Ingredient Categories", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all ingredient categories",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredient-categories",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch all ingredient categories",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredient-categories",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should be able to fetch ingredient categories",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredient-categories",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch an Ingredient Category by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch an ingredient category by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredient-categories/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch an ingredient category by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredient-categories/1",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Unauthorized user should NOT fetch an ingredient category",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredient-categories/1",
				ExpectedCode: http.StatusUnauthorized,
			},
			{
				Description:  "Should return 500 for non-existent ingredient category",
				Method:       http.MethodGet,
				URL:          "/api/test/ingredient-categories/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update an Ingredient Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update an ingredient category",
				Method:      http.MethodPut,
				URL:         "/api/test/ingredient-categories/1",
				Body: map[string]interface{}{
					"name":        "Herbs",
					"description": "Fresh and dried herbs",
				},
				AuthRole:     data.RoleAdmin, // Only Admin allowed
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Store Manager should NOT be able to update an ingredient category",
				Method:      http.MethodPut,
				URL:         "/api/test/ingredient-categories/1",
				Body: map[string]interface{}{
					"name":        "Herbs",
					"description": "Fresh and dried herbs",
				},
				AuthRole:     data.RoleStoreManager, // Not permitted to update
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete an Ingredient Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should delete an ingredient category",
				Method:       http.MethodDelete,
				URL:          "/api/test/ingredient-categories/1",
				AuthRole:     data.RoleAdmin, // Only Admin allowed
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete an ingredient category",
				Method:       http.MethodDelete,
				URL:          "/api/test/ingredient-categories/1",
				AuthRole:     data.RoleBarista, // Unauthorized role
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
