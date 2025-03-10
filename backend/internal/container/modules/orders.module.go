package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/asynqTasks"
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
)

type OrdersModule struct {
	*common.BaseModule
	Repo    orders.OrderRepository
	Service orders.OrderService
	Handler *orders.OrderHandler
}

func NewOrdersModule(base *common.BaseModule, asynqManager *asynqTasks.AsynqManager, productRepo storeProducts.StoreProductRepository, additiveRepo storeAdditives.StoreAdditiveRepository, storeStockRepo storeStocks.StoreStockRepository, notificationService notifications.NotificationService) *OrdersModule {
	repo := orders.NewOrderRepository(base.DB)
	service := orders.NewOrderService(asynqManager, repo, productRepo, additiveRepo, storeStockRepo, notificationService, base.Logger)
	handler := orders.NewOrderHandler(service)

	ordersAsynqTasks := asynqTasks.NewOrderAsynqTasks(repo, base.Logger)
	asynqManager.RegisterTask(orders.OrderPaymentFailure, ordersAsynqTasks.HandleOrderPaymentFailureTask)
	base.Router.RegisterOrderRoutes(handler)

	return &OrdersModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
