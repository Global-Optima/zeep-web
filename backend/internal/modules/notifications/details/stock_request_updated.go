package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

// StockRequestStatusUpdatedDetails specific to stock request updates
type StockRequestStatusUpdatedDetails struct {
	BaseNotificationDetails
	StockRequestID uint   `json:"stockRequestID"`
	RequestStatus  string `json:"requestStatus"`
}

func (s *StockRequestStatusUpdatedDetails) ToDetails() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StockRequestStatusUpdatedDetails) GetBaseDetails() *BaseNotificationDetails {
	return &s.BaseNotificationDetails
}

func BuildStockRequestStatusUpdatedDetails(facilityID, stockRequestID uint, facilityName, requestStatus string) (*StockRequestStatusUpdatedDetails, error) {
	if facilityID == 0 || stockRequestID == 0 || facilityName == "" || requestStatus == "" {
		return nil, fmt.Errorf("invalid input: facilityID, stockRequestID, facilityName, and requestStatus are required")
	}

	return &StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		StockRequestID: stockRequestID,
		RequestStatus:  requestStatus,
	}, nil
}

func BuildStockRequestStatusUpdatedMessage(details *StockRequestStatusUpdatedDetails) (localization.LocalizedMessages, error) {
	if details == nil {
		return localization.LocalizedMessages{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.STOCK_REQUEST_STATUS_UPDATED.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName":   details.FacilityName,
		"ID":             details.ID,
		"RequestStatus":  details.RequestStatus,
		"StockRequestID": details.StockRequestID,
	})
	if err != nil {
		return localization.LocalizedMessages{}, fmt.Errorf("failed to build StockRequestStatusUpdated message: %w", err)
	}

	return *messages, nil
}
