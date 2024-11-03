package data

type FacilityAddress struct {
	ID        uint    `gorm:"primaryKey"`
	Address   string  `gorm:"size:255;not null"`
	Longitude float64 `gorm:"type:decimal(9,6)"`
	Latitude  float64 `gorm:"type:decimal(9,6)"`
}
