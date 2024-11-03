package data

type StoreAdditive struct {
	ID         uint    `gorm:"primaryKey"`
	AdditiveID uint    `gorm:"column:additive_id;not null"`
	StoreID    uint    `gorm:"column:store_id;not null"`
	Price      float64 `gorm:"type:decimal(10,2);default:0"`
}
