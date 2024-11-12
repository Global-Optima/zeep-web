package types

type AdditiveDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
}

type AdditiveCategoryDTO struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	Additives []AdditiveDTO `json:"additives"`
}
