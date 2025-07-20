package core

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type Client interface {
	NewRequest(ctx context.Context, method, path string, opt any, options []RequestOptionFunc) (*retryablehttp.Request, error)
	Do(req *retryablehttp.Request, v any) (*http.Response, *ApiError)
}
