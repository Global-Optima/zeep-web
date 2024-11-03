package data

type StoreWarehouseStock struct {
	ID               uint   `gorm:"primaryKey"`
	IngredientID     uint   `gorm:"column:ingredient_id;not null"`
	StoreWarehouseID uint   `gorm:"column:store_warehouse_id;not null"`
	Quantity         int    `gorm:"default:0;check:quantity >= 0"`
	Status           string `gorm:"size:50;not null"`
}
