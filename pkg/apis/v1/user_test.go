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

func TestUserService_GetUser(t *testing.T) {
	testCases := []struct {
		name           string
		userID         string
		expectedUser   *User
		expectedError  error
		mockClient     *testutil.MockClient
		expectedMethod string
		expectedURL    string
	}{
		{
			name:          "successful get user",
			userID:        "test-user-id",
			expectedUser:  &User{ID: "test-user-id", Name: "Test User"},
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
					if req.URL.Path != "/user/test-user-id" {
						t.Errorf("expected URL path %s, got %s", "/user/test-user-id", req.URL.Path)
					}

					user := User{ID: "test-user-id", Name: "Test User"}
					userBytes, _ := json.Marshal(user)

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(userBytes)),
					}

					err := json.NewDecoder(resp.Body).Decode(v)
					return resp, core.ApiErrorFromError(err)
				},
			},
			expectedMethod: http.MethodGet,
			expectedURL:    "/user/test-user-id",
		},
		{
			name:           "failed get user - API error",
			userID:         "test-user-id",
			expectedUser:   nil,
			expectedError:  core.ApiErrorFromMessage("API error"),
			expectedMethod: http.MethodGet,
			expectedURL:    "/user/test-user-id",
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					req, err := retryablehttp.NewRequest(method, url, nil)
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
			userService := NewUserService(tc.mockClient)

			user, _, err := userService.Get(tc.userID)

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

			if diff := deep.Equal(user, tc.expectedUser); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestUserService_Whoami(t *testing.T) {
	testCases := []struct {
		name           string
		expectedUser   *User
		expectedError  error
		mockClient     *testutil.MockClient
		expectedMethod string
		expectedURL    string
	}{
		{
			name:          "successful whoami",
			expectedUser:  &User{ID: "current-user-id", Name: "Current User"},
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
					if req.URL.Path != "/whoami" {
						t.Errorf("expected URL path %s, got %s", "/whoami", req.URL.Path)
					}

					user := User{ID: "current-user-id", Name: "Current User"}
					userBytes, _ := json.Marshal(user)

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(userBytes)),
					}

					err := json.NewDecoder(resp.Body).Decode(v)
					return resp, core.ApiErrorFromError(err)
				},
			},
			expectedMethod: http.MethodGet,
			expectedURL:    "/whoami",
		},
		{
			name:           "failed whoami - API error",
			expectedUser:   nil,
			expectedError:  core.ApiErrorFromMessage("API error"),
			expectedMethod: http.MethodGet,
			expectedURL:    "/whoami",
			mockClient: &testutil.MockClient{
				NewRequestFunc: func(method, url string, body any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
					req, err := retryablehttp.NewRequest(method, url, nil)
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
			userService := NewUserService(tc.mockClient)

			user, _, err := userService.Whoami()

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

			if diff := deep.Equal(user, tc.expectedUser); diff != nil {
				t.Error(diff)
			}
		})
	}
}
