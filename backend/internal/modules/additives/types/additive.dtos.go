package types

import "github.com/Global-Optima/zeep-web/backend/pkg/utils"

type AdditiveCategoriesFilterQuery struct {
	ProductSizeId *uint   `form:"productSizeId"`
	Search        *string `form:"search"`
}

type AdditiveFilterQuery struct {
	Search        *string  `form:"search"`
	MinPrice      *float64 `form:"minPrice"`
	MaxPrice      *float64 `form:"maxPrice"`
	CategoryID    *uint    `form:"categoryId"`
	ProductSizeID *uint    `form:"productSizeId"`
	Pagination    *utils.Pagination
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
