package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
)

func ValidateStoreProductSizes(inputIDs []uint, validProductSizes []data.ProductSize) error {
	validProductSizeIDs := make(map[uint]struct{}, len(validProductSizes))
	for _, ps := range validProductSizes {
		validProductSizeIDs[ps.ID] = struct{}{}
	}

	for _, inputID := range inputIDs {
		if _, exists := validProductSizeIDs[inputID]; !exists {
			return ErrInappropriateProductSizeID
		}
	}

	return nil
}
