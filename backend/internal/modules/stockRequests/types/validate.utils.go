package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ValidateStatus(status string) error {
	validStatuses := []string{
		string(data.StockRequestCreated),
		string(data.StockRequestProcessed),
		string(data.StockRequestInDelivery),
		string(data.StockRequestCompleted),
		string(data.StockRequestRejectedByStore),
		string(data.StockRequestRejectedByWarehouse),
		string(data.StockRequestAcceptedWithChange),
	}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return errors.New("invalid status value")
}

func ValidateStockRequestRate(lastRequestDate *time.Time) error {
	if lastRequestDate != nil && time.Since(*lastRequestDate) < 24*time.Hour {
		return errors.New("only one stock request is allowed to send per day")
	}
	return nil
}

func IsValidTransition(currentStatus, targetStatus data.StockRequestStatus) bool {
	validTransitions := map[data.StockRequestStatus][]data.StockRequestStatus{
		data.StockRequestCreated:             {data.StockRequestProcessed},
		data.StockRequestProcessed:           {data.StockRequestInDelivery, data.StockRequestRejectedByWarehouse},
		data.StockRequestInDelivery:          {data.StockRequestCompleted, data.StockRequestAcceptedWithChange, data.StockRequestRejectedByStore},
		data.StockRequestRejectedByWarehouse: {data.StockRequestProcessed},
		data.StockRequestRejectedByStore:     {data.StockRequestInDelivery},
		data.StockRequestCompleted:           {},
		data.StockRequestAcceptedWithChange:  {},
	}

	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedTransitions {
		if status == targetStatus {
			return true
		}
	}

	return false
}

var warehouseAllowedStatuses = map[data.StockRequestStatus]bool{
	data.StockRequestCreated:             false,
	data.StockRequestProcessed:           true,
	data.StockRequestInDelivery:          true,
	data.StockRequestCompleted:           true,
	data.StockRequestRejectedByStore:     true,
	data.StockRequestRejectedByWarehouse: true,
	data.StockRequestAcceptedWithChange:  true,
}

func ValidateWarehouseStatuses(inputStatuses []data.StockRequestStatus) error {
	for _, status := range inputStatuses {
		if !warehouseAllowedStatuses[status] {
			return fmt.Errorf("invalid status for warehouse: %s", status)
		}
	}
	return nil
}

func DefaultWarehouseStasuses() []data.StockRequestStatus {
	return []data.StockRequestStatus{
		data.StockRequestProcessed,
		data.StockRequestInDelivery,
		data.StockRequestCompleted,
		data.StockRequestRejectedByStore,
		data.StockRequestRejectedByWarehouse,
		data.StockRequestAcceptedWithChange,
	}
}
