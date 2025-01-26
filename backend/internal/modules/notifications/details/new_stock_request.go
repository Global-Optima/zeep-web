package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

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

func BuildNewStockRequestDetails(facilityID uint, facilityName, requesterName string, requestID uint) (*NewStockRequestDetails, error) {
	if facilityID == 0 || facilityName == "" || requesterName == "" || requestID == 0 {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, requesterName, and requestID are required")
	}

	return &NewStockRequestDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		RequesterName: requesterName,
		RequestID:     requestID,
	}, nil
}

func BuildNewStockRequestMessage(details *NewStockRequestDetails) (localization.LocalizedMessages, error) {
	if details == nil {
		return localization.LocalizedMessages{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.NEW_STOCK_REQUEST.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName":  details.FacilityName,
		"ID":            details.ID,
		"RequesterName": details.RequesterName,
		"RequestID":     details.RequestID,
	})
	if err != nil {
		return localization.LocalizedMessages{}, fmt.Errorf("failed to build NewStockRequest message: %w", err)
	}

	return *messages, nil
}
