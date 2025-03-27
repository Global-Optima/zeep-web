package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/pkg/errors"
)

func UpdateToStoreStockModel(dto *UpdateStoreStockDTO, storeStock *data.StoreStock) error {
	if dto == nil {
		return errors.New("nil UpdateStoreStockDTO")
	}

	if dto.LowStockThreshold != nil {
		storeStock.LowStockThreshold = *dto.LowStockThreshold
	}

	if dto.Quantity != nil {
		storeStock.Quantity = *dto.Quantity
	}

	return nil
}
