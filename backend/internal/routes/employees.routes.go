package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock"
)

func (r *Router) RegisterProductRoutes(handler *product.ProductHandler) {
	router := r.EmployeeRoutes.Group("/products")
	{
		router.GET("", handler.GetProducts)
		router.GET("/:id", handler.GetProductDetails)
		router.POST("", handler.CreateProduct)
		router.PUT("/:id", handler.UpdateProduct)
		router.DELETE("/:id", handler.DeleteProduct)
		router.GET(":id/sizes", handler.GetProductSizesByProductID)
		router.POST("/sizes", handler.CreateProductSize)
		router.PUT("/sizes/:id", handler.UpdateProductSize)
		router.GET("/sizes/:id", handler.GetProductSizeByID)
	}
}

func (r *Router) RegisterRecipeRoutes(handler *recipes.RecipeHandler) {
	router := r.EmployeeRoutes.Group("/products/recipe-steps")
	{
		router.GET("/product/:product-id", handler.GetRecipeSteps)
		router.GET("/step/:id", handler.GetRecipeStepDetails)
		router.POST("/product/:product-id", handler.CreateRecipeSteps)
		router.PUT("/step/:id", handler.UpdateRecipeSteps)
		router.DELETE("/step/:id", handler.DeleteRecipeSteps)
	}
}

func (r *Router) RegisterStoreProductRoutes(handler *storeProducts.StoreProductHandler) {
	router := r.EmployeeRoutes.Group("/store-products")
	{
		router.GET("", handler.GetStoreProducts)
		router.GET("/:id", handler.GetStoreProduct)
		router.POST("", handler.CreateStoreProduct)
		router.POST("/multiple", handler.CreateMultipleStoreProducts)
		router.PUT("/:id", handler.UpdateStoreProduct)
		router.DELETE("/:id", handler.DeleteStoreProduct)

		router.GET("/sizes/:id", handler.GetStoreProductSizeByID)
	}
}

func (r *Router) RegisterIngredientRoutes(handler *ingredients.IngredientHandler) {
	router := r.EmployeeRoutes.Group("/ingredients")
	{
		router.POST("", handler.CreateIngredient)
		router.PUT("/:id", handler.UpdateIngredient)
		router.DELETE("/:id", handler.DeleteIngredient)
		router.GET("/:id", handler.GetIngredientByID)
		router.GET("", handler.GetIngredients)
	}
}

func (r *Router) RegisterIngredientCategoriesRoutes(handler *ingredientCategories.IngredientCategoryHandler) {
	router := r.EmployeeRoutes.Group("/ingredient-categories")
	{
		router.GET("", handler.GetAll)
		router.GET("/:id", handler.GetByID)
		router.POST("", handler.Create)
		router.PUT("/:id", handler.Update)
		router.DELETE("/:id", handler.Delete)

	}
}

func (r *Router) RegisterStoresRoutes(handler *stores.StoreHandler) {
	router := r.EmployeeRoutes.Group("/stores")
	{
		router.GET("/:id", handler.GetStoreByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.CreateStore)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.UpdateStore)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.DeleteStore)
	}
}

func (r *Router) RegisterProductCategoriesRoutes(handler *categories.CategoryHandler) {
	router := r.EmployeeRoutes.Group("/product-categories")
	{
		router.GET("", handler.GetAllCategories)
		router.GET("/:id", handler.GetCategoryByID)
		router.POST("", handler.CreateCategory)
		router.PUT("/:id", handler.UpdateCategory)
		router.DELETE("/:id", handler.DeleteCategory)
	}
}

func (r *Router) RegisterAdditivesRoutes(handler *additives.AdditiveHandler) {
	router := r.EmployeeRoutes.Group("/additives")
	{
		router.GET("", handler.GetAdditives)
		router.POST("", handler.CreateAdditive)
		router.PUT("/:id", handler.UpdateAdditive)
		router.DELETE("/:id", handler.DeleteAdditive)
		router.GET("/:id", handler.GetAdditiveByID)

		additiveCategories := router.Group("/categories")
		{
			additiveCategories.GET("", handler.GetAdditiveCategories)
			additiveCategories.POST("", handler.CreateAdditiveCategory)
			additiveCategories.PUT("/:id", handler.UpdateAdditiveCategory)
			additiveCategories.DELETE("/:id", handler.DeleteAdditiveCategory)
			additiveCategories.GET("/:id", handler.GetAdditiveCategoryByID)
		}
	}
}

func (r *Router) RegisterStoreAdditivesRoutes(handler *storeAdditives.StoreAdditiveHandler) {
	router := r.EmployeeRoutes.Group("/store-additives")
	{
		router.GET("", handler.GetStoreAdditives)
		router.GET("/categories/:productSizeId", handler.GetStoreAdditiveCategories)
		router.POST("", handler.CreateStoreAdditives)
		router.PUT("/:id", handler.UpdateStoreAdditive)
		router.DELETE("/:id", handler.DeleteStoreAdditive)
		router.GET("/:id", handler.GetStoreAdditiveByID)
	}
}

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.EmployeeRoutes.Group("/employees")
	{
		storeEmployees := router.Group("/stores")
		{
			storeEmployees.GET("", handler.GetStoreEmployees)
			storeEmployees.POST("", handler.CreateStoreEmployee)
			storeEmployees.GET("/:id", handler.GetStoreEmployeeByID)
			storeEmployees.PUT("/:id", handler.UpdateStoreEmployee)
		}
		warehouseEmployees := router.Group("/warehouses")
		{
			warehouseEmployees.GET("", handler.GetWarehouseEmployees)
			warehouseEmployees.POST("", handler.CreateWarehouseEmployee)
			warehouseEmployees.GET("/:id", handler.GetWarehouseEmployeeByID)
			warehouseEmployees.PUT("/:id", handler.UpdateWarehouseEmployee)
		}

		router.GET("/current", handler.GetCurrentEmployee)
		router.DELETE("/:id", handler.DeleteEmployee)
		router.GET("/roles", handler.GetAllRoles)
		router.PUT("/:id/password", handler.UpdatePassword)

		workdays := router.Group("/workdays")
		{
			workdays.POST("", handler.CreateEmployeeWorkday)
			workdays.GET("/:id", handler.GetEmployeeWorkday)
			workdays.GET("", handler.GetEmployeeWorkdays)
			workdays.PUT("/:id", handler.UpdateEmployeeWorkday)
			workdays.DELETE("/:id", handler.DeleteEmployeeWorkday)
		}
	}
}

func (r *Router) RegisterOrderRoutes(handler *orders.OrderHandler) {
	router := r.EmployeeRoutes.Group("/orders")
	{
		router.POST("", handler.CreateOrder)
		router.GET("", handler.GetOrders)
		router.GET("/ws/:storeId", handler.ServeWS)
		router.PUT("/:orderId/suborders/:subOrderId/complete", handler.CompleteSubOrder)
		router.GET("/kiosk", handler.GetAllBaristaOrders)
		router.GET("/:orderId/suborders", handler.GetSubOrders)
		router.GET("/statuses/count", handler.GetStatusesCount)
		router.GET("/:orderId/receipt", handler.GeneratePDFReceipt)
		router.GET("/:orderId", handler.GetOrderDetails)
	}
}

func (r *Router) RegisterSupplierRoutes(handler *supplier.SupplierHandler) {
	router := r.EmployeeRoutes.Group("/suppliers")
	{
		router.GET("", handler.GetSuppliers)
		router.GET("/:id", handler.GetSupplierByID)
		router.POST("", handler.CreateSupplier)
		router.PUT("/:id", handler.UpdateSupplier)
		router.DELETE("/:id", handler.DeleteSupplier)

		router.POST("/:id/materials", handler.AssociateMaterialToSupplier)
		router.GET("/:id/materials", handler.GetMaterialsBySupplier)
	}
}

func (r *Router) RegisterStoreWarehouseRoutes(handler *storeWarehouses.StoreWarehouseHandler) {
	router := r.EmployeeRoutes.Group("/store-warehouse-stock")
	{
		router.GET("", handler.GetStoreWarehouseStockList)
		router.GET("/:id", handler.GetStoreWarehouseStockById)
		router.POST("", handler.AddStoreWarehouseStock)
		router.POST("/multiple", handler.AddMultipleStoreWarehouseStock)
		router.PUT("/:id", handler.UpdateStoreWarehouseStockById)
		router.DELETE("/:id", handler.DeleteStoreWarehouseStockById)
	}
}

func (r *Router) RegisterStockMaterialRoutes(handler *stockMaterial.StockMaterialHandler) {
	router := r.EmployeeRoutes.Group("/stock-materials")
	{
		router.GET("", handler.GetAllStockMaterials)
		router.GET("/:id", handler.GetStockMaterialByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.CreateStockMaterial)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.UpdateStockMaterial)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.DeleteStockMaterial)
		router.PATCH("/:id/deactivate", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.DeactivateStockMaterial)
	}
}

func (r *Router) RegisterStockMaterialPackageRoutes(handler *stockMaterialPackage.StockMaterialPackageHandler) {
	router := r.EmployeeRoutes.Group("/stock-material-packages")
	{
		router.GET("", handler.GetAll)
		router.GET("/:id", handler.GetByID)
		router.POST("", handler.Create)
		router.PUT("/:id", handler.Update)
		router.DELETE("/:id", handler.Delete)
	}
}

func (r *Router) RegisterStockMaterialCategoryRoutes(handler *stockMaterialCategory.StockMaterialCategoryHandler) {
	router := r.EmployeeRoutes.Group("/stock-material-categories")
	{
		router.GET("", handler.GetAll)
		router.GET("/:id", handler.GetByID)
		router.POST("", handler.Create)
		router.PUT("/:id", handler.Update)
		router.DELETE("/:id", handler.Delete)
	}
}

func (r *Router) RegisterUnitRoutes(handler *units.UnitHandler) {
	router := r.EmployeeRoutes.Group("/units")
	{
		router.GET("", handler.GetAllUnits)
		router.GET("/:id", handler.GetUnitByID)
		router.POST("", handler.CreateUnit)
		router.PUT("/:id", handler.UpdateUnit)
		router.DELETE("/:id", handler.DeleteUnit)
	}
}

func (r *Router) RegisterBarcodeRoutes(handler *barcode.BarcodeHandler) {
	router := r.EmployeeRoutes.Group("/barcode")
	{
		router.POST("/generate", handler.GenerateBarcode)
		router.GET("/:barcode", handler.RetrieveStockMaterialByBarcode)
		router.POST("/print", handler.PrintAdditionalBarcodes)

		router.POST("/by-material", handler.GetBarcodesForStockMaterials)  // Retrieve multiple barcodes
		router.GET("/by-material/:id", handler.GetBarcodeForStockMaterial) // Retrieve a single barcode by ID
	}
}

func (r *Router) RegisterWarehouseRoutes(handler *warehouse.WarehouseHandler, warehouseStockHandler *warehouseStock.WarehouseStockHandler) {
	router := r.EmployeeRoutes.Group("/warehouses")
	{
		warehouseRoutes := router.Group("")
		{
			warehouseRoutes.POST("", handler.CreateWarehouse)                // Create a new warehouse
			warehouseRoutes.GET("/:warehouseId", handler.GetWarehouseByID)   // Get a specific warehouse by ID
			warehouseRoutes.PUT("/:warehouseId", handler.UpdateWarehouse)    // Update warehouse details
			warehouseRoutes.DELETE("/:warehouseId", handler.DeleteWarehouse) // Delete a warehouse
		}

		storeRoutes := router.Group("/stores")
		{
			storeRoutes.POST("", handler.AssignStoreToWarehouse)              // Assign a store to a warehouse
			storeRoutes.PUT("/:storeId", handler.ReassignStore)               // Reassign a store to another warehouse
			storeRoutes.GET("/:warehouseId", handler.GetAllStoresByWarehouse) // Get all stores assigned to a specific warehouse
		}

		stockRoutes := router.Group("/stocks")
		{
			stockRoutes.GET("", warehouseStockHandler.GetStocks)
			stockRoutes.GET("/:stockMaterialId", warehouseStockHandler.GetStockMaterialDetails)
			stockRoutes.PUT("/:stockMaterialId", warehouseStockHandler.UpdateStock)
			stockRoutes.POST("/add", warehouseStockHandler.AddWarehouseStocks)
			stockRoutes.POST("/receive", warehouseStockHandler.ReceiveInventory)
			stockRoutes.POST("/transfer", warehouseStockHandler.TransferInventory)
			stockRoutes.GET("/deliveries", warehouseStockHandler.GetDeliveries)
		}
	}
}

func (r *Router) RegisterStockRequestRoutes(handler *stockRequests.StockRequestHandler) {
	router := r.EmployeeRoutes.Group("/stock-requests")
	{
		router.GET("", handler.GetStockRequests)
		router.GET("/:requestId", handler.GetStockRequestByID)
		router.POST("", handler.CreateStockRequest)
		router.GET("/current", handler.GetLastCreatedStockRequest)
		router.PUT("/:requestId", handler.UpdateStockRequest)
		router.DELETE("/:requestId", handler.DeleteStockRequest)
		router.POST("/add-material-to-latest-cart", handler.AddStockMaterialToCart)

		statusGroup := router.Group("/status/:requestId")
		{
			statusGroup.PATCH("/accept-with-change", handler.AcceptWithChangeStatus) // DTO with different stock material
			statusGroup.PATCH("/reject-store", handler.RejectStoreStatus)            // Comment
			statusGroup.PATCH("/reject-warehouse", handler.RejectWarehouseStatus)    // Comment
			statusGroup.PATCH("/processed", handler.SetProcessedStatus)
			statusGroup.PATCH("/in-delivery", handler.SetInDeliveryStatus)
			statusGroup.PATCH("/completed", handler.SetCompletedStatus)
		}
	}
}
