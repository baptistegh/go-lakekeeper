package permission

import (
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ServerPermissionServiceInterface interface {
		// Get the access to the server
		// opt filters the access by a specific user or role.
		// If not specified, it returns the access for the current user.
		GetAccess(opts *GetServerAccessOptions, options ...core.RequestOptionFunc) (*GetServerAccessResponse, *http.Response, error)
		// Get user and role assignments of the server
		// opt filters the assignments by relations.
		// If not specified, it returns all assignments.
		GetAssignments(opts *GetServerAssignmentsOptions, options ...core.RequestOptionFunc) (*GetServerAssignmentsResponse, *http.Response, error)
		// Update permissions for the server
		Update(opts *UpdateServerPermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// ServerPermissionService handles communication with server permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	ServerPermissionService struct {
		client core.Client
	}
)

type ServerAction string

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

// GetServerAccessOptions represents the GetAccess() options.
//
// Only one of PrincipalUser or PrincipalRole should be set at a time.
// Setting both fields simultaneously is not allowed.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_access
type GetServerAccessOptions struct {
	PrincipalUser *string `url:"principalUser,omitempty"`
	PrincipalRole *string `url:"principalRole,omitempty"`
}

// GetServerAccessResponse represents the response from the GetAccess() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_access
type GetServerAccessResponse struct {
	AllowedActions []ServerAction `json:"allowed-actions"`
}

// GetAccess retrieves user or role access to the server.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_access
func (s *ServerPermissionService) GetAccess(opt *GetServerAccessOptions, options ...core.RequestOptionFunc) (*GetServerAccessResponse, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/permissions/server/access", opt, options)
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

// GetServerAssignmentsOptions represents the GetAssignments() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_assignments
type GetServerAssignmentsOptions struct {
	Relations []ServerAssignmentType `url:"relations[],omitempty"`
}

// GetServerAssignmentsResponse represents the response from the GetAssignments() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_assignments
type GetServerAssignmentsResponse struct {
	Assignments []*ServerAssignment `json:"assignments"`
}

// GetAccess gets user and role assignments of the server.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_server_assignments
func (s *ServerPermissionService) GetAssignments(opt *GetServerAssignmentsOptions, options ...core.RequestOptionFunc) (*GetServerAssignmentsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/permissions/server/assignments", opt, options)
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

// UpdateServerPermissionsOptions represents the Update() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_server_assignments
type UpdateServerPermissionsOptions struct {
	// The list of assignments to delete.
	Deletes []*ServerAssignment `json:"deletes,omitempty"`
	// The list of assignments to create.
	Writes []*ServerAssignment `json:"writes,omitempty"`
}

// Update updates the server assignments.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_server_assignments
func (s *ServerPermissionService) Update(opt *UpdateServerPermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/permissions/server/assignments", opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
