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
		// Get allowed Authorizer actions on a role
		// opts filters the access by a specific user or role.
		// If not specified, it returns the access for the current user.
		GetAllowedAuthorizerActions(ctx context.Context, id string, opts *GetRoleAllowedAuthorizerActionsOptions, options ...core.RequestOptionFunc) (*GetRoleAllowedAuthorizerActionsResponse, *http.Response, error)
	}

	// RolePermissionService handles communication with role permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	RolePermissionService struct {
		client core.Client
	}

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

	// GetRoleAllowedAuthorizerActionsOptions represents the GetAllowedAuthorizerActions() options.
	//
	// Only one of PrincipalUser or PrincipalRole should be set at a time.
	// Setting both fields simultaneously is not allowed.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions-openfga/operation/get_authorizer_role_actions
	GetRoleAllowedAuthorizerActionsOptions struct {
		PrincipalUser *string `url:"principalUser,omitempty"`
		PrincipalRole *string `url:"principalRole,omitempty"`
	}

	// GetRoleAllowedAuthorizerActionsResponse represents the response from the GetAllowedAuthorizerActions() endpoint.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions-openfga/operation/get_authorizer_role_actions
	GetRoleAllowedAuthorizerActionsResponse struct {
		AllowedActions []OpenFGARoleAction `json:"allowed-actions"`
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

// Available actions on a role
type RoleAction string

const (
	Assume              RoleAction = "assume"
	CanGrantAssignee    RoleAction = "can_grant_assignee"
	CanChangeOwnership  RoleAction = "can_change_ownership"
	DeleteRole          RoleAction = "delete"
	UpdateRole          RoleAction = "update"
	ReadRole            RoleAction = "read"
	ReadRoleAssignments RoleAction = "read_assignments"
)

type OpenFGARoleAction string

const (
	RoleAssume             OpenFGARoleAction = "assume"
	RoleGrantAssignee      OpenFGARoleAction = "can_grant_assignee"
	RoleCanChangeOwnership OpenFGARoleAction = "can_change_ownership"
	RoleReadAssignments    OpenFGARoleAction = "read_assignments"
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

// GetAllowedAuthorizerActions gets allowed Authorizer actions on a role
//
// Returns Authorizer permissions (OpenFGA relations) for the specified role.
// For Catalog permissions, use /management/v1/role/{role_id}/actions instead.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions-openfga/operation/get_authorizer_role_actions
func (s *RolePermissionService) GetAllowedAuthorizerActions(ctx context.Context, id string, opt *GetRoleAllowedAuthorizerActionsOptions, options ...core.RequestOptionFunc) (*GetRoleAllowedAuthorizerActionsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/permissions/role/%s/authorizer-actions", id), opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetRoleAllowedAuthorizerActionsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}
