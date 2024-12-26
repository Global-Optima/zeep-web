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
		string(data.StockRequestRejected),
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
