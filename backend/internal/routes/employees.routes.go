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
		router.GET("", middleware.EmployeeRoleMiddleware(data.RoleOwner, data.RoleFranchiseOwner, data.RoleFranchiseManager, data.RoleWarehouseManager), handler.GetAudits)
	}
}

func (r *Router) RegisterFranchiseeRoutes(handler *franchisees.FranchiseeHandler) {
	router := r.EmployeeRoutes.Group("/franchisees")
	{
		router.GET("", handler.GetFranchisees)
		router.GET("/:id", handler.GetFranchiseeByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateFranchisee)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateFranchisee)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteFranchisee)
	}
}

func (r *Router) RegisterRegionRoutes(handler *regions.RegionHandler) {
	router := r.EmployeeRoutes.Group("/regions")
	{
		router.GET("", handler.GetRegions)
		router.GET("/:id", handler.GetRegionByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateRegion)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateRegion)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteRegion)
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
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateProduct)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateProduct)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteProduct)
		router.GET(":id/sizes", handler.GetProductSizesByProductID)
		router.POST("/sizes", middleware.EmployeeRoleMiddleware(), handler.CreateProductSize)
		router.PUT("/sizes/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateProductSize)
		router.GET("/sizes/:id", handler.GetProductSizeByID)
	}
}

func (r *Router) RegisterRecipeRoutes(handler *recipes.RecipeHandler) {
	router := r.EmployeeRoutes.Group("/products/recipe-steps")
	{
		router.GET("/product/:product-id", handler.GetRecipeSteps)
		router.GET("/step/:id", handler.GetRecipeStepDetails)
		router.POST("/product/:product-id", middleware.EmployeeRoleMiddleware(), handler.CreateRecipeSteps)
		router.PUT("/step/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateRecipeSteps)
		router.DELETE("/step/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteRecipeSteps)
	}
}

func (r *Router) RegisterStoreProductRoutes(handler *storeProducts.StoreProductHandler) {
	router := r.EmployeeRoutes.Group("/store-products")
	{
		router.GET("/categories", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreProductCategories)

		router.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreProducts)
		router.GET("/available-to-add", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.GetAvailableProducts)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreProduct)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateStoreProduct)
		router.POST("/multiple", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateMultipleStoreProducts)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.UpdateStoreProduct)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.DeleteStoreProduct)

		router.GET("/sizes/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreProductSizeByID)
	}
}

func (r *Router) RegisterIngredientRoutes(handler *ingredients.IngredientHandler) {
	router := r.EmployeeRoutes.Group("/ingredients")
	{
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateIngredient)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateIngredient)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteIngredient)
		router.GET("/:id", handler.GetIngredientByID)
		router.GET("", handler.GetIngredients)
	}
}

func (r *Router) RegisterIngredientCategoriesRoutes(handler *ingredientCategories.IngredientCategoryHandler) {
	router := r.EmployeeRoutes.Group("/ingredient-categories")
	{
		router.GET("", handler.GetAll)
		router.GET("/:id", handler.GetByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.Create)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.Update)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.Delete)
	}
}

func (r *Router) RegisterStoresRoutes(handler *stores.StoreHandler) {
	router := r.EmployeeRoutes.Group("/stores")
	{
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateStore)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateStore)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteStore)
	}
}

func (r *Router) RegisterProductCategoriesRoutes(handler *categories.CategoryHandler) {
	router := r.EmployeeRoutes.Group("/product-categories")
	{
		router.GET("", handler.GetAllCategories)
		router.GET("/:id", handler.GetCategoryByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateCategory)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateCategory)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteCategory)
	}
}

func (r *Router) RegisterAdditivesRoutes(handler *additives.AdditiveHandler) {
	router := r.EmployeeRoutes.Group("/additives")
	{
		router.GET("", handler.GetAdditives)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateAdditive)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateAdditive)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteAdditive)
		router.GET("/:id", handler.GetAdditiveByID)

		additiveCategories := router.Group("/categories")
		{
			additiveCategories.GET("", handler.GetAdditiveCategories)
			additiveCategories.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateAdditiveCategory)
			additiveCategories.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateAdditiveCategory)
			additiveCategories.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteAdditiveCategory)
			additiveCategories.GET("/:id", handler.GetAdditiveCategoryByID)
		}
	}
}

func (r *Router) RegisterStoreAdditivesRoutes(handler *storeAdditives.StoreAdditiveHandler) {
	router := r.EmployeeRoutes.Group("/store-additives")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreAdditives)
		router.GET("/addList", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.GetAdditivesListToAdd)
		router.GET("/categories/:productSizeId", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreAdditiveCategories)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateStoreAdditives)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.UpdateStoreAdditive)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.DeleteStoreAdditive)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.GetStoreAdditiveByID)
	}
}

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.EmployeeRoutes.Group("/employees")
	{
		router.PUT("/reassign", middleware.EmployeeRoleMiddleware(), handler.ReassignEmployeeType)
		router.GET("/current", handler.GetCurrentEmployee)
		router.GET("/roles", handler.GetAllRoles)
		router.PUT("/:id/password", middleware.EmployeeRoleMiddleware(), handler.UpdatePassword)

		workdays := router.Group("/workdays")
		{
			var workdaysManagementPermissions = []data.EmployeeRole{data.RoleStoreManager, data.RoleWarehouseManager, data.RoleRegionWarehouseManager, data.RoleFranchiseManager}
			workdays.POST("", middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...), handler.CreateEmployeeWorkday)
			workdays.GET("/:id", middleware.EmployeeRoleMiddleware(), handler.GetEmployeeWorkday)
			workdays.GET("", middleware.EmployeeRoleMiddleware(), handler.GetEmployeeWorkdays)
			workdays.PUT("/:id", middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...), handler.UpdateEmployeeWorkday)
			workdays.DELETE("/:id", middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...), handler.DeleteEmployeeWorkday)
		}
	}
}

func (r *Router) RegisterStoreEmployeeRoutes(handler storeEmployees.StoreEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/store-employees")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreEmployees)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateStoreEmployee)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreEmployeeByID)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.UpdateStoreEmployee)
		router.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.DeleteStoreEmployee)
	}
}

func (r *Router) RegisterWarehouseEmployeeRoutes(handler warehouseEmployees.WarehouseEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/warehouse-employees")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetWarehouseEmployees)
		router.POST("", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.CreateWarehouseEmployee)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetWarehouseEmployeeByID)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.UpdateWarehouseEmployee)
		router.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.DeleteWarehouseEmployee)
	}
}

func (r *Router) RegisterFranchiseeEmployeeRoutes(handler franchiseeEmployees.FranchiseeEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/store-employees")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), handler.GetFranchiseeEmployees)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateFranchiseeEmployee)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), handler.GetFranchiseeEmployeeByID)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateFranchiseeEmployee)
		router.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(), handler.DeleteFranchiseeEmployee)
	}
}

func (r *Router) RegisterRegionEmployeeRoutes(handler regionEmployees.RegionEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/region-employees")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.RegionReadPermissions...), handler.GetRegionEmployees)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateRegionEmployee)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.RegionReadPermissions...), handler.GetRegionEmployeeByID)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateRegionEmployee)
		router.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(), handler.DeleteRegionEmployee)
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
		router.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreWarehouseStockList)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreWarehouseStockById)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.AddStoreWarehouseStock)
		router.POST("/multiple", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.AddMultipleStoreWarehouseStock)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.UpdateStoreWarehouseStockById)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.DeleteStoreWarehouseStockById)
	}
}

func (r *Router) RegisterStockMaterialRoutes(handler *stockMaterial.StockMaterialHandler) {
	router := r.EmployeeRoutes.Group("/stock-materials")
	{
		router.GET("", handler.GetAllStockMaterials, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.GET("/:id", handler.GetStockMaterialByID, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.POST("", handler.CreateStockMaterial, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.PUT("/:id", handler.UpdateStockMaterial, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.DELETE("/:id", handler.DeleteStockMaterial, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.PATCH("/:id/deactivate", handler.DeactivateStockMaterial, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))

		router.GET("/:id/barcode", handler.GetStockMaterialBarcode)
		router.GET("/barcodes/:barcode", handler.RetrieveStockMaterialByBarcode)
		router.POST("/barcodes/generate", handler.GenerateBarcode)
	}
}

func (r *Router) RegisterStockMaterialCategoryRoutes(handler *stockMaterialCategory.StockMaterialCategoryHandler) {
	router := r.EmployeeRoutes.Group("/stock-material-categories")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetAll)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.Create)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.Update)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.Delete)
	}
}

func (r *Router) RegisterUnitRoutes(handler *units.UnitHandler) {
	router := r.EmployeeRoutes.Group("/units")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(), handler.GetAllUnits)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(), handler.GetUnitByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateUnit)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateUnit)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteUnit)
	}
}

func (r *Router) RegisterWarehouseRoutes(handler *warehouse.WarehouseHandler, warehouseStockHandler *warehouseStock.WarehouseStockHandler) {
	router := r.EmployeeRoutes.Group("/warehouses")
	{
		warehouseRoutes := router.Group("")
		{
			warehouseRoutes.POST("", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.CreateWarehouse)                // Create a new warehouse
			warehouseRoutes.GET("/:warehouseId", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetWarehouseByID)         // Get a specific warehouse by ID
			warehouseRoutes.PUT("/:warehouseId", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.UpdateWarehouse)    // Update warehouse details
			warehouseRoutes.DELETE("/:warehouseId", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.DeleteWarehouse) // Delete a warehouse
		}

		storeRoutes := router.Group("/stores")
		{
			storeRoutes.POST("", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.AssignStoreToWarehouse)        // Assign a store to a warehouse
			storeRoutes.PUT("/:storeId", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.ReassignStore)         // Reassign a store to another warehouse
			storeRoutes.GET("/:warehouseId", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetAllStoresByWarehouse) // Get all stores assigned to a specific warehouse
		}

		stockRoutes := router.Group("/stocks")
		{
			stockRoutes.GET("", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetStocks)
			stockRoutes.GET("/available-to-add", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), warehouseStockHandler.GetAvailableToAddStockMaterials)
			stockRoutes.GET("/:stockMaterialId", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetStockMaterialDetails)
			stockRoutes.PUT("/:stockMaterialId", middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...), warehouseStockHandler.UpdateStock)
			stockRoutes.POST("/add", middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...), warehouseStockHandler.AddWarehouseStocks)
			stockRoutes.POST("/receive", middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...), warehouseStockHandler.ReceiveInventory)
			stockRoutes.POST("/transfer", middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...), warehouseStockHandler.TransferInventory)
			stockRoutes.GET("/deliveries", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetDeliveries)
			stockRoutes.GET("/deliveries/:id", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetDeliveryByID)
		}
	}
}

func (r *Router) RegisterStockRequestRoutes(handler *stockRequests.StockRequestHandler) {
	router := r.EmployeeRoutes.Group("/stock-requests")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetStockRequests)
		router.GET("/:requestId", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetStockRequestByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateStockRequest)
		router.GET("/current", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.GetLastCreatedStockRequest)
		router.PUT("/:requestId", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.UpdateStockRequest)
		router.DELETE("/:requestId", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.DeleteStockRequest)
		router.POST("/add-material-to-latest-cart", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.AddStockMaterialToCart)

		statusGroup := router.Group("/status/:requestId")
		{
			statusGroup.PATCH("/accept-with-change", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.AcceptWithChangeStatus) // DTO with different stock material
			statusGroup.PATCH("/reject-store", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.RejectStoreStatus)            // Comment
			statusGroup.PATCH("/reject-warehouse", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.RejectWarehouseStatus)    // Comment
			statusGroup.PATCH("/processed", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.SetProcessedStatus)
			statusGroup.PATCH("/in-delivery", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.SetInDeliveryStatus)
			statusGroup.PATCH("/completed", middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...), handler.SetCompletedStatus)
		}
	}
}

func (r *Router) RegisterAnalyticRoutes(handler *analytics.AnalyticsHandler) {
	router := r.EmployeeRoutes.Group("/analytics", middleware.EmployeeRoleMiddleware(data.RoleOwner, data.RoleFranchiseOwner))
	{
		router.GET("/summary", handler.GetSummary)                  // Summary analytics
		router.GET("/sales-by-month", handler.GetSalesByMonth)      // Monthly sales analytics
		router.GET("/popular-products", handler.GetPopularProducts) // Popular products analytics
	}
}
