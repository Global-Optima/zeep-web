package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers"
)

type CustomersModule struct {
	*common.BaseModule
	Repo    customers.CustomerRepository
	Service customers.CustomerService
	Handler *customers.CustomerHandler
}

func NewCustomersModule(base *common.BaseModule) *CustomersModule {
	repo := customers.NewCustomerRepository(base.DB)
	service := customers.NewCustomerService(repo, base.Logger)
	handler := customers.NewCustomerHandler(service)

	//base.Router.RegisterCustomerRoutes(handler)

	return &CustomersModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
