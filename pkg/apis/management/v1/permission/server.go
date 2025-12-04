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
		// Get allowed Authorizer actions on the server
		// opts filters the access by a specific user or role.
		// If not specified, it returns the access for the current user.
		GetAllowedAuthorizerActions(ctx context.Context, opts *GetServerAllowedAuthorizerActionsOptions, options ...core.RequestOptionFunc) (*GetServerAllowedAuthorizerActionsResponse, *http.Response, error)
	}

	// ServerPermissionService handles communication with server permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	ServerPermissionService struct {
		client core.Client
	}

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

	// GetServerAllowedAuthorizerActionsOptions represents the GetAllowedAuthorizerActions() options.
	//
	// Only one of PrincipalUser or PrincipalRole should be set at a time.
	// Setting both fields simultaneously is not allowed.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions-openfga/operation/get_authorizer_server_actions
	GetServerAllowedAuthorizerActionsOptions struct {
		PrincipalUser *string `url:"principalUser,omitempty"`
		PrincipalRole *string `url:"principalRole,omitempty"`
	}

	// GetServerAllowedAuthorizerActionsResponse represents the response from the GetAllowedAuthorizerActions() endpoint.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions-openfga/operation/get_authorizer_server_actions
	GetServerAllowedAuthorizerActionsResponse struct {
		AllowedActions []OpenFGAServerAction `json:"allowed-actions"`
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

// Available actions on a server
type ServerAction string

const (
	CreateProject    ServerAction = "create_project"
	UpdateUsers      ServerAction = "update_users"
	DeleteUsers      ServerAction = "delete_users"
	ListUsers        ServerAction = "list_users"
	ProvisionUsers   ServerAction = "provision_users"
	GrantServerAdmin ServerAction = "grant_admin"
	ReadAssignments  ServerAction = "read_assignments"
)

// Available authorizer actions on a server
type OpenFGAServerAction string

const (
	ServerGrantAdmin      OpenFGAServerAction = "grant_admin"
	ServerReadAssignments OpenFGAServerAction = "read_assignments"
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

// GetAllowedAuthorizerActions gets allowed Authorizer actions on the server
//
// Returns Authorizer permissions (OpenFGA relations) for the server.
// For Catalog permissions, use /management/v1/server/actions instead.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions-openfga/operation/get_authorizer_server_actions
func (s *ServerPermissionService) GetAllowedAuthorizerActions(ctx context.Context, opt *GetServerAllowedAuthorizerActionsOptions, options ...core.RequestOptionFunc) (*GetServerAllowedAuthorizerActionsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/permissions/server/authorizer-actions", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetServerAllowedAuthorizerActionsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}
