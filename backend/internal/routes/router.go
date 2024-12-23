package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Prefix  string
	Version string
	Routes  *gin.RouterGroup
}

func NewRouter(engine *gin.Engine, prefix string, version string) *Router {
	router := engine.Group(prefix + version)
	return &Router{
		Prefix:  prefix,
		Version: version,
		Routes:  router,
	}
}

func (r *Router) RegisterAuthenticationRoutes(handler *auth.AuthenticationHandler) {
	router := r.Routes.Group("/auth")
	{
		customersRoutes := router.Group("/customers")
		{
			customersRoutes.POST("/register", handler.CustomerRegister)
			customersRoutes.POST("/login", handler.CustomerLogin)
			customersRoutes.POST("/refresh", handler.CustomerRefresh)
			customersRoutes.POST("/logout", handler.CustomerLogout)
		}

		employeesRoutes := router.Group("/employees")
		{
			employeesRoutes.POST("/login", handler.EmployeeLogin)
			employeesRoutes.POST("/refresh", handler.EmployeeRefresh)
			employeesRoutes.POST("/logout", handler.EmployeeLogout)
		}
	}
}

func (r *Router) RegisterProductRoutes(handler *product.ProductHandler) {
	router := r.Routes.Group("/products")
	{
		router.GET("", handler.GetProducts)
		router.GET("/:productId", handler.GetProductDetails)
	}
}

func (r *Router) RegisterStoresRoutes(handler *stores.StoreHandler) {
	router := r.Routes.Group("/stores")
	{
		router.GET("", handler.GetAllStores)
		router.GET("/:id", handler.GetStoreByID)
		router.POST("", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.CreateStore)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.UpdateStore)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.DeleteStore)
	}
}

func (r *Router) RegisterProductCategoriesRoutes(handler *categories.CategoryHandler) {
	router := r.Routes.Group("/product-categories")
	{
		router.GET("", handler.GetAllCategories)
	}
}

func (r *Router) RegisterAdditivesRoutes(handler *additives.AdditiveHandler) {
	router := r.Routes.Group("/additives")
	{
		router.GET("", handler.GetAdditives)
		router.GET("/categories", handler.GetAdditiveCategories)
	}
}

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.Routes.Group("/employees")
	{
		router.POST("", handler.CreateEmployee)
		router.GET("", handler.GetEmployees)
		router.GET("/current", handler.GetCurrentEmployee)
		router.GET("/:id", handler.GetEmployeeByID)
		router.PUT("/:id", handler.UpdateEmployee)
		router.DELETE("/:id", handler.DeleteEmployee)
		router.GET("/roles", handler.GetAllRoles)
		router.PUT("/:id/password", handler.UpdatePassword)
	}
}

func (r *Router) RegisterOrderRoutes(handler *orders.OrderHandler) {
	router := r.Routes.Group("/orders")
	{
		router.POST("", handler.CreateOrder)
		router.GET("", handler.GetOrders)
		router.GET("/ws/:storeId", handler.ServeWS)
		router.PUT("/:orderId/suborders/:subOrderId/complete", handler.CompleteSubOrder)
		router.GET("/kiosk", handler.GetAllBaristaOrders)
		router.GET("/:orderId/suborders", handler.GetSubOrders)
		router.GET("/statuses/count", handler.GetStatusesCount)
		router.GET("/:orderId/receipt", handler.GeneratePDFReceipt)
	}
}

func (r *Router) RegisterSupplierRoutes(handler *supplier.SupplierHandler) {
	router := r.Routes.Group("/suppliers")
	{
		router.GET("", handler.GetSuppliers)
		router.GET("/:id", handler.GetSupplierByID)
		router.POST("", handler.CreateSupplier)
		router.PUT("/:id", handler.UpdateSupplier)
		router.DELETE("/:id", handler.DeleteSupplier)
	}
}

func (r *Router) RegisterStoreWarehouseRoutes(handler *storeWarehouses.StoreWarehouseHandler) {
	router := r.Routes.Group("/store-warehouse-stock/:store_id")
	router.Use(middleware.EmployeeAuth(), middleware.MatchesStore())
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
	router := r.Routes.Group("/stock-material")
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
	router := r.Routes.Group("/barcode")
	{
		router.POST("/generate", handler.GenerateBarcode)
		router.GET("/:barcode", handler.RetrieveStockMaterialByBarcode)
		router.POST("/print", handler.PrintAdditionalBarcodes)
	}
}

func (r *Router) RegisterInventoryRoutes(handler *inventory.InventoryHandler) {
	router := r.Routes.Group("/inventory")
	{
		stockRoutes := router.Group("/stock")
		{
			stockRoutes.POST("/receive", handler.ReceiveInventory)
			stockRoutes.GET("/levels/:warehouseID", handler.GetInventoryLevels)
			stockRoutes.POST("/pickup", handler.PickupStock) // store
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
	router := r.Routes.Group("/warehouse")
	{
		warehouseRoutes := router.Group("")
		{
			warehouseRoutes.POST("", handler.CreateWarehouse)                // Create a new warehouse
			warehouseRoutes.GET("", handler.GetAllWarehouses)                // Get all warehouses
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
	router := r.Routes.Group("/stock-requests")
	{
		router.GET("", handler.GetStockRequests)                                                                                               // Get all stock requests with filtering
		router.GET("/low-stock", handler.GetLowStockIngredients)                                                                               // Get low-stock ingredients
		router.GET("/stock-materials", handler.GetAllStockMaterials)                                                                           // Get marketplace products
		router.POST("", handler.CreateStockRequest)                                                                                            // Create a new stock request (cart creation)
		router.PATCH("/:requestId/status", middleware.EmployeeRoleMiddleware(data.RoleAdmin), handler.UpdateStockRequestStatus)                // Update stock request status
		router.POST("/:requestId/ingredients", middleware.EmployeeRoleMiddleware(data.RoleManager), handler.AddStockRequestIngredient)         // Add ingredient to cart
		router.DELETE("/ingredients/:ingredientId", middleware.EmployeeRoleMiddleware(data.RoleManager), handler.DeleteStockRequestIngredient) // Delete ingredient from cart
	}
}
