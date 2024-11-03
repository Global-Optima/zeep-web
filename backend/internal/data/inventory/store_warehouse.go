package data

type StoreWarehouse struct {
	ID              uint  `gorm:"primaryKey"`
	StoreID         uint  `gorm:"column:store_id;not null"`
	CityWarehouseID *uint `gorm:"column:city_warehouse_id"`
}
