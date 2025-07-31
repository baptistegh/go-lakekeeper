// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package permission

import (
	"encoding/json"
	"errors"
)

type (
	RoleAssignmentType string

	// RoleAssignment represents the role assignments
	//
	// Assignee can be a role or a user
	// Assignement can be ownership or assignee
	RoleAssignment struct {
		Assignee   UserOrRole
		Assignment RoleAssignmentType
	}
)

const (
	OwnershipRoleAssignment RoleAssignmentType = "ownership"
	AssigneeRoleAssignment  RoleAssignmentType = "assignee"
)

// to be sure RoleAssignment can be JSON encoded/decoded
var (
	_ json.Unmarshaler = (*RoleAssignment)(nil)
	_ json.Marshaler   = (*RoleAssignment)(nil)

	ValidRoleAssignmentTypes = []RoleAssignmentType{
		OwnershipRoleAssignment,
		AssigneeRoleAssignment,
	}
)

func (sa *RoleAssignment) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Type RoleAssignmentType `json:"type"`
		Role *string            `json:"role,omitempty"`
		User *string            `json:"user,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	sa.Assignment = aux.Type

	if aux.Role == nil && aux.User == nil {
		return errors.New("error reading role assignment, role or user must be provided")
	}

	if aux.Role != nil && aux.User != nil {
		return errors.New("error reading role assignment, role and user can't be both provided")
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
	return errors.New("incorrect role assignment")
}

func (sa RoleAssignment) MarshalJSON() ([]byte, error) {
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
