package data

import (
	"github.com/pkg/errors"
	"strings"
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
	case RoleAdmin, RoleDirector, RoleManager, RoleBarista:
		return true
	default:
		return false
	}
}

type Weekday string

const (
	Monday    Weekday = "MONDAY"
	Tuesday   Weekday = "TUESDAY"
	Wednesday Weekday = "WEDNESDAY"
	Thursday  Weekday = "THURSDAY"
	Friday    Weekday = "FRIDAY"
	Saturday  Weekday = "SATURDAY"
	Sunday    Weekday = "SUNDAY"
)

var weekdayMapping = map[string]Weekday{
	"MONDAY":    Monday,
	"TUESDAY":   Tuesday,
	"WEDNESDAY": Wednesday,
	"THURSDAY":  Thursday,
	"FRIDAY":    Friday,
	"SATURDAY":  Saturday,
	"SUNDAY":    Sunday,
}

func ToWeekday(s string) (Weekday, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if weekday, exists := weekdayMapping[s]; exists {
		return weekday, nil
	}
	return "", errors.New("invalid weekday")
}

func (w Weekday) ToString() string {
	str := string(w)
	if len(str) == 0 {
		return ""
	}
	return str
}

func IsValidWeekday(weekday Weekday) bool {
	switch weekday {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday:
		return true
	}
	return false
}

type Employee struct {
	BaseEntity
	FirstName         string             `gorm:"size:255;not null" sort:"firstName"`
	LastName          string             `gorm:"size:255;not null" sort:"lastName"`
	Phone             string             `gorm:"size:16;unique"`
	Email             string             `gorm:"size:255;unique" sort:"email"`
	HashedPassword    string             `gorm:"size:255;not null"`
	Role              EmployeeRole       `gorm:"size:50;not null" sort:"role"`
	Type              EmployeeType       `gorm:"size:50;not null" sort:"type"`
	IsActive          bool               `gorm:"default:true" sort:"isActive"`
	StoreEmployee     *StoreEmployee     `gorm:"foreignKey:EmployeeID"`
	WarehouseEmployee *WarehouseEmployee `gorm:"foreignKey:EmployeeID"`
}

type StoreEmployee struct {
	BaseEntity
	EmployeeID  uint `gorm:"not null;uniqueIndex"`
	StoreID     uint `gorm:"not null"`
	IsFranchise bool `gorm:"default:false" sort:"isFranchise"`
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
	Employee    Employee   `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employees"`
}

type EmployeeWorkday struct {
	BaseEntity
	Day        Weekday  `gorm:"size:15;not null" sort:"day"`
	StartAt    string   `gorm:"type:time;not null" sort:"startAt"`
	EndAt      string   `gorm:"type:time;not null" sort:"endAt"`
	EmployeeID uint     `gorm:"index;not null"`
	Employee   Employee `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employees"`
}
