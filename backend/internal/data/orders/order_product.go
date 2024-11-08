package data

type OrderProduct struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"column:order_id;not null"`
	ProductID uint `gorm:"column:product_id;not null"`
	Quantity  int  `gorm:"default:1;check:quantity > 0"`
}
