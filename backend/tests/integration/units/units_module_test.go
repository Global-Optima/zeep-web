package units_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestUnitEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Unit", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should create a new unit",
				Method:      http.MethodPost,
				URL:         "/api/v1/units",
				Body: map[string]interface{}{
					"name":             "Kilogram",
					"conversionFactor": 1.0,
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Units", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of units",
				Method:       http.MethodGet,
				URL:          "/units",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Unit by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a unit by ID",
				Method:       http.MethodGet,
				URL:          "/units/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent unit",
				Method:       http.MethodGet,
				URL:          "/units/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Unit", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should update a unit",
				Method:      http.MethodPut,
				URL:         "/units/1",
				Body: map[string]interface{}{
					"name":             "Gram",
					"conversionFactor": 0.001,
				},
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Unit", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should delete a unit",
				Method:       http.MethodDelete,
				URL:          "/units/1",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
