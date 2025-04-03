package data

import "time"

type Size string

const (
	S Size = "S"
	M Size = "M"
	L Size = "L"

	DEFAULT_INGREDIENT_EXPIRATION_IN_DAYS = 1095 // 3 years
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

type RecipeStep struct {
	BaseEntity
	ProductID   uint    `gorm:"index;not null"`
	Product     Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Step        int     `gorm:"not null" sort:"step"`
	Name        string  `gorm:"size:100" sort:"name"`
	Description string  `gorm:"type:text"`
	ImageKey    string  `gorm:"size:2048"`
}

type ProductSize struct {
	BaseEntity
	Name                   string                  `gorm:"size:100;not null" sort:"name"`
	UnitID                 uint                    `gorm:"index,not null"`
	Unit                   Unit                    `gorm:"foreignKey:UnitID;constraint:OnDelete:CASCADE" sort:"unit"`
	BasePrice              float64                 `gorm:"not null" sort:"price"`
	Size                   float64                 `gorm:"not null"`
	ProductID              uint                    `gorm:"index;not null"`
	Product                Product                 `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" sort:"product"`
	DiscountID             uint                    `gorm:"index"`
	MachineId              string                  `gorm:"size:40;not null;unique" sort:"machineId"`
	AdditivesUpdatedAt     time.Time               `gorm:"index" sort:"additivesUpdatedAt"`
	Additives              []ProductSizeAdditive   `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	ProductSizeIngredients []ProductSizeIngredient `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
}

type ProductSizeIngredient struct {
	BaseEntity
	IngredientID  uint        `gorm:"index;not null"`
	Ingredient    Ingredient  `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	ProductSizeID uint        `gorm:"index;not null"`
	ProductSize   ProductSize `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	Quantity      float64     `gorm:"type:decimal(10,2);not null;check:quantity > 0"`
}

type AdditiveIngredient struct {
	BaseEntity
	IngredientID uint       `gorm:"index;not null"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	AdditiveID   uint       `gorm:"index;not null"`
	Additive     Additive   `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
	Quantity     float64    `gorm:"type:decimal(10,2);not null;check:quantity > 0"`
}

type Ingredient struct {
	BaseEntity
	Name                   string                  `gorm:"size:255;not null;index" sort:"name"`
	Calories               float64                 `gorm:"type:decimal(5,2);check:calories >= 0" sort:"calories"`
	Fat                    float64                 `gorm:"type:decimal(5,2);check:fat >= 0" sort:"fat"`
	Carbs                  float64                 `gorm:"type:decimal(5,2);check:carbs >= 0" sort:"carbs"`
	Proteins               float64                 `gorm:"type:decimal(5,2);check:proteins >= 0" sort:"proteins"`
	ExpirationInDays       int                     `gorm:"not null;default:0" sort:"expirationInDays"` // Changed to int
	UnitID                 uint                    `gorm:"not null"`                                   // Link to Unit
	Unit                   Unit                    `gorm:"foreignKey:UnitID;constraint:OnDelete:CASCADE"`
	CategoryID             uint                    `gorm:"not null"` // Link to IngredientCategory
	IsAllergen             bool                    `gorm:"default:false" sort:"isAllergen"`
	IngredientCategory     IngredientCategory      `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	StockMaterials         []StockMaterial         `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"` // New association
	ProductSizeIngredients []ProductSizeIngredient `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	StoreStocks            []StoreStock            `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	AdditiveIngredients    []AdditiveIngredient    `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
}

type IngredientCategory struct {
	BaseEntity
	Name        string       `gorm:"size:255;not null;uniqueIndex"`
	Description string       `gorm:"type:text"`
	Ingredients []Ingredient `gorm:"foreignKey:CategoryID"`
}

type Provision struct {
	BaseEntity
	Name                 string                `gorm:"size:255;not null;uniqueIndex" sort:"name"`
	Description          string                `gorm:"type:text"`
	AbsoluteVolume       float64               `gorm:"type:decimal(10,2);not null;check:absoluteVolume > 0"`
	NetCost              float64               `gorm:"type:decimal(10,2);not null;check:netCost >= 0"`
	UnitID               uint                  `gorm:"not null"`
	Unit                 Unit                  `gorm:"foreignKey:UnitID;constraint:CASCADE"`
	PreparationInMinutes uint                  `gorm:"not null" sort:"preparationInMinutes"`
	LimitPerDay          uint                  `gorm:"not null" sort:"limitPerDay"`
	ProvisionIngredients []ProvisionIngredient `gorm:"foreignKey:ProvisionID"`
}

type ProvisionIngredient struct {
	BaseEntity
	IngredientID uint       `gorm:"index;not null"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE"`
	ProvisionID  uint       `gorm:"index;not null"`
	Provision    Provision  `gorm:"foreignKey:ProvisionID;constraint:OnDelete:CASCADE"`
	Quantity     float64    `gorm:"type:decimal(10,2);not null;check:quantity > 0"`
}

type ProductSizeAdditive struct {
	BaseEntity
	ProductSizeID uint        `gorm:"index;not null"`
	AdditiveID    uint        `gorm:"index;not null"`
	IsDefault     bool        `gorm:"not null" sort:"isDefault"`
	IsHidden      bool        `gorm:"default:false"`
	ProductSize   ProductSize `gorm:"foreignKey:ProductSizeID;constraint:OnDelete:CASCADE"`
	Additive      Additive    `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
}

type ProductCategory struct {
	BaseEntity
	Name        string    `gorm:"size:100;not null" sort:"name"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"foreignKey:CategoryID"`
}

type Additive struct {
	BaseEntity
	Name                 string                `gorm:"size:255;not null;index" sort:"name"`
	Description          string                `gorm:"type:text"`
	BasePrice            float64               `gorm:"type:decimal(10,2);default:0" sort:"basePrice"`
	Size                 float64               `gorm:"not null"`
	UnitID               uint                  `gorm:"index,not null"`
	Unit                 Unit                  `gorm:"foreignKey:UnitID;constraint:OnDelete:SET NULL" sort:"unit"`
	AdditiveCategoryID   uint                  `gorm:"index"`
	Category             AdditiveCategory      `gorm:"foreignKey:AdditiveCategoryID;constraint:OnDelete:SET NULL" sort:"category"`
	ImageKey             *StorageImageKey      `gorm:"size:2048"`
	MachineId            string                `gorm:"size:40;not null;unique" sort:"machineId"`
	IngredientsUpdatedAt time.Time             `gorm:"index" sort:"ingredientsUpdatedAt"`
	ProductSizeAdditives []ProductSizeAdditive `gorm:"foreignKey:AdditiveID;constraint:OnDelete:CASCADE"`
	StoreAdditives       []StoreAdditive       `gorm:"foreignKey:AdditiveID"`
	Ingredients          []AdditiveIngredient  `gorm:"foreignKey:AdditiveID"`
}

type AdditiveCategory struct {
	BaseEntity
	Name             string     `gorm:"size:100;not null" sort:"name"`
	Description      string     `gorm:"type:text"`
	Additives        []Additive `gorm:"foreignKey:AdditiveCategoryID"`
	IsMultipleSelect bool       `gorm:"default:true" sort:"isMultipleSelect"`
	IsRequired       bool       `gorm:"not null;default:false" sort:"isRequired"`
}
