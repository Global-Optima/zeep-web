package details

import "encoding/json"

type NotificationDetails interface {
	ToDetails() ([]byte, error)
	GetBaseDetails() *BaseNotificationDetails
}

type BaseNotificationDetails struct {
	ID           uint   `json:"id"`
	FacilityName string `json:"facilityName"` // Either Store or Warehouse name
}

func (b *BaseNotificationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(b)
}

func (b *BaseNotificationDetails) GetBaseDetails() *BaseNotificationDetails {
	return b
}
