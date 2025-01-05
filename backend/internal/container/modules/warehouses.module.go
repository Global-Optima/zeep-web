package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
)

type WarehousesModule struct {
	*common.BaseModule
	Repo    warehouse.WarehouseRepository
	Service warehouse.WarehouseService
	Handler *warehouse.WarehouseHandler
}

func NewWarehousesModule(base *common.BaseModule) *WarehousesModule {
	repo := warehouse.NewWarehouseRepository(base.DB)
	service := warehouse.NewWarehouseService(repo)
	handler := warehouse.NewWarehouseHandler(service)

	base.Router.RegisterWarehouseRoutes(handler)

	return &WarehousesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
