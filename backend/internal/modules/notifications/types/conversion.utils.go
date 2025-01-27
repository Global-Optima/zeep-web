package types

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/shared"
)

func ConvertNotificationToDTO(notification *data.EmployeeNotification, isRead bool) (*NotificationDTO, error) {
	registry, err := shared.GetNotificationRegistry(notification.EventType)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification registry: %w", err)
	}

	details := registry.DetailsFactory()
	if err := json.Unmarshal(notification.Details, details); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification details: %w", err)
	}

	messages, err := registry.MessageBuilder(details)
	if err != nil {
		return nil, fmt.Errorf("failed to generate localized message: %w", err)
	}

	notificationDetails, err := details.ToDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to convert notification details to JSON: %w", err)
	}
	return &NotificationDTO{
		ID:        notification.ID,
		EventType: string(notification.EventType),
		Priority:  string(notification.Priority),
		Message:   messages,
		Details:   json.RawMessage(notificationDetails), // Dynamic payload
		IsRead:    isRead,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}, nil
}

func ConvertNotificationsToDTOs(notifications []data.EmployeeNotificationRecipient) ([]NotificationDTO, error) {
	response := make([]NotificationDTO, len(notifications))
	for i, n := range notifications {
		dto, err := ConvertNotificationToDTO(&n.Notification, n.IsRead)
		if err != nil {
			return nil, fmt.Errorf("failed to convert notification to DTO: %w", err)
		}
		response[i] = *dto
	}
	return response, nil
}
