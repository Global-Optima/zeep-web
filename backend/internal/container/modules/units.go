package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/units"
)

type UnitsModule struct {
	*common.BaseModule
	Repo    units.UnitRepository
	Service units.UnitService
	Handler *units.UnitHandler
}

func NewUnitsModule(base *common.BaseModule, auditService audit.AuditService) *UnitsModule {
	repo := units.NewUnitRepository(base.DB)
	service := units.NewUnitService(repo)
	handler := units.NewUnitHandler(service, auditService)

	base.Router.RegisterUnitRoutes(handler)

	return &UnitsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
