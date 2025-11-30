package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	WarehouseServiceInterface interface {
		// TODO: implement missing API endpoints
		// GetTaskQueueConfig (expiration)
		// SetTaskQueueConfig (expiration)
		// GetTaskQueueConfig (purge)
		// SetTaskQueueConfig (purge)

		// Returns all warehouses in the project that the current user has access to.
		// By default, deactivated warehouses are not included in the results.
		List(ctx context.Context, opt *ListWarehouseOptions, options ...core.RequestOptionFunc) (*ListWarehouseResponse, *http.Response, error)
		// Creates a new warehouse in the specified project with the provided configuration.
		// The project of a warehouse cannot be changed after creation.
		// This operation validates the storage configuration.
		Create(ctx context.Context, opt *CreateWarehouseOptions, options ...core.RequestOptionFunc) (*CreateWarehouseResponse, *http.Response, error)
		// Retrieves detailed information about a specific warehouse.
		Get(ctx context.Context, id string, options ...core.RequestOptionFunc) (*Warehouse, *http.Response, error)
		// Permanently removes a warehouse and all its associated resources.
		// Use the force parameter to delete protected warehouses.
		Delete(ctx context.Context, id string, opt *DeleteWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Re-enables access to a previously deactivated warehouse.
		Activate(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error)
		// Temporarily disables access to a warehouse without deleting its data.
		Deactivate(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error)
		// Configures the soft-delete behavior for a warehouse.
		UpdateDeleteProfile(ctx context.Context, id string, opt *UpdateDeleteProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Configures whether a warehouse should be protected from deletion.
		//
		// Deprecated: user SetWarehouseProtection instead. This will be remove in the future.
		SetProtection(ctx context.Context, id string, protected bool, options ...core.RequestOptionFunc) (*SetProtectionResponse, *http.Response, error)
		// Configures whether a warehouse should be protected from deletion.
		SetWarehouseProtection(ctx context.Context, id string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error)
		// Updates the name of a specific warehouse.
		Rename(ctx context.Context, id string, opt *RenameWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Updates both the storage profile and credentials of a warehouse.
		UpdateStorageProfile(ctx context.Context, id string, opt *UpdateStorageProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Updates only the storage credential of a warehouse without modifying the storage profile.
		// Useful for refreshing expiring credentials.
		UpdateStorageCredential(ctx context.Context, id string, opt *UpdateStorageCredentialOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Returns soft-deleted tables and views in the warehouse that are visible to the current user.
		ListSoftDeletedTabulars(ctx context.Context, id string, opt *ListSoftDeletedTabularsOptions, options ...core.RequestOptionFunc) (*ListSoftDeletedTabularsResponse, *http.Response, error)
		// Restores previously deleted tables or views to make them accessible again.
		UndropTabular(ctx context.Context, id string, opt *UndropTabularOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Retrieves whether a namespace is protected from deletion.
		GetNamespaceProtection(ctx context.Context, warehouseID, namespaceID string, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error)
		// Configures whether a namespace should be protected from deletion.
		SetNamespaceProtection(ctx context.Context, warehouseID, namespaceID string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error)
		// Retrieves statistical data about a warehouse's usage and resources over time. Statistics are aggregated hourly when changes occur.
		// We lazily create a new statistics entry every hour, in between hours, the existing entry is being updated.
		// If there's a change at created_at + 1 hour, a new entry is created. If there's been no change, no new entry is created, meaning there may be gaps.
		GetStatistics(ctx context.Context, id string, opt *GetStatisticsOptions, options ...core.RequestOptionFunc) (*GetStatisticsResponse, *http.Response, error)
		// Retrieves whether a table is protected from deletion.
		GetTableProtection(ctx context.Context, warehouseID, tableID string, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error)
		// Retrieves whether a view is protected from deletion.
		GetViewProtection(ctx context.Context, warehouseID, viewID string, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error)
		// Configures whether a table should be protected from deletion.
		SetTableProtection(ctx context.Context, warehouseID, tableID string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error)
		// Configures whether a view should be protected from deletion.
		SetViewProtection(ctx context.Context, warehouseID, viewID string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error)
		// Get user allowed actions for a warehouse
		GetAllowedActions(ctx context.Context, warehouseID string, opt *GetWarehouseAllowedActionsOptions, options ...core.RequestOptionFunc) (*GetWarehouseAllowedActionsResponse, *http.Response, error)
	}

	// WarehouseService handles communication with warehouse endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse
	WarehouseService struct {
		projectID string
		client    core.Client
	}

	WarehouseStatus string

	// Warehouse represents a lakekeeper warehouse
	Warehouse struct {
		ID             string                 `json:"id"`
		ProjectID      string                 `json:"project-id"`
		Name           string                 `json:"name"`
		Protected      bool                   `json:"protected"`
		Status         WarehouseStatus        `json:"status"`
		StorageProfile profile.StorageProfile `json:"storage-profile"`
		DeleteProfile  *profile.DeleteProfile `json:"delete-profile,omitempty"`
	}

	TabularType string

	Tabular struct {
		// Unique identifier of the tabular
		ID string `json:"id"`
		// Name of the tabular
		Name string `json:"name"`
		// Warehouse ID where the tabular is stored
		WarehouseID string `json:"warehouse-id"`
		// List of namespace parts the tabular belongs to
		Namespace []string `json:"namespace"`
		// Type of the tabular
		Type TabularType `json:"typ"`
		// Date when the tabular will not be recoverable anymore
		ExpirationDate string `json:"expiration-date"`
		// Date when the tabular was deleted
		DeletedAt string `json:"deleted-at"`
		// Date when the tabular was created
		CreatedAt string `json:"created-at"`
	}

	// ListWarehouseOptions represents List() options
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
	ListWarehouseOptions struct {
		WarehouseStatus []WarehouseStatus `url:"warehouseStatus[],omitempty"`

		// Deprecated: This field will be removed in a future version.
		// ProjectID should be obtained from the Service itself and is not intended to be used here.
		// It is temporarily kept for compatibility with the Lakekeeper API until it gets removed upstream.
		ProjectID *string `url:"projectId,omitempty"`
	}

	// listWarehouseResponse represents the response on list warehouses API action
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
	ListWarehouseResponse struct {
		Warehouses []*Warehouse `json:"warehouses"`
	}

	// CreateOptions represents Create() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
	CreateWarehouseOptions struct {
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
	CreateWarehouseResponse struct {
		ID string `json:"warehouse-id"`
	}

	// RenameWarehouseOptions represents WarehouseService.Rename() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/rename_warehouse
	RenameWarehouseOptions struct {
		NewName string `json:"new-name"`
	}

	// DeleteWarehouseOptions represents Delete() options.
	//
	// force parameters needs to be true to delete protected warehouses.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/delete_warehouse
	DeleteWarehouseOptions struct {
		Force *bool `url:"force,omitempty"`
	}

	// SetProtectionResponse represent the response sent by SetProtection()
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_warehouse_protection
	//
	// Deprecated: use SetWarehouseProtection instead. This will be remove in the future
	SetProtectionResponse struct {
		Protected bool    `json:"protected"`
		UpdatedAt *string `json:"updated_at,omitempty"`
	}

	// UpdateStorageProfileOptions represent UpdateStorageProfile() options
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_profile
	UpdateStorageProfileOptions struct {
		StorageCredential *credential.StorageCredential `json:"storage-credential,omitempty"`
		StorageProfile    profile.StorageProfile        `json:"storage-profile"`
	}

	// UpdateDeleteProfileOptions represent UpdateDeleteProfile() options
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_warehouse_delete_profile
	UpdateDeleteProfileOptions struct {
		DeleteProfile profile.DeleteProfile `json:"delete-profile"`
	}

	// UpdateStorageCredentialOptions represent UpdateStorageCredential() options
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_credential
	UpdateStorageCredentialOptions struct {
		StorageCredential *credential.StorageCredential `json:"new-storage-credential,omitempty"`
	}
	// ListSoftDeletedTabularsOptions represents ListSoftDeletedTabulars() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_deleted_tabulars
	ListSoftDeletedTabularsOptions struct {
		// Filter by Namespace ID
		NamespaceID *string `url:"namespaceId"`

		ListOptions `url:",inline"`
	}

	// ListSoftDeletedTabularsResponse represents ListSoftDeletedTabulars() response.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_deleted_tabulars
	ListSoftDeletedTabularsResponse struct {
		// List of the tabulars
		Tabulars []*Tabular `json:"tabulars"`

		ListResponse `json:",inline"`
	}

	// UndropTabular restores previously deleted tables or views to make them accessible again.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_deleted_tabulars
	UndropTabularOptions struct {
		Targets []struct {
			ID   string      `json:"id"`
			Type TabularType `json:"type"`
		} `json:"targets"`
	}

	// SetProtectionOptions represents protection-related methods options
	SetProtectionOptions struct {
		// Setting this to true will prevent the entity from being deleted unless force is used.
		Protected bool `json:"protected"`
	}

	// GetProtectionResponse represents protection-related methods response.
	GetProtectionResponse struct {
		// Indicates wether the entity is protected
		Protected bool `json:"protected"`
		// Updated At
		UpdatedAt *string `json:"updated_at,omitempty"`
	}

	// GetStatisticsResponse represents GetStatistics() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse_statistics
	GetStatisticsOptions struct {
		// Next page token
		PageToken *string `url:"page_token,omitempty"`
		// Signals an upper bound of the number of results that a client will receive
		PageSize *int64 `url:"page_size,omitempty"`
	}

	// GetStatisticsResponse represents GetStatistics() response.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse_statistics
	GetStatisticsResponse struct {
		// ID of the warehouse for which the stats were collected.
		WarehouseID string `json:"warehouse-ident"`
		// Ordered list of warehouse statistics.
		Stats []struct {
			// Number of tables in the warehouse.
			NumberOfTables int64 `json:"number-of-tables"`
			// Number of views in the warehouse.
			NumberOfView int64 `json:"number-of-views"`
			// Timestamp of when these statistics are valid until.
			// We lazily create a new statistics entry every hour, in between hours, the existing entry is being updated.
			// If there's a change at created_at + 1 hour, a new entry is created.
			// If there's no change, no new entry is created.
			Timestamp string `json:"timestamp"`
			// Timestamp of when these statistics were last updated.
			UpdatedAt string `json:"updated-at"`
		} `json:"stats"`

		ListResponse `json:",inline"`
	}

	// GetWarehouseAllowedActionsOptions represents the GetAllowedActions() options.
	//
	// Only one of PrincipalUser or PrincipalRole should be set at a time.
	// Setting both fields simultaneously is not allowed.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse_actions
	GetWarehouseAllowedActionsOptions struct {
		// The user to show actions for.
		PrincipalUser *string `url:"principalUser,omitempty"`
		// The role to show actions for.
		PrincipalRole *string `url:"principalRole,omitempty"`
	}

	// GetWarehouseAllowedAccessResponse represents the GetAllowedActions() response.
	//
	//
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse_actions
	GetWarehouseAllowedActionsResponse struct {
		AllowedActions []permission.WarehouseAction `json:"allowed-actions"`
	}
)

const (
	WarehouseStatusActive   WarehouseStatus = "active"
	WarehouseStatusInactive WarehouseStatus = "inactive"

	TableTabularType TabularType = "table"
	ViewTabularType  TabularType = "view"
)

func (w *Warehouse) IsActive() bool {
	return w.Status == WarehouseStatusActive
}

func NewWarehouseService(client core.Client, projectID string) WarehouseServiceInterface {
	return &WarehouseService{
		projectID: projectID,
		client:    client,
	}
}

// Get retrieves detailed information about a specific warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse
func (s *WarehouseService) Get(ctx context.Context, id string, options ...core.RequestOptionFunc) (*Warehouse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, "/warehouse/"+id, nil, options)
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

// Returns all warehouses in the project that the current user has access to.
// By default, deactivated warehouses are not included in the results.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
func (s *WarehouseService) List(ctx context.Context, opt *ListWarehouseOptions, options ...core.RequestOptionFunc) (*ListWarehouseResponse, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opt == nil {
		opt = &ListWarehouseOptions{}
	}
	opt.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, "/warehouse", opt, options)
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

// Create creates a new warehouse in the specified project with
// the provided configuration.
// The project of a warehouse cannot be changed after creation.
// This operation validates the storage configuration.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
func (s *WarehouseService) Create(ctx context.Context, opt *CreateWarehouseOptions, options ...core.RequestOptionFunc) (*CreateWarehouseResponse, *http.Response, error) {
	// This workaround will be removed once project-id is no longer required
	// in the request by the API.
	// https://github.com/lakekeeper/lakekeeper/issues/1234
	if opt == nil {
		opt = &CreateWarehouseOptions{}
	}
	opt.ProjectID = &s.projectID

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/warehouse", opt, options)
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

// Rename updates the name of a specific warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/rename_warehouse
func (s *WarehouseService) Rename(ctx context.Context, id string, opt *RenameWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/rename", id), opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// Delete permanently removes a warehouse and all its associated resources.
// Use the force parameter to delete protected warehouses.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/delete_warehouse
func (s *WarehouseService) Delete(ctx context.Context, id string, opt *DeleteWarehouseOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodDelete, "/warehouse/"+id, opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// SetProtection configures whether a warehouse should be protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_warehouse_protection
//
// Deprecated: use SetWarehouseProtection instead. This will be remove in the future.
func (s *WarehouseService) SetProtection(ctx context.Context, id string, protected bool, options ...core.RequestOptionFunc) (*SetProtectionResponse, *http.Response, error) {
	opt := SetProtectionOptions{
		Protected: protected,
	}

	resp, r, err := s.SetWarehouseProtection(ctx, id, &opt, options...)
	if err != nil {
		return nil, r, err
	}

	return &SetProtectionResponse{resp.Protected, resp.UpdatedAt}, r, nil
}

// SetWarehouseProtection configures whether a warehouse should be protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_warehouse_protection
func (s *WarehouseService) SetWarehouseProtection(ctx context.Context, id string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/protection", id), &opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetProtectionResponse
	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// Activate re-enables access to a previously deactivated warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/activate_warehouse
func (s *WarehouseService) Activate(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/activate", id), nil, options)
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
func (s *WarehouseService) Deactivate(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/deactivate", id), nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// Deactivate updates both the storage profile and credentials of a warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_profile
func (s *WarehouseService) UpdateStorageProfile(ctx context.Context, id string, opt *UpdateStorageProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	if opt == nil {
		return nil, errors.New("update storage profile received empty options")
	}

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/storage", id), opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// UpdateDeleteProfile configures the soft-delete behavior for a warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_warehouse_delete_profile
func (s *WarehouseService) UpdateDeleteProfile(ctx context.Context, id string, opt *UpdateDeleteProfileOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	if opt == nil {
		return nil, errors.New("update delete profile received empty options")
	}

	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/delete-profile", id), opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// Deactivate updates only the storage credential of a warehouse without modifying the storage profile.
// Useful for refreshing expiring credentials.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/update_storage_credential
func (s *WarehouseService) UpdateStorageCredential(ctx context.Context, id string, opt *UpdateStorageCredentialOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/storage-credential", id), opt, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// ListSoftDeletedTabulars returns all soft-deleted tables and views in the warehouse that are visible to the current user.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_deleted_tabulars
func (s *WarehouseService) ListSoftDeletedTabulars(ctx context.Context, id string, opt *ListSoftDeletedTabularsOptions, options ...core.RequestOptionFunc) (*ListSoftDeletedTabularsResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/warehouse/%s/deleted-tabulars", id), opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp ListSoftDeletedTabularsResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// UndropTabular restores previously deleted tables or views to make them accessible again.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_deleted_tabulars
func (s *WarehouseService) UndropTabular(ctx context.Context, id string, opt *UndropTabularOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/deleted-tabulars/undrop", id), opt, options)
	if err != nil {
		return nil, err
	}

	r, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return r, apiErr
	}

	return r, nil
}

// GetNamespaceProtection retrieves whether a namespace is protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_namespace_protection
func (s *WarehouseService) GetNamespaceProtection(ctx context.Context, warehouseID, namespaceID string, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/warehouse/%s/namespace/%s/protection", warehouseID, namespaceID), nil, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetProtectionResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// SetNamespaceProtection configures whether a namespace should be protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_namespace_protection
func (s *WarehouseService) SetNamespaceProtection(ctx context.Context, warehouseID, namespaceID string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/namespace/%s/protection", warehouseID, namespaceID), opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetProtectionResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// GetTableProtection retrieves whether a table is protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_table_protection
func (s *WarehouseService) GetTableProtection(ctx context.Context, warehouseID, tableID string, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/warehouse/%s/table/%s/protection", warehouseID, tableID), nil, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetProtectionResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// SetTableProtection configures whether a table should be protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_table_protection
func (s *WarehouseService) SetTableProtection(ctx context.Context, warehouseID, tableID string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/table/%s/protection", warehouseID, tableID), opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetProtectionResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// GetViewProtection retrieves whether a view is protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_view_protection
func (s *WarehouseService) GetViewProtection(ctx context.Context, warehouseID, viewID string, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/warehouse/%s/view/%s/protection", warehouseID, viewID), nil, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetProtectionResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// SetViewProtection configures whether a view should be protected from deletion.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/set_view_protection
func (s *WarehouseService) SetViewProtection(ctx context.Context, warehouseID, viewID string, opt *SetProtectionOptions, options ...core.RequestOptionFunc) (*GetProtectionResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf("/warehouse/%s/view/%s/protection", warehouseID, viewID), opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetProtectionResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

// GetStatistics Retrieves statistical data about a warehouse's usage and resources over time. Statistics are aggregated hourly when changes occur.
//
// We lazily create a new statistics entry every hour, in between hours, the existing entry is being updated.
// If there's a change at created_at + 1 hour, a new entry is created.
// If there's been no change, no new entry is created, meaning there may be gaps.
//
// Example:
//
// 00:16:32: warehouse created:
//
//	timestamp: 01:00:00, created_at: 00:16:32, updated_at: null, 0 tables, 0 views
//
// 00:30:00: table created:
//
//	timestamp: 01:00:00, created_at: 00:16:32, updated_at: 00:30:00, 1 table, 0 views
//
// 00:45:00: view created:
//
//	timestamp: 01:00:00, created_at: 00:16:32, updated_at: 00:45:00, 1 table, 1 view
//
// 01:00:36: table deleted:
//
//	timestamp: 02:00:00, created_at: 01:00:36, updated_at: null, 0 tables, 1 view
//	timestamp: 01:00:00, created_at: 00:16:32, updated_at: 00:45:00, 1 table, 1 view
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse_statistics
func (s *WarehouseService) GetStatistics(ctx context.Context, id string, opt *GetStatisticsOptions, options ...core.RequestOptionFunc) (*GetStatisticsResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/warehouse/%s/statistics", id), opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetStatisticsResponse

	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

func (s *WarehouseService) GetAllowedActions(ctx context.Context, id string, opt *GetWarehouseAllowedActionsOptions, options ...core.RequestOptionFunc) (*GetWarehouseAllowedActionsResponse, *http.Response, error) {
	options = append(options, WithProject(s.projectID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/warehouse/%s/actions", id), opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetWarehouseAllowedActionsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}
