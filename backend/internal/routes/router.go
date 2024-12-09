package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/sku"
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

func (r *Router) RegisterProductRoutes(handler *product.ProductHandler) {
	router := r.Routes.Group("/products")
	{
		router.GET("", handler.GetStoreProducts)
		router.GET("/:productId", handler.GetStoreProductDetails)
	}
}

func (r *Router) RegisterStoresRoutes(handler *stores.StoreHandler) {
	router := r.Routes.Group("/stores")
	{
		router.GET("", handler.GetAllStores)
		router.GET("/:id", handler.GetStoreByID)
		router.POST("", middleware.EmployeeRoleMiddleware("Admin"), handler.CreateStore)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware("Admin"), handler.UpdateStore)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware("Admin"), handler.DeleteStore)
	}
}

func (r *Router) RegisterProductCategoriesRoutes(handler *categories.CategoryHandler) {
	router := r.Routes.Group("/categories")
	{
		router.GET("", handler.GetAllCategories)
	}
}

func (r *Router) RegisterAdditivesRoutes(handler *additives.AdditiveHandler) {
	router := r.Routes.Group("/additives")
	{
		router.GET("", handler.GetAdditivesByStoreAndProduct)
	}
}

func (r *Router) RegisterEmployeesRoutes(handler *employees.EmployeeHandler) {
	router := r.Routes.Group("/employees")
	{
		router.POST("", middleware.EmployeeRoleMiddleware(types.RoleDirector), handler.CreateEmployee)
		router.GET("", handler.GetEmployeesByStore)
		router.GET("/current", handler.GetCurrentEmployee)
		router.GET("/:id", handler.GetEmployeeByID)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(types.RoleDirector), handler.UpdateEmployee)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(types.RoleDirector), handler.DeleteEmployee)
		router.GET("/roles", middleware.EmployeeRoleMiddleware(types.RoleDirector, types.RoleManager), handler.GetAllRoles)
		router.PUT("/:id/password", middleware.EmployeeRoleMiddleware(types.RoleDirector, types.RoleManager, types.RoleBarista), handler.UpdatePassword)
		router.POST("/login", handler.EmployeeLogin)
		router.POST("/logout", handler.EmployeeLogout)
	}
}

func (r *Router) RegisterOrderRoutes(handler *orders.OrderHandler) {
	router := r.Routes.Group("/orders")
	{
		router.POST("", handler.CreateOrder)
		router.PUT("/:orderId/suborders/:subOrderId/complete", middleware.EmployeeRoleMiddleware(types.RoleBarista), handler.CompleteSubOrder)
		router.GET("", handler.GetAllOrders)
		router.GET("/:orderId/suborders", handler.GetSubOrders)
		router.GET("/:orderId", handler.GetActiveOrder)
		router.GET("/statuses/count", handler.GetStatusesCount)
		router.GET("/:orderId/suborders/count", handler.GetSubOrderCount)
		router.GET("/:orderId/receipt", handler.GeneratePDFReceipt)
	}
}

func (r *Router) RegisterSupplierRoutes(handler *supplier.SupplierHandler) {
	router := r.Routes.Group("/supplier")
	{
		router.GET("", handler.ListSuppliers)
		router.GET("/:id", handler.GetSupplierByID)
		router.POST("", handler.CreateSupplier)
		router.PUT("/:id", handler.UpdateSupplier)
		router.DELETE("/:id", handler.DeleteSupplier)
	}
}

func (r *Router) RegisterSKURoutes(handler *sku.SKUHandler) {
	router := r.Routes.Group("/warehouse/sku")
	{
		router.GET("", handler.GetAllSKUs)
		router.GET("/:id", handler.GetSKUByID)
		router.POST("", middleware.EmployeeRoleMiddleware(types.RoleAdmin), handler.CreateSKU)
		router.PUT("/:id", middleware.EmployeeRoleMiddleware(types.RoleAdmin), handler.UpdateSKU)
		router.DELETE("/:id", middleware.EmployeeRoleMiddleware(types.RoleAdmin), handler.DeleteSKU)
		router.PATCH("/:id/deactivate", middleware.EmployeeRoleMiddleware(types.RoleAdmin), handler.DeactivateSKU)
	}
}

func (r *Router) RegisterBarcodeRouter(handler *barcode.BarcodeHandler) {
	router := r.Routes.Group("/barcode")
	{
		router.POST("/generate", handler.GenerateBarcode)
		router.GET("/:barcode", handler.RetrieveSKUByBarcode)
		router.POST("/print", handler.PrintAdditionalBarcodes)
	}
}
