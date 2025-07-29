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
	"context"
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	RolePermissionServiceInterface interface {
		// Get the access to a role
		// opt filters the access by a specific user or role.
		// If not specified, it returns the access for the current user.
		GetAccess(ctx context.Context, id string, opts *GetRoleAccessOptions, options ...core.RequestOptionFunc) (*GetRoleAccessResponse, *http.Response, error)
		// Get a role assignments
		// opt filters the assignments by relations.
		// If not specified, it returns all assignments.
		GetAssignments(ctx context.Context, id string, opts *GetRoleAssignmentsOptions, options ...core.RequestOptionFunc) (*GetRoleAssignmentsResponse, *http.Response, error)
		// Update permissions for a role
		Update(ctx context.Context, id string, opts *UpdateRolePermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// RolePermissionService handles communication with role permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	RolePermissionService struct {
		client core.Client
	}

	// Available actions on a role
	RoleAction string

	// GetRoleAccessOptions represents the GetAccess() options.
	//
	// Only one of PrincipalUser or PrincipalRole should be set at a time.
	// Setting both fields simultaneously is not allowed.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_access
	GetRoleAccessOptions struct {
		PrincipalUser *string `url:"principalUser,omitempty"`
		PrincipalRole *string `url:"principalRole,omitempty"`
	}

	// GetRoleAccessResponse represents the response from the GetAccess() endpoint.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_access
	GetRoleAccessResponse struct {
		AllowedActions []RoleAction `json:"allowed-actions"`
	}

	// GetRoleAssignmentsOptions represents the GetAssignments() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_assignments
	GetRoleAssignmentsOptions struct {
		Relations []RoleAssignmentType `url:"relations[],omitempty"`
	}

	// GetRoleAssignmentsResponse represents the response from the GetAssignments() endpoint.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_assignments
	GetRoleAssignmentsResponse struct {
		Assignments []*RoleAssignment `json:"assignments"`
	}

	// UpdateRolePermissionsOptions represents the Update() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_role_assignments
	UpdateRolePermissionsOptions struct {
		// The list of assignments to delete.
		Deletes []*RoleAssignment `json:"deletes,omitempty"`
		// The list of assignments to create.
		Writes []*RoleAssignment `json:"writes,omitempty"`
	}
)

const (
	Assume              RoleAction = "assume"
	CanGrantAssignee    RoleAction = "can_grant_assignee"
	CanChangeOwnership  RoleAction = "can_change_ownership"
	DeleteRole          RoleAction = "delete"
	UpdateRole          RoleAction = "update"
	ReadRole            RoleAction = "read"
	ReadRoleAssignments RoleAction = "read_assignments"
)

func NewRolePermissionService(client core.Client) RolePermissionServiceInterface {
	return &RolePermissionService{
		client: client,
	}
}

// GetAccess retrieves user or role access to a role.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_access
func (s *RolePermissionService) GetAccess(ctx context.Context, id string, opt *GetRoleAccessOptions, options ...core.RequestOptionFunc) (*GetRoleAccessResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/role/%s/access", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetRoleAccessResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// GetAccess gets user and role assignments of the role.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_assignments
func (s *RolePermissionService) GetAssignments(ctx context.Context, id string, opt *GetRoleAssignmentsOptions, options ...core.RequestOptionFunc) (*GetRoleAssignmentsResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/role/%s/assignments", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetRoleAssignmentsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// Update updates the role assignments.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_role_assignments
func (s *RolePermissionService) Update(ctx context.Context, id string, opt *UpdateRolePermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	path := fmt.Sprintf("/permissions/role/%s/assignments", id)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
