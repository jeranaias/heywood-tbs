package auth

// Role constants used throughout the application.
const (
	RoleXO           = "xo"
	RoleStaff        = "staff"
	RoleSPC          = "spc"
	RoleStudent      = "student"
	RoleUnauthorized = "unauthorized"
)

// ValidRoles is the set of assignable roles.
var ValidRoles = map[string]bool{
	RoleXO:      true,
	RoleStaff:   true,
	RoleSPC:     true,
	RoleStudent: true,
}

// IsPrivileged returns true for roles with admin-level data access.
func IsPrivileged(role string) bool {
	return role == RoleXO || role == RoleStaff
}
