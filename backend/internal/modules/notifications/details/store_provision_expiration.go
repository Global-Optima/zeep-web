package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type StoreProvisionExpirationDetails struct {
	BaseNotificationDetails
	ItemName       string `json:"itemName"`
	CompletionDate string `json:"completionDate"`
}

func (s *StoreProvisionExpirationDetails) ToDetails() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StoreProvisionExpirationDetails) GetBaseDetails() *BaseNotificationDetails {
	return &s.BaseNotificationDetails
}

func BuildStoreProvisionExpirationMessage(details *StoreProvisionExpirationDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.STORE_PROVISION_EXPIRATION.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName":   details.FacilityName,
		"ID":             details.ID,
		"ItemName":       details.ItemName,
		"CompletionDate": details.CompletionDate,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build StoreProvisionExpiration message: %w", err)
	}

	return *messages, nil
}
