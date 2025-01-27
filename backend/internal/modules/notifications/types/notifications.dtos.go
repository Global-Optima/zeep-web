package types

import (
	"encoding/json"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type NotificationDTO struct {
	ID        uint                           `json:"id"`
	EventType string                         `json:"eventType"`
	Priority  string                         `json:"priority"`
	Message   localization.LocalizedMessages `json:"messages"`
	Details   json.RawMessage                `json:"details"`
	IsRead    bool                           `json:"isRead"`
	CreatedAt time.Time                      `json:"createdAt"`
	UpdatedAt time.Time                      `json:"updatedAt"`
}

type GetNotificationsFilter struct {
	Priority  *data.NotificationPriority  `form:"priority"`
	EventType *data.NotificationEventType `form:"eventType"`
	Search    *string                     `form:"search"`
	IsRead    *bool                       `form:"isRead"`
	StartDate *time.Time                  `form:"startDate" time_format:"2006-01-02"`
	EndDate   *time.Time                  `form:"endDate" time_format:"2006-01-02"`
	utils.BaseFilter
}

type MarkNotificationAsReadDTO struct {
	NotificationID uint `json:"notificationId" binding:"required"`
}

type MarkNotificationsAsReadDTO struct {
	NotificationIDs []uint `json:"notificationIds" binding:"required"`
}
