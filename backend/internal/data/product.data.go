package data

import (
	"time"
)

type Size string

const (
	S Size = "S"
	M Size = "M"
	L Size = "L"
)

func (s Size) ToString() string {
	return string(s)
}

func IsValidSize(size Size) bool {
	switch size {
	case S, M, L:
		return true
	}
	return false
}

type Product struct {
	BaseEntity
	Name         string           `gorm:"size:100;not null" sort:"name"`
	Description  string           `gorm:"type:text"`
	ImageURL     string           `gorm:"size:2048"`
	VideoURL     string           `gorm:"size:2048"`
	CategoryID   *uint            `gorm:"index;not null"`
	Category     *ProductCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" sort:"categories"`
	RecipeSteps  []RecipeStep     `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	ProductSizes []ProductSize    `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type RecipeStep struct {
	BaseEntity
	ProductID   uint    `gorm:"index;not null"`
	Product     Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Step        int     `gorm:"not null" sort:"step"`
	Name        string  `gorm:"size:100" sort:"name"`
	Description string  `gorm:"type:text"`
	ImageURL    string  `gorm:"size:2048"`
}

type ProductSize struct {
	BaseEntity
	Name               string  `gorm:"size:100;not null" sort:"name"`
	Measure            string  `gorm:"size:50" sort:"measure"`
	BasePrice          float64 `gorm:"not null" sort:"price"`
	Size               int     `gorm:"not null" sort:"size"`
	IsDefault          bool    `gorm:"default:false" sort:"isDefault"`
	ProductID          uint    `gorm:"index;not null"`
	Product            Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	DiscountID         uint
	Additives          []ProductSizeAdditive `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	ProductIngredients []ProductIngredient   `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
}

type ProductIngredient struct {
	BaseEntity
	ItemIngredientID uint           `gorm:"index;not null"`
	ItemIngredient   ItemIngredient `gorm:"foreignKey:ItemIngredientID;constraint:OnDelete:CASCADE"`
	ProductSizeID    uint           `gorm:"index;not null"`
	ProductSize      ProductSize    `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	Ingredient       Ingredient     `gorm:"foreignKey:ItemIngredientID;constraint:OnDelete:CASCADE"`
}

type AdditiveIngredient struct {
	BaseEntity
	ItemIngredientID uint           `gorm:"index;not null"`
	ItemIngredient   ItemIngredient `gorm:"foreignKey:ItemIngredientID;constraint:OnDelete:CASCADE"`
	AdditiveID       uint           `gorm:"index;not null"`
	Additive         Additive       `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	Ingredient       Ingredient     `gorm:"foreignKey:ItemIngredientID;constraint:OnDelete:CASCADE"`
}

type ItemIngredient struct {
	BaseEntity
	IngredientID uint       `gorm:"not null;index"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	ItemID       uint       `gorm:"not null;index"`
	Product      Product    `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Quantity     float64    `gorm:"type:decimal(10,2);not null;check:quantity > 0"`
}

type Ingredient struct {
	BaseEntity
	Name      string              `gorm:"size:255;not null;index" sort:"name"`
	Calories  float64             `gorm:"type:decimal(5,2);check:calories >= 0" sort:"calories"`
	Fat       float64             `gorm:"type:decimal(5,2);check:fat >= 0" sort:"fat"`
	Carbs     float64             `gorm:"type:decimal(5,2);check:carbs >= 0" sort:"carbs"`
	Proteins  float64             `gorm:"type:decimal(5,2);check:proteins >= 0" sort:"proteins"`
	ExpiresAt *time.Time          `gorm:"type:timestamp" sort:"expiresAt"`
	Products  []ProductIngredient `gorm:"foreignKey:ItemIngredientID"`
}

type ProductSizeAdditive struct {
	BaseEntity
	ProductSizeID uint        `gorm:"index;not null"`
	AdditiveID    uint        `gorm:"index;not null"`
	IsDefault     bool        `gorm:"default:true" sort:"isDefault"`
	ProductSize   ProductSize `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	Additive      Additive    `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
}

type ProductCategory struct {
	BaseEntity
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"foreignKey:CategoryID"`
}

type Additive struct {
	BaseEntity
	Name               string            `gorm:"size:255;not null;index" sort:"name"`
	Description        string            `gorm:"type:text"`
	BasePrice          float64           `gorm:"type:decimal(10,2);default:0"`
	Size               string            `gorm:"size:200"`
	AdditiveCategoryID uint              `gorm:"index"`
	Category           *AdditiveCategory `gorm:"foreignKey:AdditiveCategoryID;constraint:OnDelete:SET NULL"`
	ImageURL           string            `gorm:"size:2048"`
	StoreAdditives     []StoreAdditive   `gorm:"foreignKey:AdditiveID"`
}

type AdditiveCategory struct {
	BaseEntity
	Name             string     `gorm:"size:100;not null"`
	Description      string     `gorm:"type:text"`
	Additives        []Additive `gorm:"foreignKey:AdditiveCategoryID"`
	IsMultipleSelect bool       `gorm:"default:true" json:"is_multiple_select"`
}
