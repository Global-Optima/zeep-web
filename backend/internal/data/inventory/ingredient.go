package data

import "time"

type Ingredient struct {
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"size:255;not null"`
	Calories  float64    `gorm:"type:decimal(5,2);check:calories >= 0"`
	Fat       float64    `gorm:"type:decimal(5,2);check:fat >= 0"`
	Carbs     float64    `gorm:"type:decimal(5,2);check:carbs >= 0"`
	Proteins  float64    `gorm:"type:decimal(5,2);check:proteins >= 0"`
	ExpiresAt *time.Time `gorm:"type:timestamp"`
}
