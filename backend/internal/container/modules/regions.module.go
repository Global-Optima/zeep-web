package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
)

type RegionsModule struct {
	*common.BaseModule
	Repo    regions.RegionRepository
	Service regions.RegionService
	Handler *regions.RegionHandler
}

func NewRegionsModule(base *common.BaseModule) *RegionsModule {
	repo := regions.NewRegionRepository(base.DB)
	service := regions.NewRegionService(repo)
	handler := regions.NewRegionHandler(service)

	base.Router.RegisterRegionRoutes(handler)

	return &RegionsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
