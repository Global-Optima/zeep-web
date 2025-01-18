package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
)

type BarcodeModule struct {
	*common.BaseModule
	Repo    barcode.BarcodeRepository
	Service barcode.BarcodeService
	Handler *barcode.BarcodeHandler
}

func NewBarcodeModule(base *common.BaseModule, stockMaterialRepo stockMaterial.StockMaterialRepository) *BarcodeModule {
	repo := barcode.NewBarcodeRepository(base.DB)
	printerService := barcode.NewPrinterService()
	service := barcode.NewBarcodeService(repo, stockMaterialRepo, printerService)
	handler := barcode.NewBarcodeHandler(service)

	base.Router.RegisterBarcodeRoutes(handler)

	return &BarcodeModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
