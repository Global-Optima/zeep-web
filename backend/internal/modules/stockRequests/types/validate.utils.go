package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/logger"
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
		return ErrOneRequestPerDay
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

func DefaultWarehouseStatuses() []data.StockRequestStatus {
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
	OriginalMaterialName string  `json:"originalMaterialName,omitempty"`
	MaterialName         string  `json:"materialName,omitempty"`
	Quantity             float64 `json:"quantity,omitempty"`
	ActualQuantity       float64 `json:"actualQuantity,omitempty"`
}

func GenerateUnexpectedCommentFromDetails(details StockRequestDetails) *localization.LocalizedMessage {
	key := "stockRequestComments.unexpectedMaterial"
	translations, err := localization.Translate(key, map[string]interface{}{
		"MaterialName":   details.MaterialName,
		"ActualQuantity": details.ActualQuantity,
	})
	if err != nil {
		return nil
	}

	return translations
}

func GenerateMismatchCommentFromDetails(details StockRequestDetails) *localization.LocalizedMessage {
	key := "stockRequestComments.quantityMismatch"
	translations, err := localization.Translate(key, map[string]interface{}{
		"OriginalMaterialName": details.OriginalMaterialName,
		"Quantity":             details.Quantity,
		"ActualQuantity":       details.ActualQuantity,
	})
	if err != nil {
		return nil
	}

	return translations
}

func CombineComments(
	mismatchComments []localization.LocalizedMessage,
	unexpectedComments []localization.LocalizedMessage,
) *localization.LocalizedMessage {
	if len(mismatchComments) == 0 && len(unexpectedComments) == 0 {
		logger.GetZapSugaredLogger().Warn("No comments to combine")
		return nil
	}

	var combinedComments localization.LocalizedMessage

	for _, mismismatchComment := range mismatchComments {
		combinedComments.En += mismismatchComment.En + "\n"
		combinedComments.Ru += mismismatchComment.Ru + "\n"
		combinedComments.Kk += mismismatchComment.Kk + "\n"
	}

	if len(mismatchComments) > 0 {
		combinedComments.En += "\n"
		combinedComments.Ru += "\n"
		combinedComments.Kk += "\n"
	}

	for _, unexpectedComment := range unexpectedComments {
		combinedComments.En += unexpectedComment.En + "\n"
		combinedComments.Ru += unexpectedComment.Ru + "\n"
		combinedComments.Kk += unexpectedComment.Kk + "\n"
	}

	return &combinedComments
}
