package permission

import (
	"encoding/json"
	"errors"
)

type (
	ServerAssignmentType string

	// ServerAssignment represents an assignment a role or a user can
	// have to the server
	//
	// Assignee can be a role or a user
	// Assignement can be Operator or Admin
	ServerAssignment struct {
		Assignee   UserOrRole
		Assignment ServerAssignmentType
	}
)

const (
	OperatorServerAssignment ServerAssignmentType = "operator"
	AdminServerAssignment    ServerAssignmentType = "admin"
)

// to be sure ServerAssignment can be JSON encoded/decoded
var (
	_ json.Unmarshaler = (*ServerAssignment)(nil)
	_ json.Marshaler   = (*ServerAssignment)(nil)

	_ Assignment = (*ServerAssignment)(nil)
)

func (sa *ServerAssignment) GetAssignment() string {
	return string(sa.Assignment)
}

func (sa *ServerAssignment) GetPrincipalID() string {
	return sa.Assignee.Value
}

func (sa *ServerAssignment) GetPrincipalType() UserOrRoleType {
	return sa.Assignee.Type
}

func (sa *ServerAssignment) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Type ServerAssignmentType `json:"type"`
		Role *string              `json:"role,omitempty"`
		User *string              `json:"user,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	sa.Assignment = aux.Type

	if aux.Role == nil && aux.User == nil {
		return errors.New("error reading server assignment, role or user must be provided")
	}

	if aux.Role != nil && aux.User != nil {
		return errors.New("error reading server assignment, role and user can't be both provided")
	}

	if aux.Role != nil {
		sa.Assignee = UserOrRole{
			RoleType,
			*aux.Role,
		}
		return nil
	}

	if aux.User != nil {
		sa.Assignee = UserOrRole{
			UserType,
			*aux.User,
		}
		return nil
	}
	return errors.New("incorrect server assignment")
}

func (sa ServerAssignment) MarshalJSON() ([]byte, error) {
	aux := make(map[string]string)

	switch sa.Assignee.Type {
	case RoleType:
		aux["role"] = sa.Assignee.Value
	case UserType:
		aux["user"] = sa.Assignee.Value
	}

	aux["type"] = string(sa.Assignment)

	return json.Marshal(aux)
}
