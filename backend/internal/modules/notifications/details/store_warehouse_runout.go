package details

import (
	"encoding/json"
	"fmt"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

type StoreWarehouseRunOutDetails struct {
	BaseNotificationDetails
	StockItem   string `json:"stockItem"`
	StockItemID uint   `json:"stockItemId"`
}

func (s *StoreWarehouseRunOutDetails) ToDetails() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StoreWarehouseRunOutDetails) GetBaseDetails() *BaseNotificationDetails {
	return &s.BaseNotificationDetails
}

func BuildStoreWarehouseRunOutDetails(facilityID, stockItemID uint, facilityName, stockItem string) (*StoreWarehouseRunOutDetails, error) {
	if facilityID == 0 || stockItemID == 0 || facilityName == "" || stockItem == "" {
		return nil, fmt.Errorf("invalid input: facilityID, stockItemID, facilityName, and stockItem are required")
	}

	return &StoreWarehouseRunOutDetails{
		BaseNotificationDetails: BaseNotificationDetails{
			ID:           facilityID,
			FacilityName: facilityName,
		},
		StockItem:   stockItem,
		StockItemID: stockItemID,
	}, nil
}

func BuildStoreWarehouseRunOutMessage(details *StoreWarehouseRunOutDetails) (localization.LocalizedMessage, error) {
	if details == nil {
		return localization.LocalizedMessage{}, fmt.Errorf("details cannot be nil")
	}

	key := localization.FormTranslationKey("notification", data.STORE_WAREHOUSE_RUN_OUT.ToString())
	messages, err := localization.Translate(key, map[string]interface{}{
		"FacilityName": details.FacilityName,
		"ID":           details.ID,
		"StockItem":    details.StockItem,
		"StockItemID":  details.StockItemID,
	})
	if err != nil {
		return localization.LocalizedMessage{}, fmt.Errorf("failed to build StoreWarehouseRunOut message: %w", err)
	}

	return *messages, err
}
