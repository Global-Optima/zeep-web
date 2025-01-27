package data

import "gorm.io/datatypes"

type NotificationPriority string
type NotificationEventType string

const (
	HIGH   NotificationPriority = "HIGH"
	MEDIUM NotificationPriority = "MEDIUM"
	LOW    NotificationPriority = "LOW"
)

const (
	STOCK_REQUEST_STATUS_UPDATED NotificationEventType = "STOCK_REQUEST_STATUS_UPDATED"
	NEW_ORDER                    NotificationEventType = "NEW_ORDER"
	STORE_WAREHOUSE_RUN_OUT      NotificationEventType = "STORE_WAREHOUSE_RUN_OUT"
	CENTRAL_CATALOG_UPDATE       NotificationEventType = "CENTRAL_CATALOG_UPDATE"
	WAREHOUSE_STOCK_EXPIRATION   NotificationEventType = "WAREHOUSE_STOCK_EXPIRATION"
	WAREHOUSE_OUT_OF_STOCK       NotificationEventType = "WAREHOUSE_OUT_OF_STOCK"
	NEW_STOCK_REQUEST            NotificationEventType = "NEW_STOCK_REQUEST"
	PRICE_CHANGE                 NotificationEventType = "PRICE_CHANGE"
)

func (nt NotificationEventType) ToString() string {
	return string(nt)
}

func (np NotificationPriority) ToString() string {
	return string(np)
}

type EmployeeNotification struct {
	BaseEntity
	EventType  NotificationEventType           `gorm:"type:varchar(255);not null"`
	Priority   NotificationPriority            `gorm:"type:varchar(50);not null"`
	Details    datatypes.JSON                  `gorm:"type:jsonb"`
	Recipients []EmployeeNotificationRecipient `gorm:"foreignKey:NotificationID"`
}

type EmployeeNotificationRecipient struct {
	BaseEntity
	NotificationID uint                 `gorm:"not null;index"`
	Notification   EmployeeNotification `gorm:"foreignKey:NotificationID;constraint:OnDelete:CASCADE;"`
	EmployeeID     uint                 `gorm:"not null;index"`
	Employee       Employee             `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE;"`
	IsRead         bool                 `gorm:"default:false;not null"`
}
