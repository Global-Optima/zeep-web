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
	StockRequestID uint                    `json:"stockRequestID"`
	RequestStatus  data.StockRequestStatus `json:"requestStatus"`
}

func (s *StockRequestStatusUpdatedDetails) ToDetails() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StockRequestStatusUpdatedDetails) GetBaseDetails() *BaseNotificationDetails {
	return &s.BaseNotificationDetails
}

func TranslateStockRequestStatus(status data.StockRequestStatus, lang string) (string, error) {
	key := "notification.stockRequestStatus." + localization.ToCamelCase(string(status))

	message, err := localization.Translate(key, nil)
	if err != nil {
		return "", fmt.Errorf("failed to translate stock request status for key %s: %w", key, err)
	}

	switch lang {
	case "en":
		return message.En, nil
	case "ru":
		return message.Ru, nil
	case "kk":
		return message.Kk, nil
	default:
		return message.Ru, nil
	}
}

func BuildStockRequestStatusUpdatedDetails(facilityID, stockRequestID uint, facilityName string, requestStatus data.StockRequestStatus) (*StockRequestStatusUpdatedDetails, error) {
	if facilityID == 0 || stockRequestID == 0 || facilityName == "" || requestStatus == "" {
		return nil, fmt.Errorf("invalid input: facilityID, stockRequestID, facilityName, and requestStatus are required")
	}

	return &StockRequestStatusUpdatedDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		StockRequestID: stockRequestID,
		RequestStatus:  data.StockRequestStatus(requestStatus),
	}, nil
}

func BuildStockRequestStatusUpdatedMessage(details *StockRequestStatusUpdatedDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	localizedMessages := localization.LocalizedMessage{}

	languages := map[string]*string{
		"en": &localizedMessages.En,
		"ru": &localizedMessages.Ru,
		"kk": &localizedMessages.Kk,
	}

	for lang, msg := range languages {
		translatedStatus, err := TranslateStockRequestStatus(data.StockRequestStatus(details.RequestStatus), lang)
		if err != nil {
			return localization.LocalizedMessage{}, fmt.Errorf("failed to translate stock request status for language %s: %w", lang, err)
		}

		key := localization.FormTranslationKey("notification", data.STOCK_REQUEST_STATUS_UPDATED.ToString())
		translatedMessage, err := localization.Translate(key, map[string]interface{}{
			"FacilityName":   details.FacilityName,
			"ID":             details.ID,
			"RequestStatus":  translatedStatus,
			"StockRequestID": details.StockRequestID,
		})
		if err != nil {
			return localization.LocalizedMessage{}, fmt.Errorf("failed to build %s message: %w", lang, err)
		}

		*msg = getMessageForLang(lang, *translatedMessage)
	}

	return localizedMessages, nil
}
