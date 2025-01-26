package types

import (
	"encoding/json"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type NotificationDTO struct {
	ID        uint                           `json:"id"`
	EventType string                         `json:"eventType"`
	Priority  string                         `json:"priority"`
	Message   localization.LocalizedMessages `json:"messages"`
	Details   json.RawMessage                `json:"details"`
	Timestamp time.Time                      `json:"timestamp"`
	IsRead    bool                           `json:"isRead"`
	CreatedAt time.Time                      `json:"createdAt"`
	UpdatedAt time.Time                      `json:"updatedAt"`
}

type GetNotificationsFilter struct {
	Priority  *string    `form:"priority"`
	EventType *string    `form:"eventType"`
	Search    *string    `form:"search"`
	IsRead    *bool      `form:"isRead"`
	StartDate *time.Time `form:"startDate"`
	EndDate   *time.Time `form:"endDate"`
	utils.BaseFilter
}

type MarkNotificationAsReadDTO struct {
	NotificationID uint `json:"notificationId" binding:"required"`
}

type MarkNotificationsAsReadDTO struct {
	NotificationIDs []uint `json:"notificationIds" binding:"required"`
}
