package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock"
	"github.com/Global-Optima/zeep-web/backend/internal/scheduler"
)

type WarehousesModule struct {
	*common.BaseModule
	Repo    warehouse.WarehouseRepository
	Service warehouse.WarehouseService
	Handler *warehouse.WarehouseHandler
}

func NewWarehousesModule(
	base *common.BaseModule,
	stockMaterialRepo stockMaterial.StockMaterialRepository,
	notificationService notifications.NotificationService,
	cronManager *scheduler.CronManager,
	regionService regions.RegionService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
) *WarehousesModule {
	warehouseRepo := warehouse.NewWarehouseRepository(base.DB)
	warehouseStockRepo := warehouseStock.NewWarehouseStockRepository(base.DB)

	warehouseStockService := warehouseStock.NewWarehouseStockService(warehouseStockRepo, stockMaterialRepo, notificationService, base.Logger)
	warehouseService := warehouse.NewWarehouseService(warehouseRepo, base.Logger)

	warehouseHandler := warehouse.NewWarehouseHandler(warehouseService, regionService, auditService)
	warehouseStockHandler := warehouseStock.NewWarehouseStockHandler(warehouseStockService, franchiseeService, warehouseService, auditService)

	base.Router.RegisterWarehouseRoutes(warehouseHandler, warehouseStockHandler)
	base.Router.RegisterCommonWarehousesRoutes(warehouseHandler)

	warehouseStockCronTasks := scheduler.NewWarehouseStockCronTasks(warehouseRepo, warehouseStockService, warehouseStockRepo, base.Logger)
	err := cronManager.RegisterJob(scheduler.DailyJob, func() {
		warehouseStockCronTasks.CheckWarehouseStockNotifications()
	})

	if err != nil {
		base.Logger.Errorf("Failed to register warehouse stock cron job: %v", err)
	}

	return &WarehousesModule{
		BaseModule: base,
		Repo:       warehouseRepo,
		Service:    warehouseService,
		Handler:    warehouseHandler,
	}
}
