package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stockRequests"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
)

type StockRequestsModule struct {
	*common.BaseModule
	Repo    stockRequests.StockRequestRepository
	Service stockRequests.StockRequestService
	Handler *stockRequests.StockRequestHandler
}

func NewStockRequestsModule(
	base *common.BaseModule,
	franchiseeService franchisees.FranchiseeService,
	regionService regions.RegionService,
	stockMaterialRepo stockMaterial.StockMaterialRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	notificationService notifications.NotificationService,
	auditService audit.AuditService,
) *StockRequestsModule {
	repo := stockRequests.NewStockRequestRepository(base.DB)
	service := stockRequests.NewStockRequestService(
		repo,
		stockMaterialRepo,
		storeInventoryManagerRepo,
		stockRequests.NewTransactionManager(base.DB, repo, stockMaterialRepo),
		notificationService,
		base.Logger,
	)
	handler := stockRequests.NewStockRequestHandler(service, franchiseeService, regionService, auditService)

	base.Router.RegisterStockRequestRoutes(handler)

	return &StockRequestsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
