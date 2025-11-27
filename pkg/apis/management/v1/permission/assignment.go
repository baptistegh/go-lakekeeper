package permission

type Assignment interface {
	GetPrincipalType() UserOrRoleType
	GetPrincipalID() string
	GetAssignment() string
}
