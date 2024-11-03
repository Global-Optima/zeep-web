package data

import "time"

type Referral struct {
	ID         uint      `gorm:"primaryKey"`
	CustomerID uint      `gorm:"column:customer_id;not null"`
	RefereeID  uint      `gorm:"column:referral_id;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
