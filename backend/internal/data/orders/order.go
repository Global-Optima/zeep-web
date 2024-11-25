package data

import "time"

type Order struct {
	ID                uint    `gorm:"primaryKey"`
	CustomerName      string  `gorm:"size:255;not null"`
	CustomerID        *uint   `gorm:"column:customer_id"`
	StoreID           uint    `gorm:"column:store_id"`
	Price             float64 `gorm:"type:decimal(10,2);check:price >= 0"`
	Date              *time.Time
	Status            string    `gorm:"size:50"`
	EmployeeID        *uint     `gorm:"column:employee_id"`
	DeliveryAddressID *uint     `gorm:"column:delivery_address_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}
