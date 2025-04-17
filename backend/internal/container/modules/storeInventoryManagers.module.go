package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
)

type StoreInventoryManagerModule struct {
	*common.BaseModule
	Repo    storeInventoryManagers.StoreInventoryManagerRepository
	Service storeInventoryManagers.StoreInventoryManagerService
	Handler *storeInventoryManagers.StoreInventoryManagerHandler
}

func NewStoreInventoryManagersModule(
	base *common.BaseModule,
	notificationService notifications.NotificationService,
) *StoreInventoryManagerModule {
	repo := storeInventoryManagers.NewStoreInventoryManagerRepository(base.DB)
	service := storeInventoryManagers.NewStoreInventoryManagerService(
		repo,
		notificationService,
		base.Logger,
	)
	handler := storeInventoryManagers.NewStoreInventoryManagerHandler(service)

	base.Router.RegisterStoreInventoryManagerRoutes(handler)

	return &StoreInventoryManagerModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
