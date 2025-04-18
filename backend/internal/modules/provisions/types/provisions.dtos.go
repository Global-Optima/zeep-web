package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type BaseProvisionDTO struct {
	Name                       string             `json:"name"`
	AbsoluteVolume             float64            `json:"absoluteVolume"`
	NetCost                    float64            `json:"netCost"`
	Unit                       unitTypes.UnitsDTO `json:"unit"`
	PreparationInMinutes       uint               `json:"preparationInMinutes"`
	DefaultExpirationInMinutes uint               `json:"defaultExpirationInMinutes"`
	LimitPerDay                uint               `json:"limitPerDay"`
}

type ProvisionDTO struct {
	ID uint `json:"id"`
	BaseProvisionDTO
}

type ProvisionDetailsDTO struct {
	ProvisionDTO
	Ingredients []ProvisionIngredientDTO `json:"ingredients"`
}

type ProvisionIngredientDTO struct {
	Ingredient ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity   float64                       `json:"quantity"`
}

type CreateProvisionDTO struct {
	Name                       string                                  `json:"name" binding:"required,min=2,max=255"`
	AbsoluteVolume             float64                                 `json:"absoluteVolume" binding:"required,gt=0"`
	NetCost                    float64                                 `json:"netCost" binding:"required,gt=0"`
	UnitID                     uint                                    `json:"unitId" binding:"required,gt=0"`
	PreparationInMinutes       uint                                    `json:"preparationInMinutes" binding:"required"`
	DefaultExpirationInMinutes uint                                    `json:"defaultExpirationInMinutes" binding:"required"`
	LimitPerDay                uint                                    `json:"limitPerDay" binding:"required"`
	Ingredients                []ingredientTypes.SelectedIngredientDTO `json:"ingredients" binding:"omitempty,dive"`
}

type UpdateProvisionDTO struct {
	Name                       *string                                 `json:"name" binding:"omitempty,min=2,max=255"`
	AbsoluteVolume             *float64                                `json:"absoluteVolume" binding:"omitempty,gt=0"`
	NetCost                    *float64                                `json:"netCost" binding:"omitempty,gte=0"`
	UnitID                     *uint                                   `json:"unitId" binding:"omitempty,gt=0"`
	PreparationInMinutes       *uint                                   `json:"preparationInMinutes" binding:"omitempty"`
	DefaultExpirationInMinutes *uint                                   `json:"defaultExpirationInMinutes" binding:"omitempty"`
	LimitPerDay                *uint                                   `json:"limitPerDay" binding:"omitempty"`
	Ingredients                []ingredientTypes.SelectedIngredientDTO `json:"ingredients" binding:"omitempty,dive"`
}

type ProvisionFilterDTO struct {
	utils.BaseFilter
	Search *string `form:"search"`
}

type SelectedProvisionDTO struct {
	ProvisionID uint    `json:"provisionId" binding:"required,gt=0"`
	Volume      float64 `json:"volume" binding:"required,gt=0"`
}
