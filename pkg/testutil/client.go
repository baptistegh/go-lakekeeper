package testutil

import (
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/go-retryablehttp"
)

type MockClient struct {
	DoFunc         func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError)
	NewRequestFunc func(method, url string, body any, options []core.RequestOptionFunc) (*retryablehttp.Request, error)
}

func (m *MockClient) Do(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
	return m.DoFunc(req, v)
}

func (m *MockClient) NewRequest(method, url string, body any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
	return m.NewRequestFunc(method, url, body, options)
}

func BoolPtr(b bool) *bool {
	return &b
}

func StringPtr(s string) *string {
	return &s
}
