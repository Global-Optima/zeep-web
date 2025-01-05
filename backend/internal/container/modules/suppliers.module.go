package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/supplier"
)

type SuppliersModule struct {
	*common.BaseModule
	Repo    supplier.SupplierRepository
	Service supplier.SupplierService
	Handler *supplier.SupplierHandler
}

func NewSuppliersModule(base *common.BaseModule) *SuppliersModule {
	repo := supplier.NewSupplierRepository(base.DB)
	service := supplier.NewSupplierService(repo)
	handler := supplier.NewSupplierHandler(service)

	base.Router.RegisterSupplierRoutes(handler)

	return &SuppliersModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
