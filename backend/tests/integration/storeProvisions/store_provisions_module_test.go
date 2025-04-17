package storeProvisions_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestStoreProvisionEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	t.Run("Create Store Provision - "+data.RoleStoreManager.ToString(), func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Create a valid store provision",
				Method:      http.MethodPost,
				URL:         "/api/test/store-provisions",
				Body: map[string]interface{}{
					"provisionId":         1,
					"volume":              1.0,
					"expirationInMinutes": 120,
				},
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Fail with missing provision",
				Method:      http.MethodPost,
				URL:         "/api/test/store-provisions",
				Body: map[string]interface{}{
					"provisionId":         9999,
					"volume":              1.0,
					"expirationInMinutes": 120,
				},
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusNotFound,
			},
			{
				Description: "Fail with 0 volume",
				Method:      http.MethodPost,
				URL:         "/api/test/store-provisions",
				Body: map[string]interface{}{
					"provisionId":         1,
					"volume":              0,
					"expirationInMinutes": 120,
				},
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusBadRequest,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Get Store Provisions - "+data.RoleStoreManager.ToString(), func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Fetch all store provisions",
				Method:       http.MethodGet,
				URL:          "/api/test/store-provisions",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Filter provisions by search keyword",
				Method:       http.MethodGet,
				URL:          "/api/test/store-provisions?search=Test", // assumes provision name has "Test"
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Filter provisions with minCompletedAt",
				Method:       http.MethodGet,
				URL:          "/api/test/store-provisions?minCompletedAt=2024-01-01T00:00:00Z",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Filter provisions with maxCompletedAt",
				Method:       http.MethodGet,
				URL:          "/api/test/store-provisions?maxCompletedAt=2030-01-01T00:00:00Z",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Filter provisions within a completedAt range",
				Method:       http.MethodGet,
				URL:          "/api/test/store-provisions?minCompletedAt=2024-01-01T00:00:00Z&maxCompletedAt=2030-01-01T00:00:00Z",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Get Store Provision by ID - "+data.RoleStoreManager.ToString(), func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Fetch a store provision by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/store-provisions/1",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Return 404 for non-existent provision",
				Method:       http.MethodGet,
				URL:          "/api/test/store-provisions/9999",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update Store Provision - "+data.RoleStoreManager.ToString(), func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Update volume and expiration",
				Method:      http.MethodPut,
				URL:         "/api/test/store-provisions/1",
				Body: map[string]interface{}{
					"volume":              1.5,
					"expirationInMinutes": 180,
				},
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Complete Store Provision - "+data.RoleStoreManager.ToString(), func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Complete a store provision",
				Method:       http.MethodPost,
				URL:          "/api/test/store-provisions/1/complete",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Fail completing already completed provision",
				Method:       http.MethodPost,
				URL:          "/api/test/store-provisions/4/complete",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusConflict,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete Store Provision - "+data.RoleStoreManager.ToString(), func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Delete a store provision",
				Method:       http.MethodDelete,
				URL:          "/api/test/store-provisions/1",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Fail deleting completed provision",
				Method:       http.MethodDelete,
				URL:          "/api/test/store-provisions/4",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusConflict,
			},
		}
		env.RunTests(t, testCases)
	})
}
