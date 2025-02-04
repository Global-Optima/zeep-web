package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial"
)

type StockMaterialsModule struct {
	*common.BaseModule
	Repo    stockMaterial.StockMaterialRepository
	Service stockMaterial.StockMaterialService
	Handler *stockMaterial.StockMaterialHandler
}

func NewStockMaterialsModule(base *common.BaseModule, auditService audit.AuditService) *StockMaterialsModule {
	repo := stockMaterial.NewStockMaterialRepository(base.DB)
	service := stockMaterial.NewStockMaterialService(repo)
	handler := stockMaterial.NewStockMaterialHandler(service, auditService)

	base.Router.RegisterStockMaterialRoutes(handler)

	return &StockMaterialsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
