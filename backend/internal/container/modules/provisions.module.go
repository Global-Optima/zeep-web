package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/technicalMap"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"github.com/Global-Optima/zeep-web/backend/internal/scheduler"
)

type ProvisionsModule struct {
	*common.BaseModule
	Repo                         provisions.ProvisionRepository
	Service                      provisions.ProvisionService
	Handler                      *provisions.ProvisionHandler
	StoreProvisionsModule        *StoreProvisionsModule
	ProvisionsTechnicalMapModule *ProvisionsTechnicalMapModule
}

func NewProvisionsModule(
	base *common.BaseModule,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	storeService stores.StoreService,
	notificationService notifications.NotificationService,
	ingredientRepo ingredients.IngredientRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	cronManager *scheduler.CronManager,
) *ProvisionsModule {
	repo := provisions.NewProvisionRepository(base.DB)
	service := provisions.NewProvisionService(repo, base.Logger)
	handler := provisions.NewProvisionHandler(service, auditService, base.Logger)

	storeProvisionsModule := NewStoreProvisionsModule(
		base,
		service,
		franchiseeService,
		storeService,
		auditService,
		notificationService,
		repo,
		ingredientRepo,
		storeStockRepo,
		storeInventoryManagerRepo,
		cronManager,
	)
	provisionTechMapModule := NewProvisionsTechnicalMapModule(base)

	base.Router.RegisterProvisionsRoutes(handler, provisionTechMapModule.Handler)

	return &ProvisionsModule{
		BaseModule:                   base,
		Repo:                         repo,
		Service:                      service,
		Handler:                      handler,
		StoreProvisionsModule:        storeProvisionsModule,
		ProvisionsTechnicalMapModule: provisionTechMapModule,
	}
}

func NewStoreProvisionsModule(
	base *common.BaseModule,
	provisionService provisions.ProvisionService,
	franchiseeService franchisees.FranchiseeService,
	storeService stores.StoreService,
	auditService audit.AuditService,
	notificationService notifications.NotificationService,
	provisionRepo provisions.ProvisionRepository,
	ingredientRepo ingredients.IngredientRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	cronManager *scheduler.CronManager,
) *StoreProvisionsModule {
	repo := storeProvisions.NewStoreProvisionRepository(base.DB)
	service := storeProvisions.NewStoreProvisionService(
		repo,
		ingredientRepo,
		provisionRepo,
		storeInventoryManagerRepo,
		storeProvisions.NewTransactionManager(base.DB, repo, storeStockRepo, ingredientRepo, storeInventoryManagerRepo),
		base.Logger,
	)
	handler := storeProvisions.NewStoreProvisionHandler(service, provisionService, franchiseeService, auditService, base.Logger)

	base.Router.RegisterStoreProvisionsRoutes(handler)

	storeProvisionCronTasks := scheduler.NewStoreProvisionCronTasks(
		service,
		repo,
		storeInventoryManagerRepo,
		storeService,
		notificationService,
		base.Logger,
	)
	err := cronManager.RegisterJob(scheduler.HalfHourlyJob, func() {
		storeProvisionCronTasks.CheckStoreProvisionNotifications()
	})
	if err != nil {
		base.Logger.Errorf("Failed to register warehouse stock cron job: %v", err)
	}

	return &StoreProvisionsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type StoreProvisionsModule struct {
	*common.BaseModule
	Repo    storeProvisions.StoreProvisionRepository
	Service storeProvisions.StoreProvisionService
	Handler *storeProvisions.StoreProvisionHandler
}

type ProvisionsTechnicalMapModule struct {
	*common.BaseModule
	Repo    technicalMap.TechnicalMapRepository
	Service technicalMap.TechnicalMapService
	Handler *technicalMap.TechnicalMapHandler
}

func NewProvisionsTechnicalMapModule(base *common.BaseModule) *ProvisionsTechnicalMapModule {
	repo := technicalMap.NewTechnicalMapRepository(base.DB)
	service := technicalMap.NewTechnicalMapService(repo, base.Logger)
	handler := technicalMap.NewTechnicalMapHandler(service)

	return &ProvisionsTechnicalMapModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
