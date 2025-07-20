package permission

import (
	"context"
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	WarehousePermissionServiceInterface interface {
		GetAuthzProperties(ctx context.Context, id string, options ...core.RequestOptionFunc) (*GetWarehouseAuthzPropertiesResponse, *http.Response, error)
		// Get the access to a warehouse
		// opt filters the access by a specific user or role.
		// If not specified, it returns the access for the current user.
		GetAccess(ctx context.Context, id string, opt *GetWarehouseAccessOptions, options ...core.RequestOptionFunc) (*GetWarehouseAccessResponse, *http.Response, error)
		// Get a warehouse assignments
		// opt filters the assignments by relations.
		// If not specified, it returns all assignments.
		GetAssignments(ctx context.Context, id string, opt *GetWarehouseAssignmentsOptions, options ...core.RequestOptionFunc) (*GetWarehouseAssignmentsResponse, *http.Response, error)
		// Update permissions for a warehouse
		Update(ctx context.Context, id string, opts *UpdateWarehousePermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Set managed access property of a warehouse
		SetManagedAccess(ctx context.Context, id string, opts *SetWarehouseManagedAccessOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// WarehousePermissionService handles communication with warehouse permissions endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	WarehousePermissionService struct {
		client core.Client
	}
)

func NewWarehousePermissionService(client core.Client) WarehousePermissionServiceInterface {
	return &WarehousePermissionService{
		client: client,
	}
}

type WarehouseAction string

const (
	CreateNamespace          WarehouseAction = "create_namespace"
	DeleteWarehouse          WarehouseAction = "delete"
	ModifyStorage            WarehouseAction = "modify_storage"
	ModifyStorageCredential  WarehouseAction = "modify_storage_credential"
	GetConfig                WarehouseAction = "get_config"
	GetMetadata              WarehouseAction = "get_metadata"
	ListNamespaces           WarehouseAction = "list_namespaces"
	IncludeInList            WarehouseAction = "include_in_list"
	Deactivate               WarehouseAction = "deactivate"
	Activate                 WarehouseAction = "activate"
	Rename                   WarehouseAction = "rename"
	ListDeletedTabulars      WarehouseAction = "list_deleted_tabulars"
	ReadWarehouseAssignments WarehouseAction = "read_assignments"
	GrantCreate              WarehouseAction = "grant_create"
	GrantDescribe            WarehouseAction = "grant_describe"
	GrantModify              WarehouseAction = "grant_modify"
	GrantSelect              WarehouseAction = "grant_select"
	GrantPassGrants          WarehouseAction = "grant_pass_grants"
	GrantManageGrants        WarehouseAction = "grant_manage_grants"
	ChangeOwnership          WarehouseAction = "change_ownership"
)

// GetWarehouseAuthzPropertiesResponse represents the response from the GetAuthzProperties() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access
type GetWarehouseAuthzPropertiesResponse struct {
	ManagedAccess bool `json:"managed-access"`
}

// GetAuthzProperties retrieves authorization properties of a warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access
func (s *WarehousePermissionService) GetAuthzProperties(ctx context.Context, id string, options ...core.RequestOptionFunc) (*GetWarehouseAuthzPropertiesResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/warehouse/%s", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetWarehouseAuthzPropertiesResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// GetWarehouseAccessOptions represents the GetAccess() options.
//
// Only one of PrincipalUser or PrincipalRole should be set at a time.
// Setting both fields simultaneously is not allowed.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access
type GetWarehouseAccessOptions struct {
	PrincipalUser *string `url:"principalUser,omitempty"`
	PrincipalRole *string `url:"principalRole,omitempty"`
}

// GetWarehouseAccessResponse represents the response from the GetAccess() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access
type GetWarehouseAccessResponse struct {
	AllowedActions []WarehouseAction `json:"allowed-actions"`
}

// GetAccess retrieves user or role access to a warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access
func (s *WarehousePermissionService) GetAccess(ctx context.Context, id string, opt *GetWarehouseAccessOptions, options ...core.RequestOptionFunc) (*GetWarehouseAccessResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/warehouse/%s/access", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetWarehouseAccessResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// GetWarehouseAssignmentsOptions represents the GetAssignments() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_assignments
type GetWarehouseAssignmentsOptions struct {
	Relations []WarehouseAssignmentType `url:"relations[],omitempty"`
}

// GetWarehouseAssignmentsResponse represents the response from the GetAssignments() endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_assignments
type GetWarehouseAssignmentsResponse struct {
	Assignments []*WarehouseAssignment `json:"assignments"`
}

// GetAssignments gets user and role assignments of the warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_assignments
func (s *WarehousePermissionService) GetAssignments(ctx context.Context, id string, opt *GetWarehouseAssignmentsOptions, options ...core.RequestOptionFunc) (*GetWarehouseAssignmentsResponse, *http.Response, error) {
	path := fmt.Sprintf("/permissions/warehouse/%s/assignments", id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetWarehouseAssignmentsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}

// UpdateWarehousePermissionsOptions represents the Update() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_warehouse_assignments
type UpdateWarehousePermissionsOptions struct {
	// The list of assignments to delete.
	Deletes []*WarehouseAssignment `json:"deletes,omitempty"`
	// The list of assignments to create.
	Writes []*WarehouseAssignment `json:"writes,omitempty"`
}

// Update updates the warehouse assignments.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_warehouse_assignments
func (s *WarehousePermissionService) Update(ctx context.Context, id string, opt *UpdateWarehousePermissionsOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	path := fmt.Sprintf("/permissions/warehouse/%s/assignments", id)

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

// SetWarehouseManagedAccessOptions represents SetManagedAccess() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/set_warehouse_managed_access
type SetWarehouseManagedAccessOptions struct {
	ManagedAccess bool `json:"managed-access"`
}

// SetManagedAccess sets managed access property of a warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/set_warehouse_managed_access
func (s *WarehousePermissionService) SetManagedAccess(ctx context.Context, id string, opt *SetWarehouseManagedAccessOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	path := fmt.Sprintf("/permissions/warehouse/%s/managed-access", id)

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
