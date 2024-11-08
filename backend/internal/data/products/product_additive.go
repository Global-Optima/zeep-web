package data

type ProductAdditive struct {
	ID            uint `gorm:"primaryKey"`
	ProductSizeID uint `gorm:"column:product_size_id;not null"`
	AdditiveID    uint `gorm:"column:additive_id;not null"`
}
