package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
)

type StockMaterialsModule struct {
	*common.BaseModule
	Repo    stockMaterial.StockMaterialRepository
	Service stockMaterial.StockMaterialService
	Handler *stockMaterial.StockMaterialHandler
}

func NewStockMaterialsModule(base *common.BaseModule) *StockMaterialsModule {
	repo := stockMaterial.NewStockMaterialRepository(base.DB)
	service := stockMaterial.NewStockMaterialService(repo)
	handler := stockMaterial.NewStockMaterialHandler(service)

	base.Router.RegisterStockMaterialRoutes(handler)

	return &StockMaterialsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
