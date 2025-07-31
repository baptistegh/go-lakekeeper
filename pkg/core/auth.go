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
	"fmt"
	"os"
	"sync"

	"golang.org/x/oauth2"
)

type (
	// AuthSource is used to obtain access tokens.
	AuthSource interface {
		// Init is called once before making any requests.
		// If the token source needs access to client to initialize itself, it should do so here.
		Init(context.Context) error

		// Header returns an authentication header. When no error is returned, the
		// key and value should never be empty.
		Header(context.Context) (key, value string, err error)

		// GetToken creates a token
		// mainly use to create the Catalog REST API
		GetToken(context.Context) (string, error)
	}

	// OAuthTokenSource wraps an oauth2.TokenSource to implement the AuthSource interface.
	OAuthTokenSource struct {
		TokenSource oauth2.TokenSource
	}

	// AccessTokenAuthSource is an AuthSource that uses a static access token.
	// The token is added to the Authorization header using the Bearer scheme.
	AccessTokenAuthSource struct {
		Token string
	}

	// K8sServiceAccountAuthSource is an AuthSource that retrieves the service account token
	// from the Kubernetes environment. This is typically used in Kubernetes pods where
	// the service account token is mounted at a specific path.
	K8sServiceAccountAuthSource struct {
		// ServiceAccountTokenPath is the path to the service account token file.
		// Default is "/var/run/secrets/kubernetes.io/serviceaccount/token".
		ServiceAccountTokenPath *string

		token  string
		doOnce sync.Once
	}
)

// check the implementations
var (
	_ AuthSource = (*OAuthTokenSource)(nil)
	_ AuthSource = (*AccessTokenAuthSource)(nil)
	_ AuthSource = (*K8sServiceAccountAuthSource)(nil)
)

func (*OAuthTokenSource) Init(context.Context) error {
	return nil
}

func (as *OAuthTokenSource) Header(context.Context) (string, string, error) {
	t, err := as.TokenSource.Token()
	if err != nil {
		return "", "", err
	}

	return "Authorization", fmt.Sprintf("%s %s", t.TokenType, t.AccessToken), nil
}

func (as *OAuthTokenSource) GetToken(_ context.Context) (string, error) {
	t, err := as.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	return t.AccessToken, nil
}

func (*AccessTokenAuthSource) Init(context.Context) error {
	return nil
}

func (as *AccessTokenAuthSource) Header(context.Context) (string, string, error) {
	return "Authorization", "Bearer " + as.Token, nil
}

func (as *AccessTokenAuthSource) GetToken(context.Context) (string, error) {
	return as.Token, nil
}

func (as *K8sServiceAccountAuthSource) Init(context.Context) error {
	// Get service account token
	// This is typically done by reading the token from a file mounted in the pod.
	// For example, the token is usually available at /var/run/secrets/kubernetes.io/serviceaccount/token.
	var err error
	as.doOnce.Do(func() {
		if as.ServiceAccountTokenPath == nil {
			as.ServiceAccountTokenPath = Ptr("/var/run/secrets/kubernetes.io/serviceaccount/token")
		}

		token, e := os.ReadFile(*as.ServiceAccountTokenPath)
		if e != nil {
			err = fmt.Errorf("failed to read service account token: %w", e)
		}

		as.token = string(token)
		if as.token == "" {
			err = fmt.Errorf("service account token is empty, please ensure the file at %s contains a valid token", *as.ServiceAccountTokenPath)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (as *K8sServiceAccountAuthSource) Header(context.Context) (header, value string, err error) {
	return "Authorization", "Bearer " + as.token, nil
}

func (as *K8sServiceAccountAuthSource) GetToken(ctx context.Context) (string, error) {
	if err := as.Init(ctx); err != nil {
		return "", err
	}
	return as.token, nil
}
