package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/scheduler"
)

type StoreWarehouseModule struct {
	Repo    storeWarehouses.StoreWarehouseRepository
	Service storeWarehouses.StoreWarehouseService
	Handler *storeWarehouses.StoreWarehouseHandler
}

func NewStoreWarehouseModule(base *common.BaseModule,
	ingredientService ingredients.IngredientService,
	auditService audit.AuditService,
	notificationService notifications.NotificationService,
	storeService stores.StoreService,
	cronManager *scheduler.CronManager) *StoreWarehouseModule {
	repo := storeWarehouses.NewStoreWarehouseRepository(base.DB)
	service := storeWarehouses.NewStoreWarehouseService(repo, notificationService, base.Logger)
	handler := storeWarehouses.NewStoreWarehouseHandler(service, ingredientService, auditService, base.Logger)

	base.Router.RegisterStoreWarehouseRoutes(handler)

	storeWarehouseCronTasks := scheduler.NewStoreWarehouseCronTasks(service, repo, storeService, base.Logger)
	err := cronManager.RegisterJob(scheduler.DailyJob, func() {
		storeWarehouseCronTasks.CheckStockNotifications()
	})

	if err != nil {
		base.Logger.Errorf("Failed to register warehouse stock cron job: %v", err)
	}

	return &StoreWarehouseModule{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}
