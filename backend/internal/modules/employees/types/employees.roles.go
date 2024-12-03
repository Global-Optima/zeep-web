package types

type EmployeeRole string

const (
	RoleAdmin    EmployeeRole = "Admin"
	RoleDirector EmployeeRole = "Director"
	RoleManager  EmployeeRole = "Manager"
	RoleEmployee EmployeeRole = "Employee"
)

func IsEmployeeValidRole(role string) bool {
	switch EmployeeRole(role) {
	case RoleAdmin, RoleDirector, RoleManager, RoleEmployee:
		return true
	default:
		return false
	}
}
