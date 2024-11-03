package types

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
	SizeName  string  `json:"size_name"`
	Measure   string  `json:"measure"`
	Price     float64 `json:"price"`
	IsDefault bool    `json:"is_default"`
}

type NutritionDTO struct {
	Calories      int `json:"calories"`
	Fat           int `json:"fat"`
	Carbohydrates int `json:"carbohydrates"`
	Proteins      int `json:"proteins"`
}
