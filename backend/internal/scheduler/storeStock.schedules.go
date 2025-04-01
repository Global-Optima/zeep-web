package scheduler

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"go.uber.org/zap"
)

type StoreStockCronTasks struct {
	storeStockService storeStocks.StoreStockService
	storeStockRepo    storeStocks.StoreStockRepository
	storeService      stores.StoreService
	logger            *zap.SugaredLogger
}

func NewStoreStockCronTasks(storeStockService storeStocks.StoreStockService, storeStockRepo storeStocks.StoreStockRepository, storeService stores.StoreService, logger *zap.SugaredLogger) *StoreStockCronTasks {
	return &StoreStockCronTasks{
		storeStockService: storeStockService,
		storeStockRepo:    storeStockRepo,
		storeService:      storeService,
		logger:            logger,
	}
}

func (tasks *StoreStockCronTasks) CheckStockNotifications() {
	tasks.logger.Info("Running CheckStockNotifications...")

	storesList, err := tasks.storeService.GetAllStoresForNotifications()
	if err != nil {
		tasks.logger.Errorf("Failed to fetch stores: %v", err)
		return
	}

	for _, store := range storesList {
		processedStocks := make(map[uint]bool)

		stockList, err := tasks.storeStockRepo.GetAllStockList(store.ID)
		if err != nil {
			tasks.logger.Errorf("Failed to fetch stock list for store %d: %v", store.ID, err)
			continue
		}

		for _, stock := range stockList {
			if processedStocks[stock.IngredientID] {
				tasks.logger.Infof("Skipping duplicate IngredientID: %d", stock.IngredientID)
				continue
			}

			err := tasks.storeStockService.CheckStockNotifications(store.ID, stock)
			if err != nil {
				tasks.logger.Errorf("Failed to check notifications for stock ID %d: %v", stock.ID, err)
			}

			processedStocks[stock.IngredientID] = true
		}
	}

	tasks.logger.Info("Check Stock Notifications completed.")
}
