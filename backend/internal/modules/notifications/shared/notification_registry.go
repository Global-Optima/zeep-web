package shared

import (
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
)

type NotificationRegistry struct {
	DetailsFactory func() details.NotificationDetails
	MessageBuilder func(details.NotificationDetails) (localization.LocalizedMessages, error)
}

var notificationRegistry = map[data.NotificationEventType]NotificationRegistry{}

func RegisterNotification(eventType data.NotificationEventType, factory func() details.NotificationDetails, builder func(details.NotificationDetails) (localization.LocalizedMessages, error)) {
	if _, exists := notificationRegistry[eventType]; exists {
		panic(fmt.Errorf("notification type already registered: %s", eventType))
	}

	notificationRegistry[eventType] = NotificationRegistry{
		DetailsFactory: factory,
		MessageBuilder: builder,
	}
}

func GetNotificationRegistry(eventType data.NotificationEventType) (NotificationRegistry, error) {
	registry, exists := notificationRegistry[eventType]
	if !exists {
		return NotificationRegistry{}, fmt.Errorf("no registry found for event type: %s", eventType)
	}
	return registry, nil
}
