package details

import "encoding/json"

// StockRequestStatusUpdatedDetails specific to stock request updates
type StockRequestStatusUpdatedDetails struct {
	BaseNotificationDetails
	RequestStatus string `json:"requestStatus"`
}

func (s *StockRequestStatusUpdatedDetails) ToDetails() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StockRequestStatusUpdatedDetails) GetBaseDetails() *BaseNotificationDetails {
	return &s.BaseNotificationDetails
}

type NewStockRequestDetails struct {
	BaseNotificationDetails
	RequesterName string `json:"requesterName"`
	RequestID     uint   `json:"requestId"`
}

func (n *NewStockRequestDetails) ToDetails() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NewStockRequestDetails) GetBaseDetails() *BaseNotificationDetails {
	return &n.BaseNotificationDetails
}
