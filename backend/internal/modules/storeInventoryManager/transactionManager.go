package storeInventoryManager

import (
	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"gorm.io/gorm"
)

type TransactionManager interface {
}

type transactionManager struct {
	db                 *gorm.DB
	repo               StoreInventoryManagerRepository
	storeRepo          stores.StoreRepository
	storeAdditiveRepo  storeAdditives.StoreAdditiveRepository
	storeStockRepo     storeStocks.StoreStockRepository
	storeProvisionRepo storeProvisions.StoreProvisionRepository
	ingredientRepo     ingredients.IngredientRepository
	provisionRepo      provisions.ProvisionRepository
}

func NewTransactionManager(
	db *gorm.DB,
	repo StoreInventoryManagerRepository,
	storeRepo stores.StoreRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	storeProvisionRepo storeProvisions.StoreProvisionRepository,
	ingredientRepo ingredients.IngredientRepository,
	provisionRepo provisions.ProvisionRepository,
) TransactionManager {
	return &transactionManager{
		db:                 db,
		repo:               repo,
		storeRepo:          storeRepo,
		storeAdditiveRepo:  storeAdditiveRepo,
		storeStockRepo:     storeStockRepo,
		storeProvisionRepo: storeProvisionRepo,
		ingredientRepo:     ingredientRepo,
		provisionRepo:      provisionRepo,
	}
}
