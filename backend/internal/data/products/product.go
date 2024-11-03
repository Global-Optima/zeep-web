package data

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	ImageURL    string    `gorm:"size:2048"`
	VideoURL    string    `gorm:"size:2048"`
	CategoryID  *uint     `gorm:"column:category_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
