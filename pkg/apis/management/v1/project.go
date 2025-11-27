package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ProjectServiceInterface interface {
		// Retrieves information about a project.
		Get(ctx context.Context, id string, options ...core.RequestOptionFunc) (*Project, *http.Response, error)
		// Creates a new project with the specified configuration.
		Create(ctx context.Context, opts *CreateProjectOptions, options ...core.RequestOptionFunc) (*CreateProjectResponse, *http.Response, error)
		// Deletes a project.
		Delete(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error)
		// Lists all projects that the requesting user has access to.
		List(ctx context.Context, options ...core.RequestOptionFunc) (*ListProjectsResponse, *http.Response, error)
		// Renames a project.
		Rename(ctx context.Context, id string, opts *RenameProjectOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		// Retrieves detailed endpoint call statistics for your project, allowing you to monitor API usage patterns,
		// track frequency of operations, and analyze response codes.
		GetAPIStatistics(ctx context.Context, id string, opt *GetAPIStatisticsOptions, options ...core.RequestOptionFunc) (*GetAPIStatisticsResponse, *http.Response, error)
	}

	// ProjectService handles communication with project endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project
	ProjectService struct {
		client core.Client
	}

	// Project represents a lakekeeper project
	Project struct {
		ID   string `json:"project-id"`
		Name string `json:"project-name"`
	}

	// ListProjectsResponse represents ListProjects() response.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/list_projects
	ListProjectsResponse struct {
		Projects []*Project `json:"projects"`
	}
	// CreateProjectOptions represents CreateProject() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
	CreateProjectOptions struct {
		ID   *string `json:"project-id,omitempty"` // Request a specific project ID - optional. If not provided, a new project ID will be generated (recommended)
		Name string  `json:"project-name"`
	}

	// createProjectResponse represents the response on project creation.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
	CreateProjectResponse struct {
		ID string `json:"project-id"`
	}

	// RenameProjectOptions represents RenameProject() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/rename_project
	RenameProjectOptions struct {
		NewName string `json:"new-name"`
	}

	// GetAPIStatisticsOptions represents GetAPIStatistics() options
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_endpoint_statistics
	GetAPIStatisticsOptions struct {
		RangeSpecifier *struct {
			// type of the range specifier
			// can be `window` or `page-token`
			Type string `json:"type"`
			// End timestamp of the time window Specify
			// Required if type=window
			End *string `json:"end,omitempty"`
			// 	Duration/span of the time window
			// The returned statistics will be for the time window from end - interval to end.
			// Specify a ISO8601 duration string, e.g. PT1H for 1 hour, P1D for 1 day.
			Interval *string `json:"interval,omitempty"`
			// Opaque Token from previous response for paginating through time windows
			// Use the next_page_token or previous_page_token from a previous response
			// Required if type=page-token
			Token *string `json:"token,omitempty"`
		} `json:"range-specifier,omitempty"`
		StatusCodes []int32 `json:"status-codes,omitempty"`
		Warehouse   struct {
			// Type can be `warehouse-id`, `unmapped` or `all`
			Type string `json:"type"`
			// Required if `Type=warehouse-id`
			ID *string `json:"id,omitempty"`
		} `json:"warehouse"`
	}

	// GetAPIStatisticsResponse represents GetAPIStatistics() response
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_endpoint_statistics
	GetAPIStatisticsResponse struct {
		// Array of arrays of statistics detailing each called endpoint for each timestamp.
		// See docs of timestamps for more details.
		CalledEnpoints [][]struct {
			// Number of requests to this endpoint for the current time-slice.
			Count int64 `json:"count"`
			// Timestamp at which the datapoint was created in the database.
			// This is the exact time at which the current endpoint-status-warehouse combination was called for the first time in the current time-slice.
			CreatedAt string `json:"created-at"`
			// The route of the endpoint.
			// Format: METHOD /path/to/endpoint
			HTTPRoute string `json:"http-route"`
			// The status code of the response.
			StatusCode int32 `json:"status-code"`
			// Timestamp at which the datapoint was last updated.
			// This is the exact time at which the current datapoint was last updated.
			UpdatedAt *string `json:"updated-at,omitempty"`
			// The ID of the warehouse that handled the request.
			// Only present for requests that could be associated with a warehouse.
			// Some management endpoints cannot be associated with a warehouse,
			// e.g. warehouse creation or user management will not have a warehouse-id.
			WarehouseID *string `json:"warehouse-id,omitempty"`
			// The name of the warehouse that handled the request.
			// Only present for requests that could be associated with a warehouse.
			// Some management endpoints cannot be associated with a warehouse,
			// e.g. warehouse creation or user management will not have a warehouse-id
			WarehouseName *string `json:"warehouse-name,omitempty"`
		} `json:"called-endpoints"`
		// Token to get the next page of results.
		// Inverse of PreviousPageToken, see its documentation below.
		NextPageToken string `json:"next-page-token"`
		// Token to get the previous page of results.
		// Endpoint statistics are not paginated through page-limits, we paginate them by stepping through time.
		// By default, the list-statistics endpoint will return all statistics for now() - 1 day to now().
		// In the request, you can specify a range_specifier to set the end date and step interval.
		// The previous-page-token will then move to the neighboring window.
		// E.g. in the default case of now() and 1 day, it'd be now() - 2 days to now() - 1 day.
		PreviousPageToken string `json:"previous-page-token"`
		// Array of timestamps indicating the time at which each entry in the called_endpoints array is valid.
		// We lazily create a new statistics entry every hour, in between hours, the existing entry is being updated.
		// If any endpoint is called in the following hour, there'll be an entry in timestamps for the following hour.
		// If not, then there'll be no entry.
		Timestamps []string `json:"timestamps"`
	}
)

func NewProjectService(client core.Client) ProjectServiceInterface {
	return &ProjectService{
		client: client,
	}
}

// GetProject retrieves information about a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_default_project
func (s *ProjectService) Get(ctx context.Context, id string, options ...core.RequestOptionFunc) (*Project, *http.Response, error) {
	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(ctx, http.MethodGet, "/project", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var prj Project

	resp, apiErr := s.client.Do(req, &prj)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &prj, resp, nil
}

// ListProjects lists all projects that the requesting user has access to.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/list_projects
func (s *ProjectService) List(ctx context.Context, options ...core.RequestOptionFunc) (*ListProjectsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/project-list", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var prjs ListProjectsResponse

	resp, apiErr := s.client.Do(req, &prjs)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &prjs, resp, nil
}

// Create creates a new project with the specified configuration.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
func (s *ProjectService) Create(ctx context.Context, opts *CreateProjectOptions, options ...core.RequestOptionFunc) (*CreateProjectResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/project", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var prjResp CreateProjectResponse

	resp, apiErr := s.client.Do(req, &prjResp)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &prjResp, resp, nil
}

// RenameProject renames a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/rename_project
func (s *ProjectService) Rename(ctx context.Context, id string, opts *RenameProjectOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/project/rename", opts, options)
	if err != nil {
		return nil, err
	}

	r, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return r, apiErr
	}

	return r, nil
}

// DeleteProject delete a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/delete_default_project
func (s *ProjectService) Delete(ctx context.Context, id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(ctx, http.MethodDelete, "/project", nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// GetAPIStatistics retrieves detailed endpoint call statistics for your project, allowing you to monitor API usage patterns,
// track frequency of operations, and analyze response codes.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_endpoint_statistics
func (s *ProjectService) GetAPIStatistics(ctx context.Context, id string, opt *GetAPIStatisticsOptions, options ...core.RequestOptionFunc) (*GetAPIStatisticsResponse, *http.Response, error) {
	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(ctx, http.MethodPost, "/endpoint-statistics", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp GetAPIStatisticsResponse
	r, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, r, apiErr
	}

	return &resp, r, nil
}

func (p *Project) String() string {
	return fmt.Sprintf("Project{ID: %s, Name: %s}", p.ID, p.Name)
}
