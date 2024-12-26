package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type AdditiveCategoriesFilterQuery struct {
	utils.BaseFilter
	ProductSizeId *uint   `form:"productSizeId"`
	Search        *string `form:"search"`
}

type AdditiveFilterQuery struct {
	utils.BaseFilter
	Search        *string  `form:"search"`
	MinPrice      *float64 `form:"minPrice"`
	MaxPrice      *float64 `form:"maxPrice"`
	CategoryID    *uint    `form:"categoryId"`
	ProductSizeID *uint    `form:"productSizeId"`
}

type AdditiveDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
	Size        string  `json:"size"`
	Category    struct {
		ID               uint   `json:"id"`
		Name             string `json:"name"`
		IsMultipleSelect bool   `json:"isMultipleSelect"`
	} `json:"category"`
}

type AdditiveCategoryItemDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
	Size        string  `json:"size"`
	CategoryID  uint    `json:"categoryId"`
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
	ID               uint   `json:"id" binding:"required"`
	Name             string `json:"name" binding:"omitempty"`
	Description      string `json:"description" binding:"omitempty"`
	IsMultipleSelect *bool  `json:"isMultipleSelect"`
}

type UpdateAdditiveDTO struct {
	ID                 uint     `json:"id" binding:"required"`
	Name               string   `json:"name" binding:"omitempty"`
	Description        string   `json:"description" binding:"omitempty"`
	Price              *float64 `json:"price" binding:"omitempty,gte=0"`
	ImageURL           *string  `json:"imageUrl" binding:"omitempty"`
	Size               *string  `json:"size" binding:"omitempty"`
	AdditiveCategoryID *uint    `json:"additiveCategoryId" binding:"omitempty"`
}

type AdditiveCategoryResponseDTO struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsMultipleSelect bool   `json:"isMultipleSelect"`
}

type CreateAdditiveDTO struct {
	Name               string  `json:"name" binding:"required"`
	Description        string  `json:"description" binding:"required"`
	Price              float64 `json:"price" binding:"required,gte=0"`
	ImageURL           string  `json:"imageUrl" binding:"omitempty"`
	Size               string  `json:"size" binding:"required"`
	AdditiveCategoryID uint    `json:"additiveCategoryId" binding:"required"`
}
