package data

import "time"

type CityWarehouse struct {
	ID                uint      `gorm:"primaryKey"`
	Name              string    `gorm:"size:255;not null"`
	FacilityAddressID *uint     `gorm:"column:facility_address_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}
