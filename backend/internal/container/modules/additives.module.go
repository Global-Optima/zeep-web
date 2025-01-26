package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
)

type AdditivesModule struct {
	*common.BaseModule
	Repo                 additives.AdditiveRepository
	Service              additives.AdditiveService
	Handler              *additives.AdditiveHandler
	StoreAdditivesModule *StoreAdditivesModule
}

func NewAdditivesModule(base *common.BaseModule, ingredientRepo ingredients.IngredientRepository, storeWarehouseRepo storeWarehouses.StoreWarehouseRepository) *AdditivesModule {
	repo := additives.NewAdditiveRepository(base.DB)
	service := additives.NewAdditiveService(repo, base.Logger)
	handler := additives.NewAdditiveHandler(service)

	storeAdditivesModule := NewStoreAdditivesModule(base, ingredientRepo, storeWarehouseRepo)

	base.Router.RegisterAdditivesRoutes(handler)

	return &AdditivesModule{
		BaseModule:           base,
		Repo:                 repo,
		Service:              service,
		Handler:              handler,
		StoreAdditivesModule: storeAdditivesModule,
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
	ingredientRepo ingredients.IngredientRepository,
	storeWarehouseRepo storeWarehouses.StoreWarehouseRepository,
) *StoreAdditivesModule {
	repo := storeAdditives.NewStoreAdditiveRepository(base.DB)
	service := storeAdditives.NewStoreAdditiveService(
		repo,
		ingredientRepo,
		storeAdditives.NewTransactionManager(base.DB, repo, storeWarehouseRepo),
		base.Logger)
	handler := storeAdditives.NewStoreAdditiveHandler(service)

	base.Router.RegisterStoreAdditivesRoutes(handler)

	return &StoreAdditivesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
