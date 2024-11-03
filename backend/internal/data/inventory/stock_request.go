package data

import "time"

type StockRequest struct {
	ID               uint      `gorm:"primaryKey"`
	CityWarehouseID  *uint     `gorm:"column:city_warehouse_id"`
	StoreWarehouseID uint      `gorm:"column:store_warehouse_id;not null"`
	Status           string    `gorm:"size:50;not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}
