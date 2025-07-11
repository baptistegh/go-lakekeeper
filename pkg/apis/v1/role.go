package v1

import (
	"errors"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	RoleServiceInterface interface {
		List(opts *ListRolesOptions, options ...core.RequestOptionFunc) (*ListRolesResponse, *http.Response, error)
		Get(id string, options ...core.RequestOptionFunc) (*Role, *http.Response, error)
		Create(opts *CreateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error)
		Update(id string, opts *UpdateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error)
		Delete(id string, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// RoleService handles communication with role endpoints of the Lakekeeper API.
	//
	//
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role
	RoleService struct {
		projectID string
		client    core.Client
	}
)

var _ RoleServiceInterface = (*RoleService)(nil)

func NewRoleService(client core.Client, projectID string) RoleServiceInterface {
	return &RoleService{
		projectID: projectID,
		client:    client,
	}
}

// Project represents a lakekeeper role
type Role struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"project-id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	CreatedAt string  `json:"created-at"`
	UpdatedAt *string `json:"updated-at,omitempty"`
}

// Get retrieves information about a role.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/get_role
func (s *RoleService) Get(id string, options ...core.RequestOptionFunc) (*Role, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodGet, "/role/"+id, nil, options)
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

// ListRolesOptions represents List() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
type ListRolesOptions struct {
	Name      *string `url:"name,omitempty"`
	PageToken *string `url:"pageToken,omitempty"`
	PageSize  *string `url:"pageSize,omitempty"`
	// Deprecated: This field will be removed in a future version.
	// ProjectID should be obtained from the Service itself and is not intended to be used here.
	// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
	ProjectID *string `url:"projectId,omitempty"`
}

// ListRolesResponse represents a response from list_roles API endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/list_roles
type ListRolesResponse struct {
	NextPageToken *string `json:"next-page-token,omitempty"`
	Roles         []*Role `json:"roles"`
}

// List returns all roles in the project that the current user has access to view.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/list_roles
func (s *RoleService) List(opts *ListRolesOptions, options ...core.RequestOptionFunc) (*ListRolesResponse, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opts == nil {
		opts = &ListRolesOptions{}
	}
	opts.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodGet, "/role", opts, options)
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

// CreateRoleOptions represents Create() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/create_role
type CreateRoleOptions struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	// Deprecated: This field will be removed in a future version.
	// ProjectID should be obtained from the Service itself and is not intended to be used here.
	// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
	ProjectID *string `json:"project-id,omitempty"`
}

// Create creates a role with the specified name and description.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/create_role
func (s *RoleService) Create(opts *CreateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opts == nil {
		opts = &CreateRoleOptions{}
	}
	opts.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, "/role", opts, options)
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

// UpdateRoleOptions represents Update() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/update_role
type UpdateRoleOptions struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// Update update a role with the specified name and description.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/update_role
func (s *RoleService) Update(id string, opts *UpdateRoleOptions, options ...core.RequestOptionFunc) (*Role, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("Role ID must be defined to be updated")
	}

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, "/role/"+id, opts, options)
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
func (s *RoleService) Delete(id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodDelete, "/role/"+id, nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
