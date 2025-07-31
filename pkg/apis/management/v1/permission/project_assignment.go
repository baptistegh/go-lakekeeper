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
	ProjectAssignmentType string

	// ProjectAssignment represents an assignment a role or a user can
	// have to a project
	//
	// Assignee can be a role or a user
	// Assignement can be Operator or Admin
	ProjectAssignment struct {
		Assignee   UserOrRole
		Assignment ProjectAssignmentType
	}
)

const (
	AdminProjectAssignment         ProjectAssignmentType = "project_admin"
	SecurityAdminProjectAssignment ProjectAssignmentType = "security_admin"
	DataAdminProjectAssignment     ProjectAssignmentType = "data_admin"
	RoleCreatorProjectAssignment   ProjectAssignmentType = "role_creator"
	DescribeProjectAssignment      ProjectAssignmentType = "describe"
	SelectProjectAssignment        ProjectAssignmentType = "select"
	CreateProjectAssignment        ProjectAssignmentType = "create"
	ModifyProjectAssignment        ProjectAssignmentType = "modify"
)

// ProjectAssignment can be JSON encoded/decoded
var (
	_ json.Unmarshaler = (*ProjectAssignment)(nil)
	_ json.Marshaler   = (*ProjectAssignment)(nil)

	ValidProjectAssignmentTypes = []ProjectAssignmentType{
		AdminProjectAssignment,
		SecurityAdminProjectAssignment,
		DataAdminProjectAssignment,
		RoleCreatorProjectAssignment,
		DescribeProjectAssignment,
		SelectProjectAssignment,
		CreateProjectAssignment,
		ModifyProjectAssignment,
	}
)

func (sa *ProjectAssignment) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Type ProjectAssignmentType `json:"type"`
		Role *string               `json:"role,omitempty"`
		User *string               `json:"user,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	sa.Assignment = aux.Type

	if aux.Role == nil && aux.User == nil {
		return errors.New("error reading project assignment, role or user must be provided")
	}

	if aux.Role != nil && aux.User != nil {
		return errors.New("error reading project assignment, role and user can't be both provided")
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
	return errors.New("incorrect project assignment")
}

func (sa ProjectAssignment) MarshalJSON() ([]byte, error) {
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
