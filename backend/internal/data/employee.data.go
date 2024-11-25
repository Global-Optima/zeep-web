package data

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type Employee struct {
	BaseEntity
	Name           string             `gorm:"size:255;not null"`
	Phone          string             `gorm:"size:15;unique"`
	Email          string             `gorm:"size:255;unique"`
	Role           types.EmployeeRole `gorm:"size:50;not null"`
	StoreID        uint               `gorm:"index"`
	Store          Store              `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE"`
	IsActive       bool               `gorm:"default:true"`
	HashedPassword string             `gorm:"size:255;not null"`
	Audits         []EmployeeAudit    `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
	Workdays       []EmployeeWorkday  `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

type EmployeeAudit struct {
	BaseEntity
	StartWorkAt *time.Time `gorm:"type:timestamp"`
	EndWorkAt   *time.Time `gorm:"type:timestamp"`
	EmployeeID  uint       `gorm:"index;not null"`
	Employee    Employee   `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

type EmployeeWorkday struct {
	BaseEntity
	Day        string   `gorm:"size:15;not null"`
	StartAt    string   `gorm:"type:time;not null"`
	EndAt      string   `gorm:"type:time;not null"`
	EmployeeID uint     `gorm:"index;not null"`
	Employee   Employee `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}
