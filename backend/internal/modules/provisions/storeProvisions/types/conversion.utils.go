package types

import (
	"fmt"
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type StoreProvisionModels struct {
	StoreProvision                      *data.StoreProvision
	StoreProvisionIngredientsMultiplier *float64
}

func MapToStoreProvisionDTO(sp *data.StoreProvision) *StoreProvisionDTO {
	storeProvisionDTO := &StoreProvisionDTO{
		ID:                  sp.ID,
		Provision:           *provisionsTypes.MapToProvisionDTO(&sp.Provision),
		ExpirationInMinutes: sp.ExpirationInMinutes,
		Volume:              sp.Volume,
		CompletedAt:         sp.CompletedAt,
		ExpiresAt:           sp.ExpiresAt,
		CreatedAt:           sp.CreatedAt,
	}

	if sp.ExpiresAt != nil && sp.ExpiresAt.UTC().Before(time.Now().UTC()) {
		storeProvisionDTO.Status = data.STORE_PROVISION_VISUAL_STATUS_EXPIRED
	} else {
		storeProvisionDTO.Status = sp.Status
	}

	return storeProvisionDTO
}

func MapToStoreProvisionDetailsDTO(sp *data.StoreProvision) *StoreProvisionDetailsDTO {
	ingredients := make([]StoreProvisionIngredientDTO, len(sp.StoreProvisionIngredients))
	for i, spi := range sp.StoreProvisionIngredients {
		ingredients[i] = StoreProvisionIngredientDTO{
			Ingredient:      *ingredientTypes.ConvertToIngredientResponseDTO(&spi.Ingredient),
			Quantity:        spi.Quantity,
			InitialQuantity: spi.InitialQuantity,
		}
	}

	return &StoreProvisionDetailsDTO{
		StoreProvisionDTO: *MapToStoreProvisionDTO(sp),
		Ingredients:       ingredients,
	}
}

func CreateToStoreProvisionModel(storeID uint, dto *CreateStoreProvisionDTO, centralCatalogProvision *data.Provision) (*data.StoreProvision, error) {
	ingredients, err := mapIngredientsToStoreProvisionIngredients(dto.Volume, centralCatalogProvision)
	if err != nil {
		return nil, err
	}

	return &data.StoreProvision{
		ProvisionID:               dto.ProvisionID,
		Volume:                    dto.Volume,
		ExpirationInMinutes:       dto.ExpirationInMinutes,
		Status:                    data.STORE_PROVISION_STATUS_PREPARING,
		StoreID:                   storeID,
		StoreProvisionIngredients: ingredients,
	}, nil
}

func UpdateToStoreProvisionModels(storeProvision *data.StoreProvision, dto *UpdateStoreProvisionDTO) (*StoreProvisionModels, error) {
	if dto == nil {
		return nil, fmt.Errorf("dto is nil")
	}

	if storeProvision == nil || storeProvision.ID == 0 || storeProvision.ProvisionID == 0 || storeProvision.Provision.AbsoluteVolume == 0 {
		return nil, fmt.Errorf("invalid arguments for storeProvision or provision model")
	}

	updateModels := &StoreProvisionModels{
		StoreProvision: storeProvision,
	}

	if dto.Volume != nil {
		multiplier, err := CalculateStoreProvisionIngredientsMultiplier(*dto.Volume, storeProvision.Provision.AbsoluteVolume)
		if err != nil {
			return nil, ErrInvalidStoreProvisionIngredientsVolume
		}
		updateModels.StoreProvisionIngredientsMultiplier = &multiplier
		storeProvision.Volume = *dto.Volume
	}

	if dto.ExpirationInMinutes != nil {
		storeProvision.ExpirationInMinutes = *dto.ExpirationInMinutes
	}

	return updateModels, nil
}

func mapIngredientsToStoreProvisionIngredients(volume float64, provision *data.Provision) ([]data.StoreProvisionIngredient, error) {
	if provision == nil || provision.AbsoluteVolume == 0 {
		return nil, fmt.Errorf("invalid input data")
	}

	multiplier, err := CalculateStoreProvisionIngredientsMultiplier(volume, provision.AbsoluteVolume)
	if err != nil {
		return nil, ErrInvalidStoreProvisionIngredientsVolume
	}

	result := make([]data.StoreProvisionIngredient, 0, len(provision.ProvisionIngredients))
	for _, ingredient := range provision.ProvisionIngredients {
		calculated := ingredient.Quantity * multiplier
		rounded := utils.RoundToDecimal(calculated, 2)

		if rounded == 0 {
			return nil, ErrInvalidStoreProvisionIngredientsVolume
		}

		result = append(result, data.StoreProvisionIngredient{
			IngredientID:    ingredient.IngredientID,
			Quantity:        rounded,
			InitialQuantity: ingredient.Quantity,
		})
	}
	return result, nil
}

func CalculateStoreProvisionIngredientsMultiplier(volume, absoluteVolume float64) (float64, error) {
	if absoluteVolume == 0 {
		return 0, fmt.Errorf("cannot divide by 0")
	}
	return utils.RoundToDecimal(volume/absoluteVolume, 2), nil
}
