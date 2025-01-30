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
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock"
)

func (r *Router) RegisterAuditRoutes(handler *audit.AuditHandler) {
	router := r.EmployeeRoutes.Group("/audits")
	{
		router.GET("", handler.GetAudits, middleware.EmployeeRoleMiddleware(data.RoleOwner, data.RoleFranchiseOwner, data.RoleFranchiseManager, data.RoleWarehouseManager))
	}
}

func (r *Router) RegisterFranchiseeRoutes(handler *franchisees.FranchiseeHandler) {
	router := r.EmployeeRoutes.Group("/franchisees")
	{
		router.GET("", handler.GetFranchisees)
		router.GET("/:id", handler.GetFranchiseeByID)
		router.POST("", handler.CreateFranchisee, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateFranchisee, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteFranchisee, middleware.EmployeeRoleMiddleware())
	}
}

func (r *Router) RegisterRegionRoutes(handler *regions.RegionHandler) {
	router := r.EmployeeRoutes.Group("/regions")
	{
		router.GET("", handler.GetRegions)
		router.GET("/:id", handler.GetRegionByID)
		router.POST("", handler.CreateRegion, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateRegion, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteRegion, middleware.EmployeeRoleMiddleware())
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
		router.POST("", handler.CreateProduct, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateProduct, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteProduct, middleware.EmployeeRoleMiddleware())
		router.GET(":id/sizes", handler.GetProductSizesByProductID)
		router.POST("/sizes", handler.CreateProductSize, middleware.EmployeeRoleMiddleware())
		router.PUT("/sizes/:id", handler.UpdateProductSize, middleware.EmployeeRoleMiddleware())
		router.GET("/sizes/:id", handler.GetProductSizeByID)
	}
}

func (r *Router) RegisterRecipeRoutes(handler *recipes.RecipeHandler) {
	router := r.EmployeeRoutes.Group("/products/recipe-steps")
	{
		router.GET("/product/:product-id", handler.GetRecipeSteps)
		router.GET("/step/:id", handler.GetRecipeStepDetails)
		router.POST("/product/:product-id", handler.CreateRecipeSteps, middleware.EmployeeRoleMiddleware())
		router.PUT("/step/:id", handler.UpdateRecipeSteps, middleware.EmployeeRoleMiddleware())
		router.DELETE("/step/:id", handler.DeleteRecipeSteps, middleware.EmployeeRoleMiddleware())
	}
}

func (r *Router) RegisterStoreProductRoutes(handler *storeProducts.StoreProductHandler) {
	router := r.EmployeeRoutes.Group("/store-products")
	{
		router.GET("/categories", handler.GetStoreProductCategories, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))

		router.GET("", handler.GetStoreProducts, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/addList", handler.GetProductsListToAdd, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.GET("/:id", handler.GetStoreProduct, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.POST("", handler.CreateStoreProduct, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.POST("/multiple", handler.CreateMultipleStoreProducts, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.PUT("/:id", handler.UpdateStoreProduct, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.DELETE("/:id", handler.DeleteStoreProduct, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))

		router.GET("/sizes/:id", handler.GetStoreProductSizeByID, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
	}
}

func (r *Router) RegisterIngredientRoutes(handler *ingredients.IngredientHandler) {
	router := r.EmployeeRoutes.Group("/ingredients")
	{
		router.POST("", handler.CreateIngredient, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateIngredient, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteIngredient, middleware.EmployeeRoleMiddleware())
		router.GET("/:id", handler.GetIngredientByID)
		router.GET("", handler.GetIngredients)
	}
}

func (r *Router) RegisterIngredientCategoriesRoutes(handler *ingredientCategories.IngredientCategoryHandler) {
	router := r.EmployeeRoutes.Group("/ingredient-categories")
	{
		router.GET("", handler.GetAll)
		router.GET("/:id", handler.GetByID)
		router.POST("", handler.Create, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.Update, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.Delete, middleware.EmployeeRoleMiddleware())

	}
}

func (r *Router) RegisterStoresRoutes(handler *stores.StoreHandler) {
	router := r.EmployeeRoutes.Group("/stores")
	{
		router.GET("/:id", handler.GetStoreByID, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.POST("", handler.CreateStore, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateStore, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteStore, middleware.EmployeeRoleMiddleware())
	}
}

func (r *Router) RegisterProductCategoriesRoutes(handler *categories.CategoryHandler) {
	router := r.EmployeeRoutes.Group("/product-categories")
	{
		router.GET("", handler.GetAllCategories)
		router.GET("/:id", handler.GetCategoryByID)
		router.POST("", handler.CreateCategory, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateCategory, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteCategory, middleware.EmployeeRoleMiddleware())
	}
}

func (r *Router) RegisterAdditivesRoutes(handler *additives.AdditiveHandler) {
	router := r.EmployeeRoutes.Group("/additives")
	{
		router.GET("", handler.GetAdditives)
		router.POST("", handler.CreateAdditive, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateAdditive, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteAdditive, middleware.EmployeeRoleMiddleware())
		router.GET("/:id", handler.GetAdditiveByID)

		additiveCategories := router.Group("/categories")
		{
			additiveCategories.GET("", handler.GetAdditiveCategories)
			additiveCategories.POST("", handler.CreateAdditiveCategory, middleware.EmployeeRoleMiddleware())
			additiveCategories.PUT("/:id", handler.UpdateAdditiveCategory, middleware.EmployeeRoleMiddleware())
			additiveCategories.DELETE("/:id", handler.DeleteAdditiveCategory, middleware.EmployeeRoleMiddleware())
			additiveCategories.GET("/:id", handler.GetAdditiveCategoryByID)
		}
	}
}

func (r *Router) RegisterStoreAdditivesRoutes(handler *storeAdditives.StoreAdditiveHandler) {
	router := r.EmployeeRoutes.Group("/store-additives")
	{
		router.GET("", handler.GetStoreAdditives, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/addList", handler.GetAdditivesListToAdd, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.GET("/categories/:productSizeId", handler.GetStoreAdditiveCategories, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.POST("", handler.CreateStoreAdditives, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.PUT("/:id", handler.UpdateStoreAdditive, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.DELETE("/:id", handler.DeleteStoreAdditive, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.GET("/:id", handler.GetStoreAdditiveByID, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
	}
}

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.EmployeeRoutes.Group("/employees")
	{

		router.PUT("/reassign", handler.ReassignEmployeeType, middleware.EmployeeRoleMiddleware())
		router.GET("/current", handler.GetCurrentEmployee)
		router.GET("/roles", handler.GetAllRoles)
		router.PUT("/:id/password", handler.UpdatePassword, middleware.EmployeeRoleMiddleware())

		workdays := router.Group("/workdays")

		{
			var workdaysManagementPermissions = []data.EmployeeRole{data.RoleStoreManager, data.RoleWarehouseManager, data.RoleRegionWarehouseManager, data.RoleFranchiseManager}
			workdays.POST("", handler.CreateEmployeeWorkday, middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...))
			workdays.GET("/:id", handler.GetEmployeeWorkday)
			workdays.GET("", handler.GetEmployeeWorkdays)
			workdays.PUT("/:id", handler.UpdateEmployeeWorkday, middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...))
			workdays.DELETE("/:id", handler.DeleteEmployeeWorkday, middleware.EmployeeRoleMiddleware(workdaysManagementPermissions...))
		}
	}
}

func (r *Router) RegisterStoreEmployeeRoutes(handler storeEmployees.StoreEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/store-employees")
	{
		router.GET("", handler.GetStoreEmployees, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.POST("", handler.CreateStoreEmployee, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.GET("/:id", handler.GetStoreEmployeeByID, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.PUT("/:id", handler.UpdateStoreEmployee, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.DELETE("/:employeeId", handler.DeleteStoreEmployee, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
	}
}

func (r *Router) RegisterWarehouseEmployeeRoutes(handler warehouseEmployees.WarehouseEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/warehouse-employees")
	{
		router.GET("", handler.GetWarehouseEmployees, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.POST("", handler.CreateWarehouseEmployee, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.GET("/:id", handler.GetWarehouseEmployeeByID, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.PUT("/:id", handler.UpdateWarehouseEmployee, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.DELETE("/:employeeId", handler.DeleteWarehouseEmployee, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
	}
}

func (r *Router) RegisterFranchiseeEmployeeRoutes(handler franchiseeEmployees.FranchiseeEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/store-employees")
	{
		router.GET("", handler.GetFranchiseeEmployees, middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...))
		router.POST("", handler.CreateFranchiseeEmployee, middleware.EmployeeRoleMiddleware())
		router.GET("/:id", handler.GetFranchiseeEmployeeByID, middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...))
		router.PUT("/:id", handler.UpdateFranchiseeEmployee, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:employeeId", handler.DeleteFranchiseeEmployee, middleware.EmployeeRoleMiddleware())
	}
}

func (r *Router) RegisterRegionEmployeeRoutes(handler regionEmployees.RegionEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/region-employees")
	{
		router.GET("", handler.GetRegionEmployees, middleware.EmployeeRoleMiddleware(data.RegionReadPermissions...))
		router.POST("", handler.CreateRegionEmployee, middleware.EmployeeRoleMiddleware())
		router.GET("/:id", handler.GetRegionEmployeeByID, middleware.EmployeeRoleMiddleware(data.RegionReadPermissions...))
		router.PUT("/:id", handler.UpdateRegionEmployee, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:employeeId", handler.DeleteRegionEmployee, middleware.EmployeeRoleMiddleware())
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
		router.GET("", handler.GetStoreWarehouseStockList, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.GET("/:id", handler.GetStoreWarehouseStockById, middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...))
		router.POST("", handler.AddStoreWarehouseStock, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.POST("/multiple", handler.AddMultipleStoreWarehouseStock, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.PUT("/:id", handler.UpdateStoreWarehouseStockById, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.DELETE("/:id", handler.DeleteStoreWarehouseStockById, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
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
	}
}

func (r *Router) RegisterStockMaterialCategoryRoutes(handler *stockMaterialCategory.StockMaterialCategoryHandler) {
	router := r.EmployeeRoutes.Group("/stock-material-categories")
	{
		router.GET("", handler.GetAll, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.GET("/:id", handler.GetByID, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.POST("", handler.Create, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.PUT("/:id", handler.Update, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.DELETE("/:id", handler.Delete, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
	}
}

func (r *Router) RegisterUnitRoutes(handler *units.UnitHandler) {
	router := r.EmployeeRoutes.Group("/units")
	{
		router.GET("", handler.GetAllUnits, middleware.EmployeeRoleMiddleware())
		router.GET("/:id", handler.GetUnitByID, middleware.EmployeeRoleMiddleware())
		router.POST("", handler.CreateUnit, middleware.EmployeeRoleMiddleware())
		router.PUT("/:id", handler.UpdateUnit, middleware.EmployeeRoleMiddleware())
		router.DELETE("/:id", handler.DeleteUnit, middleware.EmployeeRoleMiddleware())
	}
}

func (r *Router) RegisterBarcodeRoutes(handler *barcode.BarcodeHandler) {
	router := r.EmployeeRoutes.Group("/barcode")
	{
		router.POST("/generate", handler.GenerateBarcode, middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...))
		router.GET("/:barcode", handler.RetrieveStockMaterialByBarcode, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.POST("/print", handler.PrintAdditionalBarcodes, middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...))

		router.POST("/by-material", handler.GetBarcodesForStockMaterials, middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...)) // Retrieve multiple barcodes
		router.GET("/by-material/:id", handler.GetBarcodeForStockMaterial, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))  // Retrieve a single barcode by ID
	}
}

func (r *Router) RegisterWarehouseRoutes(handler *warehouse.WarehouseHandler, warehouseStockHandler *warehouseStock.WarehouseStockHandler) {
	router := r.EmployeeRoutes.Group("/warehouses")
	{
		warehouseRoutes := router.Group("")
		{
			warehouseRoutes.POST("", handler.CreateWarehouse, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))                // Create a new warehouse
			warehouseRoutes.GET("/:warehouseId", handler.GetWarehouseByID, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))         // Get a specific warehouse by ID
			warehouseRoutes.PUT("/:warehouseId", handler.UpdateWarehouse, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))    // Update warehouse details
			warehouseRoutes.DELETE("/:warehouseId", handler.DeleteWarehouse, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...)) // Delete a warehouse
		}

		storeRoutes := router.Group("/stores")
		{
			storeRoutes.POST("", handler.AssignStoreToWarehouse, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))        // Assign a store to a warehouse
			storeRoutes.PUT("/:storeId", handler.ReassignStore, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))         // Reassign a store to another warehouse
			storeRoutes.GET("/:warehouseId", handler.GetAllStoresByWarehouse, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...)) // Get all stores assigned to a specific warehouse
		}

		stockRoutes := router.Group("/stocks")
		{
			stockRoutes.GET("", warehouseStockHandler.GetStocks, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
			stockRoutes.GET("/:stockMaterialId", warehouseStockHandler.GetStockMaterialDetails, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
			stockRoutes.PUT("/:stockMaterialId", warehouseStockHandler.UpdateStock, middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...))
			stockRoutes.POST("/add", warehouseStockHandler.AddWarehouseStocks, middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...))
			stockRoutes.POST("/receive", warehouseStockHandler.ReceiveInventory, middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...))
			stockRoutes.POST("/transfer", warehouseStockHandler.TransferInventory, middleware.EmployeeRoleMiddleware(data.WarehouseWorkerPermissions...))
			stockRoutes.GET("/deliveries", warehouseStockHandler.GetDeliveries, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
			stockRoutes.GET("/deliveries/:id", warehouseStockHandler.GetDeliveryByID, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		}
	}
}

func (r *Router) RegisterStockRequestRoutes(handler *stockRequests.StockRequestHandler) {
	router := r.EmployeeRoutes.Group("/stock-requests")
	{
		router.GET("", handler.GetStockRequests, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.GET("/:requestId", handler.GetStockRequestByID, middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...))
		router.POST("", handler.CreateStockRequest, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.GET("/current", handler.GetLastCreatedStockRequest, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))
		router.PUT("/:requestId", handler.UpdateStockRequest, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.DELETE("/:requestId", handler.DeleteStockRequest, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
		router.POST("/add-material-to-latest-cart", handler.AddStockMaterialToCart, middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...))

		statusGroup := router.Group("/status/:requestId")
		{
			statusGroup.PATCH("/accept-with-change", handler.AcceptWithChangeStatus, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...)) // DTO with different stock material
			statusGroup.PATCH("/reject-store", handler.RejectStoreStatus, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))            // Comment
			statusGroup.PATCH("/reject-warehouse", handler.RejectWarehouseStatus, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))    // Comment
			statusGroup.PATCH("/processed", handler.SetProcessedStatus, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
			statusGroup.PATCH("/in-delivery", handler.SetInDeliveryStatus, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
			statusGroup.PATCH("/completed", handler.SetCompletedStatus, middleware.EmployeeRoleMiddleware(data.WarehouseManagementPermissions...))
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
