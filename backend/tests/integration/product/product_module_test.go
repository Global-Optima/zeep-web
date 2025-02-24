package product_test

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/tests/mockFiles"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestProductEndpoints(t *testing.T) {
	env := utils.NewTestEnvironment(t)
	defer env.Close()

	imageFileHeader := mockFiles.GetMockImageByFileName("test-image-valid.png")
	videoFileHeader := mockFiles.GetMockVideoByFileName("test-video-valid.mp4")
	imageFileHeaders := []*multipart.FileHeader{imageFileHeader}
	videoFileHeaders := []*multipart.FileHeader{videoFileHeader}

	//products tests
	t.Run("Create a Product", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Creating new product with an image",
				Method:      http.MethodPost,
				URL:         "/api/test/products",
				FormData: map[string]string{
					"name":        "Latte",
					"description": "A smooth coffee with milk",
					"categoryId":  "2",
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Creating new product with an image and a video",
				Method:      http.MethodPost,
				URL:         "/api/test/products",
				FormData: map[string]string{
					"name":        "Americano",
					"description": "A smooth coffee with milk",
					"categoryId":  "2",
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
					"video": videoFileHeaders,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusCreated,
			},
			{
				Description: "Barista should NOT be able to create a product",
				Method:      http.MethodPost,
				URL:         "/api/test/products",
				FormData: map[string]string{
					"name":        "Latte",
					"description": "A smooth coffee with milk",
					"categoryId":  "2",
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
				},
				AuthRole:     data.RoleBarista,
				ExpectedCode: http.StatusForbidden,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch products list", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch all products",
				Method:       http.MethodGet,
				URL:          "/api/test/products",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Store Manager should fetch all products",
				Method:       http.MethodGet,
				URL:          "/api/test/products",
				AuthRole:     data.RoleStoreManager,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Fetch a Product by ID", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should fetch a product by ID",
				Method:       http.MethodGet,
				URL:          "/api/test/products/2",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description:  "Should return 404 for non-existent product",
				Method:       http.MethodGet,
				URL:          "/api/test/products/9999",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusNotFound,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Update a Product", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description: "Update product with an image",
				Method:      http.MethodPut,
				URL:         "/api/test/products/1",
				FormData: map[string]string{
					"name":        "Updated Latte",
					"description": "A smooth coffee with milk",
					"categoryId":  "1",
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
			{
				Description: "Update product with an image and a video",
				Method:      http.MethodPut,
				URL:         "/api/test/products/1",
				FormData: map[string]string{
					"name":        "Updated Latte",
					"description": "A smooth coffee with milk",
					"categoryId":  "1",
				},
				Files: map[string][]*multipart.FileHeader{
					"image": imageFileHeaders,
					"video": videoFileHeaders,
				},
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})

	t.Run("Delete a Product", func(t *testing.T) {
		testCases := []utils.TestCase{
			{
				Description:  "Admin should delete a product",
				Method:       http.MethodDelete,
				URL:          "/api/test/products/1",
				AuthRole:     data.RoleAdmin,
				ExpectedCode: http.StatusOK,
			},
		}
		env.RunTests(t, testCases)
	})
}
