package data

type StoreProductSize struct {
	ID            uint    `gorm:"primaryKey"`
	ProductSizeID uint    `gorm:"column:product_size_id;not null"`
	StoreID       uint    `gorm:"column:store_id;not null"`
	Price         float64 `gorm:"type:decimal(10,2);default:0"`
}
