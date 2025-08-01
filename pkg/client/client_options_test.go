// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialBootstrapEnabled(t *testing.T) {
	t.Parallel()

	t.Run("Enable Bootstrap Default", func(t *testing.T) {
		t.Parallel()
		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		err = WithInitialBootstrapV1Enabled(true, true, core.Ptr(managementv1.ApplicationUserType))(c)
		require.NoError(t, err)
		assert.True(t, c.bootstrap)
		assert.True(t, c.bootstrapAsOperator)
		assert.Equal(t, managementv1.ApplicationUserType, c.bootstrapUserType)
	})

	t.Run("Enable Bootstrap / No Accept", func(t *testing.T) {
		t.Parallel()
		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		err = WithInitialBootstrapV1Enabled(false, true, core.Ptr(managementv1.ApplicationUserType))(c)
		require.NoError(t, err)
		assert.False(t, c.bootstrap)
		assert.False(t, c.bootstrapAsOperator)
	})
}

func TestCustomHOptions(t *testing.T) {
	t.Parallel()
	t.Run("CustomHTTPClient", func(t *testing.T) {
		t.Parallel()
		// Arrange
		customHTTPClient := &http.Client{}

		// Act
		opt := WithHTTPClient(customHTTPClient)

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to add custom http client, %v", err)
		}

		assert.Equal(t, customHTTPClient, c.client.HTTPClient)
	})

	t.Run("CustomRetryWaitMinMax", func(t *testing.T) {
		t.Parallel()
		// Act
		opt := WithCustomRetryWaitMinMax(30*time.Second, 40*time.Second)

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to set custom retry wait, %v", err)
		}

		assert.Equal(t, 30*time.Second, c.client.RetryWaitMin)
		assert.Equal(t, 40*time.Second, c.client.RetryWaitMax)
	})

	t.Run("WithtoutRetries", func(t *testing.T) {
		t.Parallel()
		// Act
		opt := WithoutRetries()

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to set without retry, %v", err)
		}

		assert.True(t, c.disableRetries)
	})

	t.Run("ErrorHandler", func(t *testing.T) {
		t.Parallel()
		var handler retryablehttp.ErrorHandler = func(_ *http.Response, _ error, _ int) (*http.Response, error) {
			return nil, errors.New("custom error handler")
		}

		opt := WithErrorHandler(handler)

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to set custom error handler, %v", err)
		}

		r, err := c.client.ErrorHandler(nil, nil, 0)
		assert.Nil(t, r)
		require.Error(t, err, "custom error handler")
	})

	t.Run("RequestOptions", func(t *testing.T) {
		t.Parallel()
		opt := WithRequestOptions(core.WithHeader("test", "test"))

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to set custom request option, %v", err)
		}

		assert.Len(t, c.defaultRequestOptions, 1)
	})

	t.Run("CustomBackOff", func(t *testing.T) {
		t.Parallel()
		var custom retryablehttp.Backoff = func(_, _ time.Duration, _ int, _ *http.Response) time.Duration {
			return 67 * time.Second
		}

		opt := WithCustomBackoff(custom)

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to set custom retry, %v", err)
		}

		assert.Equal(t, 67*time.Second, c.client.Backoff(time.Second, time.Second, 0, nil))
	})

	t.Run("CustomRetry", func(t *testing.T) {
		t.Parallel()
		var custom retryablehttp.CheckRetry = func(_ context.Context, _ *http.Response, _ error) (bool, error) {
			return true, errors.New("custom check retry")
		}

		opt := WithCustomRetry(custom)

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to set custom retry, %v", err)
		}

		b, err := c.client.CheckRetry(context.Background(), nil, nil)
		assert.True(t, b)
		require.Error(t, err, "custom check retry")
	})

	t.Run("CustomRetryMax", func(t *testing.T) {
		t.Parallel()
		custom := 258

		opt := WithCustomRetryMax(custom)

		c, err := NewClient(t.Context(), "", "http://localhost:8080")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}

		if err := opt(c); err != nil {
			t.Fatalf("Failed to set custom retry max, %v", err)
		}

		assert.Equal(t, int(258), c.client.RetryMax)
	})
}
