package permission

import (
	"encoding/json"
	"fmt"
)

// WarehouseAssignment represents an assignment a role or a user can
// have to a warehouse
//
// Assignee can be a role or a user
// Assignement can be Operator or Admin
type WarehouseAssignment struct {
	Assignee   UserOrRole
	Assignment WarehouseAssignmentType
}

// to be sure WarehouseAssignment can be JSON encoded/decoded
var (
	_ json.Unmarshaler = (*WarehouseAssignment)(nil)
	_ json.Marshaler   = (*WarehouseAssignment)(nil)
)

type WarehouseAssignmentType string

const (
	OwnershipWarehouseAssignment         WarehouseAssignmentType = "ownership"
	PassGrantsAdminWarehouseAssignment   WarehouseAssignmentType = "pass_grants"
	ManageGrantsAdminWarehouseAssignment WarehouseAssignmentType = "manage_grants"
	DescribeWarehouseAssignment          WarehouseAssignmentType = "describe"
	SelectWarehouseAssignment            WarehouseAssignmentType = "select"
	CreateWarehouseAssignment            WarehouseAssignmentType = "create"
	ModifyWarehouseAssignment            WarehouseAssignmentType = "modify"
)

func (sa *WarehouseAssignment) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Type WarehouseAssignmentType `json:"type"`
		Role *string                 `json:"role,omitempty"`
		User *string                 `json:"user,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	sa.Assignment = aux.Type

	if aux.Role == nil && aux.User == nil {
		return fmt.Errorf("error reading warehouse assignment, role or user must be provided")
	}

	if aux.Role != nil && aux.User != nil {
		return fmt.Errorf("error reading warehouse assignment, role and user can't be both provided")
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
	return fmt.Errorf("incorrect warehouse assignment")
}

func (sa WarehouseAssignment) MarshalJSON() ([]byte, error) {
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
