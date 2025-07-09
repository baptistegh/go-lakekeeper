package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-retryablehttp"
)

func TestProjectService_Get(t *testing.T) {
	testCases := []struct {
		name            string
		projectID       string
		expectedProject *Project
		expectedError   error
		mockClient      *testutil.MockClient
	}{
		{
			name:      "successful get project",
			projectID: "test-project-id",
			expectedProject: &Project{
				ID:   "test-project-id",
				Name: "test-project",
			},
			expectedError: nil,
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/project"
					if url != expectedURL {
						return nil, fmt.Errorf("expected URL %s, got %s", expectedURL, url)
					}
					if len(opts) != 1 {
						return nil, fmt.Errorf("expected 1 option, got %d", len(opts))
					}
					req, err := retryablehttp.NewRequest(method, url, body)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					project := Project{
						ID:   "test-project-id",
						Name: "test-project",
					}
					projectBytes, _ := json.Marshal(project)
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(projectBytes)),
					}
					err := json.NewDecoder(resp.Body).Decode(v)
					return resp, core.ApiErrorFromError(err)
				},
			},
		},
		{
			name:            "failed get project - API error",
			projectID:       "test-project-id",
			expectedProject: nil,
			expectedError:   core.ApiErrorFromMessage("API error"),
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/project"
					if url != expectedURL {
						return nil, fmt.Errorf("expected URL %s, got %s", expectedURL, url)
					}
					if len(opts) != 1 {
						return nil, fmt.Errorf("expected 1 option, got %d", len(opts))
					}
					req, err := retryablehttp.NewRequest(method, url, body)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					return nil, core.ApiErrorFromMessage("API error")
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectService := NewProjectService(tc.mockClient)

			project, _, err := projectService.Get(tc.projectID)

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if diff := deep.Equal(project, tc.expectedProject); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestProjectService_Default(t *testing.T) {
	testCases := []struct {
		name            string
		expectedProject *Project
		expectedError   error
		mockClient      *testutil.MockClient
	}{
		{
			name: "successful get default project",
			expectedProject: &Project{
				ID:   "default-project-id",
				Name: "default-project",
			},
			expectedError: nil,
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/default-project"
					if url != expectedURL {
						return nil, fmt.Errorf("expected URL %s, got %s", expectedURL, url)
					}
					if len(opts) != 0 {
						return nil, fmt.Errorf("expected 0 options, got %d", len(opts))
					}
					req, err := retryablehttp.NewRequest(method, url, body)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					project := Project{
						ID:   "default-project-id",
						Name: "default-project",
					}
					projectBytes, _ := json.Marshal(project)
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(projectBytes)),
					}
					err := json.NewDecoder(resp.Body).Decode(v)
					return resp, core.ApiErrorFromError(err)
				},
			},
		},
		{
			name:            "failed get default project - API error",
			expectedProject: nil,
			expectedError:   core.ApiErrorFromMessage("API error"),
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/default-project"
					if url != expectedURL {
						return nil, fmt.Errorf("expected URL %s, got %s", expectedURL, url)
					}
					if len(opts) != 0 {
						return nil, fmt.Errorf("expected 0 options, got %d", len(opts))
					}
					req, err := retryablehttp.NewRequest(method, url, body)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					return nil, core.ApiErrorFromMessage("API error")
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectService := NewProjectService(tc.mockClient)

			project, _, err := projectService.Default()

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if diff := deep.Equal(project, tc.expectedProject); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestProjectService_List(t *testing.T) {
	testCases := []struct {
		name             string
		expectedProjects *ListProjectsResponse
		expectedError    error
		mockClient       *testutil.MockClient
	}{
		{
			name: "successful list projects",
			expectedProjects: &ListProjectsResponse{
				Projects: []*Project{
					{
						ID:   "project-1",
						Name: "Project 1",
					},
					{
						ID:   "project-2",
						Name: "Project 2",
					},
				},
			},
			expectedError: nil,
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/project-list"
					if url != expectedURL {
						return nil, fmt.Errorf("expected URL %s, got %s", expectedURL, url)
					}
					if len(opts) != 0 {
						return nil, fmt.Errorf("expected 0 options, got %d", len(opts))
					}
					req, err := retryablehttp.NewRequest(method, url, body)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					projects := ListProjectsResponse{
						Projects: []*Project{
							{
								ID:   "project-1",
								Name: "Project 1",
							},
							{
								ID:   "project-2",
								Name: "Project 2",
							},
						},
					}
					projectsBytes, _ := json.Marshal(projects)
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(projectsBytes)),
					}
					err := json.NewDecoder(resp.Body).Decode(v)
					return resp, core.ApiErrorFromError(err)
				},
			},
		},
		{
			name:             "failed list projects - API error",
			expectedProjects: nil,
			expectedError:    core.ApiErrorFromMessage("API error"),
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/project-list"
					if url != expectedURL {
						return nil, fmt.Errorf("expected URL %s, got %s", expectedURL, url)
					}
					if len(opts) != 0 {
						return nil, fmt.Errorf("expected 0 options, got %d", len(opts))
					}
					req, err := retryablehttp.NewRequest(method, url, body)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					return nil, core.ApiErrorFromMessage("API error")
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectService := NewProjectService(tc.mockClient)

			projects, _, err := projectService.List()

			if tc.expectedError != nil {
				if err == nil || err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if diff := deep.Equal(projects, tc.expectedProjects); diff != nil {
				t.Error(diff)
			}
		})
	}
}
