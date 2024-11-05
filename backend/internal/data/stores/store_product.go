package data

type StoreProduct struct {
	ID          uint `gorm:"primaryKey"`
	ProductID   uint `gorm:"column:product_id;not null"`
	StoreID     uint `gorm:"column:store_id;not null"`
	IsAvailable bool `gorm:"default:true"`
}
