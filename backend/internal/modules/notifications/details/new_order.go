package details

import "encoding/json"

// NewOrderNotificationDetails specific to new order notifications
type NewOrderNotificationDetails struct {
	BaseNotificationDetails
	EmployeeName string `json:"employeeName"`
	OrderID      uint   `json:"orderId"`
}

func (n *NewOrderNotificationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NewOrderNotificationDetails) GetBaseDetails() *BaseNotificationDetails {
	return &n.BaseNotificationDetails
}
