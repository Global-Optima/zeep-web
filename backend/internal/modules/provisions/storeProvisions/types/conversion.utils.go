package types

import (
	"fmt"
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
	return &StoreProvisionDTO{
		ID:                sp.ID,
		BaseProvisionDTO:  *provisionsTypes.MapToBaseProvisionDTO(&sp.Provision),
		ProvisionID:       sp.ProvisionID,
		ExpirationInHours: sp.ExpirationInHours,
		Volume:            sp.Volume,
		Status:            sp.Status,
		CompletedAt:       sp.CompletedAt,
		ExpiresAt:         sp.ExpiresAt,
		CreatedAt:         sp.CreatedAt,
	}
}

func MapToStoreProvisionDetailsDTO(sp *data.StoreProvision) *StoreProvisionDetailsDTO {
	ingredients := make([]StoreProvisionIngredientDTO, len(sp.StoreProvisionIngredients))
	for i, spi := range sp.StoreProvisionIngredients {
		ingredients[i] = StoreProvisionIngredientDTO{
			Ingredient: *ingredientTypes.ConvertToIngredientResponseDTO(&spi.Ingredient),
			Quantity:   spi.Quantity,
		}
	}

	return &StoreProvisionDetailsDTO{
		StoreProvisionDTO: *MapToStoreProvisionDTO(sp),
		Ingredients:       ingredients,
	}
}

func CreateToStoreProvisionModel(storeID uint, dto *CreateStoreProvisionDTO, centralCatalogProvision *data.Provision) *data.StoreProvision {
	return &data.StoreProvision{
		ProvisionID:               dto.ProvisionID,
		Volume:                    dto.Volume,
		ExpirationInHours:         dto.ExpirationInHours,
		Status:                    data.PROVISION_STATUS_PREPARING,
		StoreID:                   storeID,
		StoreProvisionIngredients: mapIngredientsToStoreProvisionIngredients(dto.Volume, centralCatalogProvision),
	}
}

func UpdateToStoreProvisionModel(storeProvision *data.StoreProvision, dto *UpdateStoreProvisionDTO) (*StoreProvisionModels, error) {
	if dto == nil {
		return nil, fmt.Errorf("dto is nil")
	}

	if storeProvision == nil || storeProvision.ID == 0 || storeProvision.ProvisionID == 0 {
		return nil, fmt.Errorf("invalid argument for ID paramters fetched while validating existing provision")
	}

	updateModels := &StoreProvisionModels{}

	if dto.Volume != nil {
		multiplier := CalculateStoreProvisionIngredientsMultiplier(*dto.Volume, storeProvision.Volume)
		if multiplier != 1 {
			storeProvision.Volume = *dto.Volume
			updateModels.StoreProvisionIngredientsMultiplier = &multiplier
		}
	}
	if dto.ExpirationInHours != nil {
		storeProvision.ExpirationInHours = *dto.ExpirationInHours
	}

	return updateModels, nil
}

func mapIngredientsToStoreProvisionIngredients(volume float64, provision *data.Provision) []data.StoreProvisionIngredient {
	result := make([]data.StoreProvisionIngredient, len(provision.ProvisionIngredients))
	multiplier := CalculateStoreProvisionIngredientsMultiplier(volume, provision.AbsoluteVolume)
	for i, ingredient := range provision.ProvisionIngredients {
		result[i] = data.StoreProvisionIngredient{
			IngredientID: ingredient.IngredientID,
			Quantity:     ingredient.Quantity * multiplier,
		}
	}
	return result
}

func CalculateStoreProvisionIngredientsMultiplier(volume, absoluteVolume float64) float64 {
	return utils.RoundToDecimal(volume/absoluteVolume, 1)
}
