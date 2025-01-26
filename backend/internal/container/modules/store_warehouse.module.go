package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
)

type StoreWarehouseModule struct {
	Repo    storeWarehouses.StoreWarehouseRepository
	Service storeWarehouses.StoreWarehouseService
	Handler *storeWarehouses.StoreWarehouseHandler
}

func NewStoreWarehouseModule(base *common.BaseModule) *StoreWarehouseModule {
	repo := storeWarehouses.NewStoreWarehouseRepository(base.DB)
	service := storeWarehouses.NewStoreWarehouseService(repo, base.Logger)
	handler := storeWarehouses.NewStoreWarehouseHandler(service)

	base.Router.RegisterStoreWarehouseRoutes(handler)

	return &StoreWarehouseModule{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}
