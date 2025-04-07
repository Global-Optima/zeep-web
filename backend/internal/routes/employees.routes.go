package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	additivesTechnicalMap "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/technicalMap"
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
	productTechnicalMap "github.com/Global-Optima/zeep-web/backend/internal/modules/product/technicalMap"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers"
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
		router.GET("", middleware.EmployeeRoleMiddleware(data.RoleOwner), handler.GetFranchisees)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.RoleOwner), handler.GetFranchiseeByID)
		router.GET("/my", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), handler.GetMyFranchisee)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateFranchisee)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), handler.UpdateFranchisee)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.RoleOwner, data.RoleFranchiseOwner), handler.DeleteFranchisee)
	}
}

func (r *Router) RegisterRegionRoutes(handler *regions.RegionHandler) {
	router := r.EmployeeRoutes.Group("/regions")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.RoleOwner), handler.GetRegions)        // owner
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.RoleOwner), handler.GetRegionByID) // owner
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

func (r *Router) RegisterProductRoutes(handler *product.ProductHandler, productTechMapHandler *productTechnicalMap.TechnicalMapHandler) {
	router := r.EmployeeRoutes.Group("/products")
	{
		router.GET("", handler.GetProducts)
		router.GET("/:id", handler.GetProductDetails)
		router.GET(":id/sizes", handler.GetProductSizesByProductID)
		router.GET("/sizes/:id", handler.GetProductSizeByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateProduct)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateProduct)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteProduct)
		router.DELETE("sizes/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteProductSize)
		router.POST("/sizes", middleware.EmployeeRoleMiddleware(), handler.CreateProductSize)
		router.PUT("/sizes/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateProductSize)
		router.GET("/sizes/:id/technical-map", middleware.EmployeeRoleMiddleware(), productTechMapHandler.GetProductSizeTechnicalMapByID)
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
		router.GET("/available-to-add", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetAvailableProductsToAdd)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreProduct)
		router.GET("/recommended", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetRecommendedStoreProducts)
		router.GET("/sizes/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreProductSizeByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateStoreProduct)
		router.POST("/multiple", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateMultipleStoreProducts)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.UpdateStoreProduct)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.DeleteStoreProduct)
	}
}

func (r *Router) RegisterIngredientRoutes(handler *ingredients.IngredientHandler) {
	router := r.EmployeeRoutes.Group("/ingredients")
	{
		router.GET("/:id", handler.GetIngredientByID)
		router.GET("", handler.GetIngredients)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateIngredient)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateIngredient)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteIngredient)
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
		router.GET("/:id", handler.GetStoreByID)
		router.GET("", handler.GetStoresByFranchisee)
		router.POST("", middleware.EmployeeRoleMiddleware(data.FranchiseePermissions...), handler.CreateStore)       // franchise owner, manager
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.FranchiseePermissions...), handler.UpdateStore)    // franchise owner, manager
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.FranchiseePermissions...), handler.DeleteStore) // franchise owner, manager
	}
}

func (r *Router) RegisterProductCategoriesRoutes(handler *categories.CategoryHandler) {
	router := r.EmployeeRoutes.Group("/product-categories") // merge with products routes, like in additives (/products/categories)
	{
		router.GET("", handler.GetAllCategories)
		router.GET("/:id", handler.GetCategoryByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateCategory)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateCategory)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteCategory)
	}
}

func (r *Router) RegisterAdditivesRoutes(handler *additives.AdditiveHandler, additivesTechMapHandler *additivesTechnicalMap.TechnicalMapHandler) {
	router := r.EmployeeRoutes.Group("/additives")
	{
		router.GET("", handler.GetAdditives)
		router.GET("/:id", handler.GetAdditiveByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateAdditive)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateAdditive)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteAdditive)
		router.GET("/:id/technical-map", middleware.EmployeeRoleMiddleware(), additivesTechMapHandler.GetAdditiveTechnicalMapByID)

		additiveCategories := router.Group("/categories")
		{
			additiveCategories.GET("", handler.GetAdditiveCategories)
			additiveCategories.GET("/:id", handler.GetAdditiveCategoryByID)
			additiveCategories.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateAdditiveCategory)
			additiveCategories.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateAdditiveCategory)
			additiveCategories.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteAdditiveCategory)
		}
	}
}

func (r *Router) RegisterStoreAdditivesRoutes(handler *storeAdditives.StoreAdditiveHandler) {
	router := r.EmployeeRoutes.Group("/store-additives")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreAdditives)
		router.GET("/available-to-add", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetAdditivesListToAdd)
		router.GET("/categories/:storeProductSizeId", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreAdditiveCategories)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreAdditiveByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.CreateStoreAdditives)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.UpdateStoreAdditive)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.DeleteStoreAdditive)
	}
}

func (r *Router) RegisterEmployeesRoutes(
	handler *employees.EmployeeHandler,
	storeEmployeeHandler *storeEmployees.StoreEmployeeHandler,
	warehouseEmployeeHandler *warehouseEmployees.WarehouseEmployeeHandler,
	franchiseeEmployeeHandler *franchiseeEmployees.FranchiseeEmployeeHandler,
	regionEmployeeHandler *regionEmployees.RegionEmployeeHandler,
	adminEmployeeHandler *adminEmployees.AdminEmployeeHandler,
) {
	router := r.EmployeeRoutes.Group("/employees")
	{
		router.GET("/:id", middleware.EmployeeRoleMiddleware(), handler.GetEmployeeByID)
		router.GET("", middleware.EmployeeRoleMiddleware(), handler.GetEmployees)
		router.GET("/current", handler.GetCurrentEmployee)
		router.GET("/roles", handler.GetAllRoles)
		router.PUT("/:id/password", middleware.EmployeeRoleMiddleware(), handler.UpdatePassword)
		router.PUT("/:id/reassign", middleware.EmployeeRoleMiddleware(), handler.ReassignEmployeeType)

		storeEmployeeRouter := router.Group("/store")
		{
			storeEmployeeRouter.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), storeEmployeeHandler.GetStoreEmployees)
			storeEmployeeRouter.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), storeEmployeeHandler.GetStoreEmployeeByID)
			storeEmployeeRouter.POST("", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), storeEmployeeHandler.CreateStoreEmployee)
			storeEmployeeRouter.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), storeEmployeeHandler.UpdateStoreEmployee)
			storeEmployeeRouter.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), storeEmployeeHandler.DeleteStoreEmployee)
		}

		warehouseEmployeeRouter := router.Group("/warehouse") // warehouse and region managers
		{
			warehouseEmployeeRouter.GET("", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseEmployeeHandler.GetWarehouseEmployees)
			warehouseEmployeeRouter.GET("/:id", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseEmployeeHandler.GetWarehouseEmployeeByID)
			warehouseEmployeeRouter.POST("", middleware.EmployeeRoleMiddleware(data.RegionPermissions...), warehouseEmployeeHandler.CreateWarehouseEmployee)
			warehouseEmployeeRouter.PUT("/:id", middleware.EmployeeRoleMiddleware(data.RegionPermissions...), warehouseEmployeeHandler.UpdateWarehouseEmployee)
			warehouseEmployeeRouter.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(data.RegionPermissions...), warehouseEmployeeHandler.DeleteWarehouseEmployee)
		}

		franchiseeEmployeeRouter := router.Group("/franchisee")
		{
			franchiseeEmployeeRouter.GET("", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), franchiseeEmployeeHandler.GetFranchiseeEmployees)        // owner, all franchisee
			franchiseeEmployeeRouter.GET("/:id", middleware.EmployeeRoleMiddleware(data.FranchiseeReadPermissions...), franchiseeEmployeeHandler.GetFranchiseeEmployeeByID) // owner, all franchisee
			franchiseeEmployeeRouter.POST("", middleware.EmployeeRoleMiddleware(data.RoleFranchiseOwner), franchiseeEmployeeHandler.CreateFranchiseeEmployee)               // franchise owner
			franchiseeEmployeeRouter.PUT("/:id", middleware.EmployeeRoleMiddleware(data.RoleFranchiseOwner), franchiseeEmployeeHandler.UpdateFranchiseeEmployee)            // franchise owner
			franchiseeEmployeeRouter.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(data.RoleFranchiseOwner), franchiseeEmployeeHandler.DeleteFranchiseeEmployee) // franchise owner
		}

		regionEmployeeRouter := router.Group("/region")
		{
			regionEmployeeRouter.GET("", middleware.EmployeeRoleMiddleware(data.RoleOwner), regionEmployeeHandler.GetRegionEmployees)
			regionEmployeeRouter.GET("/:id", middleware.EmployeeRoleMiddleware(data.RoleOwner), regionEmployeeHandler.GetRegionEmployeeByID)
			regionEmployeeRouter.POST("", middleware.EmployeeRoleMiddleware(), regionEmployeeHandler.CreateRegionEmployee)
			regionEmployeeRouter.PUT("/:id", middleware.EmployeeRoleMiddleware(), regionEmployeeHandler.UpdateRegionEmployee)
			regionEmployeeRouter.DELETE("/:employeeId", middleware.EmployeeRoleMiddleware(), regionEmployeeHandler.DeleteRegionEmployee)
		}

		adminEmployeeRouter := router.Group("/admin")
		{
			adminEmployeeRouter.GET("", middleware.EmployeeRoleMiddleware(), adminEmployeeHandler.GetAdminEmployees)
			adminEmployeeRouter.POST("", middleware.EmployeeRoleMiddleware(), adminEmployeeHandler.CreateAdminEmployee)
			adminEmployeeRouter.GET("/:id", middleware.EmployeeRoleMiddleware(), adminEmployeeHandler.GetAdminEmployeeByID)
		}

		workdays := router.Group("/workdays")
		{
			workdays.GET("/:id", handler.GetEmployeeWorkday)
			workdays.GET("", handler.GetEmployeeWorkdays)
		}
	}
}

func (r *Router) RegisterOrderRoutes(handler *orders.OrderHandler) {
	router := r.EmployeeRoutes.Group("/orders")
	{
		router.POST("", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.CreateOrder)
		router.GET("/:orderId", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetOrderDetails) // all franchise, stores
		router.POST("/check-name", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.CheckCustomerName)
		router.POST("/:orderId/payment/fail", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.FailOrderPayment)
		router.POST("/:orderId/payment/success", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.SuccessOrderPayment)
		router.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetOrders)                   // all franchise, stores
		router.GET("/ws", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.ServeWS)                      // Store manager and barista
		router.GET("/kiosk", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.GetAllBaristaOrders)       // Store manager and barista
		router.GET("/:orderId/suborders", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.GetSubOrders) // Store manager and barista

		router.GET("/export", middleware.EmployeeRoleMiddleware(data.StoreManagementPermissions...), handler.ExportOrders) // franchise and store management
		router.PUT("/suborders/:subOrderId/status-change", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.SetNextSubOrderStatus)
	}
}

func (r *Router) RegisterSupplierRoutes(handler *supplier.SupplierHandler) {
	router := r.EmployeeRoutes.Group("/suppliers")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetSuppliers)        // Region, warehouse
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetSupplierByID) // Region, warehouse
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateSupplier)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateSupplier)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteSupplier)

		router.GET("/:id/materials", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), handler.GetMaterialsBySupplier) // Region, warehouse
		router.PUT("/:id/materials", middleware.EmployeeRoleMiddleware(), handler.UpsertMaterialsForSupplier)
	}
}

func (r *Router) RegisterStoreWarehouseRoutes(handler *storeStocks.StoreStockHandler) {
	router := r.EmployeeRoutes.Group("/store-stocks") // Franchise and store all roles
	{
		router.GET("/available-to-add", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetAvailableIngredientsToAdd)
		router.GET("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreStockList)
		router.GET("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.GetStoreStockById)
		router.POST("", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.AddStoreStock)
		router.POST("/multiple", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.AddMultipleStoreStock)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.UpdateStoreStockById)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.DeleteStoreStockById)
	}
}

func (r *Router) RegisterStockMaterialRoutes(handler *stockMaterial.StockMaterialHandler) {
	router := r.EmployeeRoutes.Group("/stock-materials")
	{
		router.GET("", handler.GetAllStockMaterials)
		router.GET("/:id", handler.GetStockMaterialByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.CreateStockMaterial)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.UpdateStockMaterial)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.DeleteStockMaterial)
		router.PATCH("/:id/deactivate", middleware.EmployeeRoleMiddleware(), handler.DeactivateStockMaterial)

		router.GET("/:id/barcode", handler.GetStockMaterialBarcode)
		router.GET("/barcodes/:barcode", handler.RetrieveStockMaterialByBarcode)
		router.POST("/barcodes/generate", middleware.EmployeeRoleMiddleware(), handler.GenerateBarcode)
	}
}

func (r *Router) RegisterStockMaterialCategoryRoutes(handler *stockMaterialCategory.StockMaterialCategoryHandler) {
	router := r.EmployeeRoutes.Group("/stock-material-categories")
	{
		router.GET("", handler.GetAll)
		router.GET("/:id", handler.GetByID)
		router.POST("", middleware.EmployeeRoleMiddleware(), handler.Create)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(), handler.Update)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(), handler.Delete)
	}
}

func (r *Router) RegisterUnitRoutes(handler *units.UnitHandler) {
	router := r.EmployeeRoutes.Group("/units")
	{
		router.GET("", handler.GetAllUnits)
		router.GET("/:id", handler.GetUnitByID)
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
			warehouseRoutes.GET("/:warehouseId", handler.GetWarehouseByID)
			warehouseRoutes.GET("", handler.GetWarehouses)
			warehouseRoutes.POST("", middleware.EmployeeRoleMiddleware(data.RegionPermissions...), handler.CreateWarehouse)                // region
			warehouseRoutes.PUT("/:warehouseId", middleware.EmployeeRoleMiddleware(data.RegionPermissions...), handler.UpdateWarehouse)    // region
			warehouseRoutes.DELETE("/:warehouseId", middleware.EmployeeRoleMiddleware(data.RegionPermissions...), handler.DeleteWarehouse) // region
		}

		storeRoutes := router.Group("/stores")
		{
			storeRoutes.POST("", handler.AssignStoreToWarehouse)
			storeRoutes.GET("/:warehouseId", middleware.EmployeeRoleMiddleware(data.StoreAndWarehousePermissions...), handler.GetAllStoresByWarehouse) // Store and warehouse
		}

		stockRoutes := router.Group("/stocks")
		{
			stockRoutes.GET("/available-to-add", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), warehouseStockHandler.GetAvailableToAddStockMaterials)
			stockRoutes.GET("", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetStocks) // Region and warehouses all roles
			stockRoutes.GET("/:stockMaterialId", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetStockMaterialDetails)
			stockRoutes.PUT("/:stockMaterialId", middleware.EmployeeRoleMiddleware(data.WarehousePermissions...), warehouseStockHandler.UpdateStock)       // Warehouse all roles
			stockRoutes.POST("/add", middleware.EmployeeRoleMiddleware(data.WarehousePermissions...), warehouseStockHandler.AddWarehouseStocks)            // Warehouse all roles
			stockRoutes.POST("/receive", middleware.EmployeeRoleMiddleware(data.WarehousePermissions...), warehouseStockHandler.ReceiveInventory)          // Warehouse all roles
			stockRoutes.GET("/deliveries", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetDeliveries)       // Region and warehouses
			stockRoutes.GET("/deliveries/:id", middleware.EmployeeRoleMiddleware(data.WarehouseReadPermissions...), warehouseStockHandler.GetDeliveryByID) // Region and warehouses
		}
	}
}

func (r *Router) RegisterStockRequestRoutes(handler *stockRequests.StockRequestHandler) {
	stockRequestReadPermissions := append(data.StoreAndWarehousePermissions, data.FranchiseeAndRegionPermissions...)

	router := r.EmployeeRoutes.Group("/stock-requests")
	{
		router.GET("", middleware.EmployeeRoleMiddleware(stockRequestReadPermissions...), handler.GetStockRequests)                              // Store and warehouses all roles
		router.GET("/:requestId", middleware.EmployeeRoleMiddleware(stockRequestReadPermissions...), handler.GetStockRequestByID)                // Store and warehouses all roles
		router.GET("/current", middleware.EmployeeRoleMiddleware(data.StoreAndWarehousePermissions...), handler.GetLastCreatedStockRequest)      // Store and warehouses all roles
		router.POST("", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.CreateStockRequest)                                 // Store all roles
		router.POST("/add-material-to-latest-cart", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.AddStockMaterialToCart) // Store all roles
		router.PUT("/:requestId", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.UpdateStockRequest)                       // Store all roles
		router.DELETE("/:requestId", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.DeleteStockRequest)                    // Store all roles

		statusGroup := router.Group("/status/:requestId")
		{
			statusGroup.PATCH("/processed", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.SetProcessedStatus)               // Store
			statusGroup.PATCH("/reject-warehouse", middleware.EmployeeRoleMiddleware(data.WarehousePermissions...), handler.RejectWarehouseStatus) // Warehouse
			statusGroup.PATCH("/in-delivery", middleware.EmployeeRoleMiddleware(data.WarehousePermissions...), handler.SetInDeliveryStatus)        // Warehouse
			statusGroup.PATCH("/accept-with-change", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.AcceptWithChangeStatus)  // Store
			statusGroup.PATCH("/reject-store", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.RejectStoreStatus)             // Store
			statusGroup.PATCH("/completed", middleware.EmployeeRoleMiddleware(data.StorePermissions...), handler.SetCompletedStatus)               // Store
		}
	}
}

func (r *Router) RegisterStoreSynchronizerSynchronizerRoutes(handler *storeSynchronizers.StoreSynchronizerHandler) {
	router := r.EmployeeRoutes.Group("/sync")
	{
		router.GET("/store", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.IsSynchronizedStore)
		router.POST("/store", middleware.EmployeeRoleMiddleware(data.StoreReadPermissions...), handler.SynchronizeStore)
	}
}

func (r *Router) RegisterAnalyticRoutes(handler *analytics.AnalyticsHandler) {
	router := r.EmployeeRoutes.Group("/analytics")
	{
		router.GET("/summary", handler.GetSummary)
		router.GET("/sales-by-month", handler.GetSalesByMonth)
		router.GET("/popular-products", handler.GetPopularProducts)
	}
}
