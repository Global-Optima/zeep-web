package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type CentralCatalogUpdateDetails struct {
	BaseNotificationDetails
	UpdatedBy string `json:"updatedBy"`
	Changes   string `json:"changes"` // Summary of changes
}

func (c *CentralCatalogUpdateDetails) ToDetails() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CentralCatalogUpdateDetails) GetBaseDetails() *BaseNotificationDetails {
	return &c.BaseNotificationDetails
}

func BuildCentralCatalogUpdateDetails(facilityID uint, facilityName, updatedBy, changes string) (*CentralCatalogUpdateDetails, error) {
	if facilityID == 0 || facilityName == "" || updatedBy == "" || changes == "" {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, updatedBy, and changes are required")
	}

	return &CentralCatalogUpdateDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		UpdatedBy: updatedBy,
		Changes:   changes,
	}, nil
}

func BuildCentralCatalogUpdateMessage(details *CentralCatalogUpdateDetails) (localization.LocalizedMessages, error) {
	if details == nil {
		return localization.LocalizedMessages{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.CENTRAL_CATALOG_UPDATE.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName": details.FacilityName,
		"ID":           details.ID,
		"UpdatedBy":    details.UpdatedBy,
		"Changes":      details.Changes,
	})
	if err != nil {
		return localization.LocalizedMessages{}, fmt.Errorf("failed to build CentralCatalogUpdate message: %w", err)
	}

	return *messages, nil
}
