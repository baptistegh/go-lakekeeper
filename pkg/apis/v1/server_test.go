package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-retryablehttp"
)

func TestServerService_Info(t *testing.T) {
	testCases := []struct {
		name            string
		expectedInfo    *ServerInfo
		expectedError   error
		mockClient      *testutil.MockClient
		expectedMethod  string
		expectedURL     string
		expectedOptions []core.RequestOptionFunc
	}{
		{
			name: "successful get server info",
			expectedInfo: &ServerInfo{
				AuthzBackend:                 "opa",
				Bootstrapped:                 true,
				DefaultProjectID:             "default",
				AWSSystemIdentitiesEnabled:   true,
				AzureSystemIdentitiesEnabled: false,
				GCPSystemIdentitiesEnabled:   false,
				ServerID:                     "test-server",
				Version:                      "v0.1.0",
				Queues:                       []string{"queue1", "queue2"},
			},
			expectedError: nil,
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					req, err := retryablehttp.NewRequest(method, url, nil)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					if req.Method != http.MethodGet {
						t.Errorf("expected method %s, got %s", http.MethodGet, req.Method)
					}
					if req.URL.Path != "/info" {
						t.Errorf("expected URL path %s, got %s", "/info", req.URL.Path)
					}

					info := ServerInfo{
						AuthzBackend:                 "opa",
						Bootstrapped:                 true,
						DefaultProjectID:             "default",
						AWSSystemIdentitiesEnabled:   true,
						AzureSystemIdentitiesEnabled: false,
						GCPSystemIdentitiesEnabled:   false,
						ServerID:                     "test-server",
						Version:                      "v0.1.0",
						Queues:                       []string{"queue1", "queue2"},
					}
					infoBytes, _ := json.Marshal(info)

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(infoBytes)),
					}

					err := json.NewDecoder(resp.Body).Decode(v)
					return resp, core.ApiErrorFromError(err)
				},
			},
			expectedMethod:  http.MethodGet,
			expectedURL:     "/info",
			expectedOptions: []core.RequestOptionFunc{},
		},
		{
			name:          "failed get server info - API error",
			expectedInfo:  nil,
			expectedError: core.ApiErrorFromMessage("API error"),
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					req, err := retryablehttp.NewRequest(method, url, nil)
					return req, err
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					return nil, core.ApiErrorFromMessage("API error")
				},
			},
			expectedMethod:  http.MethodGet,
			expectedURL:     "/info",
			expectedOptions: []core.RequestOptionFunc{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serverService := NewServerService(tc.mockClient)

			info, resp, err := serverService.Info(tc.expectedOptions...)

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

			if resp == nil {
				t.Error("expected response, got nil")
				return
			}

			if diff := deep.Equal(info, tc.expectedInfo); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestServerService_Bootstrap(t *testing.T) {
	// Define a dummy BootstrapServerOptions for testing.
	human := HumanUserType
	dummyOpts := &BootstrapServerOptions{
		AcceptTermsOfUse: true,
		IsOperator:       testutil.BoolPtr(true),
		UserEmail:        testutil.StringPtr("test@example.com"),
		UserName:         testutil.StringPtr("Test User"),
		UserType:         &human,
	}

	testCases := []struct {
		name           string
		opts           *BootstrapServerOptions
		expectedError  error
		mockClient     *testutil.MockClient
		expectedMethod string
		expectedURL    string
	}{
		// Add test cases for successful and failed bootstrap scenarios here.
		// Example:
		{
			name:          "successful bootstrap",
			opts:          dummyOpts,
			expectedError: nil,
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					j, err := json.Marshal(body)
					if err != nil {
						return nil, err
					}

					req, err := retryablehttp.NewRequest(method, url, j)
					if err != nil {
						return nil, err
					}

					if diff := deep.Equal(body, dummyOpts); diff != nil {
						t.Error(diff)
					}

					return req, nil
				},
				DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
					return &http.Response{StatusCode: http.StatusOK}, nil
				},
			},
			expectedMethod: http.MethodPost,
			expectedURL:    "/bootstrap",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serverService := NewServerService(tc.mockClient)

			resp, err := serverService.Bootstrap(tc.opts)

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

			if resp == nil {
				t.Error("expected response, got nil")
				return
			}
		})
	}
}
