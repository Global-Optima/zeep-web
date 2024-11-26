package data

import (
	"time"
)

type Product struct {
	BaseEntity
	Name               string                   `gorm:"size:100;not null"`
	Description        string                   `gorm:"type:text"`
	ImageURL           string                   `gorm:"size:2048"`
	VideoURL           string                   `gorm:"size:2048"`
	CategoryID         *uint                    `gorm:"index;not null"`
	Category           *ProductCategory         `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	RecipeSteps        []RecipeStep             `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	ProductSizes       []ProductSize            `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	DefaultAdditives   []DefaultProductAdditive `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	ProductIngredients []ProductIngredient      `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type RecipeStep struct {
	BaseEntity
	ProductID   uint    `gorm:"index;not null"`
	Product     Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Step        int     `gorm:"not null"`
	Name        string  `gorm:"size:100"`
	Description string  `gorm:"type:text"`
	ImageURL    string  `gorm:"size:2048"`
}

type ProductSize struct {
	BaseEntity
	Name       string  `gorm:"size:100;not null"`
	Measure    string  `gorm:"size:50"`
	BasePrice  float64 `gorm:"not null"`
	Size       int     `gorm:"not null"`
	IsDefault  bool    `gorm:"default:false"`
	ProductID  uint    `gorm:"index;not null"`
	Product    Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	DiscountID uint
	Additives  []ProductAdditive `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
}

type ProductIngredient struct {
	BaseEntity
	ItemIngredientID uint           `gorm:"index;not null"`
	ItemIngredient   ItemIngredient `gorm:"foreignKey:ItemIngredientID;constraint:OnDelete:CASCADE"`
	ProductID        uint           `gorm:"index;not null"`
	Product          Product        `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Ingredient       Ingredient     `gorm:"foreignKey:ItemIngredientID;constraint:OnDelete:CASCADE"`
}

type ProductAdditive struct {
	BaseEntity
	ProductSizeID uint        `gorm:"index;not null"`
	AdditiveID    uint        `gorm:"index;not null"`
	Additive      Additive    `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
	ProductSize   ProductSize `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
}

type Ingredient struct {
	BaseEntity
	Name      string              `gorm:"size:255;not null;index"`
	Calories  float64             `gorm:"type:decimal(5,2);check:calories >= 0"`
	Fat       float64             `gorm:"type:decimal(5,2);check:fat >= 0"`
	Carbs     float64             `gorm:"type:decimal(5,2);check:carbs >= 0"`
	Proteins  float64             `gorm:"type:decimal(5,2);check:proteins >= 0"`
	ExpiresAt *time.Time          `gorm:"type:timestamp"`
	Products  []ProductIngredient `gorm:"foreignKey:ItemIngredientID"`
}

type DefaultProductAdditive struct {
	BaseEntity
	ProductID  uint     `gorm:"index;not null"`
	AdditiveID uint     `gorm:"index;not null"`
	Product    Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Additive   Additive `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
}

type ProductCategory struct {
	BaseEntity
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"foreignKey:CategoryID"`
}

type Additive struct {
	BaseEntity
	Name               string            `gorm:"size:255;not null;index"`
	Description        string            `gorm:"type:text"`
	BasePrice          float64           `gorm:"type:decimal(10,2);default:0"`
	Size               string            `gorm:"size:200"`
	AdditiveCategoryID uint              `gorm:"index"`
	Category           *AdditiveCategory `gorm:"foreignKey:AdditiveCategoryID;constraint:OnDelete:SET NULL"`
	ImageURL           string            `gorm:"size:2048"`

	StorePrice float64 `gorm:"-"`
}

type AdditiveCategory struct {
	BaseEntity
	Name             string     `gorm:"size:100;not null"`
	Description      string     `gorm:"type:text"`
	Additives        []Additive `gorm:"foreignKey:AdditiveCategoryID"`
	IsMultipleSelect bool       `gorm:"default:true" json:"is_multiple_select"`
}
