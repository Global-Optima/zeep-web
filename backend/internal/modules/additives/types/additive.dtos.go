package types

import (
	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

type AdditiveCategoriesFilterQuery struct {
	utils.BaseFilter
	IncludeEmpty     *bool   `form:"includeEmpty"`
	ProductSizeId    *uint   `form:"productSizeId"`
	IsMultipleSelect *bool   `form:"isMultipleSelect"`
	Search           *string `form:"search"`
}

type AdditiveFilterQuery struct {
	utils.BaseFilter
	Search        *string  `form:"search"`
	MinPrice      *float64 `form:"minPrice"`
	MaxPrice      *float64 `form:"maxPrice"`
	CategoryID    *uint    `form:"categoryId"`
	ProductSizeID *uint    `form:"productSizeId"`
}

type BaseAdditiveCategoryDTO struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsMultipleSelect bool   `json:"isMultipleSelect"`
}

// BaseAdditiveDTO should not be returned directly as a response,
// instead wrap it into another struct with more info like ID and etc
type BaseAdditiveDTO struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	BasePrice   float64                 `json:"basePrice"`
	ImageURL    string                  `json:"imageUrl"`
	Size        int                     `json:"size"`
	Unit        unitTypes.UnitsDTO      `json:"unit"`
	Category    BaseAdditiveCategoryDTO `json:"category"`
}

type AdditiveDTO struct {
	ID uint `json:"id"`
	BaseAdditiveDTO
}

type AdditiveDetailsDTO struct {
	AdditiveDTO
	Ingredients []AdditiveIngredientDTO `json:"ingredients"`
}

type AdditiveIngredientDTO struct {
	Quantity   float64                       `json:"quantity"`
	Ingredient ingredientTypes.IngredientDTO `json:"ingredient"`
}

type BaseAdditiveCategoryItemDTO struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	BasePrice   float64            `json:"basePrice"`
	ImageURL    string             `json:"imageUrl"`
	Size        int                `json:"size"`
	Unit        unitTypes.UnitsDTO `json:"unit"`
	CategoryID  uint               `json:"categoryId"`
}

type AdditiveCategoryItemDTO struct {
	ID uint `json:"id"`
	BaseAdditiveCategoryItemDTO
}

type AdditiveCategoryDTO struct {
	ID               uint                      `json:"id"`
	Name             string                    `json:"name"`
	Description      string                    `json:"description"`
	Additives        []AdditiveCategoryItemDTO `json:"additives"`
	IsMultipleSelect bool                      `json:"isMultipleSelect"`
}

type CreateAdditiveCategoryDTO struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description" binding:"omitempty"`
	IsMultipleSelect bool   `json:"isMultipleSelect"`
}

type UpdateAdditiveCategoryDTO struct {
	Name             *string `json:"name" binding:"omitempty"`
	Description      *string `json:"description" binding:"omitempty"`
	IsMultipleSelect *bool   `json:"isMultipleSelect"`
}

type UpdateAdditiveDTO struct {
	Name               string                  `form:"name" binding:"omitempty"`
	Description        string                  `form:"description" binding:"omitempty"`
	BasePrice          *float64                `form:"basePrice" binding:"omitempty,gte=0"`
	Size               *int                    `form:"size" binding:"omitempty,gt=0"`
	UnitID             *uint                   `form:"unitId" binding:"omitempty,gt=0"`
	AdditiveCategoryID *uint                   `form:"additiveCategoryId" binding:"omitempty,gt=0"`
	Ingredients        []SelectedIngredientDTO `form:"ingredients" binding:"omitempty,dive"`
}

type AdditiveCategoryResponseDTO struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsMultipleSelect bool   `json:"isMultipleSelect"`
}

type CreateAdditiveDTO struct {
	Name               string                  `form:"name" binding:"required"`
	Description        string                  `form:"description" binding:"required"`
	BasePrice          float64                 `form:"basePrice" binding:"required,gte=0"`
	Size               int                     `form:"size" binding:"required,gt=0"`
	UnitID             uint                    `form:"unitId" binding:"required,gt=0"`
	AdditiveCategoryID uint                    `form:"additiveCategoryId" binding:"required,gt=0"`
	Ingredients        []SelectedIngredientDTO `json:"-"`
}

type SelectedIngredientDTO struct {
	IngredientID uint    `json:"ingredientId" binding:"required,gt=0"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
}
