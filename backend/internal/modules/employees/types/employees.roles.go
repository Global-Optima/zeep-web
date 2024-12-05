package types

type EmployeeRole string

const (
	RoleAdmin    EmployeeRole = "Admin"
	RoleDirector EmployeeRole = "Director"
	RoleManager  EmployeeRole = "Manager"
	RoleBarista  EmployeeRole = "Barista"
)

func IsEmployeeValidRole(role string) bool {
	switch EmployeeRole(role) {
	case RoleAdmin, RoleDirector, RoleManager, RoleBarista:
		return true
	default:
		return false
	}
}
