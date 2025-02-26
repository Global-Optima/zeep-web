package data

import (
	"time"
)

type Customer struct {
	BaseEntity
	FirstName  string            `gorm:"size:255;not null" sort:"firstName"`
	LastName   string            `gorm:"size:255;not null" sort:"lastName"`
	Password   string            `gorm:"size:255;not null"`
	Phone      string            `gorm:"size:16;not null"`
	IsVerified *bool             `gorm:"default:false" sort:"isVerified"`
	IsBanned   *bool             `gorm:"default:false" sort:"isBanned"`
	Addresses  []CustomerAddress `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Bonuses    []Bonus           `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Referrals  []Referral        `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
}

type Referral struct {
	BaseEntity
	CustomerID uint     `gorm:"index;not null"`
	RefereeID  uint     `gorm:"index;not null"`
	Customer   Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" sort:"customers"`
	Referee    Customer `gorm:"foreignKey:RefereeID;constraint:OnDelete:CASCADE" sort:"referees"`
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
	Customer   Customer `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" sort:"customers"`
	Address    string   `gorm:"size:255;not null"`
	Longitude  string   `gorm:"size:20"`
	Latitude   string   `gorm:"size:20"`
}

type Bonus struct {
	BaseEntity
	Bonuses    float64    `gorm:"type:decimal(10,2);check:bonuses >= 0"`
	CustomerID uint       `gorm:"index;not null"`
	Customer   Customer   `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" sort:"customers"`
	ExpiresAt  *time.Time `sort:"expiresAt"`
}
