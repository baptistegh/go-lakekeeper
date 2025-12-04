package v1

import (
	"context"
	"encoding/json"
	"net/http"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ServerServiceInterface interface {
		Info(ctx context.Context, options ...core.RequestOptionFunc) (*ServerInfo, *http.Response, error)
		Bootstrap(ctx context.Context, opts *BootstrapServerOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		GetAllowedActions(ctx context.Context, opts *GetServerAllowedActionsOptions, options ...core.RequestOptionFunc) (*GetServerAllowedActionsResponse, *http.Response, error)
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

	// GetServerAllowedActionsOptions represents the GetAllowedActions() options.
	//
	// Only one of PrincipalUser or PrincipalRole should be set at a time.
	// Setting both fields simultaneously is not allowed.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/get_server_actions
	GetServerAllowedActionsOptions struct {
		// The user to show actions for.
		PrincipalUser *string `url:"principalUser,omitempty"`
		// The role to show actions for.
		PrincipalRole *string `url:"principalRole,omitempty"`
	}

	// GetServerAllowedActionsResponse represents the GetAllowedActions() response.
	//
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/get_server_actions
	GetServerAllowedActionsResponse struct {
		AllowedActions []permissionv1.ServerAction `json:"allowed-actions"`
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

// GetAllowedActions retrieves the allowed actions for a user or role on the server.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/get_server_actions
func (s *ServerService) GetAllowedActions(ctx context.Context, opt *GetServerAllowedActionsOptions, options ...core.RequestOptionFunc) (*GetServerAllowedActionsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/server/actions", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var response GetServerAllowedActionsResponse
	resp, apiErr := s.client.Do(req, &response)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &response, resp, nil
}
