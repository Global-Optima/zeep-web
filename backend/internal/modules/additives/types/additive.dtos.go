package types

import (
	provisionsTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/types"
	"mime/multipart"

	ingredientTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/ingredients/types"
	unitTypes "github.com/Global-Optima/zeep-web/backend/internal/modules/units/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
)

const (
	PROVISIONS_FORM_DATA_KEY  = "provisions"
	INGREDIENTS_FORM_DATA_KEY = "ingredients"
)

type AdditiveCategoriesFilterQuery struct {
	utils.BaseFilter
	IncludeEmpty     *bool   `form:"includeEmpty"`
	ProductSizeId    *uint   `form:"productSizeId"`
	IsMultipleSelect *bool   `form:"isMultipleSelect"`
	IsRequired       *bool   `form:"isRequired"`
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
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsMultipleSelect bool   `json:"isMultipleSelect"`
	IsRequired       bool   `json:"isRequired"`
}

// BaseAdditiveDTO should not be returned directly as a response,
// instead wrap it into another struct with more info like ID and etc
type BaseAdditiveDTO struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	BasePrice   float64             `json:"basePrice"`
	ImageURL    string              `json:"imageUrl"`
	Size        float64             `json:"size"`
	Unit        unitTypes.UnitsDTO  `json:"unit"`
	Category    AdditiveCategoryDTO `json:"category"`
	MachineId   string              `json:"machineId"`
}

type AdditiveDTO struct {
	ID uint `json:"id"`
	BaseAdditiveDTO
}

type AdditiveDetailsDTO struct {
	AdditiveDTO
	Ingredients []AdditiveIngredientDTO `json:"ingredients"`
	Provisions  []AdditiveProvisionDTO  `json:"provisions"`
}

type AdditiveProvisionDTO struct {
	Volume    float64                      `json:"volume"`
	Provision provisionsTypes.ProvisionDTO `json:"provision"`
}

type AdditiveIngredientDTO struct {
	Quantity   float64                       `json:"quantity"`
	Ingredient ingredientTypes.IngredientDTO `json:"ingredient"`
}

type AdditiveCategoryItemDTO struct {
	ID uint `json:"id"`
	BaseAdditiveDTO
}

type AdditiveCategoryDetailsDTO struct {
	AdditiveCategoryDTO
	AdditivesCount int `json:"additivesCount"`
}

type CreateAdditiveCategoryDTO struct {
	Name             string  `json:"name" binding:"required"`
	Description      *string `json:"description" binding:"omitempty"`
	IsMultipleSelect bool    `json:"isMultipleSelect"`
	IsRequired       bool    `json:"isRequired"`
}

type UpdateAdditiveCategoryDTO struct {
	Name             *string `json:"name" binding:"min=0,omitempty"`
	Description      *string `json:"description" binding:"omitempty"`
	IsMultipleSelect *bool   `json:"isMultipleSelect"`
	IsRequired       *bool   `json:"isRequired"`
}

type UpdateAdditiveDTO struct {
	Name               *string                                 `form:"name" binding:"min=0,omitempty"`
	Description        *string                                 `form:"description" binding:"omitempty"`
	BasePrice          *float64                                `form:"basePrice" binding:"omitempty,gte=0"`
	Size               *float64                                `form:"size" binding:"omitempty,gt=0"`
	UnitID             *uint                                   `form:"unitId" binding:"omitempty,gt=0"`
	AdditiveCategoryID *uint                                   `form:"additiveCategoryId" binding:"omitempty,gt=0"`
	MachineId          *string                                 `form:"machineId" binding:"omitempty"`
	Ingredients        []ingredientTypes.SelectedIngredientDTO `json:"-"`
	Provisions         []provisionsTypes.SelectedProvisionDTO  `json:"-"`
	Image              *multipart.FileHeader
	DeleteImage        bool `form:"deleteImage"`
}

type AdditiveCategoryDTO struct {
	ID uint `json:"id"`
	BaseAdditiveCategoryDTO
}

type CreateAdditiveDTO struct {
	Name               string                                  `form:"name" binding:"required"`
	Description        *string                                 `form:"description" binding:"omitempty"`
	BasePrice          float64                                 `form:"basePrice" binding:"gte=0"`
	Size               float64                                 `form:"size" binding:"gt=0"`
	UnitID             uint                                    `form:"unitId" binding:"gt=0"`
	AdditiveCategoryID uint                                    `form:"additiveCategoryId" binding:"gt=0"`
	MachineId          string                                  `form:"machineId" binding:"required"`
	Ingredients        []ingredientTypes.SelectedIngredientDTO `json:"-"`
	Provisions         []provisionsTypes.SelectedProvisionDTO  `json:"-"`
	Image              *multipart.FileHeader
}
