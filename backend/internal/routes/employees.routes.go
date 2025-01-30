package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/analytics"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	adminEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/adminEmployees"
	franchiseeEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees"
	regionEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees"
	storeEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees"
	warehouseEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/warehouseEmployees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/ingredientCategories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/recipes"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock"
)

func (r *Router) RegisterAuditRoutes(handler *audit.AuditHandler) {
	router := r.EmployeeRoutes.Group("/audits")
	{
		router.GET("", handler.GetAudits)
	}
}

func (r *Router) RegisterFranchiseeRoutes(handler *franchisees.FranchiseeHandler) {
	router := r.EmployeeRoutes.Group("/franchisees")
	{
		router.GET("", handler.GetFranchisees)
		router.GET("/:id", handler.GetFranchiseeByID)
		router.POST("", handler.CreateFranchisee)
		router.PUT("/:id", handler.UpdateFranchisee)
		router.DELETE("/:id", handler.DeleteFranchisee)
	}
}

func (r *Router) RegisterRegionRoutes(handler *regions.RegionHandler) {
	router := r.EmployeeRoutes.Group("/regions")
	{
		router.GET("", handler.GetRegions)
		router.GET("/:id", handler.GetRegionByID)
		router.POST("", handler.CreateRegion)
		router.PUT("/:id", handler.UpdateRegion)
		router.DELETE("/:id", handler.DeleteRegion)
	}
}

func (r *Router) RegisterNotificationsRoutes(handler *notifications.NotificationHandler) {
	router := r.EmployeeRoutes.Group("/notifications")
	{
		router.GET("", handler.GetNotificationsByEmployee)
		router.GET("/:id", handler.GetNotificationByID)
		router.POST("/:id/mark-as-read", handler.MarkNotificationAsRead)
		router.POST("/mark-multiple-as-read", handler.MarkMultipleNotificationsAsRead)
		router.DELETE("/:id", handler.DeleteNotification)
	}
}

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
		router.GET("/categories", handler.GetStoreProductCategories)

		router.GET("", handler.GetStoreProducts)
		router.GET("/available-to-add", handler.GetAvailableProducts)
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
		router.POST("", handler.CreateStore)
		router.PUT("/:id", handler.UpdateStore)
		router.DELETE("/:id", handler.DeleteStore)
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
		router.GET("/addList", handler.GetAdditivesListToAdd)
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
		router.PUT("/reassign", handler.ReassignEmployeeType)
		router.GET("/current", handler.GetCurrentEmployee)
		router.GET("/roles", handler.GetAllRoles)
		router.PUT("/:id/password", handler.UpdatePassword)

		workdays := router.Group("/workdays")
		{
			var workdaysManagementPermissions = []data.EmployeeRole{data.RoleStoreManager, data.RoleWarehouseManager, data.RoleRegionWarehouseManager, data.RoleFranchiseManager}
			workdays.POST("", middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...), handler.CreateEmployeeWorkday)
			workdays.GET("/:id", handler.GetEmployeeWorkday)
			workdays.GET("", handler.GetEmployeeWorkdays)
			workdays.PUT("/:id", middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...), handler.UpdateEmployeeWorkday)
			workdays.DELETE("/:id", middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...), handler.DeleteEmployeeWorkday)
		}
	}
}

func (r *Router) RegisterStoreEmployeeRoutes(handler storeEmployees.StoreEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/store-employees")
	{
		router.GET("", handler.GetStoreEmployees)
		router.POST("", handler.CreateStoreEmployee)
		router.GET("/:id", handler.GetStoreEmployeeByID)
		router.PUT("/:id", handler.UpdateStoreEmployee)
		router.DELETE("/:employeeId", handler.DeleteStoreEmployee)
	}
}

func (r *Router) RegisterWarehouseEmployeeRoutes(handler warehouseEmployees.WarehouseEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/warehouse-employees")
	{
		router.GET("", handler.GetWarehouseEmployees)
		router.POST("", handler.CreateWarehouseEmployee)
		router.GET("/:id", handler.GetWarehouseEmployeeByID)
		router.PUT("/:id", handler.UpdateWarehouseEmployee)
		router.DELETE("/:employeeId", handler.DeleteWarehouseEmployee)
	}
}

func (r *Router) RegisterFranchiseeEmployeeRoutes(handler franchiseeEmployees.FranchiseeEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/store-employees")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), handler.GetFranchiseeEmployees)
		router.POST("", handler.CreateFranchiseeEmployee)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), handler.GetFranchiseeEmployeeByID)
		router.PUT("/:id", handler.UpdateFranchiseeEmployee)
		router.DELETE("/:employeeId", handler.DeleteFranchiseeEmployee)
	}
}

func (r *Router) RegisterRegionEmployeeRoutes(handler regionEmployees.RegionEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/region-employees")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.RegionReadPermissions...), handler.GetRegionEmployees)
		router.POST("", handler.CreateRegionEmployee)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.RegionReadPermissions...), handler.GetRegionEmployeeByID)
		router.PUT("/:id", handler.UpdateRegionEmployee)
		router.DELETE("/:employeeId", handler.DeleteRegionEmployee)
	}
}

func (r *Router) RegisterAdminEmployeeRoutes(handler adminEmployees.AdminEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/warehouse-employees")
	{
		router.GET("", handler.GetAdminEmployees, middleware.EmployeeRoleMiddleware())
		router.POST("", handler.CreateAdminEmployee, middleware.EmployeeRoleMiddleware())
		router.GET("/:id", handler.GetAdminEmployeeByID, middleware.EmployeeRoleMiddleware())
	}
}

func (r *Router) RegisterOrderRoutes(handler *orders.OrderHandler) {
	router := r.EmployeeRoutes.Group("/orders")
	{
		router.POST("", handler.CreateOrder)
		router.GET("", handler.GetOrders, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/ws", handler.ServeWS, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.PUT("/:orderId/suborders/:subOrderId/complete", handler.CompleteSubOrder, middleware.EmployeeRoleMiddleware(data.StoreWorkerPermissions...))
		router.GET("/kiosk", handler.GetAllBaristaOrders, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/:orderId/suborders", handler.GetSubOrders, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/statuses/count", handler.GetStatusesCount, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/:orderId/receipt", handler.GeneratePDFReceipt, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/:orderId", handler.GetOrderDetails, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))

		router.GET("/export", handler.ExportOrders, middleware.EmployeeRoleMiddleware(data.RoleOwner, data.RoleFranchiseOwner, data.RoleFranchiseManager))
		router.PUT("/suborders/:subOrderId/complete", handler.CompleteSubOrderByBarcode, middleware.EmployeeRoleMiddleware(data.StoreWorkerPermissions...))
		router.GET("/suborders/:subOrderId/barcode", handler.GetSuborderBarcode)
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

		router.PUT("/:id/materials", handler.UpsertMaterialsForSupplier)
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
		router.POST("", handler.CreateStockMaterial)
		router.PUT("/:id", handler.UpdateStockMaterial)
		router.DELETE("/:id", handler.DeleteStockMaterial)
		router.PATCH("/:id/deactivate", handler.DeactivateStockMaterial)

		router.GET("/:id/barcode", handler.GetStockMaterialBarcode)
		router.GET("/barcodes/:barcode", handler.RetrieveStockMaterialByBarcode)
		router.POST("/barcodes/generate", handler.GenerateBarcode)
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
			stockRoutes.GET("/available-to-add", warehouseStockHandler.GetAvailableToAddStockMaterials)
			stockRoutes.GET("/:stockMaterialId", warehouseStockHandler.GetStockMaterialDetails)
			stockRoutes.PUT("/:stockMaterialId", warehouseStockHandler.UpdateStock)
			stockRoutes.POST("/add", warehouseStockHandler.AddWarehouseStocks)
			stockRoutes.POST("/receive", warehouseStockHandler.ReceiveInventory)
			stockRoutes.POST("/transfer", warehouseStockHandler.TransferInventory)
			stockRoutes.GET("/deliveries", warehouseStockHandler.GetDeliveries)
			stockRoutes.GET("/deliveries/:id", warehouseStockHandler.GetDeliveryByID)
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

func (r *Router) RegisterAnalyticRoutes(handler *analytics.AnalyticsHandler) {
	router := r.EmployeeRoutes.Group("/analytics")
	{
		router.GET("/summary", handler.GetSummary)                  // Summary analytics
		router.GET("/sales-by-month", handler.GetSalesByMonth)      // Monthly sales analytics
		router.GET("/popular-products", handler.GetPopularProducts) // Popular products analytics
	}
}
