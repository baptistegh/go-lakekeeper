package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/v1/permission" // Import the permission package
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/version"
	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
)

var userAgent = fmt.Sprintf("go-lakekeeper/%s", version.Version)

type Client struct {
	// HTTP client used to communicate with the API.
	client *retryablehttp.Client

	// Base URL for API requests. baseURL should always
	// be specified without a trailing slash.
	baseURL *url.URL

	// disableRetries is used to disable the default retry logic.
	disableRetries bool

	// authSource is used to obtain authentication headers.
	authSource core.AuthSource

	// authSourceInit is used to ensure that AuthSources are initialized only
	// once.
	authSourceInit sync.Once

	// Default request options applied to every request.
	defaultRequestOptions []core.RequestOptionFunc

	// User agent used when communicating with the Lakekeeper API.
	UserAgent string

	// bootstrap is used to check if client needs to bootstrap
	// server at startup.
	bootstrap bool

	// bootstrapInit is used to ensure that the bootstrap flow
	// is executed once
	bootstrapInit sync.Once
}

var _ core.Client = (*Client)(nil)

func (c *Client) ServerV1() v1.ServerServiceInterface {
	return v1.NewServerService(c)
}

func (c *Client) ProjectV1() v1.ProjectServiceInterface {
	return v1.NewProjectService(c)
}

func (c *Client) UserV1() v1.UserServiceInterface {
	return v1.NewUserService(c)
}

func (c *Client) RoleV1(projectID string) v1.RoleServiceInterface {
	return v1.NewRoleService(c, projectID)
}

func (c *Client) WarehouseV1(projectID string) v1.WarehouseServiceInterface {
	return v1.NeWarehouseService(c, projectID)
}

func (c *Client) ServerPermissionV1() permission.ServerPermissionServiceInterface {
	return c.ServerV1().ServerPermission()
}

// NewClient returns a new Lakekeeper API client.
// You must provide a valid access token.
func NewClient(token string, baseURL string, options ...ClientOptionFunc) (*Client, error) {
	as := core.AccessTokenAuthSource{Token: token}
	return NewAuthSourceClient(as, baseURL, options...)
}

// NewAuthSourceClient returns a new Lakekeeper API client that uses the AuthSource for authentication.
func NewAuthSourceClient(as core.AuthSource, baseURL string, options ...ClientOptionFunc) (*Client, error) {
	var err error

	c := &Client{
		UserAgent:  userAgent,
		authSource: as,
		bootstrap:  false,
	}

	// Configure the HTTP client.
	c.client = &retryablehttp.Client{
		Backoff:      c.retryHTTPBackoff,
		CheckRetry:   c.retryHTTPCheck,
		ErrorHandler: retryablehttp.PassthroughErrorHandler,
		HTTPClient:   cleanhttp.DefaultPooledClient(),
		RetryWaitMin: 100 * time.Millisecond,
		RetryWaitMax: 400 * time.Millisecond,
		RetryMax:     5,
	}

	// Set the default base URL.
	if err := c.setBaseURL(baseURL); err != nil {
		return nil, err
	}

	// Apply any given client options.
	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(c); err != nil {
			return nil, err
		}
	}

	c.bootstrapInit.Do(func() {
		if !c.bootstrap {
			return
		}

		var info *v1.ServerInfo
		info, _, err = c.ServerV1().Info()
		if err != nil {
			return
		}

		if info != nil && info.Bootstrapped {
			return
		}

		isOperator := true
		userType := v1.ApplicationUserType

		bootstrapOpts := v1.BootstrapServerOptions{
			AcceptTermsOfUse: true,
			IsOperator:       &isOperator,
			UserType:         &userType,
		}
		_, err = c.ServerV1().Bootstrap(&bootstrapOpts)
	})
	if err != nil {
		return nil, fmt.Errorf("error bootstraping the server, %w", err)
	}

	return c, nil
}

// BaseURL return a copy of the baseURL.
func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

// setBaseURL sets the base URL for API requests.
func (c *Client) setBaseURL(urlStr string) error {
	if urlStr == "" {
		return errors.New("base URL must be provided")
	}

	// Make sure the given URL does not end with "/"
	urlStr = strings.TrimSuffix(urlStr, "/")

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(baseURL.Path, v1.ApiManagementVersionPath) {
		baseURL.Path += v1.ApiManagementVersionPath
	}

	// Update the base URL of the client.
	c.baseURL = baseURL

	return nil
}

// NewRequest creates a new API request. The method expects a relative URL
// path that will be resolved relative to the base URL of the Client.
// Relative URL paths should always be specified with a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included
// as the request body.
func (c *Client) NewRequest(method, path string, opt any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
	u := *c.baseURL
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	// Create a request specific headers map.
	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")

	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	var body any
	switch {
	case method == http.MethodPatch || method == http.MethodPost || method == http.MethodPut:
		reqHeaders.Set("Content-Type", "application/json")

		if opt != nil {
			body, err = json.Marshal(opt)
			if err != nil {
				return nil, err
			}
		}
	case opt != nil:
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for _, fn := range append(c.defaultRequestOptions, options...) {
		if fn == nil {
			continue
		}
		if err := fn(req); err != nil {
			return nil, err
		}
	}

	// Set the request specific headers.
	maps.Copy(req.Header, reqHeaders)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
	var err error

	c.authSourceInit.Do(func() {
		err = c.authSource.Init(req.Context(), c)
	})
	if err != nil {
		return nil, core.ApiErrorFromMessage("initializing token source failed:").WithCause(err)
	}

	authKey, authValue, err := c.authSource.Header(req.Context())
	if err != nil {
		return nil, core.ApiErrorFromError(err)
	}

	if v := req.Header.Values(authKey); len(v) == 0 {
		req.Header.Set(authKey, authValue)
	}

	client := c.client

	resp, err := client.Do(req)
	if err != nil {
		return nil, core.ApiErrorFromError(err)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	apiErr := CheckResponse(resp)
	if apiErr != nil {
		// Even though there was an error, we still return the response
		// in case the caller wants to inspect it further.
		return resp, apiErr
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, core.ApiErrorFromError(err)
}

// CheckResponse checks the API response for errors, and returns them if present.
func CheckResponse(r *http.Response) *core.ApiError {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	}

	return core.ApiErrorFromResponse(r)
}

// retryHTTPCheck provides a callback for Client.CheckRetry which
// will retry both rate limit (429) and server (>= 500) errors.
func (c *Client) retryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return false, err
	}
	if !c.disableRetries && (resp.StatusCode == 429 || resp.StatusCode >= 500) {
		return true, nil
	}
	return false, nil
}

// retryHTTPBackoff provides a generic callback for Client.Backoff which
// will pass through all calls based on the status code of the response.
func (c *Client) retryHTTPBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
}
