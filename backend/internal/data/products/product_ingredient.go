package data

type ProductIngredient struct {
	ID               uint `gorm:"primaryKey"`
	ItemIngredientID uint `gorm:"column:item_ingredient_id;not null"`
	ProductID        uint `gorm:"column:product_id;not null"`
}
