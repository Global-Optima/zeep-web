package data

import (
	"time"
)

type Employee struct {
	BaseEntity
	Name     string            `gorm:"size:255;not null"`
	Phone    string            `gorm:"size:15;unique"`
	Email    string            `gorm:"size:255;unique"`
	RoleID   *uint             `gorm:"index"`
	Role     *EmployeeRole     `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	StoreID  *uint             `gorm:"index"`
	Store    *Store            `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsActive bool              `gorm:"default:true"`
	Audits   []EmployeeAudit   `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
	Workdays []EmployeeWorkday `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
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

type EmployeeRole struct {
	BaseEntity
	Name      string     `gorm:"size:50;not null;unique"`
	Employees []Employee `gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL"`
}
