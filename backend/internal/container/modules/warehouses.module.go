package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/warehouseStock"
)

type WarehousesModule struct {
	*common.BaseModule
	Repo    warehouse.WarehouseRepository
	Service warehouse.WarehouseService
	Handler *warehouse.WarehouseHandler
}

func NewWarehousesModule(base *common.BaseModule, stockMaterialRepo stockMaterial.StockMaterialRepository, barcodeRepo barcode.BarcodeRepository) *WarehousesModule {
	repo := warehouse.NewWarehouseRepository(base.DB)
	warehouseStockRepo := warehouseStock.NewWarehouseStockRepository(base.DB)
	warehouseStockService := warehouseStock.NewWarehouseStockService(warehouseStockRepo, stockMaterialRepo, barcodeRepo)
	warehouseStockHandler := warehouseStock.NewWarehouseStockHandler(warehouseStockService)
	service := warehouse.NewWarehouseService(repo)
	handler := warehouse.NewWarehouseHandler(service)

	base.Router.RegisterWarehouseRoutes(handler, warehouseStockHandler)
	base.Router.RegisterCommonWarehousesRoutes(handler)

	return &WarehousesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
