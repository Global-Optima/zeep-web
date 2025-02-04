package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type PriceChangeNotificationDetails struct {
	BaseNotificationDetails
	ProductSizeID uint    `json:"productSizeId"`
	ProductName   string  `json:"productName"`
	OldPrice      float64 `json:"oldPrice"`
	NewPrice      float64 `json:"newPrice"`
}

func (p *PriceChangeNotificationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PriceChangeNotificationDetails) GetBaseDetails() *BaseNotificationDetails {
	return &p.BaseNotificationDetails
}

func BuildPriceChangeDetails(facilityID uint, facilityName, productName string, oldPrice, newPrice float64) (*PriceChangeNotificationDetails, error) {
	if facilityID == 0 || facilityName == "" || productName == "" || oldPrice <= 0 || newPrice <= 0 {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, productName, oldPrice, and newPrice are required")
	}

	return &PriceChangeNotificationDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		ProductName: productName,
		OldPrice:    oldPrice,
		NewPrice:    newPrice,
	}, nil
}

func BuildPriceChangeMessage(details *PriceChangeNotificationDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.PRICE_CHANGE.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName": details.FacilityName,
		"ID":           details.ID,
		"ProductName":  details.ProductName,
		"OldPrice":     details.OldPrice,
		"NewPrice":     details.NewPrice,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build PriceChange message: %w", err)
	}

	return *messages, nil
}
