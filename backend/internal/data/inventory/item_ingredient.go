package data

type ItemIngredient struct {
	ID           uint    `gorm:"primaryKey"`
	IngredientID uint    `gorm:"column:ingredient_id;not null"`
	Name         string  `gorm:"size:255;not null"`
	Weight       float64 `gorm:"type:decimal(5,2);check:weight > 0"`
	Label        string  `gorm:"size:20"`
}
