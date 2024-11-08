package data

type DefaultProductAdditive struct {
	ID         uint `gorm:"primaryKey"`
	ProductID  uint `gorm:"column:product_id;not null"`
	AdditiveID uint `gorm:"column:additive_id;not null"`
}
