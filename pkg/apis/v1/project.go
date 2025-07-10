package v1

import (
	"fmt"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ProjectServiceInterface interface {
		List(options ...core.RequestOptionFunc) (*ListProjectsResponse, *http.Response, error)
		Get(id string, options ...core.RequestOptionFunc) (*Project, *http.Response, error)
		Delete(id string, options ...core.RequestOptionFunc) (*http.Response, error)
		Default(options ...core.RequestOptionFunc) (*Project, *http.Response, error)
		Create(opts *CreateProjectOptions, options ...core.RequestOptionFunc) (*CreateProjectResponse, *http.Response, error)
		Rename(id string, opts *RenameProjectOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// ProjectService handles communication with project endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project
	ProjectService struct {
		client core.Client
	}
)

var _ ProjectServiceInterface = (*ProjectService)(nil)

func NewProjectService(client core.Client) ProjectServiceInterface {
	return &ProjectService{
		client: client,
	}
}

// Project represents a lakekeeper project
type Project struct {
	ID   string `json:"project-id"`
	Name string `json:"project-name"`
}

// GetProject retrieves information about a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_default_project
func (s *ProjectService) Get(id string, options ...core.RequestOptionFunc) (*Project, *http.Response, error) {
	options = append(options, WithProject(id))
	req, err := s.client.NewRequest(http.MethodGet, "/project", nil, options)
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

// GetDefaultProject retrieves information about the user's default project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_default_project
//
// Deprecated: This endpoint is deprecated and will be removed in a future version.
func (s *ProjectService) Default(options ...core.RequestOptionFunc) (*Project, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/default-project", nil, options)
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

// ListProjectsResponse represents ListProjects() response.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/list_projects
type ListProjectsResponse struct {
	Projects []*Project `json:"projects"`
}

// ListProjects lists all projects that the requesting user has access to.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/list_projects
func (s *ProjectService) List(options ...core.RequestOptionFunc) (*ListProjectsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/project-list", nil, options)
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

// CreateProjectOptions represents CreateProject() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
type CreateProjectOptions struct {
	ID   *string `json:"project-id,omitempty"` // Request a specific project ID - optional. If not provided, a new project ID will be generated (recommended)
	Name string  `json:"project-name"`
}

// createProjectResponse represents the response on project creation.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
type CreateProjectResponse struct {
	ID string `json:"project-id"`
}

// CreateProject creates a new project with the specified configuration.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
func (s *ProjectService) Create(opts *CreateProjectOptions, options ...core.RequestOptionFunc) (*CreateProjectResponse, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/project", opts, options)
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

// RenameProjectOptions represents RenameProject() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/rename_project
type RenameProjectOptions struct {
	NewName string `json:"new-name"`
}

// RenameProject renames a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/rename_project
func (s *ProjectService) Rename(id string, opts *RenameProjectOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	if opts == nil {
		return nil, fmt.Errorf("RenameProjectOptions cannot be nil")
	}

	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(http.MethodPost, "/project/rename", opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// DeleteProject delete a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/delete_default_project
func (s *ProjectService) Delete(id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(http.MethodDelete, "/project", nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
