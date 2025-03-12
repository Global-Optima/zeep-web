package data

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type EmployeeType string

func (t EmployeeType) ToString() string {
	return string(t)
}

const (
	StoreEmployeeType      EmployeeType = "STORE"
	WarehouseEmployeeType  EmployeeType = "WAREHOUSE"
	FranchiseeEmployeeType EmployeeType = "FRANCHISEE"
	RegionEmployeeType     EmployeeType = "REGION"
	AdminEmployeeType      EmployeeType = "ADMIN"
)

type EmployeeRole string

func (e EmployeeRole) ToString() string {
	return string(e)
}

func ToEmployeeRole(role string) (EmployeeRole, error) {
	if IsValidEmployeeRole(EmployeeRole(role)) {
		return EmployeeRole(role), nil
	}
	return "", fmt.Errorf("invalid employee role: %s", role)
}

type AdminRole = EmployeeRole

const (
	RoleAdmin AdminRole = "ADMIN"
	RoleOwner AdminRole = "OWNER"
)

type StoreEmployeeRole = EmployeeRole

const (
	RoleStoreManager StoreEmployeeRole = "STORE_MANAGER"
	RoleBarista      StoreEmployeeRole = "BARISTA"
)

type WarehouseEmployeeRole = EmployeeRole

const (
	RoleWarehouseManager  WarehouseEmployeeRole = "WAREHOUSE_MANAGER"
	RoleWarehouseEmployee WarehouseEmployeeRole = "WAREHOUSE_EMPLOYEE"
)

type FranchiseeEmployeeRole = EmployeeRole

const (
	RoleFranchiseManager FranchiseeEmployeeRole = "FRANCHISEE_MANAGER"
	RoleFranchiseOwner   FranchiseeEmployeeRole = "FRANCHISEE_OWNER"
)

type RegionManagerRole = EmployeeRole

const (
	RoleRegionWarehouseManager RegionManagerRole = "REGION_WAREHOUSE_MANAGER"
)

var EmployeeTypeRoleMap = map[EmployeeType][]EmployeeRole{
	StoreEmployeeType:      {RoleStoreManager, RoleBarista},
	WarehouseEmployeeType:  {RoleWarehouseManager, RoleWarehouseEmployee},
	FranchiseeEmployeeType: {RoleFranchiseManager, RoleFranchiseOwner},
	RegionEmployeeType:     {RoleRegionWarehouseManager},
	AdminEmployeeType:      {RoleAdmin, RoleOwner},
}

var RoleManagePermissions = map[EmployeeRole][]EmployeeRole{
	RoleStoreManager:           {RoleBarista},
	RoleWarehouseManager:       {RoleWarehouseEmployee},
	RoleFranchiseManager:       {RoleStoreManager, RoleBarista},
	RoleRegionWarehouseManager: {RoleWarehouseManager, RoleWarehouseEmployee},
	RoleAdmin: {
		RoleOwner, RoleFranchiseOwner, RoleFranchiseManager,
		RoleRegionWarehouseManager, RoleStoreManager, RoleWarehouseManager,
		RoleBarista, RoleWarehouseEmployee,
	},
}

func CanManageRole(currentRole, targetRole EmployeeRole) bool {
	allowedRoles, exists := RoleManagePermissions[currentRole]
	if !exists {
		return false
	}

	for _, allowedRole := range allowedRoles {
		if allowedRole == targetRole {
			return true
		}
	}
	return false
}

var (
	AdminPermissions = []EmployeeRole{
		RoleAdmin,
		RoleOwner,
	}
	FranchiseePermissions = []EmployeeRole{
		RoleFranchiseManager,
		RoleFranchiseOwner,
	}
	FranchiseeReadPermissions = append(
		FranchiseePermissions,
		RoleOwner,
	)
	RegionPermissions = []EmployeeRole{
		RoleRegionWarehouseManager,
	}
	RegionReadPermissions = append(
		RegionPermissions,
		RoleOwner,
	)
	WarehouseManagementPermissions = []EmployeeRole{
		RoleRegionWarehouseManager,
		RoleWarehouseManager,
	}
	WarehouseReadPermissions = append(
		WarehouseManagementPermissions,
		RoleOwner,
		RoleWarehouseEmployee,
	)
	WarehousePermissions = []EmployeeRole{
		RoleWarehouseManager,
		RoleWarehouseEmployee,
	}
	StoreManagementPermissions = append(
		FranchiseePermissions,
		RoleStoreManager,
	)
	StorePermissions = []EmployeeRole{
		RoleStoreManager,
		RoleBarista,
	}
	StoreReadPermissions = append(
		StoreManagementPermissions,
		RoleBarista,
	)
	StoreAndWarehousePermissions = append(
		StorePermissions,
		WarehousePermissions...,
	)
	FranchiseeAndRegionPermissions = append(
		FranchiseePermissions,
		RegionPermissions...,
	)
)

func GetEmployeeTypeByRole(role EmployeeRole) EmployeeType {
	switch {
	case slices.Contains(FranchiseePermissions, role):
		return FranchiseeEmployeeType
	case slices.Contains(RegionPermissions, role):
		return RegionEmployeeType
	case slices.Contains(WarehousePermissions, role):
		return WarehouseEmployeeType
	case slices.Contains(StorePermissions, role):
		return StoreEmployeeType
	case slices.Contains(AdminPermissions, role):
		return AdminEmployeeType
	}
	return ""
}

func IsAllowableRole(employeeType EmployeeType, role EmployeeRole) bool {
	roles, exists := EmployeeTypeRoleMap[employeeType]
	if !exists {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func IsValidEmployeeRole(role EmployeeRole) bool {
	switch role {
	case RoleAdmin, RoleOwner, RoleRegionWarehouseManager, RoleFranchiseManager, RoleFranchiseOwner,
		RoleWarehouseManager, RoleStoreManager, RoleBarista, RoleWarehouseEmployee:
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
	// English
	"MONDAY":    Monday,
	"TUESDAY":   Tuesday,
	"WEDNESDAY": Wednesday,
	"THURSDAY":  Thursday,
	"FRIDAY":    Friday,
	"SATURDAY":  Saturday,
	"SUNDAY":    Sunday,

	// Russian
	"ПОНЕДЕЛЬНИК": Monday,
	"ВТОРНИК":     Tuesday,
	"СРЕДА":       Wednesday,
	"ЧЕТВЕРГ":     Thursday,
	"ПЯТНИЦА":     Friday,
	"СУББОТА":     Saturday,
	"ВОСКРЕСЕНЬЕ": Sunday,

	// Kazakh
	"ДҮЙСЕНБІ": Monday,
	"СЕЙСЕНБІ": Tuesday,
	"СӘРСЕНБІ": Wednesday,
	"БЕЙСЕНБІ": Thursday,
	"ЖҰМА":     Friday,
	"СЕНБІ":    Saturday,
	"ЖЕКСЕНБІ": Sunday,
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
	_, exists := weekdayMapping[string(weekday)]
	return exists
}

type Employee struct {
	BaseEntity
	FirstName          string              `gorm:"size:255;not null" sort:"firstName"`
	LastName           string              `gorm:"size:255;not null" sort:"lastName"`
	Phone              string              `gorm:"size:16;not null"`
	Email              string              `gorm:"size:255;not null" sort:"email"`
	HashedPassword     string              `gorm:"size:255;not null"`
	IsActive           bool                `gorm:"not null" sort:"isActive"`
	StoreEmployee      *StoreEmployee      `gorm:"foreignKey:EmployeeID"`
	WarehouseEmployee  *WarehouseEmployee  `gorm:"foreignKey:EmployeeID"`
	RegionEmployee     *RegionEmployee     `gorm:"foreignKey:EmployeeID"`
	FranchiseeEmployee *FranchiseeEmployee `gorm:"foreignKey:EmployeeID"`
	AdminEmployee      *AdminEmployee      `gorm:"foreignKey:EmployeeID"`
	Workdays           []EmployeeWorkday   `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

func (e *Employee) GetType() EmployeeType {
	switch {
	case e.StoreEmployee != nil:
		return StoreEmployeeType
	case e.FranchiseeEmployee != nil:
		return FranchiseeEmployeeType
	case e.AdminEmployee != nil:
		return AdminEmployeeType
	case e.RegionEmployee != nil:
		return RegionEmployeeType
	case e.WarehouseEmployee != nil:
		return WarehouseEmployeeType
	default:
		return ""
	}
}

type StoreEmployee struct {
	BaseEntity
	EmployeeID uint              `gorm:"index,not null"`
	StoreID    uint              `gorm:"index,not null"`
	Role       StoreEmployeeRole `gorm:"type:store_employee_role;not null" sort:"role"`
	Employee   Employee          `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employee"`
	Store      Store             `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE" sort:"store"`
}

type WarehouseEmployee struct {
	BaseEntity
	EmployeeID  uint                  `gorm:"index,not null"`
	WarehouseID uint                  `gorm:"index,not null"`
	Role        WarehouseEmployeeRole `gorm:"type:warehouse_employee_role;not null" sort:"role"`
	Employee    Employee              `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employee"`
	Warehouse   Warehouse             `gorm:"foreignKey:WarehouseID;constraint:OnDelete:CASCADE" sort:"warehouse"`
}

type FranchiseeEmployee struct {
	BaseEntity
	FranchiseeID uint                   `gorm:"index,not null"`
	EmployeeID   uint                   `gorm:"index,not null"`
	Role         FranchiseeEmployeeRole `gorm:"type:franchisee_employee_role;not null" sort:"role"`
	Employee     Employee               `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employee"`
	Franchisee   Franchisee             `gorm:"foreignKey:FranchiseeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" sort:"franchisee"`
}

type RegionEmployee struct {
	BaseEntity
	EmployeeID uint              `gorm:"index;not null"`
	RegionID   uint              `gorm:"index;not null"`
	Role       RegionManagerRole `gorm:"type:warehouse_region_manager_role;not null" sort:"role"`
	Employee   Employee          `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employee"`
	Region     Region            `gorm:"foreignKey:RegionID;constraint:OnDelete:CASCADE" sort:"region"`
}

type AdminEmployee struct {
	BaseEntity
	EmployeeID uint      `gorm:"index;not null"`
	Role       AdminRole `gorm:"type:admin_role;not null" sort:"role"`
	Employee   Employee  `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" sort:"employee"`
}

type EmployeeWorkTrack struct {
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
