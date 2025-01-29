package scheduler

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock"
	"go.uber.org/zap"
)

type WarehouseStockCronTasks struct {
	warehouseRepo         warehouse.WarehouseRepository
	warehouseStockService warehouseStock.WarehouseStockService
	warehouseStockRepo    warehouseStock.WarehouseStockRepository
	logger                *zap.SugaredLogger
}

func NewWarehouseStockCronTasks(warehouseRepo warehouse.WarehouseRepository, warehouseStockService warehouseStock.WarehouseStockService, warehouseStockRepo warehouseStock.WarehouseStockRepository, logger *zap.SugaredLogger) *WarehouseStockCronTasks {
	return &WarehouseStockCronTasks{
		warehouseRepo:         warehouseRepo,
		warehouseStockService: warehouseStockService,
		warehouseStockRepo:    warehouseStockRepo,
		logger:                logger,
	}
}

func (tasks *WarehouseStockCronTasks) CheckWarehouseStockNotifications() {
	tasks.logger.Info("Running CheckWarehouseStockNotifications...")

	warehouses, err := tasks.warehouseRepo.GetAllWarehousesForNotifications()
	if err != nil {
		tasks.logger.Errorf("Failed to fetch warehouses: %v", err)
		return
	}

	processedStocks := make(map[uint]bool)

	for _, warehouse := range warehouses {
		tasks.logger.Infof("Processing warehouse ID: %d, Name: %s", warehouse.ID, warehouse.Name)

		stocks, err := tasks.warehouseStockRepo.GetWarehouseStocksForNotifications(warehouse.ID)
		if err != nil {
			tasks.logger.Errorf("Failed to fetch stock list for warehouse ID %d: %v", warehouse.ID, err)
			continue
		}

		for _, stock := range stocks {
			if processedStocks[stock.StockMaterialID] {
				tasks.logger.Infof("Skipping duplicate stock material ID: %d", stock.StockMaterialID)
				continue
			}

			err = tasks.warehouseStockService.CheckStockNotifications(warehouse.ID, stock)
			if err != nil {
				tasks.logger.Errorf("Failed to check notifications for stock ID %d in warehouse ID %d: %v", stock.StockMaterialID, warehouse.ID, err)
			}

			processedStocks[stock.StockMaterialID] = true
		}
	}

	tasks.logger.Info("CheckWarehouseStockNotifications completed.")
}
