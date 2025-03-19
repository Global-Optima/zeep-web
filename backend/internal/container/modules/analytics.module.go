package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/analytics"
)

type AnalyticsModule struct {
	*common.BaseModule
	Repo    analytics.AnalyticsRepo
	Service analytics.AnalyticsService
	Handler *analytics.AnalyticsHandler
}

func NewAnalyticsModule(
	base *common.BaseModule,
) *AnalyticsModule {
	repo := analytics.NewAnalyticsRepo(base.DB)
	service := analytics.NewAnalyticsService(repo, base.Logger)
	handler := analytics.NewAnalyticsHandler(service)

	base.Router.RegisterAnalyticRoutes(handler)

	return &AnalyticsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
