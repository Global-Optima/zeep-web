package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

// StockExpirationDetails specific to stock expiration
type StoreStockExpirationDetails struct {
	BaseNotificationDetails
	ItemName       string `json:"itemName"`
	ExpirationDate string `json:"expirationDate"`
}

func (s *StoreStockExpirationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StoreStockExpirationDetails) GetBaseDetails() *BaseNotificationDetails {
	return &s.BaseNotificationDetails
}

func BuildStockExpirationDetails(facilityID uint, facilityName, itemName, expirationDate string) (*StoreStockExpirationDetails, error) {
	if facilityID == 0 || facilityName == "" || itemName == "" || expirationDate == "" {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, itemName, and expirationDate are required")
	}

	return &StoreStockExpirationDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		ItemName:       itemName,
		ExpirationDate: expirationDate,
	}, nil
}

func BuildStockExpirationMessage(details *StoreStockExpirationDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.STORE_STOCK_EXPIRATION.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName":   details.FacilityName,
		"ID":             details.ID,
		"ItemName":       details.ItemName,
		"ExpirationDate": details.ExpirationDate,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build StockExpiration message: %w", err)
	}

	return *messages, nil
}
