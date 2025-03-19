package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

// OutOfStockDetails specific to items out of stock
type OutOfStockDetails struct {
	BaseNotificationDetails
	ItemName string `json:"itemName"`
}

func (o *OutOfStockDetails) ToDetails() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OutOfStockDetails) GetBaseDetails() *BaseNotificationDetails {
	return &o.BaseNotificationDetails
}

func BuildOutOfStockDetails(facilityID uint, facilityName, itemName string) (*OutOfStockDetails, error) {
	if facilityID == 0 || facilityName == "" || itemName == "" {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, and itemName are required")
	}

	return &OutOfStockDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		ItemName: itemName,
	}, nil
}

func BuildOutOfStockMessage(details *OutOfStockDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.WAREHOUSE_OUT_OF_STOCK.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName": details.FacilityName,
		"ID":           details.ID,
		"ItemName":     details.ItemName,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build OutOfStock message: %w", err)
	}

	return *messages, nil
}
