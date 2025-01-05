package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product"
)

type OrdersModule struct {
	*common.BaseModule
	Repo    orders.OrderRepository
	Service orders.OrderService
	Handler *orders.OrderHandler
}

func NewOrdersModule(base *common.BaseModule, productRepo product.ProductRepository, additiveRepo additives.AdditiveRepository) *OrdersModule {
	repo := orders.NewOrderRepository(base.DB)
	service := orders.NewOrderService(repo, productRepo, additiveRepo, base.Logger)
	handler := orders.NewOrderHandler(service)

	base.Router.RegisterOrderRoutes(handler)

	return &OrdersModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
