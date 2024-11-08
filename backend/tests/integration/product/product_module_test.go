package product

import (
	"net/http"
	"testing"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/tests/integration/utils"
)

func TestGetStoreProducts(t *testing.T) {
	handler, db := utils.SetupProductHandler(t)
	router := utils.SetupTestRouter(handler)
	TruncateTables(db)
	SetupMockData(db)

	tests := []utils.TestCase{
		{
			Description:  "Valid Request with Category Filter",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products?category=Coffee",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 1, ProductName: "Latte", ProductDescription: "A creamy coffee drink", Category: "Coffee", Price: 5.0, IsAvailable: true, IsOutOfStock: false},
				{ProductID: 2, ProductName: "Espresso", ProductDescription: "Strong coffee", Category: "Coffee", Price: 4.0, IsAvailable: true, IsOutOfStock: false},
			},
		},
		{
			Description:  "Valid Request without Category Filter",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 1, ProductName: "Latte", ProductDescription: "A creamy coffee drink", Category: "Coffee", Price: 5.0, IsAvailable: true, IsOutOfStock: false},
				{ProductID: 2, ProductName: "Espresso", ProductDescription: "Strong coffee", Category: "Coffee", Price: 4.0, IsAvailable: true, IsOutOfStock: false},
				{ProductID: 4, ProductName: "Green Tea", ProductDescription: "Fresh green tea", Category: "Tea", Price: 0, IsAvailable: false, IsOutOfStock: true},
			},
		},
		{
			Description:  "Pagination",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products?offset=1&limit=1",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 2, ProductName: "Espresso", ProductDescription: "Strong coffee", Category: "Coffee", Price: 4.0, IsAvailable: true, IsOutOfStock: false},
			},
		},
		{
			Description:  "Pagination at Offset with No Results",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products?offset=10&limit=5",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{}, // No products should be returned as offset is beyond total products
		},
		{
			Description:  "No Products in Category",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products?category=Nonexistent",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{},
		},
		{
			Description:  "Store with No Products",
			Method:       http.MethodGet,
			URL:          "/api/stores/99/products",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{},
		},
		{
			Description:  "Unavailable Products Only",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products?category=Tea",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 4, ProductName: "Green Tea", ProductDescription: "Fresh green tea", Category: "Tea", Price: 0, IsAvailable: false, IsOutOfStock: true},
			},
		},
		{
			Description:  "Available and Out-of-Stock Products",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 1, ProductName: "Latte", ProductDescription: "A creamy coffee drink", Category: "Coffee", Price: 5.0, IsAvailable: true, IsOutOfStock: false},
				{ProductID: 2, ProductName: "Espresso", ProductDescription: "Strong coffee", Category: "Coffee", Price: 4.0, IsAvailable: true, IsOutOfStock: false},
				{ProductID: 4, ProductName: "Green Tea", ProductDescription: "Fresh green tea", Category: "Tea", Price: 0, IsAvailable: false, IsOutOfStock: true},
			},
		},
		{
			Description:  "Invalid Store ID",
			Method:       http.MethodGet,
			URL:          "/api/stores/abc/products",
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Description:  "Invalid Pagination Parameters",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products?offset=-1&limit=-10",
			ExpectedCode: http.StatusBadRequest,
		},
	}

	utils.TestRunner(t, router, tests)
}

func TestSearchStoreProducts(t *testing.T) {
	handler, db := utils.SetupProductHandler(t)
	router := utils.SetupTestRouter(handler)
	TruncateTables(db)
	SetupMockData(db)

	tests := []utils.TestCase{
		{
			Description:  "Valid Search Query with Category Filter",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=Latte&category=Coffee",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 1, ProductName: "Latte", ProductDescription: "A creamy coffee drink", Category: "Coffee", Price: 5.0, IsAvailable: true, IsOutOfStock: false},
			},
		},
		{
			Description:  "Valid Search Query without Category Filter",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=Espresso",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 2, ProductName: "Espresso", ProductDescription: "Strong coffee", Category: "Coffee", Price: 4.0, IsAvailable: true, IsOutOfStock: false},
			},
		},
		{
			Description:  "Search Query with No Results",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=Mocha&category=Tea",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{}, // No matching products
		},
		{
			Description:  "Store with No Matching Products",
			Method:       http.MethodGet,
			URL:          "/api/stores/99/products/search?q=Latte",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{}, // Store 99 does not exist or has no products
		},
		{
			Description:  "Empty Search Query",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=",
			ExpectedCode: http.StatusBadRequest, // Search query cannot be empty
		},
		{
			Description:  "Invalid Store ID",
			Method:       http.MethodGet,
			URL:          "/api/stores/abc/products/search?q=Latte",
			ExpectedCode: http.StatusBadRequest, // Invalid store ID format
		},
		{
			Description:  "Search Query with Multiple Matches",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=Tea",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 3, ProductName: "Black Tea", ProductDescription: "Classic black tea", Category: "Tea", Price: 0, IsAvailable: true, IsOutOfStock: false},
				{ProductID: 4, ProductName: "Green Tea", ProductDescription: "Fresh green tea", Category: "Tea", Price: 0, IsAvailable: false, IsOutOfStock: true},
			},
		},
		{
			Description:  "Pagination on Search Results",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=Tea&offset=1&limit=1",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 4, ProductName: "Green Tea", ProductDescription: "Fresh green tea", Category: "Tea", Price: 0, IsAvailable: false, IsOutOfStock: true},
			},
		},
		{
			Description:  "Invalid Pagination Parameters",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=Latte&offset=-1&limit=-10",
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Description:  "Search for Out of Stock Product",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/search?q=Green",
			ExpectedCode: http.StatusOK,
			ExpectedBody: []types.ProductCatalogDTO{
				{ProductID: 4, ProductName: "Green Tea", ProductDescription: "Fresh green tea", Category: "Tea", Price: 0, IsAvailable: false, IsOutOfStock: true},
			},
		},
	}

	utils.TestRunner(t, router, tests)
}

func TestGetStoreProductDetails(t *testing.T) {
	handler, db := utils.SetupProductHandler(t)
	router := utils.SetupTestRouter(handler)
	TruncateTables(db)
	SetupMockData(db)

	tests := []utils.TestCase{
		{
			Description:  "Valid Product Details",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/1",
			ExpectedCode: http.StatusOK,
			ExpectedBody: types.ProductDTO{
				ProductID:          1,
				ProductName:        "Latte",
				ProductDescription: "A creamy coffee drink",
				Category:           "Coffee",
				ProductImageURL:    "https://example.com/latte.jpg",
				ProductVideoURL:    "https://example.com/latte.mp4",
				Price:              5.0,
				IsAvailable:        true,
				IsOutOfStock:       false,
				Sizes: []types.SizeDTO{
					{SizeID: 1, SizeName: "Small", Size: 200, Measure: "ml", Price: 5.0, IsDefault: true},
					{SizeID: 2, SizeName: "Large", Size: 300, Measure: "ml", Price: 6.5, IsDefault: false},
				},
				Additives: []types.AdditivesDTO{
					{AdditiveID: 1, AdditiveName: "Vanilla Syrup", AdditiveDescription: "Sweet vanilla flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.5},
					{AdditiveID: 2, AdditiveName: "Caramel Syrup", AdditiveDescription: "Rich caramel flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.75},
				},
				DefaultAdditives: []types.AdditivesDTO{
					{AdditiveID: 1, AdditiveName: "Vanilla Syrup", AdditiveDescription: "Sweet vanilla flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.5},
					{AdditiveID: 2, AdditiveName: "Caramel Syrup", AdditiveDescription: "Rich caramel flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.75},
				},
				Nutrition: types.NutritionDTO{
					Calories:      100,
					Fat:           5,
					Carbohydrates: 15,
					Proteins:      3,
				},
			},
		},
		{
			Description:  "Product Exists but is Unavailable",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/4",
			ExpectedCode: http.StatusOK,
			ExpectedBody: types.ProductDTO{
				ProductID:          4,
				ProductName:        "Green Tea",
				ProductDescription: "Fresh green tea",
				Category:           "Tea",
				ProductImageURL:    "https://example.com/greentea.jpg",
				Price:              0,
				IsAvailable:        false,
				IsOutOfStock:       true,
				Sizes:              []types.SizeDTO{},
				Additives:          []types.AdditivesDTO{},
				DefaultAdditives:   []types.AdditivesDTO{},
				Nutrition:          types.NutritionDTO{},
			},
		},
		{
			Description:  "Product Exists but is Out of Stock",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/3",
			ExpectedCode: http.StatusOK,
			ExpectedBody: types.ProductDTO{
				ProductID:          3,
				ProductName:        "Black Tea",
				ProductDescription: "Classic black tea",
				Category:           "Tea",
				ProductImageURL:    "https://example.com/blacktea.jpg",
				Price:              0,
				IsAvailable:        true,
				IsOutOfStock:       true,
				Sizes:              []types.SizeDTO{},
				Additives:          []types.AdditivesDTO{},
				DefaultAdditives:   []types.AdditivesDTO{},
				Nutrition:          types.NutritionDTO{},
			},
		},
		{
			Description:  "Product with Multiple Sizes and Additives",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/1",
			ExpectedCode: http.StatusOK,
			ExpectedBody: types.ProductDTO{
				ProductID:          1,
				ProductName:        "Latte",
				ProductDescription: "A creamy coffee drink",
				Category:           "Coffee",
				ProductImageURL:    "https://example.com/latte.jpg",
				ProductVideoURL:    "https://example.com/latte.mp4",
				Price:              5.0,
				IsAvailable:        true,
				IsOutOfStock:       false,
				Sizes: []types.SizeDTO{
					{SizeID: 1, SizeName: "Small", Size: 200, Measure: "ml", Price: 5.0, IsDefault: true},
					{SizeID: 2, SizeName: "Large", Size: 300, Measure: "ml", Price: 6.5, IsDefault: false},
				},
				Additives: []types.AdditivesDTO{
					{AdditiveID: 1, AdditiveName: "Vanilla Syrup", AdditiveDescription: "Sweet vanilla flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.5},
					{AdditiveID: 2, AdditiveName: "Caramel Syrup", AdditiveDescription: "Rich caramel flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.75},
				},
				DefaultAdditives: []types.AdditivesDTO{
					{AdditiveID: 1, AdditiveName: "Vanilla Syrup", AdditiveDescription: "Sweet vanilla flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.5},
					{AdditiveID: 2, AdditiveName: "Caramel Syrup", AdditiveDescription: "Rich caramel flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.75},
				},
				Nutrition: types.NutritionDTO{
					Calories:      100,
					Fat:           5,
					Carbohydrates: 15,
					Proteins:      3,
				},
			},
		},
		{
			Description:  "Invalid Product ID",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/999",
			ExpectedCode: http.StatusNotFound,
		},
		{
			Description:  "Invalid Store ID",
			Method:       http.MethodGet,
			URL:          "/api/stores/99/products/1",
			ExpectedCode: http.StatusNotFound,
		},
		{
			Description:  "Product with Missing Nutrition Information",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/5",
			ExpectedCode: http.StatusOK,
			ExpectedBody: types.ProductDTO{
				ProductID:          5,
				ProductName:        "Mocha",
				ProductDescription: "Chocolate-flavored coffee",
				Category:           "Coffee",
				ProductImageURL:    "https://example.com/mocha.jpg",
				Price:              0,
				IsAvailable:        true,
				IsOutOfStock:       false,
				Sizes:              []types.SizeDTO{},
				Additives:          []types.AdditivesDTO{},
				DefaultAdditives:   []types.AdditivesDTO{},
				Nutrition:          types.NutritionDTO{},
			},
		},
		{
			Description:  "Database Error During Query Execution",
			Method:       http.MethodGet,
			URL:          "/api/stores/1/products/1",
			ExpectedCode: http.StatusInternalServerError,
		},
	}

	utils.TestRunner(t, router, tests)
}
