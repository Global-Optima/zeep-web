package data

type CustomerAddress struct {
	ID         uint   `gorm:"primaryKey"`
	CustomerID uint   `gorm:"column:customer_id;not null"`
	Address    string `gorm:"size:255;not null"`
	Longitude  string `gorm:"size:20"`
	Latitude   string `gorm:"size:20"`
}
