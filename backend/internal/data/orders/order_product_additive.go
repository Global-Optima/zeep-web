package data

type OrderProductAdditive struct {
	ID             uint `gorm:"primaryKey"`
	AdditiveID     uint `gorm:"column:additive_id;not null"`
	OrderProductID uint `gorm:"column:order_products_id;not null"`
}
