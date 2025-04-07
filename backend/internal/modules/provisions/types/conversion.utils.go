package types

import (
	"encoding/json"
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/pkg/errors"
)

type ProvisionModels struct {
	Provision   *data.Provision
	Ingredients []data.ProvisionIngredient
}

func MapToBaseProvisionDTO(provision *data.Provision) *BaseProvisionDTO {
	return &BaseProvisionDTO{
		Name:                 provision.Name,
		AbsoluteVolume:       provision.AbsoluteVolume,
		NetCost:              provision.NetCost,
		Unit:                 unitTypes.ToUnitResponse(provision.Unit),
		PreparationInMinutes: provision.PreparationInMinutes,
		LimitPerDay:          provision.LimitPerDay,
	}
}

func MapToProvisionDTO(provision *data.Provision) *ProvisionDTO {
	return &ProvisionDTO{
		ID:               provision.ID,
		BaseProvisionDTO: *MapToBaseProvisionDTO(provision),
	}
}

func MapToProvisionDetailsDTO(provision *data.Provision) *ProvisionDetailsDTO {
	ingredients := make([]ProvisionIngredientDTO, len(provision.ProvisionIngredients))
	for i, pi := range provision.ProvisionIngredients {
		ingredients[i] = ProvisionIngredientDTO{
			Ingredient: *ingredientTypes.ConvertToIngredientResponseDTO(&pi.Ingredient),
			Quantity:   pi.Quantity,
		}
	}

	return &ProvisionDetailsDTO{
		ProvisionDTO: *MapToProvisionDTO(provision),
		Ingredients:  ingredients,
	}
}

func CreateToProvisionModel(dto *CreateProvisionDTO) *data.Provision {
	return &data.Provision{
		Name:                 dto.Name,
		AbsoluteVolume:       dto.AbsoluteVolume,
		NetCost:              dto.NetCost,
		UnitID:               dto.UnitID,
		PreparationInMinutes: dto.PreparationInMinutes,
		LimitPerDay:          dto.LimitPerDay,
		ProvisionIngredients: mapIngredientsToProvisionIngredients(dto.Ingredients),
	}
}

func UpdateToProvisionModels(provision *data.Provision, dto *UpdateProvisionDTO) (*ProvisionModels, error) {
	if dto == nil {
		return nil, errors.New("dto is nil")
	}

	if provision == nil || provision.ID == 0 {
		return nil, errors.New("empty existing provision fetched")
	}

	if dto.Name != nil {
		provision.Name = *dto.Name
	}
	if dto.AbsoluteVolume != nil {
		provision.AbsoluteVolume = *dto.AbsoluteVolume
	}
	if dto.NetCost != nil {
		provision.NetCost = *dto.NetCost
	}
	if dto.PreparationInMinutes != nil {
		provision.PreparationInMinutes = *dto.PreparationInMinutes
	}
	if dto.LimitPerDay != nil {
		provision.LimitPerDay = *dto.LimitPerDay
	}

	var ingredients []data.ProvisionIngredient

	if dto.Ingredients != nil {
		if len(dto.Ingredients) == 0 {
			ingredients = []data.ProvisionIngredient{}
		} else {
			for _, ingredient := range dto.Ingredients {
				temp := data.ProvisionIngredient{
					IngredientID: ingredient.IngredientID,
					Quantity:     ingredient.Quantity,
				}
				ingredients = append(ingredients, temp)
			}
		}
	}

	return &ProvisionModels{
		Provision:   provision,
		Ingredients: ingredients,
	}, nil
}

func mapIngredientsToProvisionIngredients(dto []ingredientTypes.SelectedIngredientDTO) []data.ProvisionIngredient {
	result := make([]data.ProvisionIngredient, len(dto))
	for i, item := range dto {
		result[i] = data.ProvisionIngredient{
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
		}
	}
	return result
}

func ParseJSONProvisionsFromString(provisionsJSON string, provisions []SelectedProvisionDTO) error {
	if provisionsJSON != "" {
		err := json.Unmarshal([]byte(provisionsJSON), &provisions)
		if err != nil {
			return err
		}

		for _, provision := range provisions {
			if provision.ProvisionID == 0 || provision.Volume <= 0 {
				return fmt.Errorf("invalid provision json input")
			}
		}
	}
	return nil
}
