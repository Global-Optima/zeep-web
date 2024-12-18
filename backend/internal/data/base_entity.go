package data

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"autoCreateTime" sort:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" sort:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
