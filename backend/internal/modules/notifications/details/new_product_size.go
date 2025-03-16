package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type NewProductSizeDetails struct {
	BaseNotificationDetails
	ProductName     string  `json:"productName"`
	ProductSizeName string  `json:"productSizeName"`
	Size            float64 `json:"size"`
}

func (n *NewProductSizeDetails) ToDetails() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NewProductSizeDetails) GetBaseDetails() *BaseNotificationDetails {
	return &n.BaseNotificationDetails
}

func BuildNewProductSizeDetails(facilityID uint, facilityName, productName, productSizeName string, size float64) (*NewProductSizeDetails, error) {
	if facilityID == 0 || facilityName == "" || productName == "" || productSizeName == "" || size == 0 {
		return nil, fmt.Errorf("invalid input: all fields are required and size != 0")
	}
	return &NewProductSizeDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		ProductName:     productName,
		ProductSizeName: productSizeName,
		Size:            size,
	}, nil
}

func BuildNewProductSizeMessage(details *NewProductSizeDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}
	key := localization.FormTranslationKey("notification", data.NEW_PRODUCT_SIZE.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName":    details.FacilityName,
		"ProductName":     details.ProductName,
		"ProductSizeName": details.ProductSizeName,
		"Size":            details.Size,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build NewProductSize message: %w", err)
	}
	return *messages, nil
}
