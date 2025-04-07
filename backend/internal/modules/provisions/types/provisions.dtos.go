package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type BaseProvisionDTO struct {
	Name                 string                   `json:"name"`
	Description          string                   `json:"description"`
	AbsoluteVolume       float64                  `json:"absoluteVolume"`
	NetCost              float64                  `json:"netCost"`
	Unit                 unitTypes.UnitsDTO       `json:"unit"`
	PreparationInMinutes uint                     `json:"preparationInMinutes"`
	LimitPerDay          uint                     `json:"limitPerDay"`
	Ingredients          []ProvisionIngredientDTO `json:"ingredients"`
}

type ProvisionDTO struct {
	ID uint `json:"id"`
	BaseProvisionDTO
}

type ProvisionIngredientDTO struct {
	Ingredient ingredientTypes.IngredientDTO `json:"ingredient"`
	Quantity   float64                       `json:"quantity"`
}

type CreateProvisionDTO struct {
	Name                 string                                  `json:"name" binding:"required,min=2,max=255"`
	Description          *string                                 `json:"description" binding:"omitempty,max=1000"`
	AbsoluteVolume       float64                                 `json:"absoluteVolume" binding:"required,gt=0"`
	NetCost              float64                                 `json:"netCost" binding:"required,gte=0"`
	UnitID               uint                                    `json:"unitId" binding:"required,gt=0"`
	PreparationInMinutes uint                                    `json:"preparationInMinutes" binding:"required"`
	LimitPerDay          uint                                    `json:"limitPerDay" binding:"required"`
	Ingredients          []ingredientTypes.SelectedIngredientDTO `json:"ingredients" binding:"omitempty,dive"`
}

type UpdateProvisionDTO struct {
	Name                 *string                                 `json:"name" binding:"omitempty,min=2,max=255"`
	Description          *string                                 `json:"description" binding:"omitempty,max=1000"`
	AbsoluteVolume       *float64                                `json:"absoluteVolume" binding:"omitempty,gt=0"`
	NetCost              *float64                                `json:"netCost" binding:"omitempty,gte=0"`
	UnitID               *uint                                   `json:"unitId" binding:"omitempty,gt=0"`
	PreparationInMinutes *uint                                   `json:"preparationInMinutes" binding:"omitempty"`
	LimitPerDay          *uint                                   `json:"limitPerDay" binding:"omitempty"`
	Ingredients          []ingredientTypes.SelectedIngredientDTO `json:"ingredients" binding:"omitempty,dive"`
}

type ProvisionFilterDTO struct {
	utils.BaseFilter
	Search *string `form:"search"`
}
