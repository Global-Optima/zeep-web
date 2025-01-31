package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
)

type OrdersModule struct {
	*common.BaseModule
	Repo    orders.OrderRepository
	Service orders.OrderService
	Handler *orders.OrderHandler
}

func NewOrdersModule(base *common.BaseModule, productRepo storeProducts.StoreProductRepository, additiveRepo storeAdditives.StoreAdditiveRepository, storeWarehouseRepo storeWarehouses.StoreWarehouseRepository, notificationService notifications.NotificationService) *OrdersModule {
	repo := orders.NewOrderRepository(base.DB)
	service := orders.NewOrderService(repo, productRepo, additiveRepo, storeWarehouseRepo, notificationService, base.Logger)
	handler := orders.NewOrderHandler(service)

	base.Router.RegisterOrderRoutes(handler)

	return &OrdersModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
