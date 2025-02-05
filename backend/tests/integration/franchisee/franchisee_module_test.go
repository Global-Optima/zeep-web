package franchisees_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestFranchiseeEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Franchisee", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should create a new franchisee",
				Method:      http.MethodPost,
				URL:         "/franchisees",
				Body: map[string]interface{}{
					"name":        "Franchisee A",
					"description": "Main branch",
				},
				ExpectedCode: http.StatusCreated,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Franchisees", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a list of franchisees",
				Method:       http.MethodGet,
				URL:          "/franchisees",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Franchisee by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should return a franchisee by ID",
				Method:       http.MethodGet,
				URL:          "/franchisees/1",
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent franchisee",
				Method:       http.MethodGet,
				URL:          "/franchisees/9999",
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Franchisee", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Should update a franchisee",
				Method:      http.MethodPut,
				URL:         "/franchisees/1",
				Body: map[string]interface{}{
					"name": "Updated Franchisee A",
				},
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Franchisee", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Should delete a franchisee",
				Method:       http.MethodDelete,
				URL:          "/franchisees/1",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
