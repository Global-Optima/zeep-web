package additives_test

import (
	"github.com/Global-Optima/zeep-web/backend/tests/mockFiles"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestAdditiveEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	imageFileHeader := mockFiles.GetMockImageByFileName("test-image-valid.png")
	imageFileHeaders := []*multipart.FileHeader{imageFileHeader}

	t.Run("Create an Additive", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should create a new additive WITH an image",
				Method:      http.MethodPost,
				URL:         "/api/test/additives",
				FormData: map[string]string{
					"name":               "New Test Additive 1",
					"description":        "Sweet vanilla flavor",
					"basePrice":          "3.99",
					"size":               "250",
					"unitId":             "1",
					"additiveCategoryId": "1",
					"ingredients":        `[{"ingredientId":2,"quantity":5.0}]`,
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Additive should not be created WITHOUT an image",
				Method:      http.MethodPost,
				URL:         "/api/test/additives",
				FormData: map[string]string{
					"name":               "New Test Additive 2",
					"description":        "Rich caramel flavor",
					"basePrice":          "4.49",
					"size":               "300",
					"unitId":             "1",
					"additiveCategoryId": "1",
					"ingredients":        `[{"ingredientId":3,"quantity":4.0}]`,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusBadRequest,
			},
			{
				Description: "Barista should NOT be able to create an additive",
				Method:      http.MethodPost,
				URL:         "/api/test/additives",
				FormData: map[string]string{
					"name":               "New Test Additive 3",
					"description":        "Sweet vanilla flavor",
					"basePrice":          "3.99",
					"size":               "250",
					"unitId":             "1",
					"additiveCategoryId": "1",
					"ingredients":        `[{"ingredientId":2,"quantity":5.0}]`,
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders, // Attach image
				},
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch All Additives", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all additives",
				Method:       http.MethodGet,
				URL:          "/api/test/additives",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch all additives",
				Method:       http.MethodGet,
				URL:          "/api/test/additives",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should fetch all additives",
				Method:       http.MethodGet,
				URL:          "/api/test/additives",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch an Additive by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch an additive by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/additives/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch an additive by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/additives/1",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should fetch an additive by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/additives/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent additive",
				Method:       http.MethodGet,
				URL:          "/api/test/additives/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update an Additive", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Admin should update an additive WITH an image",
				Method:      http.MethodPut,
				URL:         "/api/test/additives/1",
				FormData: map[string]string{
					"name":               "Vanilla Syrup - Updated",
					"description":        "Updated description",
					"basePrice":          "4.50",
					"size":               "350",
					"unitId":             "1",
					"additiveCategoryId": "1",
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Admin should update an additive WITHOUT an image",
				Method:      http.MethodPut,
				URL:         "/api/test/additives/1",
				FormData: map[string]string{
					"name":               "Caramel Syrup - Updated",
					"description":        "Updated description",
					"basePrice":          "5.00",
					"size":               "400",
					"unitId":             "1",
					"additiveCategoryId": "1",
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Store Manager should NOT be able to update an additive",
				Method:      http.MethodPut,
				URL:         "/api/test/additives/1",
				FormData: map[string]string{
					"name":               "Vanilla Syrup - Updated",
					"description":        "Updated description",
					"basePrice":          "4.50",
					"size":               "350",
					"unitId":             "1",
					"additiveCategoryId": "1",
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
				},
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete an Additive", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should delete an additive",
				Method:       http.MethodDelete,
				URL:          "/api/test/additives/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Barista should NOT be able to delete an additive",
				Method:       http.MethodDelete,
				URL:          "/api/test/additives/1",
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})
}
