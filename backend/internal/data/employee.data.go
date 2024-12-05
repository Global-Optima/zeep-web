package data

import (
	"time"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees/types"
)

type Employee struct {
	BaseEntity
	Name              string             `gorm:"size:255;not null"`
	Phone             string             `gorm:"size:15;unique"`
	Email             string             `gorm:"size:255;unique"`
	HashedPassword    string             `gorm:"size:255;not null"`
	Role              types.EmployeeRole `gorm:"size:50;not null"` // Admin, Manager, etc.
	Type              types.EmployeeType `gorm:"size:50;not null"`
	IsActive          bool               `gorm:"default:true"`
	StoreEmployee     *StoreEmployee     `gorm:"foreignKey:EmployeeID"`
	WarehouseEmployee *WarehouseEmployee `gorm:"foreignKey:EmployeeID"`
}

type StoreEmployee struct {
	BaseEntity
	EmployeeID  uint `gorm:"not null;uniqueIndex"`
	StoreID     uint `gorm:"not null"`
	IsFranchise bool `gorm:"default:false"`
}

type WarehouseEmployee struct {
	BaseEntity
	EmployeeID  uint `gorm:"not null;uniqueIndex"`
	WarehouseID uint `gorm:"not null"`
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
