package franchisees_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestFranchiseeEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create a Franchisee", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new franchisee",
				Method:      http.MethodPost,
				URL:         "/api/test/franchisees",
				Body: map[string]interface{}{
					"name":        "Franchisee C",
					"description": "Main branch",
				},
				AuthRole:     data.RoleAdmin, // Only Admin can create
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a franchisee",
				Method:      http.MethodPost,
				URL:         "/api/test/franchisees",
				Body: map[string]interface{}{
					"name":        "Franchisee A",
					"description": "Main branch",
				},
				AuthRole:     data.RoleBarista, // Unauthorized role
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description: "Admin should NOT create a franchisee with an empty name",
				Method:      http.MethodPost,
				URL:         "/api/test/franchisees",
				Body: map[string]interface{}{
					"name":        " ", // Empty name
					"description": "Branch with empty name",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest, // Expected failure due to empty name
			},
			{
				Description: "Admin should NOT create a franchisee with a duplicate name",
				Method:      http.MethodPost,
				URL:         "/api/test/franchisees",
				Body: map[string]interface{}{
					"name":        "Franchisee A", // Duplicate of the first test case
					"description": "Duplicate branch",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError, // Expected failure due to duplicate name
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch all Franchisees", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all franchisees",
				Method:       http.MethodGet,
				URL:          "/api/test/franchisees",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Franchise Manager NOT(why?) should fetch all franchisees",
				Method:       http.MethodGet,
				URL:          "/api/test/franchisees",
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description:  "Barista should NOT be able to fetch franchisees",
				Method:       http.MethodGet,
				URL:          "/api/test/franchisees",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Franchisee by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a franchisee by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/franchisees/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Franchise Manager should NOT(why?) fetch a franchisee by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/franchisees/1",
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusForbidden,
			},
			{
				Description:  "Unauthorized user should NOT fetch a franchisee",
				Method:       http.MethodGet,
				URL:          "/api/test/franchisees/1",
				ExpectedCode: http.StatusUnauthorized,
			},
			{
				Description:  "Should return 500 for non-existent franchisee",
				Method:       http.MethodGet,
				URL:          "/api/test/franchisees/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusInternalServerError,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Franchisee", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update a franchisee",
				Method:      http.MethodPut,
				URL:         "/api/test/franchisees/1",
				Body: map[string]interface{}{
					"name": "Updated Franchisee A",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Franchise Manager should be able to update a franchisee",
				Method:      http.MethodPut,
				URL:         "/api/test/franchisees/1",
				Body: map[string]interface{}{
					"name": "Updated Franchisee A",
				},
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Admin should NOT update a franchisee with an empty name",
				Method:      http.MethodPut,
				URL:         "/api/test/franchisees/1",
				Body: map[string]interface{}{
					"name":        " ",                 // Empty name
					"description": "Valid description", // Valid description provided
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest, // Expected failure due to empty name
			},
			{
				Description: "Admin should update a franchisee with an empty description",
				Method:      http.MethodPut,
				URL:         "/api/test/franchisees/1",
				Body: map[string]interface{}{
					"name":        "Valid Franchisee Name", // Valid name provided
					"description": " ",                     // Empty description
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK, // Expected failure due to empty description
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Franchisee", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should delete a franchisee",
				Method:       http.MethodDelete,
				URL:          "/api/test/franchisees/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete a franchisee",
				Method:       http.MethodDelete,
				URL:          "/api/test/franchisees/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
