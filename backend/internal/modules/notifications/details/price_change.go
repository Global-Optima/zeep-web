package details

import "encoding/json"

type PriceChangeNotificationDetails struct {
	BaseNotificationDetails
	ProductID   uint    `json:"productId"`
	ProductName string  `json:"productName"`
	OldPrice    float64 `json:"oldPrice"`
	NewPrice    float64 `json:"newPrice"`
}

func (p *PriceChangeNotificationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PriceChangeNotificationDetails) GetBaseDetails() *BaseNotificationDetails {
	return &p.BaseNotificationDetails
}
