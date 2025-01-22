package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
)

type AuditsModule struct {
	*common.BaseModule
	Repo    audit.AuditRepository
	Service audit.AuditService
	Handler *audit.AuditHandler
}

func NewAuditsModule(
	base *common.BaseModule,
) *AuditsModule {

	repo := audit.NewAuditRepository(base.DB)
	service := audit.NewAuditService(repo, base.Logger)
	handler := audit.NewAuditHandler(service)

	base.Router.RegisterAuditRoutes(handler)

	return &AuditsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
