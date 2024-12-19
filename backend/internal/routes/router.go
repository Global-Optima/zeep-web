package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
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

func (r *Router) RegisterIngredientRoutes(handler *ingredients.IngredientHandler) {
	router := r.Routes.Group("/ingredients")
	{
		router.POST("", handler.CreateIngredient)
		router.PUT("/:id", handler.UpdateIngredient)
		router.DELETE("/:id", handler.DeleteIngredient)
		router.GET("/:id", handler.GetIngredientByID)
		router.GET("", handler.GetIngredients)
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
		router.POST("", handler.CreateAdditive)
		router.PUT("", handler.UpdateAdditive)
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

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.Routes.Group("/employees")
	{
		storeEmployees := router.Group("/store")
		{
			storeEmployees.POST("", handler.CreateStoreEmployee)
			storeEmployees.GET("", handler.GetStoreEmployees)
			storeEmployees.GET("/:id", handler.GetStoreEmployeeByID)
			storeEmployees.PUT("/:id", handler.UpdateStoreEmployee)
		}
		warehouseEmployees := router.Group("/warehouse")
		{
			warehouseEmployees.POST("", handler.CreateWarehouseEmployee)
			warehouseEmployees.GET("", handler.GetWarehouseEmployees)
			warehouseEmployees.GET("/:id", handler.GetWarehouseEmployeeByID)
			warehouseEmployees.PUT("/:id", handler.UpdateWarehouseEmployee)
		}

		router.GET("/current", handler.GetCurrentEmployee)
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
		router.GET("/:orderId", handler.GetOrderDetails)
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
	router := r.Routes.Group("/store-warehouse-stock")
	router.Use(middleware.EmployeeAuth())
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
		router.POST("/receive", handler.ReceiveInventory)
		router.GET("/levels/:warehouseID", handler.GetInventoryLevels)
		router.POST("/pickup", handler.PickupStock) // store
		router.POST("/transfer", handler.TransferInventory)

		router.GET("/expiration/upcoming/:warehouseID", handler.GetExpiringItems)
		router.POST("/expiration/extend", handler.ExtendExpiration)

		router.GET("/deliveries", handler.GetDeliveries)
	}
}

func (r *Router) RegisterWarehouseRoutes(handler *warehouse.WarehouseHandler) {
	router := r.Routes.Group("/warehouse") // store
	{
		router.POST("/stores", handler.AssignStoreToWarehouse)              // store
		router.PUT("/stores/:storeId", handler.ReassignStore)               // store
		router.GET("/:warehouseId/stores", handler.GetAllStoresByWarehouse) // store
	}
}
