package data

import (
	"time"
)

type Customer struct {
	BaseEntity
	Name       string            `gorm:"size:255;not null"`
	Password   string            `gorm:"size:255;not null"`
	Phone      string            `gorm:"size:15;unique"`
	IsVerified bool              `gorm:"default:false"`
	IsBanned   bool              `gorm:"default:false"`
	Addresses  []CustomerAddress `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Bonuses    []Bonus           `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Referrals  []Referral        `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}

type Referral struct {
	BaseEntity
	CustomerID uint     `gorm:"index;not null"`
	RefereeID  uint     `gorm:"index;not null"`
	Customer   Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Referee    Customer `gorm:"foreignKey:RefereeID;constraint:OnDelete:CASCADE"`
}

type VerificationCode struct {
	BaseEntity
	CustomerID uint      `gorm:"index;not null"`
	Customer   Customer  `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Code       string    `gorm:"size:6;not null"`
	ExpiresAt  time.Time `gorm:"not null"`
}

type CustomerAddress struct {
	BaseEntity
	CustomerID uint     `gorm:"index;not null"`
	Customer   Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Address    string   `gorm:"size:255;not null"`
	Longitude  string   `gorm:"size:20"`
	Latitude   string   `gorm:"size:20"`
}

type Bonus struct {
	BaseEntity
	Bonuses    float64  `gorm:"type:decimal(10,2);check:bonuses >= 0"`
	CustomerID uint     `gorm:"index;not null"`
	Customer   Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	ExpiresAt  *time.Time
}
