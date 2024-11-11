package types

type StoreProductDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
	BasePrice   float64 `json:"basePrice"`
}

type StoreProductDetailsDTO struct {
	ID               uint                 `json:"id"`
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	ImageURL         string               `json:"imageUrl"`
	Sizes            []ProductSizeDTO     `json:"sizes"`
	DefaultAdditives []ProductAdditiveDTO `json:"defaultAdditives"`
	RecipeSteps      []RecipeStepDTO      `json:"recipeSteps"`
}

type ProductSizeDTO struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	BasePrice float64 `json:"basePrice"`
	Measure   string  `json:"measure"`
}

type ProductAdditiveDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type RecipeStepDTO struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	Step        int    `json:"step"`
}
