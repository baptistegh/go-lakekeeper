package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("Default Configuration", func(t *testing.T) {
		t.Parallel()
		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		expectedBaseURL := "http://localhost:8080" + managementv1.APIManagementVersionPath

		if c.BaseURL().String() != expectedBaseURL {
			t.Errorf("NewClient BaseURL is %s, want %s", c.BaseURL().String(), expectedBaseURL)
		}
		if c.UserAgent != userAgent {
			t.Errorf("NewClient UserAgent is %s, want %s", c.UserAgent, userAgent)
		}
	})

	t.Run("Custom UserAgent", func(t *testing.T) {
		t.Parallel()
		c, err := NewClient(t.Context(), "", "http://localhost:8080", WithUserAgent("any-custom-user-agent"))
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		expectedBaseURL := "http://localhost:8080" + managementv1.APIManagementVersionPath

		if c.BaseURL().String() != expectedBaseURL {
			t.Errorf("NewClient BaseURL is %s, want %s", c.BaseURL().String(), expectedBaseURL)
		}
		if c.UserAgent != "any-custom-user-agent" {
			t.Errorf("NewClient UserAgent is %s, want any-custom-user-agent", c.UserAgent)
		}
	})

	t.Run("Invalid Base URL", func(t *testing.T) {
		t.Parallel()
		_, err := NewClient(t.Context(), "", ":invalid:")
		require.Error(t, err)
	})
}

func TestSendingUserAgent_Default(t *testing.T) {
	t.Parallel()

	c, err := NewClient(t.Context(), "", "http://localhost:8080")
	require.NoError(t, err)

	req, err := c.NewRequest(t.Context(), http.MethodGet, "test", nil, nil)
	require.NoError(t, err)

	assert.Equal(t, userAgent, req.Header.Get("User-Agent"))
}

func TestSendingUserAgent_Custom(t *testing.T) {
	t.Parallel()

	c, err := NewClient(t.Context(), "", "http://localhost:8080", WithUserAgent("any-custom-user-agent"))
	require.NoError(t, err)

	req, err := c.NewRequest(t.Context(), http.MethodGet, "test", nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "any-custom-user-agent", req.Header.Get("User-Agent"))
}

func TestCheckResponse(t *testing.T) {
	t.Parallel()
	c, err := NewClient(t.Context(), "", "http://localhost:8181")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := c.NewRequest(t.Context(), http.MethodGet, "test", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp := &http.Response{
		Request:    req.Request,
		StatusCode: http.StatusBadRequest,
		Body: io.NopCloser(strings.NewReader(`
		{
			"error": {
				"code": 3254,
				"message": "This is an error message.",
				"stack":[
					"message 1",
					"message 2"
				],
				"type":"ErrorMessage"
			}
		}`)),
	}

	errResp := CheckResponse(resp)
	if errResp == nil {
		t.Fatal("Expected error response.")
	}

	want := "api error, code=3254 message=This is an error message. type=ErrorMessage"

	if errResp.Error() != want {
		t.Errorf("Expected error: %s, got %s", want, errResp.Error())
	}
}

func TestCheckResponseOnUnknownErrorFormat(t *testing.T) {
	t.Parallel()
	c, err := NewClient(t.Context(), "", "http://localhost:8181")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := c.NewRequest(t.Context(), http.MethodGet, "test", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp := &http.Response{
		Request:    req.Request,
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader("some error message but not JSON")),
	}

	errResp := CheckResponse(resp)
	if errResp == nil {
		t.Fatal("Expected error response.")
	}

	want := "unexpected error response, some error message but not JSON"

	if errResp.Error() != want {
		t.Errorf("Expected error: %s, got %s", want, errResp.Error())
	}
}

func TestRequestWithContext(t *testing.T) {
	t.Parallel()
	c, err := NewClient(t.Context(), "", "http://localhost:8181")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	req, err := c.NewRequest(t.Context(), http.MethodGet, "test", nil, []core.RequestOptionFunc{core.WithContext(ctx)}) //nolint:staticcheck // we let the unit test use the deprecated method for coverage
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	defer cancel()

	if req.Context() != ctx {
		t.Fatal("Context was not set correctly")
	}
}
