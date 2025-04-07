package types

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
)

func MapToStoreProvisionDTO(sp *data.StoreProvision) *StoreProvisionDTO {
	ingredients := make([]StoreProvisionIngredientDTO, len(sp.StoreProvisionIngredients))
	for i, spi := range sp.StoreProvisionIngredients {
		ingredients[i] = StoreProvisionIngredientDTO{
			Ingredient: *ingredientTypes.ConvertToIngredientResponseDTO(&spi.Ingredient),
			Quantity:   spi.Quantity,
		}
	}

	return &StoreProvisionDTO{
		ID:               sp.ID,
		BaseProvisionDTO: *provisionsTypes.MapToBaseProvisionDTO(&sp.Provision),
		Volume:           sp.Volume,
		Status:           sp.Status.ToString(),
		CompletedAt:      sp.CompletedAt,
		ExpiresAt:        sp.ExpiresAt,
		Ingredients:      ingredients,
	}
}

func CreateToStoreProvisionModel(storeID uint, dto *CreateStoreProvisionDTO) *data.StoreProvision {
	return &data.StoreProvision{
		ProvisionID:               dto.ProvisionID,
		Volume:                    dto.Volume,
		ExpirationInHours:         dto.ExpirationInHours,
		Status:                    data.PROVISION_STATUS_PREPARING,
		StoreID:                   storeID,
		StoreProvisionIngredients: mapIngredientsToStoreProvisionIngredients(dto.Ingredients),
	}
}

func UpdateToStoreProvisionModel(storeProvision *data.StoreProvision, dto *UpdateStoreProvisionDTO) error {
	if dto == nil {
		return fmt.Errorf("dto is nil")
	}

	if storeProvision == nil || storeProvision.ID == 0 || storeProvision.ProvisionID == 0 {
		return fmt.Errorf("invalid argument for ID paramters fetched while validating existing provision")
	}

	if dto.Volume != nil {
		storeProvision.Volume = *dto.Volume
	}
	if dto.ExpirationInHours != nil {
		storeProvision.ExpirationInHours = *dto.ExpirationInHours
	}
	storeProvision.StoreProvisionIngredients = mapIngredientsToStoreProvisionIngredients(dto.Ingredients)

	return nil
}

func mapIngredientsToStoreProvisionIngredients(dto []ingredientTypes.SelectedIngredientDTO) []data.StoreProvisionIngredient {
	result := make([]data.StoreProvisionIngredient, len(dto))
	for i, item := range dto {
		result[i] = data.StoreProvisionIngredient{
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
		}
	}
	return result
}
