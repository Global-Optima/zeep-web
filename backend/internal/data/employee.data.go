package data

import (
	"time"
)

type EmployeeType string

const (
	StoreEmployeeType     EmployeeType = "STORE"
	WarehouseEmployeeType EmployeeType = "WAREHOUSE"
)

type EmployeeRole string

const (
	RoleAdmin     EmployeeRole = "ADMIN"
	RoleDirector  EmployeeRole = "DIRECTOR"
	RoleManager   EmployeeRole = "MANAGER"
	RoleBarista   EmployeeRole = "BARISTA"
	RoleWarehouse EmployeeRole = "WAREHOUSE_EMPLOYEE"
)

func IsValidEmployeeRole(role EmployeeRole) bool {
	switch EmployeeRole(role) {
	case RoleAdmin, RoleDirector, RoleManager, RoleBarista, RoleWarehouse:
		return true
	default:
		return false
	}
}

type Employee struct {
	BaseEntity
	Name              string             `gorm:"size:255;not null"`
	Phone             string             `gorm:"size:15;unique"`
	Email             string             `gorm:"size:255;unique"`
	HashedPassword    string             `gorm:"size:255;not null"`
	Role              EmployeeRole       `gorm:"size:50;not null"`
	Type              EmployeeType       `gorm:"size:50;not null"`
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
