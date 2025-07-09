package core

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type Client interface {
	NewRequest(method, path string, opt any, options []RequestOptionFunc) (*retryablehttp.Request, error)
	Do(req *retryablehttp.Request, v any) (*http.Response, *ApiError)
}
