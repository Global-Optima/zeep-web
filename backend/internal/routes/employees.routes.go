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
		router.GET("", handler.GetAudits) // store manager, warehouse manager, franchise owner
	}
}

func (r *Router) RegisterFranchiseeRoutes(handler *franchisees.FranchiseeHandler) {
	router := r.EmployeeRoutes.Group("/franchisees")
	{
		router.GET("", handler.GetFranchisees)        // owner
		router.GET("/:id", handler.GetFranchiseeByID) // owner
		router.POST("", handler.CreateFranchisee)
		router.PUT("/:id", handler.UpdateFranchisee)    // franchise owner, manager
		router.DELETE("/:id", handler.DeleteFranchisee) // franchise owner, manager

		// Add new endpoint to get MY franchise (/franchisees/my)
	}
}

func (r *Router) RegisterRegionRoutes(handler *regions.RegionHandler) {
	router := r.EmployeeRoutes.Group("/regions")
	{
		router.GET("", handler.GetRegions)        // owner
		router.GET("/:id", handler.GetRegionByID) // owner
		router.POST("", handler.CreateRegion)
		router.PUT("/:id", handler.UpdateRegion)
		router.DELETE("/:id", handler.DeleteRegion)
	}
}

func (r *Router) RegisterNotificationsRoutes(handler *notifications.NotificationHandler) {
	router := r.EmployeeRoutes.Group("/notifications")
	{
		router.GET("", handler.GetNotificationsByEmployee) // rename to my endpoints
		router.GET("/:id", handler.GetNotificationByID)
		router.POST("/:id/mark-as-read", handler.MarkNotificationAsRead)
		router.POST("/mark-multiple-as-read", handler.MarkMultipleNotificationsAsRead)
		router.DELETE("/:id", handler.DeleteNotification)
	}
}

func (r *Router) RegisterProductRoutes(handler *product.ProductHandler) {
	router := r.EmployeeRoutes.Group("/products")
	{
		router.GET("", handler.GetProducts)           // all
		router.GET("/:id", handler.GetProductDetails) // all
		router.POST("", handler.CreateProduct)
		router.PUT("/:id", handler.UpdateProduct)
		router.DELETE("/:id", handler.DeleteProduct)
		router.GET(":id/sizes", handler.GetProductSizesByProductID) // all
		router.POST("/sizes", handler.CreateProductSize)
		router.PUT("/sizes/:id", handler.UpdateProductSize)
		router.GET("/sizes/:id", handler.GetProductSizeByID)
	}
}

func (r *Router) RegisterRecipeRoutes(handler *recipes.RecipeHandler) {
	router := r.EmployeeRoutes.Group("/products/recipe-steps")
	{
		router.GET("/product/:product-id", handler.GetRecipeSteps) // all
		router.GET("/step/:id", handler.GetRecipeStepDetails)      // all
		router.POST("/product/:product-id", handler.CreateRecipeSteps)
		router.PUT("/step/:id", handler.UpdateRecipeSteps)
		router.DELETE("/step/:id", handler.DeleteRecipeSteps)
	}
}

func (r *Router) RegisterStoreProductRoutes(handler *storeProducts.StoreProductHandler) {
	router := r.EmployeeRoutes.Group("/store-products") // to all below store, franchisee
	{
		router.GET("/categories", handler.GetStoreProductCategories)
		router.GET("", handler.GetStoreProducts)
		router.GET("/available-to-add", handler.GetAvailableProducts)
		router.GET("/:id", handler.GetStoreProduct)
		router.POST("", handler.CreateStoreProduct)                   // Franchise + store manager
		router.POST("/multiple", handler.CreateMultipleStoreProducts) // Franchise + store manager
		router.PUT("/:id", handler.UpdateStoreProduct)                // Franchise + store manager
		router.DELETE("/:id", handler.DeleteStoreProduct)             // Franchise + store manager
		router.GET("/sizes/:id", handler.GetStoreProductSizeByID)
	}
}

func (r *Router) RegisterIngredientRoutes(handler *ingredients.IngredientHandler) {
	router := r.EmployeeRoutes.Group("/ingredients")
	{
		router.POST("", handler.CreateIngredient)
		router.PUT("/:id", handler.UpdateIngredient)
		router.DELETE("/:id", handler.DeleteIngredient)
		router.GET("/:id", handler.GetIngredientByID) // all
		router.GET("", handler.GetIngredients)        // all
	}
}

func (r *Router) RegisterIngredientCategoriesRoutes(handler *ingredientCategories.IngredientCategoryHandler) {
	router := r.EmployeeRoutes.Group("/ingredient-categories")
	{
		router.GET("", handler.GetAll)      // all
		router.GET("/:id", handler.GetByID) // all
		router.POST("", handler.Create)
		router.PUT("/:id", handler.Update)
		router.DELETE("/:id", handler.Delete)
	}
}

func (r *Router) RegisterStoresRoutes(handler *stores.StoreHandler) {
	router := r.EmployeeRoutes.Group("/stores")
	{
		router.GET("/:id", handler.GetStoreByID)   // all
		router.POST("", handler.CreateStore)       // franchise owner, manager
		router.PUT("/:id", handler.UpdateStore)    // franchise owner, manager
		router.DELETE("/:id", handler.DeleteStore) // franchise owner, manager
	}
}

func (r *Router) RegisterProductCategoriesRoutes(handler *categories.CategoryHandler) {
	router := r.EmployeeRoutes.Group("/product-categories") // merge with products routes, like in additives (/products/categories)
	{
		router.GET("", handler.GetAllCategories)    // all
		router.GET("/:id", handler.GetCategoryByID) // all
		router.POST("", handler.CreateCategory)
		router.PUT("/:id", handler.UpdateCategory)
		router.DELETE("/:id", handler.DeleteCategory)
	}
}

func (r *Router) RegisterAdditivesRoutes(handler *additives.AdditiveHandler) {
	router := r.EmployeeRoutes.Group("/additives")
	{
		router.GET("", handler.GetAdditives)        // all
		router.GET("/:id", handler.GetAdditiveByID) // all
		router.POST("", handler.CreateAdditive)
		router.PUT("/:id", handler.UpdateAdditive)
		router.DELETE("/:id", handler.DeleteAdditive)

		additiveCategories := router.Group("/categories")
		{
			additiveCategories.GET("", handler.GetAdditiveCategories)       // all
			additiveCategories.GET("/:id", handler.GetAdditiveCategoryByID) // all
			additiveCategories.POST("", handler.CreateAdditiveCategory)
			additiveCategories.PUT("/:id", handler.UpdateAdditiveCategory)
			additiveCategories.DELETE("/:id", handler.DeleteAdditiveCategory)
		}
	}
}

func (r *Router) RegisterStoreAdditivesRoutes(handler *storeAdditives.StoreAdditiveHandler) {
	router := r.EmployeeRoutes.Group("/store-additives") // to all below store, franchisee
	{
		router.GET("", handler.GetStoreAdditives)
		router.GET("/available-to-add", handler.GetAdditivesListToAdd)
		router.GET("/categories/:productSizeId", handler.GetStoreAdditiveCategories)
		router.GET("/:id", handler.GetStoreAdditiveByID)
		router.POST("", handler.CreateStoreAdditives)      // franchisee and store manager
		router.PUT("/:id", handler.UpdateStoreAdditive)    // franchisee and store manager
		router.DELETE("/:id", handler.DeleteStoreAdditive) // franchisee and store manager
	}
}

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.EmployeeRoutes.Group("/employees")
	{
		router.PUT("/reassign", handler.ReassignEmployeeType)
		router.GET("/current", handler.GetCurrentEmployee)  // all
		router.GET("/roles", handler.GetAllRoles)           // all
		router.PUT("/:id/password", handler.UpdatePassword) // all

		// leave only GET methods, make update and create inside the employee updates
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

func (r *Router) RegisterStoreEmployeeRoutes(handler *storeEmployees.StoreEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/store-employees") // franchisee and store managers
	{
		router.GET("", handler.GetStoreEmployees)        // fr owner and owner
		router.GET("/:id", handler.GetStoreEmployeeByID) // fr owner and owner
		router.POST("", handler.CreateStoreEmployee)
		router.PUT("/:id", handler.UpdateStoreEmployee)
		router.DELETE("/:employeeId", handler.DeleteStoreEmployee)
	}
}

func (r *Router) RegisterWarehouseEmployeeRoutes(handler warehouseEmployees.WarehouseEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/warehouse-employees") // warehouse and region managers
	{
		router.GET("", handler.GetWarehouseEmployees)        // owner
		router.GET("/:id", handler.GetWarehouseEmployeeByID) // owner
		router.POST("", handler.CreateWarehouseEmployee)
		router.PUT("/:id", handler.UpdateWarehouseEmployee)
		router.DELETE("/:employeeId", handler.DeleteWarehouseEmployee)
	}
}

func (r *Router) RegisterFranchiseeEmployeeRoutes(handler franchiseeEmployees.FranchiseeEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/franchisee-employees")
	{
		router.GET("", handler.GetFranchiseeEmployees)                  // owner, all franchisee
		router.GET("/:id", handler.GetFranchiseeEmployeeByID)           // owner, all franchisee
		router.POST("", handler.CreateFranchiseeEmployee)               // franchise owner
		router.PUT("/:id", handler.UpdateFranchiseeEmployee)            // franchise owner
		router.DELETE("/:employeeId", handler.DeleteFranchiseeEmployee) // franchise owner
	}
}

func (r *Router) RegisterRegionEmployeeRoutes(handler regionEmployees.RegionEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/region-employees")
	{
		router.GET("", handler.GetRegionEmployees)
		router.GET("/:id", handler.GetRegionEmployeeByID)
		router.POST("", handler.CreateRegionEmployee)
		router.PUT("/:id", handler.UpdateRegionEmployee)
		router.DELETE("/:employeeId", handler.DeleteRegionEmployee)
	}
}

func (r *Router) RegisterAdminEmployeeRoutes(handler adminEmployees.AdminEmployeeHandler) {
	router := r.EmployeeRoutes.Group("/admin-employees")
	{
		router.GET("", handler.GetAdminEmployees)
		router.POST("", handler.CreateAdminEmployee)
		router.GET("/:id", handler.GetAdminEmployeeByID)
	}
}

func (r *Router) RegisterOrderRoutes(handler *orders.OrderHandler) {
	router := r.EmployeeRoutes.Group("/orders")
	{
		router.POST("", handler.CreateOrder)                                             // Store manager and barista
		router.GET("", handler.GetOrders)                                                // all franchise, stores
		router.GET("/ws", handler.ServeWS)                                               // Store manager and barista
		router.PUT("/:orderId/suborders/:subOrderId/complete", handler.CompleteSubOrder) // Store manager and barista
		router.GET("/kiosk", handler.GetAllBaristaOrders)                                // Store manager and barista
		router.GET("/:orderId/suborders", handler.GetSubOrders)                          // Store manager and barista
		router.GET("/statuses/count", handler.GetStatusesCount)                          // Store manager and barista
		router.GET("/:orderId/receipt", handler.GeneratePDFReceipt)                      // Store manager and barista
		router.GET("/:orderId", handler.GetOrderDetails)                                 // all franchise, stores

		router.GET("/export", handler.ExportOrders)                                      // franchise and store management
		router.PUT("/suborders/:subOrderId/complete", handler.CompleteSubOrderByBarcode) // Store manager and barista
		router.GET("/suborders/:subOrderId/barcode", handler.GetSuborderBarcode)         // Store manager and barista
	}
}

func (r *Router) RegisterSupplierRoutes(handler *supplier.SupplierHandler) {
	router := r.EmployeeRoutes.Group("/suppliers")
	{
		router.GET("", handler.GetSuppliers)        // Region, warehouse
		router.GET("/:id", handler.GetSupplierByID) // Region, warehouse
		router.POST("", handler.CreateSupplier)
		router.PUT("/:id", handler.UpdateSupplier)
		router.DELETE("/:id", handler.DeleteSupplier)

		router.PUT("/:id/materials", handler.UpsertMaterialsForSupplier)
		router.GET("/:id/materials", handler.GetMaterialsBySupplier) // Region, warehouse
	}
}

func (r *Router) RegisterStoreWarehouseRoutes(handler *storeWarehouses.StoreWarehouseHandler) {
	router := r.EmployeeRoutes.Group("/store-warehouse-stock") // Franchise and store all roles
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
		router.GET("", handler.GetAllStockMaterials)     // all
		router.GET("/:id", handler.GetStockMaterialByID) // all
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
			warehouseRoutes.POST("", handler.CreateWarehouse) // region
			warehouseRoutes.GET("/:warehouseId", handler.GetWarehouseByID)
			warehouseRoutes.PUT("/:warehouseId", handler.UpdateWarehouse)    // region
			warehouseRoutes.DELETE("/:warehouseId", handler.DeleteWarehouse) // region
		}

		storeRoutes := router.Group("/stores")
		{
			storeRoutes.POST("", handler.AssignStoreToWarehouse)              // On considerations
			storeRoutes.PUT("/:storeId", handler.ReassignStore)               // On considerations
			storeRoutes.GET("/:warehouseId", handler.GetAllStoresByWarehouse) // Store and warehouse
		}

		stockRoutes := router.Group("/stocks")
		{
			stockRoutes.GET("", warehouseStockHandler.GetStocks) // Region and warehouses all roles
			stockRoutes.GET("/:stockMaterialId", warehouseStockHandler.GetStockMaterialDetails)
			stockRoutes.PUT("/:stockMaterialId", warehouseStockHandler.UpdateStock)   // Warehouse all roles
			stockRoutes.POST("/add", warehouseStockHandler.AddWarehouseStocks)        // Warehouse all roles
			stockRoutes.POST("/receive", warehouseStockHandler.ReceiveInventory)      // Warehouse all roles
			stockRoutes.GET("/deliveries", warehouseStockHandler.GetDeliveries)       // Region and warehouses
			stockRoutes.GET("/deliveries/:id", warehouseStockHandler.GetDeliveryByID) // Region and warehouses

			// Consider if it needed, consider requests
			stockRoutes.POST("/transfer", warehouseStockHandler.TransferInventory) // region manager
		}
	}
}

func (r *Router) RegisterStockRequestRoutes(handler *stockRequests.StockRequestHandler) {
	router := r.EmployeeRoutes.Group("/stock-requests")
	{
		router.GET("", handler.GetStockRequests)                                    // Store and warehouses all roles
		router.GET("/:requestId", handler.GetStockRequestByID)                      // Store and warehouses all roles
		router.GET("/current", handler.GetLastCreatedStockRequest)                  // Store and warehouses all roles
		router.POST("", handler.CreateStockRequest)                                 // Store all roles
		router.POST("/add-material-to-latest-cart", handler.AddStockMaterialToCart) // Store all roles
		router.PUT("/:requestId", handler.UpdateStockRequest)                       // Store all roles
		router.DELETE("/:requestId", handler.DeleteStockRequest)                    // Store all roles

		statusGroup := router.Group("/status/:requestId")
		{
			statusGroup.PATCH("/processed", handler.SetProcessedStatus)              // Store
			statusGroup.PATCH("/reject-warehouse", handler.RejectWarehouseStatus)    // Warehouse
			statusGroup.PATCH("/in-delivery", handler.SetInDeliveryStatus)           // Warehouse
			statusGroup.PATCH("/accept-with-change", handler.AcceptWithChangeStatus) // Store
			statusGroup.PATCH("/reject-store", handler.RejectStoreStatus)            // Store
			statusGroup.PATCH("/completed", handler.SetCompletedStatus)              // Store
		}
	}
}

func (r *Router) RegisterAnalyticRoutes(handler *analytics.AnalyticsHandler) {
	router := r.EmployeeRoutes.Group("/analytics")
	// Pasha don't do this, this is bullshit
	{
		router.GET("/summary", handler.GetSummary)
		router.GET("/sales-by-month", handler.GetSalesByMonth)
		router.GET("/popular-products", handler.GetPopularProducts)
	}
}
