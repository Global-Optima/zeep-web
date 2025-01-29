package scheduler

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeWarehouses"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
	"go.uber.org/zap"
)

type StoreWarehouseCronTasks struct {
	storeWarehouseService storeWarehouses.StoreWarehouseService
	storeWarehouseRepo    storeWarehouses.StoreWarehouseRepository
	storeService          stores.StoreService
	logger                *zap.SugaredLogger
}

func NewStoreWarehouseCronTasks(storeWarehouseService storeWarehouses.StoreWarehouseService, storeWarehouseRepo storeWarehouses.StoreWarehouseRepository, storeService stores.StoreService, logger *zap.SugaredLogger) *StoreWarehouseCronTasks {
	return &StoreWarehouseCronTasks{
		storeWarehouseService: storeWarehouseService,
		storeWarehouseRepo:    storeWarehouseRepo,
		storeService:          storeService,
		logger:                logger,
	}
}

func (tasks *StoreWarehouseCronTasks) CheckStockNotifications() {
	tasks.logger.Info("Running CheckStockNotifications...")

	stores, err := tasks.storeService.GetAllStoresForNotifications()
	if err != nil {
		tasks.logger.Errorf("Failed to fetch stores: %v", err)
		return
	}

	processedStocks := make(map[uint]bool)

	for _, store := range stores {
		stockList, err := tasks.storeWarehouseRepo.GetAllStockList(store.ID)
		if err != nil {
			tasks.logger.Errorf("Failed to fetch stock list for store %d: %v", store.ID, err)
			continue
		}

		for _, stock := range stockList {
			if processedStocks[stock.IngredientID] {
				tasks.logger.Infof("Skipping duplicate IngredientID: %d", stock.IngredientID)
				continue
			}

			err := tasks.storeWarehouseService.CheckStockNotifications(store.ID, stock)
			if err != nil {
				tasks.logger.Errorf("Failed to check notifications for stock ID %d: %v", stock.ID, err)
			}

			processedStocks[stock.IngredientID] = true
		}
	}

	tasks.logger.Info("Check Stock Notifications completed.")
}
