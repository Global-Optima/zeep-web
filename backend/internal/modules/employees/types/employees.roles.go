package types

type Role string

const (
	RoleAdmin    Role = "Admin"
	RoleDirector Role = "Director"
	RoleManager  Role = "Manager"
	RoleEmployee Role = "Employee"
)

func IsValidRole(role string) bool {
	switch Role(role) {
	case RoleAdmin, RoleDirector, RoleManager, RoleEmployee:
		return true
	default:
		return false
	}
}
