package scheduler

import (
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock/types"
	"go.uber.org/zap"
)

type WarehouseStockCronTasks struct {
	warehouseRepo         warehouse.WarehouseRepository
	warehouseStockService warehouseStock.WarehouseStockService
	warehouseStockRepo    warehouseStock.WarehouseStockRepository
	logger                *zap.SugaredLogger
}

func NewWarehouseStockCronTasks(warehouseRepo warehouse.WarehouseRepository, warehouseStockService warehouseStock.WarehouseStockService, logger *zap.SugaredLogger) *WarehouseStockCronTasks {
	return &WarehouseStockCronTasks{
		warehouseRepo:         warehouseRepo,
		warehouseStockService: warehouseStockService,
		logger:                logger,
	}
}

func (tasks *WarehouseStockCronTasks) CheckWarehouseStockNotifications() {
	tasks.logger.Info("Running CheckWarehouseStockNotifications...")

	warehouses, err := tasks.warehouseRepo.GetAllWarehouses(nil)
	if err != nil {
		tasks.logger.Errorf("Failed to fetch warehouses: %v", err)
		return
	}

	for _, warehouse := range warehouses {
		stocks, err := tasks.warehouseStockRepo.GetWarehouseStock(&types.GetWarehouseStockFilterQuery{WarehouseID: &warehouse.ID})
		if err != nil {
			tasks.logger.Errorf("Failed to fetch stock list for warehouse %d: %v", warehouse.ID, err)
			continue
		}

		for _, stock := range stocks {
			modelStock, err := tasks.warehouseStockRepo.GetWarehouseStockByID(warehouse.ID, stock.StockMaterialID)
			if err != nil {
				tasks.logger.Errorf("Failed to fetch stock model for stock ID %d: %v", stock.StockMaterialID, err)
				continue
			}
			err = tasks.warehouseStockService.CheckStockNotifications(warehouse.ID, *modelStock)
			if err != nil {
				tasks.logger.Errorf("Failed to check notifications for stock ID %d: %v", stock.StockMaterialID, err)
			}
		}
	}

	tasks.logger.Info("CheckWarehouseStockNotifications completed.")
}
