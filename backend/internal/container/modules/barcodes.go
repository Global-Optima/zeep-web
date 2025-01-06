package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/barcode"
)

type BarcodeModule struct {
	*common.BaseModule
	Repo    barcode.BarcodeRepository
	Service barcode.BarcodeService
	Handler *barcode.BarcodeHandler
}

func NewBarcodeModule(base *common.BaseModule, stockMaterialModule StockMaterialsModule) *BarcodeModule {
	repo := barcode.NewBarcodeRepository(base.DB)
	printerService := barcode.NewPrinterService()
	service := barcode.NewBarcodeService(repo, stockMaterialModule.Repo, printerService)
	handler := barcode.NewBarcodeHandler(service)

	base.Router.RegisterBarcodeRoutes(handler)

	return &BarcodeModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
