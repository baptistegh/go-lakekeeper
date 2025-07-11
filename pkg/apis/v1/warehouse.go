package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/apis/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/v1/storage/profile"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	WarehouseServiceInterface interface {
		Get(id string, options ...core.RequestOptionFunc) (*Warehouse, *http.Response, error)
		List(opts *ListWarehouseOptions, options ...core.RequestOptionFunc) (*ListWarehouseResponse, *http.Response, error)
		Create(opts *CreateWarehouseOptions, options ...core.RequestOptionFunc) (*CreateWarehouseResponse, *http.Response, error)
		Delete(id string, opts *DeleteWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		SetProtection(id string, protected bool, options ...core.RequestOptionFunc) (*SetProtectionResponse, *http.Response, error)
		Activate(id string, options ...core.RequestOptionFunc) (*http.Response, error)
		Deactivate(id string, options ...core.RequestOptionFunc) (*http.Response, error)
		Rename(id string, opts *RenameWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		UpdateStorageProfile(id string, opts *UpdateStorageProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		UpdateDeleteProfile(id string, opts *UpdateDeleteProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		UpdateStorageCredential(id string, opts *UpdateStorageCredentialOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// WarehouseService handles communication with warehouse endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse
	WarehouseService struct {
		projectID string
		client    core.Client
	}
)

var _ WarehouseServiceInterface = (*WarehouseService)(nil)

// Warehouse represents a lakekeeper warehouse
type Warehouse struct {
	ID             string                 `json:"id"`
	ProjectID      string                 `json:"project-id"`
	Name           string                 `json:"name"`
	Protected      bool                   `json:"protected"`
	Status         WarehouseStatus        `json:"status"`
	StorageProfile profile.StorageProfile `json:"storage-profile"`
	DeleteProfile  *profile.DeleteProfile `json:"delete-profile,omitempty"`
}

type WarehouseStatus string

const (
	WarehouseStatusActive   WarehouseStatus = "active"
	WarehouseStatusInactive WarehouseStatus = "inactive"
)

func (w *Warehouse) IsActive() bool {
	return w.Status == WarehouseStatusActive
}

func NeWarehouseService(client core.Client, projectID string) WarehouseServiceInterface {
	return &WarehouseService{
		projectID: projectID,
		client:    client,
	}
}

// Get retrieves detailed information about a specific warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse
func (s *WarehouseService) Get(id string, options ...core.RequestOptionFunc) (*Warehouse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodGet, "/warehouse/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var wh Warehouse

	resp, apiErr := s.client.Do(req, &wh)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &wh, resp, nil
}

// ListWarehouseOptions represents List() options
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
type ListWarehouseOptions struct {
	WarehouseStatus *WarehouseStatus `url:"warehouseStatus,omitempty"`

	// Deprecated: This field will be removed in a future version.
	// ProjectID should be obtained from the Service itself and is not intended to be used here.
	// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
	ProjectID *string `url:"projectId,omitempty"`
}

// listWarehouseResponse represents the response on list warehouses API action
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
type ListWarehouseResponse struct {
	Warehouses []*Warehouse `json:"warehouses"`
}

// Returns all warehouses in the project that the current user has access to.
// By default, deactivated warehouses are not included in the results.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
func (s *WarehouseService) List(opts *ListWarehouseOptions, options ...core.RequestOptionFunc) (*ListWarehouseResponse, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opts == nil {
		opts = &ListWarehouseOptions{}
	}
	opts.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodGet, "/warehouse", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var whs ListWarehouseResponse

	resp, apiErr := s.client.Do(req, &whs)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &whs, resp, nil
}

// CreateOptions represents Create() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
type CreateWarehouseOptions struct {
	Name string `json:"warehouse-name"`
	// Deprecated: This field will be removed in a future version.
	// ProjectID should be obtained from the Service itself and is not intended to be used here.
	// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
	ProjectID         *string                      `json:"project-id,omitempty"`
	StorageProfile    profile.StorageProfile       `json:"storage-profile"`
	StorageCredential credential.StorageCredential `json:"storage-credential"`
	DeleteProfile     *profile.DeleteProfile       `json:"delete-profile,omitempty"`
}

// CreateOptions represents the response from the API
// on a create_warehouse action.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
type CreateWarehouseResponse struct {
	ID string `json:"warehouse-id"`
}

// Create creates a new warehouse in the specified project with
// the provided configuration.
// The project of a warehouse cannot be changed after creation.
// This operation validates the storage configuration.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
func (s *WarehouseService) Create(opts *CreateWarehouseOptions, options ...core.RequestOptionFunc) (*CreateWarehouseResponse, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opts == nil {
		opts = &CreateWarehouseOptions{}
	}
	opts.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, "/warehouse", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var whResp CreateWarehouseResponse

	resp, apiErr := s.client.Do(req, &whResp)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &whResp, resp, nil
}

// RenameWarehouseOptions represents WarehouseService.Rename() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/rename_warehouse
type RenameWarehouseOptions struct {
	NewName string `json:"new-name"`
}

// Rename updates the name of a specific warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/rename_warehouse
func (s *WarehouseService) Rename(id string, opts *RenameWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("/warehouse/%s/rename", id), opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// DeleteWarehouseOptions represents Delete() options.
//
// force parameters needs to be true to delete protected warehouses.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/delete_warehouse
type DeleteWarehouseOptions struct {
	Force bool `url:"force"`
}

// Delete permanently removes a warehouse and all its associated resources.
// Use the force parameter to delete protected warehouses.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/delete_warehouse
func (s *WarehouseService) Delete(id string, opts *DeleteWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodDelete, "/warehouse/"+id, opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// SetProtectionResponse represent the reponse sent by SetProtection()
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_warehouse_protection
type SetProtectionResponse struct {
	Protected bool   `json:"protected"`
	UpdatedAt string `json:"updated_at"`
}

// setWarehouseProtectionOptions represent the request sent to SetProtection()
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_warehouse_protection
type setWarehouseProtectionOptions struct {
	Protected bool `json:"protected"`
}

// SetProtection configures whether a warehouse should be protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_warehouse_protection
func (s *WarehouseService) SetProtection(id string, protected bool, options ...core.RequestOptionFunc) (*SetProtectionResponse, *http.Response, error) {
	opts := setWarehouseProtectionOptions{
		Protected: protected,
	}

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("/warehouse/%s/protection", id), &opts, options)
	if err != nil {
		return nil, nil, err
	}

	var wProtec SetProtectionResponse
	resp, apiErr := s.client.Do(req, &wProtec)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &wProtec, resp, nil
}

// Activate re-enables access to a previously deactivated warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/activate_warehouse
func (s *WarehouseService) Activate(id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("/warehouse/%s/activate", id), nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// Deactivate temporarily disables access to a warehouse without deleting its data.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/deactivate_warehouse
func (s *WarehouseService) Deactivate(id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("/warehouse/%s/deactivate", id), nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// UpdateStorageProfileOptions represent UpdateStorageProfile() options
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_profile
type UpdateStorageProfileOptions struct {
	StorageCredential *credential.StorageCredential `json:"storage-credential,omitempty"`
	StorageProfile    profile.StorageProfile        `json:"storage-profile"`
}

// Deactivate updates both the storage profile and credentials of a warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_profile
func (s *WarehouseService) UpdateStorageProfile(id string, opts *UpdateStorageProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	if opts == nil {
		return nil, errors.New("update storage profile received empty options")
	}

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("/warehouse/%s/storage", id), opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// UpdateDeleteProfileOptions represent UpdateDeleteProfile() options
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_warehouse_delete_profile
type UpdateDeleteProfileOptions struct {
	DeleteProfile profile.DeleteProfile `json:"delete-profile"`
}

// UpdateDeleteProfile configures the soft-delete behavior for a warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_warehouse_delete_profile
func (s *WarehouseService) UpdateDeleteProfile(id string, opts *UpdateDeleteProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	if opts == nil {
		return nil, errors.New("update delete profile received empty options")
	}

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("/warehouse/%s/delete-profile", id), opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// UpdateStorageCredentialOptions represent UpdateStorageCredential() options
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_credential
type UpdateStorageCredentialOptions struct {
	StorageCredential *credential.StorageCredential `json:"new-storage-credential,omitempty"`
}

// Deactivate updates only the storage credential of a warehouse without modifying the storage profile.
// Useful for refreshing expiring credentials.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_credential
func (s *WarehouseService) UpdateStorageCredential(id string, opts *UpdateStorageCredentialOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf("/warehouse/%s/storage-credential", id), opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
