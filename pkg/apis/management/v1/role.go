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

package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	RoleServiceInterface interface {
		// Returns all roles in the project that the current user has access to view.
		List(ctx context.Context, opts *ListRolesOptions, options ...core.RequestOptionFunc) (*ListRolesResponse, *http.Response, error)
		// Retrieves detailed information about a specific role.
		Get(ctx context.Context, id string, options ...core.RequestOptionFunc) (*Role, *http.Response, error)
		// Creates a role with the specified name, description, and permissions.
		Create(ctx context.Context, opts *CreateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error)
		// Updates a role
		Update(ctx context.Context, id string, opts *UpdateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error)
		// Permanently removes a role and all its associated permissions.
		Delete(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error)
		// Performs a fuzzy search for roles based on the provided criteria.
		Search(ctx context.Context, opts *SearchRoleOptions, options ...core.RequestOptionFunc) (*SearchRoleResponse, *http.Response, error)
	}

	// RoleService handles communication with role endpoints of the Lakekeeper API.
	//
	//
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role
	RoleService struct {
		projectID string
		client    core.Client
	}

	// Project represents a lakekeeper role
	Role struct {
		ID          string  `json:"id"`
		ProjectID   string  `json:"project-id"`
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`

		CreatedAt string  `json:"created-at"`
		UpdatedAt *string `json:"updated-at,omitempty"`
	}

	// ListRolesOptions represents List() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
	ListRolesOptions struct {
		Name *string `url:"name,omitempty"`
		// Deprecated: This field will be removed in a future version.
		// ProjectID should be obtained from the Service itself and is not intended to be used here.
		// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
		ProjectID *string `url:"projectId,omitempty"`

		ListOptions `url:",inline"` // Embed ListOptions for pagination support
	}

	// ListRolesResponse represents a response from list_roles API endpoint.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/list_roles
	ListRolesResponse struct {
		Roles []*Role `json:"roles"`

		ListResponse `json:",inline"` // Embed ListResponse for pagination support
	}

	// CreateRoleOptions represents Create() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/create_role
	CreateRoleOptions struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		// Deprecated: This field will be removed in a future version.
		// ProjectID should be obtained from the Service itself and is not intended to be used here.
		// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
		ProjectID *string `json:"project-id,omitempty"`
	}

	// UpdateRoleOptions represents Update() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/update_role
	UpdateRoleOptions struct {
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
	}

	// SearchRoleOptions reprensents Search() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/search_role
	SearchRoleOptions struct {
		// Deprecated: This field will be removed in a future version.
		// ProjectID should be obtained from the Service itself and is not intended to be used here.
		// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
		ProjectID *string `json:"project-id,omitempty"`
		// Search string for fuzzy search. Length is truncated to 64 characters.
		Search string `json:"search"`
	}

	// SearchRoleResponse reprensents a Search() response.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/search_role
	SearchRoleResponse struct {
		Roles []*Role `json:"roles"`
	}
)

func NewRoleService(client core.Client, projectID string) RoleServiceInterface {
	return &RoleService{
		projectID: projectID,
		client:    client,
	}
}

// Get retrieves information about a role.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/get_role
func (s *RoleService) Get(ctx context.Context, id string, options ...core.RequestOptionFunc) (*Role, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, "/role/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var role Role

	resp, apiErr := s.client.Do(req, &role)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &role, resp, nil
}

// Search performs a fuzzy search for roles based on the provided criteria.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/search_role
func (s *RoleService) Search(ctx context.Context, opts *SearchRoleOptions, options ...core.RequestOptionFunc) (*SearchRoleResponse, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opts == nil {
		opts = &SearchRoleOptions{}
	}
	opts.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/search/role", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var roles SearchRoleResponse
	resp, apiErr := s.client.Do(req, &roles)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &roles, resp, nil
}

// List returns all roles in the project that the current user has access to view.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/list_roles
func (s *RoleService) List(ctx context.Context, opts *ListRolesOptions, options ...core.RequestOptionFunc) (*ListRolesResponse, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opts == nil {
		opts = &ListRolesOptions{}
	}
	opts.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, "/role", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var r ListRolesResponse
	resp, apiErr := s.client.Do(req, &r)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &r, resp, nil
}

// Create creates a role with the specified name and description.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/create_role
func (s *RoleService) Create(ctx context.Context, opts *CreateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opts == nil {
		opts = &CreateRoleOptions{}
	}
	opts.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/role", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var role Role

	resp, apiErr := s.client.Do(req, &role)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &role, resp, nil
}

// Update update a role with the specified name and description.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/update_role
func (s *RoleService) Update(ctx context.Context, id string, opts *UpdateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("Role ID must be defined to be updated")
	}

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/role/"+id, opts, options)
	if err != nil {
		return nil, nil, err
	}

	var role Role

	resp, apiErr := s.client.Do(req, &role)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &role, resp, nil
}

// Delete permanently removes a role and all its associated permissions.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/delete_role
func (s *RoleService) Delete(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodDelete, "/role/"+id, nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
