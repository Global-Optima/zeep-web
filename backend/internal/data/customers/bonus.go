package data

import "time"

type Bonus struct {
	ID         uint      `gorm:"primaryKey"`
	Bonuses    float64   `gorm:"type:decimal(10,2);check:bonuses >= 0"`
	CustomerID uint      `gorm:"column:customer_id;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	ExpiresAt  *time.Time
}
