package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
)

type RegionsModule struct {
	*common.BaseModule
	Repo    regions.RegionRepository
	Service regions.RegionService
	Handler *regions.RegionHandler
}

func NewRegionsModule(base *common.BaseModule, auditService audit.AuditService) *RegionsModule {
	repo := regions.NewRegionRepository(base.DB)
	service := regions.NewRegionService(repo, base.Logger)
	handler := regions.NewRegionHandler(service, auditService)

	base.Router.RegisterRegionRoutes(handler)
	base.Router.RegisterCommonRegionsRoutes(handler)

	return &RegionsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
