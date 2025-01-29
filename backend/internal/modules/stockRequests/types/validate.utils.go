package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
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

type StockRequestDetails struct {
	MaterialName   string  `json:"materialName"`
	Quantity       float64 `json:"quantity"`
	ActualQuantity float64 `json:"actualQuantity"`
}

func GenerateUnexpectedCommentFromDetails(details StockRequestDetails) *localization.LocalizedMessages {
	key := "stockRequestComment.unexpected"
	translations, err := localization.Translate(key, map[string]interface{}{
		"MaterialName":   details.MaterialName,
		"ActualQuantity": details.ActualQuantity,
	})

	if err != nil {
		return nil
	}

	return translations
}

func GenerateMismatchCommentFromDetails(details StockRequestDetails) *localization.LocalizedMessages {
	key := "stockRequestComment.mismatch"
	translations, err := localization.Translate(key, map[string]interface{}{
		"MaterialName":   details.MaterialName,
		"Quantity":       details.Quantity,
		"ActualQuantity": details.ActualQuantity,
	})

	if err != nil {
		return nil
	}

	return translations
}

func CombineComments(
	storeComment string,
	mismatchComments []localization.LocalizedMessages,
	unexpectedComments []localization.LocalizedMessages,
) *localization.LocalizedMessages {
	var combinedComments localization.LocalizedMessages

	for _, mismismatchComment := range mismatchComments {
		combinedComments.En += mismismatchComment.En + "\n"
		combinedComments.Ru += mismismatchComment.Ru + "\n"
		combinedComments.Kk += mismismatchComment.Kk + "\n"
	}

	combinedComments.En += "\n"
	combinedComments.Ru += "\n"
	combinedComments.Kk += "\n"

	for _, unexpectedComment := range unexpectedComments {
		unexpectedComment.En += unexpectedComment.En + "\n"
		unexpectedComment.Ru += unexpectedComment.Ru + "\n"
		unexpectedComment.Kk += unexpectedComment.Kk + "\n"
	}

	combinedComments.En += "\n" + "Store comment: " + storeComment
	combinedComments.Ru += "\n" + "Комментарий от магазина: " + storeComment
	combinedComments.Kk += "\n" + "Кафе пікірі: " + storeComment

	return &combinedComments
}
