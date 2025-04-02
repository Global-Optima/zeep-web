package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/scheduler"
)

type StoreStockModule struct {
	Repo    storeStocks.StoreStockRepository
	Service storeStocks.StoreStockService
	Handler *storeStocks.StoreStockHandler
}

func NewStoreStockModule(
	base *common.BaseModule,
	ingredientService ingredients.IngredientService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	notificationService notifications.NotificationService,
	storeService stores.StoreService,
	cronManager *scheduler.CronManager,
) *StoreStockModule {
	repo := storeStocks.NewStoreStockRepository(base.DB)
	service := storeStocks.NewStoreStockService(repo, notificationService, base.Logger)
	handler := storeStocks.NewStoreStockHandler(service, ingredientService, auditService, franchiseeService, base.Logger)

	base.Router.RegisterStoreWarehouseRoutes(handler)

	storeWarehouseCronTasks := scheduler.NewStoreStockCronTasks(service, repo, storeService, base.Logger)
	err := cronManager.RegisterJob(scheduler.DailyJob, func() {
		storeWarehouseCronTasks.CheckStockNotifications()
	})
	if err != nil {
		base.Logger.Errorf("Failed to register warehouse stock cron job: %v", err)
	}

	return &StoreStockModule{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}
