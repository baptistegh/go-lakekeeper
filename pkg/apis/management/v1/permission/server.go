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
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ServerPermissionServiceInterface interface {
		// Get the access to the server
		// opt filters the access by a specific user or role.
		// If not specified, it returns the access for the current user.
		GetAccess(ctx context.Context, opts *GetServerAccessOptions, options ...core.RequestOptionFunc) (*GetServerAccessResponse, *http.Response, error)
		// Get user and role assignments of the server
		// opt filters the assignments by relations.
		// If not specified, it returns all assignments.
		GetAssignments(ctx context.Context, opts *GetServerAssignmentsOptions, options ...core.RequestOptionFunc) (*GetServerAssignmentsResponse, *http.Response, error)
		// Update permissions for the server
		Update(ctx context.Context, opts *UpdateServerPermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// ServerPermissionService handles communication with server permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	ServerPermissionService struct {
		client core.Client
	}

	// Available actions on a server
	ServerAction string

	// GetServerAccessOptions represents the GetAccess() options.
	//
	// Only one of PrincipalUser or PrincipalRole should be set at a time.
	// Setting both fields simultaneously is not allowed.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_access
	GetServerAccessOptions struct {
		PrincipalUser *string `url:"principalUser,omitempty"`
		PrincipalRole *string `url:"principalRole,omitempty"`
	}

	// GetServerAccessResponse represents the response from the GetAccess() endpoint.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_access
	GetServerAccessResponse struct {
		AllowedActions []ServerAction `json:"allowed-actions"`
	}

	// GetServerAssignmentsOptions represents the GetAssignments() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_assignments
	GetServerAssignmentsOptions struct {
		Relations []ServerAssignmentType `url:"relations[],omitempty"`
	}

	// GetServerAssignmentsResponse represents the response from the GetAssignments() endpoint.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_assignments
	GetServerAssignmentsResponse struct {
		Assignments []*ServerAssignment `json:"assignments"`
	}

	// UpdateServerPermissionsOptions represents the Update() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_server_assignments
	UpdateServerPermissionsOptions struct {
		// The list of assignments to delete.
		Deletes []*ServerAssignment `json:"deletes,omitempty"`
		// The list of assignments to create.
		Writes []*ServerAssignment `json:"writes,omitempty"`
	}
)

const (
	CreateProject    ServerAction = "create_project"
	UpdateUsers      ServerAction = "update_users"
	DeleteUsers      ServerAction = "delete_users"
	ListUsers        ServerAction = "list_users"
	GrantServerAdmin ServerAction = "grant_admin"
	ProvisionUsers   ServerAction = "provision_users"
	ReadAssignments  ServerAction = "read_assignments"
)

func NewServerPermissionService(client core.Client) ServerPermissionServiceInterface {
	return &ServerPermissionService{
		client: client,
	}
}

// GetAccess retrieves user or role access to the server.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_access
func (s *ServerPermissionService) GetAccess(ctx context.Context, opt *GetServerAccessOptions, options ...core.RequestOptionFunc) (*GetServerAccessResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/permissions/server/access", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetServerAccessResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// GetAccess gets user and role assignments of the server.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_assignments
func (s *ServerPermissionService) GetAssignments(ctx context.Context, opt *GetServerAssignmentsOptions, options ...core.RequestOptionFunc) (*GetServerAssignmentsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/permissions/server/assignments", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetServerAssignmentsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// Update updates the server assignments.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_server_assignments
func (s *ServerPermissionService) Update(ctx context.Context, opt *UpdateServerPermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/permissions/server/assignments", opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
