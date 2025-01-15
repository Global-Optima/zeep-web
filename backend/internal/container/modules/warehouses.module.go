package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/inventory"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage"
)

type WarehousesModule struct {
	*common.BaseModule
	Repo    warehouse.WarehouseRepository
	Service warehouse.WarehouseService
	Handler *warehouse.WarehouseHandler
}

func NewWarehousesModule(base *common.BaseModule, stockMaterialRepo stockMaterial.StockMaterialRepository, barcodeRepo barcode.BarcodeRepository, packageRepo stockMaterialPackage.StockMaterialPackageRepository) *WarehousesModule {
	repo := warehouse.NewWarehouseRepository(base.DB)
	inventoryRepo := inventory.NewInventoryRepository(base.DB)
	inventoryService := inventory.NewInventoryService(inventoryRepo, stockMaterialRepo, barcodeRepo, packageRepo)
	inventoryHandler := inventory.NewInventoryHandler(inventoryService)
	service := warehouse.NewWarehouseService(repo)
	handler := warehouse.NewWarehouseHandler(service)

	base.Router.RegisterWarehouseRoutes(handler, inventoryHandler)
	base.Router.RegisterCommonWarehousesRoutes(handler)

	return &WarehousesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
