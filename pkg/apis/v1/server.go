package v1

import (
	"encoding/json"
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/apis/v1/permission" // Import the permission package
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	ServerServiceInterface interface {
		Info(options ...core.RequestOptionFunc) (*ServerInfo, *http.Response, error)
		Bootstrap(opts *BootstrapServerOptions, options ...core.RequestOptionFunc) (*http.Response, error)
		ServerPermission() permission.ServerPermissionServiceInterface // Add this line
	}

	// BootstrapService handles communication with server endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server
	ServerService struct {
		client core.Client
	}
)

var _ ServerServiceInterface = (*ServerService)(nil)

func NewServerService(client core.Client) ServerServiceInterface {
	return &ServerService{
		client: client,
	}
}

func (s *ServerService) ServerPermission() permission.ServerPermissionServiceInterface {
	return permission.NewServerPermissionService(s.client)
}

// ServerInfo represents the servier informations.
type ServerInfo struct {
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
func (s *ServerService) Info(options ...core.RequestOptionFunc) (*ServerInfo, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/info", nil, options)
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

// BootstrapServerOptions represents the available Bootstrap() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/bootstrap
type BootstrapServerOptions struct {
	AcceptTermsOfUse bool      `json:"accept-terms-of-use"`
	IsOperator       *bool     `json:"is-operator,omitempty"`
	UserEmail        *string   `json:"user-email,omitempty"`
	UserName         *string   `json:"user-name,omitempty"`
	UserType         *UserType `json:"user-type,omitempty"`
}

// Bootstrap initializes the Lakekeeper server and sets the initial administrator account.
// This operation can only be performed once.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/bootstrap
func (s *ServerService) Bootstrap(opts *BootstrapServerOptions, options ...core.RequestOptionFunc) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/bootstrap", opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil && apiErr.Type() != "CatalogAlreadyBootstrapped" {
		return nil, apiErr
	}

	return resp, nil
}
