package regions_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestRegionEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Region", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should create a new region",
				Method:      http.MethodPost,
				URL:         "/regions",
				Body: map[string]interface{}{
					"name": "North Region",
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Regions", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of regions",
				Method:       http.MethodGet,
				URL:          "/regions",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Region by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a region by ID",
				Method:       http.MethodGet,
				URL:          "/regions/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent region",
				Method:       http.MethodGet,
				URL:          "/regions/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Region", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should update a region",
				Method:      http.MethodPut,
				URL:         "/regions/1",
				Body: map[string]interface{}{
					"name": "Updated Region",
				},
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Region", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should delete a region",
				Method:       http.MethodDelete,
				URL:          "/regions/1",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
