package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
)

type StockRequestsModule struct {
	*common.BaseModule
	Repo    stockRequests.StockRequestRepository
	Service stockRequests.StockRequestService
	Handler *stockRequests.StockRequestHandler
}

func NewStockRequestsModule(base *common.BaseModule, stockMaterialRepo stockMaterial.StockMaterialRepository, notificationService notifications.NotificationService) *StockRequestsModule {
	repo := stockRequests.NewStockRequestRepository(base.DB)
	service := stockRequests.NewStockRequestService(repo, stockMaterialRepo, notificationService)
	handler := stockRequests.NewStockRequestHandler(service)

	base.Router.RegisterStockRequestRoutes(handler)

	return &StockRequestsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
