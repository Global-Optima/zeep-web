package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/employeeToken"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
)

type AuthModule struct {
	*common.BaseModule
	Repo              auth.AuthenticationRepository
	EmployeeTokenRepo employeeToken.EmployeeTokenRepository
	Service           auth.AuthenticationService
	Handler           *auth.AuthenticationHandler
}

func NewAuthModule(
	base *common.BaseModule,
	customersRepo customers.CustomerRepository,
	employeesRepo employees.EmployeeRepository,
) *AuthModule {
	employeeTokenRepo := employeeToken.NewEmployeeTokenRepository(base.DB)
	repo := auth.NewAuthenticationRepository(base.DB)
	service := auth.NewAuthenticationService(repo, customersRepo, employeesRepo, employeeTokenRepo, base.Logger)
	handler := auth.NewAuthenticationHandler(service)

	base.Router.RegisterAuthenticationRoutes(handler)

	return &AuthModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
