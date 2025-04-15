package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeInventoryManagers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeSynchronizers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
)

type StoreSynchronizerModule struct {
	*common.BaseModule
	Repo    storeSynchronizers.StoreSynchronizeRepository
	Service storeSynchronizers.StoreSynchronizeService
	Handler *storeSynchronizers.StoreSynchronizerHandler
}

func NewStoreSynchronizerSynchronizerModule(
	base *common.BaseModule,
	storeRepo stores.StoreRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	ingredientRepo ingredients.IngredientRepository,
	storeInventoryManagerRepo storeInventoryManagers.StoreInventoryManagerRepository,
) *StoreSynchronizerModule {
	repo := storeSynchronizers.NewStoreSynchronizeRepository(base.DB)
	service := storeSynchronizers.NewStoreSynchronizeService(
		repo,
		storeSynchronizers.NewTransactionManager(base.DB, repo, storeRepo, storeAdditiveRepo, storeStockRepo, ingredientRepo, storeInventoryManagerRepo),
		base.Logger,
	)
	handler := storeSynchronizers.NewStoreSynchronizeHandler(service)

	base.Router.RegisterStoreSynchronizerSynchronizerRoutes(handler)

	return &StoreSynchronizerModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
