package stockMaterialCategory_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestStockMaterialCategoryEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Stock Material Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should create a new stock material category",
				Method:      http.MethodPost,
				URL:         "/stock-material-categories",
				Body: map[string]interface{}{
					"name":        "Metal",
					"description": "Used in construction",
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Stock Material Categories", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of stock material categories",
				Method:       http.MethodGet,
				URL:          "/stock-material-categories",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Stock Material Category by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a stock material category by ID",
				Method:       http.MethodGet,
				URL:          "/stock-material-categories/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent category",
				Method:       http.MethodGet,
				URL:          "/stock-material-categories/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Stock Material Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should update a stock material category",
				Method:      http.MethodPut,
				URL:         "/stock-material-categories/1",
				Body: map[string]interface{}{
					"name":        "Plastic",
					"description": "Synthetic material",
				},
				ExpectedCode: http.StatusNoContent,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Stock Material Category", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should delete a stock material category",
				Method:       http.MethodDelete,
				URL:          "/stock-material-categories/1",
				ExpectedCode: http.StatusNoContent,
			},
		}
		env.RunTests(t, testCases)
	})
}
