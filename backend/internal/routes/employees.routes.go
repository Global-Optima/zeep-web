package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
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
			additiveCategories.PUT("", handler.UpdateAdditiveCategory)
			additiveCategories.DELETE("/:id", handler.DeleteAdditiveCategory)
			additiveCategories.GET("/:id", handler.GetAdditiveCategoryByID)
		}
	}
}

func (r *Router) RegisterStoreAdditivesRoutes(handler *storeAdditives.StoreAdditiveHandler) {
	router := r.EmployeeRoutes.Group("/store-additives")
	{
		router.GET("", handler.GetStoreAdditives)
		router.GET("/categories", handler.GetStoreAdditiveCategories)
		router.POST("", handler.CreateStoreAdditives)
		router.PUT("/:id", handler.UpdateStoreAdditive)
		router.DELETE("/:id", handler.DeleteStoreAdditive)
		router.GET("/:id", handler.GetStoreAdditiveByID)
	}
}

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.EmployeeRoutes.Group("/employees")
	{
		storeEmployees := router.Group("/store")
		{
			storeEmployees.POST("", handler.CreateStoreEmployee)
			storeEmployees.GET("/:id", handler.GetStoreEmployeeByID)
			storeEmployees.PUT("/:id", handler.UpdateStoreEmployee)
		}
		warehouseEmployees := router.Group("/warehouse")
		{
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
	router := r.EmployeeRoutes.Group("/stock-material")
	{
		router.GET("", handler.GetAllStockMaterials)
		router.GET("/:id", handler.GetStockMaterialByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.CreateStockMaterial)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.UpdateStockMaterial)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.DeleteStockMaterial)
		router.PATCH("/:id/deactivate", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.DeactivateStockMaterial)
	}
}

func (r *Router) RegisterBarcodeRouter(handler *barcode.BarcodeHandler) {
	router := r.EmployeeRoutes.Group("/barcode")
	{
		router.POST("/generate", handler.GenerateBarcode)
		router.GET("/:barcode", handler.RetrieveStockMaterialByBarcode)
		router.POST("/print", handler.PrintAdditionalBarcodes)
	}
}

func (r *Router) RegisterInventoryRoutes(handler *inventory.InventoryHandler) {
	router := r.EmployeeRoutes.Group("/warehouse-stocks")
	{
		stockRoutes := router.Group("/stock")
		{
			stockRoutes.GET("", handler.GetInventoryLevels)
			stockRoutes.POST("/receive", handler.ReceiveInventory)
			stockRoutes.POST("/pickup", handler.PickupStock)
			stockRoutes.POST("/transfer", handler.TransferInventory)
			stockRoutes.GET("/deliveries", handler.GetDeliveries)
		}

		expirationRoutes := router.Group("/expiration")
		{
			expirationRoutes.GET("/:warehouseID", handler.GetExpiringItems)
			expirationRoutes.POST("/extend", handler.ExtendExpiration)
		}
	}
}

func (r *Router) RegisterWarehouseRoutes(handler *warehouse.WarehouseHandler) {
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

		stockRoutes := router.Group("/stock")
		{
			stockRoutes.POST("/add", handler.AddToStock)
			stockRoutes.POST("/deduct", handler.DeductFromStock)
			stockRoutes.GET("", handler.GetStock)
			stockRoutes.POST("/reset", handler.ResetStock)
		}
	}
}

func (r *Router) RegisterStockRequestRoutes(handler *stockRequests.StockRequestHandler) {
	router := r.EmployeeRoutes.Group("/stock-requests")
	{
		router.GET("", handler.GetStockRequests)
		router.GET("/:id", handler.GetStockRequestByID)
		router.GET("/low-stock", handler.GetLowStockIngredients)
		router.GET("/marketplace-products", handler.GetAllStockMaterials)
		router.POST("", handler.CreateStockRequest)
		router.PUT("/:requestId/status", handler.UpdateStockRequestStatus)
		router.PUT("/:requestId/ingredients", handler.UpdateStockRequestIngredients)
	}
}
