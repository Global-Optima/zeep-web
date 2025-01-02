package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/additives"
)

type AdditivesModule struct {
	*common.BaseModule
	Repo    additives.AdditiveRepository
	Service additives.AdditiveService
	Handler *additives.AdditiveHandler
}

func NewAdditivesModule(base *common.BaseModule) *AdditivesModule {
	repo := additives.NewAdditiveRepository(base.DB)
	service := additives.NewAdditiveService(repo, base.Logger)
	handler := additives.NewAdditiveHandler(service)

	base.Router.RegisterAdditivesRoutes(handler)

	return &AdditivesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
