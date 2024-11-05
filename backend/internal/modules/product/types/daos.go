package types

type ProductDAO struct {
	ProductID          uint          `gorm:"column:id" json:"product_id"`
	ProductName        string        `gorm:"column:name" json:"product_name"`
	ProductDescription string        `gorm:"column:description" json:"product_description"`
	CategoryID         uint          `gorm:"column:category_id" json:"category_id"`
	CategoryName       string        `gorm:"column:category_name" json:"category_name"`
	ProductImageURL    string        `gorm:"column:image_url" json:"product_image_url"`
	ProductVideoURL    string        `gorm:"column:video_url" json:"product_video_url"`
	BasePrice          float64       `gorm:"column:base_price" json:"base_price"`
	IsAvailable        bool          `gorm:"column:is_available" json:"is_available"`
	IsOutOfStock       bool          `gorm:"column:is_out_of_stock" json:"is_out_of_stock"`
	Sizes              []SizeDAO     `gorm:"-" json:"sizes"`
	Additives          []AdditiveDAO `gorm:"-" json:"additives"`
	Nutrition          NutritionDAO  `gorm:"-" json:"nutrition"`
}

type SizeDAO struct {
	SizeName  string  `gorm:"column:name" json:"size_name"`
	Size      float64 `gorm:"column:size" json:"size"`
	Measure   string  `gorm:"column:measure" json:"measure"`
	Price     float64 `gorm:"column:price" json:"price"`
	IsDefault bool    `gorm:"column:is_default" json:"is_default"`
}

type AdditiveDAO struct {
	AdditiveID          uint    `gorm:"column:id" json:"additive_id"`
	AdditiveName        string  `gorm:"column:name" json:"additive_name"`
	AdditiveDescription string  `gorm:"column:description" json:"additive_description"`
	AdditiveCategory    string  `gorm:"column:category_name" json:"additive_category"`
	AdditivePrice       float64 `gorm:"column:price" json:"additive_price"`
}

type NutritionDAO struct {
	Calories      float64 `gorm:"column:calories" json:"calories"`
	Fat           float64 `gorm:"column:fat" json:"fat"`
	Carbohydrates float64 `gorm:"column:carbs" json:"carbohydrates"`
	Proteins      float64 `gorm:"column:proteins" json:"proteins"`
}
