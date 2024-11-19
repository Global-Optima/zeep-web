package types

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
)

func IsValidRole(role string) bool {
	switch Role(role) {
	case RoleAdmin, RoleManager, RoleEmployee:
		return true
	default:
		return false
	}
}
