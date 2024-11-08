package data

import (
	"gorm.io/gorm"
)

type ProductSize struct {
	gorm.Model
	Name       string  `gorm:"size:100;not null"`
	Measure    string  `gorm:"size:50"`
	BasePrice  float64 `gorm:"not null"`
	Size       int     `gorm:"not null"`
	IsDefault  bool    `gorm:"default:false"`
	ProductID  uint    `gorm:"not null"`
	DiscountID uint
}
