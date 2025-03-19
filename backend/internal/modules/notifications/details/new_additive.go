package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type NewAdditiveDetails struct {
	BaseNotificationDetails
	AdditiveName string `json:"additiveName"`
}

func (n *NewAdditiveDetails) ToDetails() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NewAdditiveDetails) GetBaseDetails() *BaseNotificationDetails {
	return &n.BaseNotificationDetails
}

func BuildNewAdditiveDetails(facilityID uint, facilityName, additiveName string) (*NewAdditiveDetails, error) {
	if facilityID == 0 || facilityName == "" || additiveName == "" {
		return nil, fmt.Errorf("invalid input: all fields are required")
	}
	return &NewAdditiveDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		AdditiveName: additiveName,
	}, nil
}

func BuildNewAdditiveMessage(details *NewAdditiveDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}
	key := localization.FormTranslationKey("notification", data.NEW_ADDITIVE.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName": details.FacilityName,
		"AdditiveName": details.AdditiveName,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build NewAdditive message: %w", err)
	}
	return *messages, nil
}
