package types

import (
	"errors"
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
		return errors.New("only one stock request is allowed per day")
	}
	return nil
}

func IsValidTransition(currentStatus, targetStatus data.StockRequestStatus) bool {
	validTransitions := map[data.StockRequestStatus][]data.StockRequestStatus{
		data.StockRequestCreated:             {data.StockRequestProcessed},
		data.StockRequestProcessed:           {data.StockRequestInDelivery, data.StockRequestRejectedByWarehouse},
		data.StockRequestInDelivery:          {data.StockRequestCompleted, data.StockRequestAcceptedWithChange, data.StockRequestRejectedByStore},
		data.StockRequestRejectedByWarehouse: {data.StockRequestProcessed}, // Terminal state, can reuse rejected cart
		data.StockRequestRejectedByStore:     {data.StockRequestProcessed}, // Terminal state, can reuse rejected cart
		data.StockRequestCompleted:           {},                           // Terminal state
		data.StockRequestAcceptedWithChange:  {},                           // Terminal state
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

func FilterWarehouseStatuses(inputStatuses []data.StockRequestStatus) []data.StockRequestStatus {
	allowedStatuses := map[data.StockRequestStatus]bool{
		data.StockRequestProcessed:           true,
		data.StockRequestInDelivery:          true,
		data.StockRequestCompleted:           true,
		data.StockRequestRejectedByStore:     true,
		data.StockRequestRejectedByWarehouse: true,
		data.StockRequestAcceptedWithChange:  true,
	}

	var filteredStatuses []data.StockRequestStatus
	for _, status := range inputStatuses {
		if allowedStatuses[status] {
			filteredStatuses = append(filteredStatuses, status)
		}
	}
	return filteredStatuses
}
