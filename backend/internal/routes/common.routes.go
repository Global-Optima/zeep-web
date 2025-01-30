package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
)

func (r *Router) RegisterAuthenticationRoutes(handler *auth.AuthenticationHandler) {
	router := r.CommonRoutes.Group("/auth")
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

func (r *Router) RegisterEmployeeAccountRoutes(handler *employees.EmployeeHandler) {
	router := r.CommonRoutes.Group("/auth/employees")
	{
		router.GET("/store/:id", handler.GetStoreAccounts)
		router.GET("region/:id", handler.GetRegionAccounts)
		router.GET("franchisee/:id", handler.GetFranchiseeAccounts)
		router.GET("/warehouse/:id", handler.GetWarehouseAccounts)
		router.GET("/admins", handler.GetAdminAccounts)
	}
}

func (r *Router) RegisterCommonStoresRoutes(handler *stores.StoreHandler) {
	router := r.CommonRoutes.Group("/stores")
	{
		router.GET("", handler.GetAllStores)
	}
}

func (r *Router) RegisterCommonWarehousesRoutes(handler *warehouse.WarehouseHandler) {
	router := r.CommonRoutes.Group("/warehouses")
	{
		router.GET("", handler.GetAllWarehouses)
	}
}
