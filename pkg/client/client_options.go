package client

import (
	"net/http"
	"time"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/go-retryablehttp"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
)

// ClientOptionFunc can be used to customize a new Lakekeeper API client.
type ClientOptionFunc func(*Client) error

// WithCustomBackoff can be used to configure a custom backoff policy.
func WithCustomBackoff(backoff retryablehttp.Backoff) ClientOptionFunc {
	return func(c *Client) error {
		c.client.Backoff = backoff
		return nil
	}
}

// WithCustomRetry can be used to configure a custom retry policy.
func WithCustomRetry(checkRetry retryablehttp.CheckRetry) ClientOptionFunc {
	return func(c *Client) error {
		c.client.CheckRetry = checkRetry
		return nil
	}
}

// WithCustomRetryMax can be used to configure a custom maximum number of retries.
func WithCustomRetryMax(retryMax int) ClientOptionFunc {
	return func(c *Client) error {
		c.client.RetryMax = retryMax
		return nil
	}
}

// WithErrorHandler can be used to configure a custom error handler.
func WithErrorHandler(handler retryablehttp.ErrorHandler) ClientOptionFunc {
	return func(c *Client) error {
		c.client.ErrorHandler = handler
		return nil
	}
}

// WithoutRetries disables the default retry logic.
func WithoutRetries() ClientOptionFunc {
	return func(c *Client) error {
		c.disableRetries = true
		return nil
	}
}

// WithCustomRetryWaitMinMax can be used to configure a custom minimum and
// maximum time to wait between retries.
func WithCustomRetryWaitMinMax(waitMin, waitMax time.Duration) ClientOptionFunc {
	return func(c *Client) error {
		c.client.RetryWaitMin = waitMin
		c.client.RetryWaitMax = waitMax
		return nil
	}
}

// WithRequestOptions can be used to configure default request options applied to every request.
func WithRequestOptions(options ...core.RequestOptionFunc) ClientOptionFunc {
	return func(c *Client) error {
		c.defaultRequestOptions = append(c.defaultRequestOptions, options...)
		return nil
	}
}

// WithUserAgent can be used to configure a custom user agent.
func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithHTTPClient can be used to configure a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *Client) error {
		c.client.HTTPClient = httpClient
		return nil
	}
}

// WithInitialBootstrapV1Enabled enables automatic server
// bootstrap on client startup.
//
// acceptTermsOfUse is here to be sure the user is aware.
// if false, this client options will do nothing.
//
// isOperator controls wether the provisioned user will
// have the operator role. default is false.
//
// userType can be human or application, default is application.
func WithInitialBootstrapV1Enabled(
	acceptTermsOfUse bool,
	isOperator bool,
	userType *managementv1.UserType,
) ClientOptionFunc {
	return func(c *Client) error {
		if !acceptTermsOfUse {
			return nil
		}

		c.bootstrapAsOperator = isOperator
		c.bootstrap = true

		if userType != nil {
			c.bootstrapUserType = *userType
		}

		return nil
	}
}
