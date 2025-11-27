package permission

type UserOrRoleType string

const (
	UserType UserOrRoleType = "user"
	RoleType UserOrRoleType = "role"
)

type UserOrRole struct {
	Type  UserOrRoleType
	Value string
}
