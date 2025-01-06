package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialPackage"
)

type StockMaterialPackagesModule struct {
	*common.BaseModule
	Repo    stockMaterialPackage.StockMaterialPackageRepository
	Service stockMaterialPackage.StockMaterialPackageService
	Handler *stockMaterialPackage.StockMaterialPackageHandler
}

func NewStockMaterialPackagesModule(base *common.BaseModule) *StockMaterialPackagesModule {
	repo := stockMaterialPackage.NewStockMaterialPackageRepository(base.DB)
	service := stockMaterialPackage.NewStockMaterialPackageService(repo)
	handler := stockMaterialPackage.NewStockMaterialPackageHandler(service)

	base.Router.RegisterStockMaterialPackageRoutes(handler)

	return &StockMaterialPackagesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
