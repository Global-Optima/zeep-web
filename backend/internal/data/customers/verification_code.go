package data

import "time"

type VerificationCode struct {
	ID         uint      `gorm:"primaryKey"`
	CustomerID uint      `gorm:"column:customer_id;not null"`
	Code       string    `gorm:"size:6;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  *time.Time
	ExpiresAt  time.Time `gorm:"not null"`
}
