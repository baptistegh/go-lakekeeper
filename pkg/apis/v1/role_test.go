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

func TestRoleService_Get(t *testing.T) {
	testCases := []struct {
		name          string
		projectID     string
		roleID        string
		expectedRole  *Role
		expectedError error
		mockClient    *testutil.MockClient
	}{
		{
			name:      "successful get role",
			projectID: "test-project",
			roleID:    "test-role-id",
			expectedRole: &Role{
				ID:          "test-role-id",
				ProjectID:   "test-project",
				Name:        "test-role",
				Description: testutil.StringPtr("test-description"),
				CreatedAt:   "2024-01-01T00:00:00Z",
				UpdatedAt:   testutil.StringPtr("2024-01-01T00:00:00Z"),
			},
			expectedError: nil,
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/role/test-role-id"
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
					role := Role{
						ID:          "test-role-id",
						ProjectID:   "test-project",
						Name:        "test-role",
						Description: testutil.StringPtr("test-description"),
						CreatedAt:   "2024-01-01T00:00:00Z",
						UpdatedAt:   testutil.StringPtr("2024-01-01T00:00:00Z"),
					}
					roleBytes, _ := json.Marshal(role)

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(roleBytes)),
					}

					err := json.NewDecoder(resp.Body).Decode(v)
					return resp, core.ApiErrorFromError(err)
				},
			},
		},
		{
			name:          "failed get role - API error",
			projectID:     "test-project",
			roleID:        "test-role-id",
			expectedRole:  nil,
			expectedError: core.ApiErrorFromMessage("API error"),
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, opts []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					if method != http.MethodGet {
						return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
					}
					expectedURL := "/role/test-role-id"
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
			roleService := NewRoleService(tc.mockClient, tc.projectID)

			role, _, err := roleService.Get(tc.roleID)

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

			if diff := deep.Equal(role, tc.expectedRole); diff != nil {
				t.Error(diff)
			}
		})
	}
}
