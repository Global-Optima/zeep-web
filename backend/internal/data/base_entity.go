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

func toUTC(t time.Time) time.Time {
	return t.UTC()
}

// GORM Hooks for BaseEntity
func (b *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = toUTC(b.CreatedAt)
	b.UpdatedAt = toUTC(b.UpdatedAt)
	return
}

func (b *BaseEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = toUTC(b.UpdatedAt)
	return
}

func (b *BaseEntity) AfterFind(tx *gorm.DB) (err error) {
	b.CreatedAt = toUTC(b.CreatedAt)
	b.UpdatedAt = toUTC(b.UpdatedAt)
	if b.DeletedAt.Valid {
		b.DeletedAt.Time = toUTC(b.DeletedAt.Time)
	}
	return
}
