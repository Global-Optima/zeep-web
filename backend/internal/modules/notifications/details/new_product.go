package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type NewProductDetails struct {
	BaseNotificationDetails
	ProductName string `json:"productName"`
}

func (n *NewProductDetails) ToDetails() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NewProductDetails) GetBaseDetails() *BaseNotificationDetails {
	return &n.BaseNotificationDetails
}

func BuildNewProductDetails(facilityID uint, facilityName, productName string) (*NewProductDetails, error) {
	if facilityID == 0 || facilityName == "" || productName == "" {
		return nil, fmt.Errorf("invalid input: facilityID, facilityName, and productName are required")
	}
	return &NewProductDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		ProductName: productName,
	}, nil
}

func BuildNewProductMessage(details *NewProductDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}
	key := localization.FormTranslationKey("notification", data.NEW_PRODUCT.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName": details.FacilityName,
		"ProductName":  details.ProductName,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build NewProduct message: %w", err)
	}
	return *messages, nil
}
