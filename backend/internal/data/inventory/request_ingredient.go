package data

type RequestIngredient struct {
	ID             uint `gorm:"primaryKey"`
	StockRequestID uint `gorm:"column:stock_request_id;not null"`
	IngredientID   uint `gorm:"column:ingredient_id;not null"`
	Quantity       int  `gorm:"default:1;check:quantity > 0"`
}
