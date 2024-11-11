package product

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/types"
	"github.com/Global-Optima/zeep-web/backend/tests/unit/utils"
	"github.com/stretchr/testify/assert"
)

func repoTestGetStoreProducts(args ...interface{}) (interface{}, error) {
	repo := args[0].(product.ProductRepository)
	storeID := args[1].(uint)
	category := args[2].(string)
	offset := args[3].(int)
	limit := args[4].(int)

	return repo.GetStoreProducts(storeID, category, offset, limit)
}

func TestProductRepository_GetStoreProducts(t *testing.T) {
	repo, mock := utils.SetupProductRepository(t)

	columns := []string{"id", "name", "description", "image_url", "category_name", "base_price", "is_available", "is_out_of_stock"}
	queryRegex := regexp.MustCompile(`SELECT DISTINCT products\.id, products\.name, products\.description, products\.image_url, c\.name as category_name, COALESCE\(store_product_sizes\.price, product_sizes\.base_price\) as base_price, store_products\.is_available, \(CASE WHEN COALESCE\(store_warehouse_stocks\.quantity, 0\) = 0 THEN true ELSE false END\) as is_out_of_stock FROM "products"`)

	testCases := []utils.TestCase{
		{
			Name:      "Basic Retrieval",
			InputArgs: []interface{}{repo, uint(1), "Coffee", 0, 10},
			Expected: []types.ProductDAO{
				{ProductID: 1, ProductName: "Latte", ProductDescription: "Creamy coffee", CategoryName: "Coffee", BasePrice: 5.0, IsAvailable: true, IsOutOfStock: false},
				{ProductID: 2, ProductName: "Espresso", ProductDescription: "Strong coffee", CategoryName: "Coffee", BasePrice: 4.0, IsAvailable: true, IsOutOfStock: false},
			},
			ShouldFail: false,
			SetupMock: func() {
				mockRows := sqlmock.NewRows(columns)
				expectedData := []types.ProductDAO{
					{ProductID: 1, ProductName: "Latte", ProductDescription: "Creamy coffee", CategoryName: "Coffee", BasePrice: 5.0, IsAvailable: true, IsOutOfStock: false},
					{ProductID: 2, ProductName: "Espresso", ProductDescription: "Strong coffee", CategoryName: "Coffee", BasePrice: 4.0, IsAvailable: true, IsOutOfStock: false},
				}
				for _, e := range expectedData {
					mockRows.AddRow(
						e.ProductID,
						e.ProductName,
						e.ProductDescription,
						e.ProductImageURL,
						e.CategoryName,
						e.BasePrice,
						e.IsAvailable,
						e.IsOutOfStock,
					)
				}
				mock.ExpectQuery(queryRegex.String()).WillReturnRows(mockRows)
			},
		},
		{
			Name:       "No Results",
			InputArgs:  []interface{}{repo, uint(1), "Tea", 0, 10},
			Expected:   []types.ProductDAO(nil),
			ShouldFail: false,
			SetupMock: func() {
				mock.ExpectQuery(queryRegex.String()).WillReturnRows(sqlmock.NewRows(columns))
			},
		},
		{
			Name:       "Database Error",
			InputArgs:  []interface{}{repo, uint(1), "Coffee", 0, 10},
			Expected:   nil,
			ShouldFail: true,
			SetupMock: func() {
				mock.ExpectQuery(queryRegex.String()).WillReturnError(sqlmock.ErrCancelled)
			},
		},
	}

	utils.TestRunner(t, repoTestGetStoreProducts, testCases)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func repoTestSearchStoreProducts(args ...interface{}) (interface{}, error) {
	repo := args[0].(product.ProductRepository)
	storeID := args[1].(uint)
	searchQuery := args[2].(string)
	category := args[3].(string)
	offset := args[4].(int)
	limit := args[5].(int)

	return repo.SearchStoreProducts(storeID, searchQuery, category, offset, limit)
}

func TestProductRepository_SearchStoreProducts(t *testing.T) {
	repo, mock := utils.SetupProductRepository(t)
	columns := []string{"id", "name", "description", "image_url", "category_name", "base_price", "is_available", "is_out_of_stock"}
	queryRegex := regexp.MustCompile(`SELECT DISTINCT products\.id, products\.name, products\.description, products\.image_url, c\.name as category_name, COALESCE\(store_product_sizes\.price, product_sizes\.base_price\) as base_price, store_products\.is_available, \(CASE WHEN COALESCE\(store_warehouse_stocks\.quantity, 0\) = 0 THEN true ELSE false END\) as is_out_of_stock FROM "products"`)

	testCases := []utils.TestCase{
		{
			Name:      "Search by Name",
			InputArgs: []interface{}{repo, uint(1), "Latte", "Coffee", 0, 10},
			Expected: []types.ProductDAO{
				{ProductID: 1, ProductName: "Latte", ProductDescription: "Creamy coffee", CategoryName: "Coffee", ProductImageURL: "https://example.com/latte.jpg", BasePrice: 5.0, IsAvailable: true, IsOutOfStock: false},
			},
			ShouldFail: false,
			SetupMock: func() {
				mockRows := sqlmock.NewRows(columns).AddRow(1, "Latte", "Creamy coffee", "https://example.com/latte.jpg", "Coffee", 5.0, true, false)
				mock.ExpectQuery(queryRegex.String()).WillReturnRows(mockRows)
			},
		},
		{
			Name:       "No Matching Results",
			InputArgs:  []interface{}{repo, uint(1), "Mocha", "Coffee", 0, 10},
			Expected:   []types.ProductDAO(nil),
			ShouldFail: false,
			SetupMock: func() {
				mock.ExpectQuery(queryRegex.String()).WillReturnRows(sqlmock.NewRows(columns)) // Empty rows
			},
		},
		{
			Name:       "Database Error",
			InputArgs:  []interface{}{repo, uint(1), "Latte", "Coffee", 0, 10},
			Expected:   nil,
			ShouldFail: true,
			SetupMock: func() {
				mock.ExpectQuery(queryRegex.String()).WillReturnError(sqlmock.ErrCancelled)
			},
		},
	}

	utils.TestRunner(t, repoTestSearchStoreProducts, testCases)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func repoTestGetStoreProductDetails(args ...interface{}) (interface{}, error) {
	repo := args[0].(product.ProductRepository)
	storeID := args[1].(uint)
	productID := args[2].(uint)

	return repo.GetStoreProductDetails(storeID, productID)
}

func TestProductRepository_GetStoreProductDetails(t *testing.T) {
	repo, mock := utils.SetupProductRepository(t)

	// Main product details query regex
	mainQueryRegex := regexp.QuoteMeta(`SELECT products.id, products.name, products.description, products.image_url, products.video_url, c.id as category_id, c.name as category_name, COALESCE(store_product_sizes.price, product_sizes.base_price) as base_price, store_products.is_available, (CASE WHEN COALESCE(store_warehouse_stocks.quantity, 0) = 0 THEN true ELSE false END) as is_out_of_stock FROM "products" JOIN store_products ON store_products.product_id = products.id AND store_products.store_id = $1 LEFT JOIN categories c ON c.id = products.category_id LEFT JOIN product_sizes ON product_sizes.product_id = products.id AND product_sizes.is_default = true LEFT JOIN store_product_sizes ON store_product_sizes.product_size_id = product_sizes.id AND store_product_sizes.store_id = $2 LEFT JOIN store_warehouse_stocks ON store_warehouse_stocks.ingredient_id = products.id AND store_warehouse_stocks.store_warehouse_id = store_products.store_id WHERE products.id = $3 ORDER BY "products"."id" LIMIT $4`)

	// Mocked rows for product details
	mainMockRows := sqlmock.NewRows([]string{"id", "name", "description", "image_url", "video_url", "category_id", "category_name", "base_price", "is_available", "is_out_of_stock"}).
		AddRow(1, "Latte", "A creamy coffee drink", "https://example.com/latte.jpg", "https://example.com/latte.mp4", 1, "Coffee", 5.0, true, false)

	// Test case
	testCases := []utils.TestCase{
		{
			Name:      "Get Product Details Success",
			InputArgs: []interface{}{repo, uint(1), uint(1)},
			Expected: &types.ProductDAO{
				ProductID:          1,
				ProductName:        "Latte",
				ProductDescription: "A creamy coffee drink",
				CategoryID:         1,
				CategoryName:       "Coffee",
				ProductImageURL:    "https://example.com/latte.jpg",
				ProductVideoURL:    "https://example.com/latte.mp4",
				BasePrice:          5.0,
				IsAvailable:        true,
				IsOutOfStock:       false,
				Sizes: []types.SizeDAO{
					{SizeID: 1, SizeName: "Small", Size: 200, Measure: "ml", Price: 5.0, IsDefault: true},
				},
				Additives: []types.AdditiveDAO{
					{AdditiveID: 1, AdditiveName: "Vanilla Syrup", AdditiveDescription: "Sweet vanilla flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.5},
				},
				DefaultAdditives: []types.AdditiveDAO{
					{AdditiveID: 2, AdditiveName: "Caramel Syrup", AdditiveDescription: "Rich caramel flavor", AdditiveCategory: "Syrups", AdditivePrice: 1.75},
				},
				Nutrition: types.NutritionDAO{
					Calories:      100,
					Fat:           5,
					Carbohydrates: 15,
					Proteins:      3,
				},
			},
			ShouldFail: false,
			SetupMock: func() {
				// Set up main product query expectation with four arguments
				mock.ExpectQuery(mainQueryRegex).WithArgs(uint(1), uint(1), uint(1), 1).WillReturnRows(mainMockRows)

				// Mock for fetching default size ID for the product in loadAdditives
				defaultSizeIDMockRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM "product_sizes" WHERE product_id = $1 AND is_default = true`)).
					WithArgs(uint(1)).WillReturnRows(defaultSizeIDMockRows)

				// Sizes query in loadSizes
				sizeMockRows := sqlmock.NewRows([]string{"id", "name", "size", "measure", "price", "is_default"}).
					AddRow(1, "Small", 200, "ml", 5.0, true)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT product_sizes.id, product_sizes.name, product_sizes.size, product_sizes.measure, COALESCE(store_product_sizes.price, product_sizes.base_price) as price, product_sizes.is_default FROM "product_sizes" LEFT JOIN store_product_sizes ON store_product_sizes.product_size_id = product_sizes.id AND store_product_sizes.store_id = $1 WHERE product_sizes.product_id = $2`)).
					WithArgs(uint(1), uint(1)).WillReturnRows(sizeMockRows)

				// Additives query in loadAdditives
				additiveMockRows := sqlmock.NewRows([]string{"id", "name", "description", "category_name", "price"}).
					AddRow(1, "Vanilla Syrup", "Sweet vanilla flavor", "Syrups", 1.5)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT DISTINCT additives.id, additives.name, additives.description, additive_categories.name as category_name, COALESCE(store_additives.price, additives.base_price) as price FROM "additives" JOIN product_additives ON product_additives.additive_id = additives.id AND product_additives.product_size_id = $1 LEFT JOIN store_additives ON store_additives.additive_id = additives.id AND store_additives.store_id = $2 LEFT JOIN additive_categories ON additive_categories.id = additives.additive_category_id LEFT JOIN store_warehouse_stocks ON store_warehouse_stocks.ingredient_id = additives.id AND store_warehouse_stocks.store_warehouse_id = $3 WHERE store_additives.additive_id IS NOT NULL OR store_warehouse_stocks.quantity > 0`)).
					WithArgs(1, uint(1), uint(1)).WillReturnRows(additiveMockRows)

				// Default additives query in loadDefaultAdditives
				defaultAdditiveMockRows := sqlmock.NewRows([]string{"id", "name", "description", "category_name", "price"}).
					AddRow(2, "Caramel Syrup", "Rich caramel flavor", "Syrups", 1.75)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT additives.id, additives.name, additives.description, additive_categories.name as category_name, additives.base_price as price FROM "additives" JOIN default_product_additives ON default_product_additives.additive_id = additives.id LEFT JOIN additive_categories ON additive_categories.id = additives.additive_category_id WHERE default_product_additives.product_id = $1`)).
					WithArgs(uint(1)).WillReturnRows(defaultAdditiveMockRows)

				// Nutrition query in loadNutrition
				nutritionMockRows := sqlmock.NewRows([]string{"calories", "fat", "carbohydrates", "proteins"}).
					AddRow(100, 5, 15, 3)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(SUM(ingredients.calories), 0) as calories, COALESCE(SUM(ingredients.fat), 0) as fat, COALESCE(SUM(ingredients.carbs), 0) as carbohydrates, COALESCE(SUM(ingredients.proteins), 0) as proteins FROM "product_ingredients" JOIN item_ingredients ON product_ingredients.item_ingredient_id = item_ingredients.id JOIN ingredients ON item_ingredients.ingredient_id = ingredients.id WHERE product_ingredients.product_id = $1 GROUP BY product_ingredients.product_id`)).
					WithArgs(uint(1)).WillReturnRows(nutritionMockRows)
			},
		},
		{
			Name:       "Product Not Found",
			InputArgs:  []interface{}{repo, uint(1), uint(5)},
			Expected:   (*types.ProductDAO)(nil),
			ShouldFail: true,
			SetupMock: func() {
				mock.ExpectQuery(mainQueryRegex).
					WithArgs(uint(1), uint(1), uint(5), 1).
					WillReturnRows(sqlmock.NewRows(nil))
			},
		},
	}

	utils.TestRunner(t, repoTestGetStoreProductDetails, testCases)
	assert.NoError(t, mock.ExpectationsWereMet())
}
