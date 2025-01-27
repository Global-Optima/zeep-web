package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/shared"
)

type NotificationModule struct {
	*common.BaseModule
	Repo    notifications.NotificationRepository
	Service notifications.NotificationService
	Handler *notifications.NotificationHandler
}

func NewNotificationModule(base *common.BaseModule) *NotificationModule {
	shared.InitNotificationRegistry()

	repo := notifications.NewNotificationRepository(base.DB)
	roleManager := shared.NewNotificationRoleMappingManager()
	service := notifications.NewNotificationService(base.DB, repo, roleManager)
	handler := notifications.NewNotificationHandler(service)

	base.Router.RegisterNotificationsRoutes(handler)

	return &NotificationModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
