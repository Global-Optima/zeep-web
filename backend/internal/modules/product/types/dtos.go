package types

type ProductCatalogDTO struct {
	ProductID          uint    `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	Category           string  `json:"category"`
	ProductImageURL    string  `json:"product_image_url"`
	Price              float64 `json:"price"`
	IsAvailable        bool    `json:"is_available"`
	IsOutOfStock       bool    `json:"is_out_of_stock"`
}

type ProductDTO struct {
	ProductID          uint           `json:"product_id"`
	ProductName        string         `json:"product_name"`
	ProductDescription string         `json:"product_description"`
	Category           string         `json:"category"`
	ProductImageURL    string         `json:"product_image_url"`
	ProductVideoURL    string         `json:"product_video_url"`
	Price              float64        `json:"price"`
	IsAvailable        bool           `json:"is_available"`
	IsOutOfStock       bool           `json:"is_out_of_stock"`
	Additives          []AdditivesDTO `json:"additives"`
	DefaultAdditives   []AdditivesDTO `json:"default_additives"`
	Sizes              []SizeDTO      `json:"sizes"`
	Nutrition          NutritionDTO   `json:"nutrition"`
}

type AdditivesDTO struct {
	AdditiveID          uint    `json:"additive_id"`
	AdditiveName        string  `json:"additive_name"`
	AdditiveDescription string  `json:"additive_description"`
	AdditiveCategory    string  `json:"additive_category"`
	AdditivePrice       float64 `json:"additive_price"`
}

type SizeDTO struct {
	SizeID    uint    `json:"size_id"`
	SizeName  string  `json:"size_name"`
	Size      float64 `json:"size"`
	Measure   string  `json:"measure"`
	Price     float64 `json:"price"`
	IsDefault bool    `json:"is_default"`
}

type NutritionDTO struct {
	Calories      float64 `json:"calories"`
	Fat           float64 `json:"fat"`
	Carbohydrates float64 `json:"carbohydrates"`
	Proteins      float64 `json:"proteins"`
}
