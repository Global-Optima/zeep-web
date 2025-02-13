package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

// WarehouseStockExpirationDetails specific to stock expiration
type WarehouseStockExpirationDetails struct {
	BaseNotificationDetails
	ItemName       string `json:"itemName"`
	ExpirationDate string `json:"expirationDate"`
}

func (s *WarehouseStockExpirationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(s)
}

func (s *WarehouseStockExpirationDetails) GetBaseDetails() *BaseNotificationDetails {
	return &s.BaseNotificationDetails
}

func BuildWarehouseStockExpirationDetails(facilityID uint, facilityName, itemName, expirationDate string) (*WarehouseStockExpirationDetails, error) {
	if facilityID == 0 || facilityName == "" || itemName == "" || expirationDate == "" {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, itemName, and expirationDate are required")
	}

	return &WarehouseStockExpirationDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		ItemName:       itemName,
		ExpirationDate: expirationDate,
	}, nil
}

func BuildWarehouseStockExpirationMessage(details *WarehouseStockExpirationDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.WAREHOUSE_STOCK_EXPIRATION.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName":   details.FacilityName,
		"ID":             details.ID,
		"ItemName":       details.ItemName,
		"ExpirationDate": details.ExpirationDate,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build WarehouseStockExpiration message: %w", err)
	}

	return *messages, nil
}
