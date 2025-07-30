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

package core

import (
	"context"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// contextKey is key of context used internal client-go
type contextKey struct{}

// checkRetryKey is context key of requestRetry.
// Value type of this key must be `retryablehttp.CheckRetry`
// This is used in [WithRequestRetry].
var checkRetryKey = &contextKey{}

// checkRetryFromContext returns checkRetry from Context.
// If checkRetry doesn't exist in context, return nil
func checkRetryFromContext(ctx context.Context) retryablehttp.CheckRetry {
	val := ctx.Value(checkRetryKey)

	// There is no checkRetry in context
	if val == nil {
		return nil
	}

	return val.(retryablehttp.CheckRetry)
}

// contextWithCheckRetry create and return new context with checkRetry
func contextWithCheckRetry(ctx context.Context, checkRetry retryablehttp.CheckRetry) context.Context {
	return context.WithValue(ctx, checkRetryKey, checkRetry)
}
