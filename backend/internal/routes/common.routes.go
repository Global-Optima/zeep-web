package routes

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth"
	adminEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/adminEmployees"
	franchiseeEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees"
	regionEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees"
	storeEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees"
	warehouseEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/warehouseEmployees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
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
			customersRoutes.POST("/logout", handler.CustomerLogout)
		}

		employeesRoutes := router.Group("/employees")
		{
			employeesRoutes.POST("/login", handler.EmployeeLogin)
			employeesRoutes.POST("/logout", handler.EmployeeLogout)
		}
	}
}

func (r *Router) RegisterEmployeeAccountRoutes(
	storeEmployeeHandler *storeEmployees.StoreEmployeeHandler,
	warehouseEmployeeHandler *warehouseEmployees.WarehouseEmployeeHandler,
	franchiseeEmployeeHandler *franchiseeEmployees.FranchiseeEmployeeHandler,
	regionEmployeeHandler *regionEmployees.RegionEmployeeHandler,
	adminEmployeeHandler *adminEmployees.AdminEmployeeHandler,
) {
	router := r.CommonRoutes.Group("/auth/employees")
	{
		router.GET("/store/:id", storeEmployeeHandler.GetStoreAccounts)
		router.GET("region/:id", regionEmployeeHandler.GetRegionAccounts)
		router.GET("franchisee/:id", franchiseeEmployeeHandler.GetFranchiseeAccounts)
		router.GET("/warehouse/:id", warehouseEmployeeHandler.GetWarehouseAccounts)
		router.GET("/admins", adminEmployeeHandler.GetAdminAccounts)
	}
}

func (r *Router) RegisterCommonStoresRoutes(handler *stores.StoreHandler) {
	router := r.CommonRoutes.Group("/stores/all")
	{
		router.GET("", handler.GetAllStores)
	}
}

func (r *Router) RegisterCommonWarehousesRoutes(handler *warehouse.WarehouseHandler) {
	router := r.CommonRoutes.Group("/warehouses/all")
	{
		router.GET("", handler.GetAllWarehouses)
	}
}

func (r *Router) RegisterCommonFranchiseesRoutes(handler *franchisees.FranchiseeHandler) {
	router := r.CommonRoutes.Group("/franchisees/all")
	{
		router.GET("", handler.GetAllFranchisees)
	}
}

func (r *Router) RegisterCommonRegionsRoutes(handler *regions.RegionHandler) {
	router := r.CommonRoutes.Group("/regions/all")
	{
		router.GET("", handler.GetAllRegions)
	}
}
