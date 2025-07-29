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

package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ServerServiceInterface interface {
		Info(ctx context.Context, options ...core.RequestOptionFunc) (*ServerInfo, *http.Response, error)
		Bootstrap(ctx context.Context, opts *BootstrapServerOptions, options ...core.RequestOptionFunc) (*http.Response, error)
	}

	// BootstrapService handles communication with server endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server
	ServerService struct {
		client core.Client
	}

	// ServerInfo represents the servier informations.
	ServerInfo struct {
		AuthzBackend                 string   `json:"authz-backend"`
		Bootstrapped                 bool     `json:"bootstrapped"`
		DefaultProjectID             string   `json:"default-project-id"`
		AWSSystemIdentitiesEnabled   bool     `json:"aws-system-identities-enabled"`
		AzureSystemIdentitiesEnabled bool     `json:"azure-system-identities-enabled"`
		GCPSystemIdentitiesEnabled   bool     `json:"gcp-system-identities-enabled"`
		ServerID                     string   `json:"server-id"`
		Version                      string   `json:"version"`
		Queues                       []string `json:"queues"`
	}

	// BootstrapServerOptions represents the available Bootstrap() options.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/bootstrap
	BootstrapServerOptions struct {
		AcceptTermsOfUse bool      `json:"accept-terms-of-use"`
		IsOperator       *bool     `json:"is-operator,omitempty"`
		UserEmail        *string   `json:"user-email,omitempty"`
		UserName         *string   `json:"user-name,omitempty"`
		UserType         *UserType `json:"user-type,omitempty"`
	}
)

func NewServerService(client core.Client) ServerServiceInterface {
	return &ServerService{
		client: client,
	}
}

func (s *ServerInfo) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// Info returns basic information about the server configuration and status.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/get_server_info
func (s *ServerService) Info(ctx context.Context, options ...core.RequestOptionFunc) (*ServerInfo, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/info", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var info ServerInfo

	resp, apiErr := s.client.Do(req, &info)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &info, resp, nil
}

// Bootstrap initializes the Lakekeeper server and sets the initial administrator account.
// This operation can only be performed once.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/bootstrap
func (s *ServerService) Bootstrap(ctx context.Context, opts *BootstrapServerOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/bootstrap", opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil && apiErr.Type() != "CatalogAlreadyBootstrapped" {
		return nil, apiErr
	}

	return resp, nil
}
