package permission

import (
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	RolePermissionServiceInterface interface {
		// Get a role assignments
		// opt filters the assignments by relations.
		// If not specified, it returns all assignments.
		GetAssignments(id string, opts *GetRoleAssignmentsOptions, options ...core.RequestOptionFunc) (*GetRoleAssignmentsResponse, *http.Response, error)
		// Update permissions for a role
		Update(id string, opts *UpdateRolePermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// RolePermissionService handles communication with role permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	RolePermissionService struct {
		client core.Client
	}
)

func NewRolePermissionService(client core.Client) RolePermissionServiceInterface {
	return &RolePermissionService{
		client: client,
	}
}

// GetRoleAssignmentsOptions represents the GetAssignments() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_assignments
type GetRoleAssignmentsOptions struct {
	Relations []RoleAssignmentType `url:"relations[],omitempty"`
}

// GetRoleAssignmentsResponse represents the response from the GetAssignments() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_assignments
type GetRoleAssignmentsResponse struct {
	Assignments []*RoleAssignment `json:"assignments"`
}

// GetAccess gets user and role assignments of the role.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_role_assignments
func (s *RolePermissionService) GetAssignments(id string, opt *GetRoleAssignmentsOptions, options ...core.RequestOptionFunc) (*GetRoleAssignmentsResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/role/%s/assignments", id)

	req, err := s.client.NewRequest(http.MethodGet, path, opt, options)
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

// UpdateRolePermissionsOptions represents the Update() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_role_assignments
type UpdateRolePermissionsOptions struct {
	// The list of assignments to delete.
	Deletes []*RoleAssignment `json:"deletes,omitempty"`
	// The list of assignments to create.
	Writes []*RoleAssignment `json:"writes,omitempty"`
}

// Update updates the role assignments.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_role_assignments
func (s *RolePermissionService) Update(id string, opt *UpdateRolePermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	path := fmt.Sprintf("/permissions/role/%s/assignments", id)

	req, err := s.client.NewRequest(http.MethodPost, path, opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
