package stores_test

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestStoreEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	// --- Test Creating a New Cafe ---
	t.Run("Create new cafe", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Valid creation of a new cafe",
				Method:      http.MethodPost,
				URL:         "/api/test/stores",
				Body: map[string]interface{}{
					"name":         "Cafe Central",
					"franchiseeId": 1, // example franchisee ID
					"warehouseId":  1,
					"facilityAddress": map[string]interface{}{
						"address": "789 Cafe Road",
					},
					"isActive":     true,
					"contactPhone": "+77776667788",
					"contactEmail": "contact@cafecentral.com",
					"storeHours":   "08:00-20:00",
				},
				AuthRole:     data.RoleAdmin, // Use appropriate role for creation
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Create new cafe with empty name",
				Method:      http.MethodPost,
				URL:         "/api/test/stores",
				Body: map[string]interface{}{
					"name":         "", // Empty name
					"franchiseeId": 1,
					"warehouseId":  1,
					"facilityAddress": map[string]interface{}{
						"address": "789 Cafe Road",
					},
					"isActive":     true,
					"contactPhone": "+77776667788",
					"contactEmail": "contact@cafecentral.com",
					"storeHours":   "08:00-20:00",
				},
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Create new cafe with no contact number",
				Method:      http.MethodPost,
				URL:         "/api/test/stores",
				Body: map[string]interface{}{
					"name":         "Cafe No Number",
					"franchiseeId": 1,
					"warehouseId":  1,
					"facilityAddress": map[string]interface{}{
						"address": "789 Cafe Road",
					},
					"isActive": true,
					// "contactPhone" is omitted intentionally
					"contactEmail": "contact@cafenonumber.com",
					"storeHours":   "08:00-20:00",
				},
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Create new cafe with no facility address",
				Method:      http.MethodPost,
				URL:         "/api/test/stores",
				Body: map[string]interface{}{
					"name":         "Cafe No Address",
					"franchiseeId": 1,
					"warehouseId":  1,
					// "facilityAddress" is omitted intentionally
					"isActive":     true,
					"contactPhone": "+77776667788",
					"contactEmail": "contact@cafenostreet.com",
					"storeHours":   "08:00-20:00",
				},
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Create new cafe with no email",
				Method:      http.MethodPost,
				URL:         "/api/test/stores",
				Body: map[string]interface{}{
					"name":         "Cafe No Email",
					"franchiseeId": 1,
					"warehouseId":  1,
					"facilityAddress": map[string]interface{}{
						"address": "789 Cafe Road",
					},
					"isActive":     true,
					"contactPhone": "+77776667788",
					// "contactEmail" is omitted intentionally
					"storeHours": "08:00-20:00",
				},
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusBadRequest,
			},
		}
		env.RunTests(t, testCases)
	})

	// --- Test Updating an Existing Cafe ---
	t.Run("Update cafe", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Valid update of an existing cafe",
				Method:      http.MethodPut,
				URL:         "/api/test/stores/1",
				Body: map[string]interface{}{
					"name":         "Updated Cafe Central",
					"contactPhone": "+77776667788",
					"contactEmail": "updated@cafecentral.com",
					"facilityAddress": map[string]interface{}{
						"address": "101 Updated Road",
					},
					"storeHours": "09:00-21:00",
				},
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	// --- Test Deleting a Cafe Store ---
	t.Run("Delete cafe", func(t *testing.T) {
		testCases := []utils.TestCase{
			// {
			// 	Description:  "Valid deletion of a cafe store",
			// 	Method:       http.MethodDelete,
			// 	URL:          "/api/test/stores/1",
			// 	AuthRole:     data.RoleAdmin,
			// 	ExpectedCode: http.StatusOK,
			// },
			{
				Description:  "Unauthorized deletion attempt for a cafe store",
				Method:       http.MethodDelete,
				URL:          "/api/test/stores/1",
				AuthRole:     data.RoleBarista, // Assuming Barista role is not permitted to delete
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	// --- Test Fetching a Cafe Store by ID ---
	t.Run("Get store by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Fetch cafe store by valid ID",
				Method:       http.MethodGet,
				URL:          "/api/test/stores/1",
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Fetch cafe store with non-existent ID",
				Method:       http.MethodGet,
				URL:          "/api/test/stores/9999",
				AuthRole:     data.RoleFranchiseManager,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	// --- Test Fetching All Cafe Stores ---
	t.Run("Get all stores", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Fetch all cafe stores",
				Method:       http.MethodGet,
				URL:          "/api/test/stores/all",
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
