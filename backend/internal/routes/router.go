package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/middleware"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
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
		router.POST("", middleware.RoleMiddleware("Admin"), handler.CreateStore)
		router.PUT("/:id", middleware.RoleMiddleware("Admin"), handler.UpdateStore)
		router.DELETE("/:id", middleware.RoleMiddleware("Admin"), handler.DeleteStore)
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
		router.POST("", middleware.RoleMiddleware("Director"), handler.CreateEmployee)
		router.GET("", middleware.RoleMiddleware("Director", "Manager"), handler.GetEmployeesByStore)
		router.GET("/:id", handler.GetEmployeeByID)
		router.PUT("/:id", middleware.RoleMiddleware("Director"), handler.UpdateEmployee)
		router.DELETE("/:id", middleware.RoleMiddleware("Director"), handler.DeleteEmployee)
		router.GET("/roles", middleware.RoleMiddleware("Director", "Manager"), handler.GetAllRoles)
		router.PUT("/:id/password", middleware.RoleMiddleware("Employee", "Manager", "Director"), handler.UpdatePassword)
		router.POST("/login", handler.EmployeeLogin)
	}
}
