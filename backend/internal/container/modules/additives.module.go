package modules

import (
	"github.com/Global-Optima/zeep-web/backend/api/storage"
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives/technicalMap"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
)

type AdditivesModule struct {
	*common.BaseModule
	Repo                   additives.AdditiveRepository
	Service                additives.AdditiveService
	Handler                *additives.AdditiveHandler
	StoreAdditivesModule   *StoreAdditivesModule
	AdditivesTechMapModule *AdditivesTechMapModule
}

func NewAdditivesModule(
	base *common.BaseModule,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	ingredientRepo ingredients.IngredientRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	storageRepo storage.StorageRepository,
	notificationService notifications.NotificationService,
) *AdditivesModule {
	repo := additives.NewAdditiveRepository(base.DB)
	service := additives.NewAdditiveService(repo, storageRepo, notificationService, base.Logger)
	handler := additives.NewAdditiveHandler(service, auditService)

	storeAdditivesModule := NewStoreAdditivesModule(
		base,
		service,
		franchiseeService,
		auditService,
		ingredientRepo,
		storeStockRepo,
		storeInventoryManagerRepo,
		storageRepo,
	)
	additivesTechMapModule := NewAdditivesTechMapModule(base)

	base.Router.RegisterAdditivesRoutes(handler, additivesTechMapModule.Handler)

	return &AdditivesModule{
		BaseModule:             base,
		Repo:                   repo,
		Service:                service,
		Handler:                handler,
		StoreAdditivesModule:   storeAdditivesModule,
		AdditivesTechMapModule: additivesTechMapModule,
	}
}

type StoreAdditivesModule struct {
	*common.BaseModule
	Repo    storeAdditives.StoreAdditiveRepository
	Service storeAdditives.StoreAdditiveService
	Handler *storeAdditives.StoreAdditiveHandler
}

func NewStoreAdditivesModule(
	base *common.BaseModule,
	additiveService additives.AdditiveService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	ingredientRepo ingredients.IngredientRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
	storageRepo storage.StorageRepository,
) *StoreAdditivesModule {
	repo := storeAdditives.NewStoreAdditiveRepository(base.DB)
	service := storeAdditives.NewStoreAdditiveService(
		repo,
		ingredientRepo,
		storageRepo,
		storeAdditives.NewTransactionManager(base.DB, repo, storeStockRepo, ingredientRepo, storeInventoryManagerRepo),
		base.Logger,
	)

	handler := storeAdditives.NewStoreAdditiveHandler(service, additiveService, franchiseeService, auditService, base.Logger)

	base.Router.RegisterStoreAdditivesRoutes(handler)

	return &StoreAdditivesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type AdditivesTechMapModule struct {
	*common.BaseModule
	Repo    technicalMap.TechnicalMapRepository
	Service technicalMap.TechnicalMapService
	Handler *technicalMap.TechnicalMapHandler
}

func NewAdditivesTechMapModule(base *common.BaseModule) *AdditivesTechMapModule {
	repo := technicalMap.NewTechnicalMapRepository(base.DB)
	service := technicalMap.NewTechnicalMapService(repo, base.Logger)
	handler := technicalMap.NewTechnicalMapHandler(service)

	return &AdditivesTechMapModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
