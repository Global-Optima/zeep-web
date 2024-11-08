package data

import (
	"time"
)

type Additive struct {
	ID                 uint      `gorm:"primaryKey"`
	Name               string    `gorm:"size:255;not null"`
	Description        string    `gorm:"type:text"`
	BasePrice          float64   `gorm:"type:decimal(10,2);default:0"`
	Size               string    `gorm:"size:200"`
	AdditiveCategoryID *uint     `gorm:"column:additive_category_id"` // Foreign key to AdditiveCategory
	ImageURL           string    `gorm:"size:2048"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}
