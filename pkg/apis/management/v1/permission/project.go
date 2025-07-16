package permission

import (
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ProjectPermissionServiceInterface interface {
		// Get the access to a project
		// opt filters the access by a specific user or role.
		// If not specified, it returns the access for the current user.
		GetAccess(id string, opts *GetProjectAccessOptions, options ...core.RequestOptionFunc) (*GetProjectAccessResponse, *http.Response, error)
		// Get user and role assignments of a project
		// opt filters the assignments by relations.
		// If not specified, it returns all assignments.
		GetAssignments(id string, opts *GetProjectAssignmentsOptions, options ...core.RequestOptionFunc) (*GetProjectAssignmentsResponse, *http.Response, error)
		// Update permissions for the Project
		Update(id string, opts *UpdateProjectPermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// ProjectPermissionService handles communication with project permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	ProjectPermissionService struct {
		client core.Client
	}
)

type ProjectAction string

const (
	CreateWarehouse         ProjectAction = "create_warehouse"
	DeleteProject           ProjectAction = "delete"
	RenameProject           ProjectAction = "rename"
	ListWarehouses          ProjectAction = "list_warehouses"
	CreateRole              ProjectAction = "create_role"
	ListRoles               ProjectAction = "list_roles"
	SearchRoles             ProjectAction = "search_roles"
	ReadProjectAssignments  ProjectAction = "read_assignments"
	GrantProjectRoleCreator ProjectAction = "grant_role_creator"
	GrantProjectCreate      ProjectAction = "grant_create"
	GrantProjectDescribe    ProjectAction = "grant_describe"
	GrantProjectModify      ProjectAction = "grant_modify"
	GrantProjectSelet       ProjectAction = "grant_select"
	GrantProjectAdmin       ProjectAction = "grant_project_admin"
	GrantSecurityAdmin      ProjectAction = "grant_security_admin"
	GrantDataAdmin          ProjectAction = "grant_data_admin"
)

func NewProjectPermissionService(client core.Client) ProjectPermissionServiceInterface {
	return &ProjectPermissionService{
		client: client,
	}
}

// GetProjectAccessOptions represents the GetAccess() options.
//
// Only one of PrincipalUser or PrincipalRole should be set at a time.
// Setting both fields simultaneously is not allowed.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_access
type GetProjectAccessOptions struct {
	PrincipalUser *string `url:"principalUser,omitempty"`
	PrincipalRole *string `url:"principalRole,omitempty"`
}

// GetProjectAccessResponse represents the response from the GetAccess() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_access
type GetProjectAccessResponse struct {
	AllowedActions []ProjectAction `json:"allowed-actions"`
}

// GetAccess retrieves user or role access to a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_access
func (s *ProjectPermissionService) GetAccess(id string, opt *GetProjectAccessOptions, options ...core.RequestOptionFunc) (*GetProjectAccessResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/project/%s/access", id)

	req, err := s.client.NewRequest(http.MethodGet, path, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetProjectAccessResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// GetProjectAssignmentsOptions represents the GetAssignments() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_assignments
type GetProjectAssignmentsOptions struct {
	Relations []ProjectAssignmentType `url:"relations[],omitempty"`
}

// GetProjectAssignmentsResponse represents the response from the GetAssignments() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_assignments
type GetProjectAssignmentsResponse struct {
	Assignments []*ProjectAssignment `json:"assignments"`
	ProjectID   string               `json:"project-id"`
}

// GetAccess gets user and role assignments of a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_assignments
func (s *ProjectPermissionService) GetAssignments(id string, opt *GetProjectAssignmentsOptions, options ...core.RequestOptionFunc) (*GetProjectAssignmentsResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/project/%s/assignments", id)

	req, err := s.client.NewRequest(http.MethodGet, path, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetProjectAssignmentsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// UpdateProjectPermissionsOptions represents the Update() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_project_assignments
type UpdateProjectPermissionsOptions struct {
	// The list of assignments to delete.
	Deletes []*ProjectAssignment `json:"deletes,omitempty"`
	// The list of assignments to create.
	Writes []*ProjectAssignment `json:"writes,omitempty"`
}

// Update updates the project assignments.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_project_assignments
func (s *ProjectPermissionService) Update(id string, opt *UpdateProjectPermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	path := fmt.Sprintf("/permissions/project/%s/assignments", id)

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
