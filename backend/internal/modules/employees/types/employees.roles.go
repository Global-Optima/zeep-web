package types

type EmployeeRole string

const (
	RoleAdmin    EmployeeRole = "ADMIN"
	RoleDirector EmployeeRole = "DIRECTOR"
	RoleManager  EmployeeRole = "MANAGER"
	RoleEmployee EmployeeRole = "EMPLOYEE"
)

func IsEmployeeValidRole(role string) bool {
	switch EmployeeRole(role) {
	case RoleAdmin, RoleDirector, RoleManager, RoleEmployee:
		return true
	default:
		return false
	}
}
