package permission

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http" // Added for http.Response

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/google/go-querystring/query"
)

// ServerAssignment represents an assignment a role or a user can
// have to the server
//
// Assignee can be a role or a user
// Assignement can be Operator or Admin
type ServerAssignment struct {
	Assignee   UserOrRole
	Assignment ServerAssignmentType
}

// to be sure ServerAssignment can be JSON encoded/decoded
var (
	_ json.Unmarshaler = (*ServerAssignment)(nil)
	_ json.Marshaler   = (*ServerAssignment)(nil)
)

type (
	UserOrRoleType string

	ServerAssignmentType string
)

const (
	UserType UserOrRoleType = "user"
	RoleType UserOrRoleType = "role"

	OperatorServerAssignment ServerAssignmentType = "operator"
	AdminServerAssignment    ServerAssignmentType = "admin"
)

// ServerAction describes the actions that can be performed on a server
type ServerAction string

const (
	CreateProjectServerAction   ServerAction = "create_project"
	UpdateUsersServerAction     ServerAction = "update_users"
	DeleteUsersServerAction     ServerAction = "delete_users"
	ListUsersServerAction       ServerAction = "list_users"
	GrantAdminServerAction      ServerAction = "grant_admin"
	ProvisionUsersServerAction  ServerAction = "provision_users"
	ReadAssignmentsServerAction ServerAction = "read_assignments"
)

// ServerRelation describes the relations that can be applied to a server
type ServerRelation string

const (
	AdminServerRelation    ServerRelation = "admin"
	OperatorServerRelation ServerRelation = "operator"
)

// GetServerAccessResponse contains the allowed actions for a server
type GetServerAccessResponse struct {
	AllowedActions []ServerAction `json:"allowed-actions"`
}

// GetServerAssignmentsResponse contains the server assignments
type GetServerAssignmentsResponse struct {
	Assignments []ServerAssignment `json:"assignments"`
}

// UpdateServerAssignmentsRequest is the request for updating server assignments
type UpdateServerAssignmentsRequest struct {
	Deletes []ServerAssignment `json:"deletes,omitempty"`
	Writes  []ServerAssignment `json:"writes,omitempty"`
}

// ServerPermissionServiceInterface defines the interface for server permission operations
type ServerPermissionServiceInterface interface {
	GetServerAccess(ctx context.Context, principal *UserOrRole) (*GetServerAccessResponse, *http.Response, *core.ApiError)
	GetServerAssignments(ctx context.Context, relations []ServerRelation) (*GetServerAssignmentsResponse, *http.Response, *core.ApiError)
	UpdateServerAssignments(ctx context.Context, body *UpdateServerAssignmentsRequest) (*http.Response, *core.ApiError)
}

// serverPermissionService implements ServerPermissionServiceInterface
type serverPermissionService struct {
	client core.Client
}

// NewServerPermissionService creates a new ServerPermissionService
func NewServerPermissionService(client core.Client) ServerPermissionServiceInterface {
	return &serverPermissionService{
		client: client,
	}
}

// GetServerAccess gets the allowed actions for a server
func (s *serverPermissionService) GetServerAccess(ctx context.Context, principal *UserOrRole) (*GetServerAccessResponse, *http.Response, *core.ApiError) {
	path := "/permissions/server/access"
	req, err := s.client.NewRequest(http.MethodGet, path, principal, nil)
	if err != nil {
		return nil, nil, core.ApiErrorFromError(err)
	}

	resp := new(GetServerAccessResponse)
	httpResp, apiErr := s.client.Do(req, resp)
	if apiErr != nil {
		return nil, httpResp, apiErr
	}

	return resp, httpResp, nil
}

// GetServerAssignments gets user and role assignments of the server
func (s *serverPermissionService) GetServerAssignments(ctx context.Context, relations []ServerRelation) (*GetServerAssignmentsResponse, *http.Response, *core.ApiError) {
	path := "/permissions/server/assignments"
	opts := struct {
		Relations []ServerRelation `url:"relations,omitempty"`
	}{
		Relations: relations,
	}
	q, err := query.Values(opts)
	if err != nil {
		return nil, nil, core.ApiErrorFromError(err)
	}
	req, err := s.client.NewRequest(http.MethodGet, path, q, nil)
	if err != nil {
		return nil, nil, core.ApiErrorFromError(err)
	}

	resp := new(GetServerAssignmentsResponse)
	httpResp, apiErr := s.client.Do(req, resp)
	if apiErr != nil {
		return nil, httpResp, apiErr
	}

	return resp, httpResp, nil
}

// UpdateServerAssignments updates permissions for this server
func (s *serverPermissionService) UpdateServerAssignments(ctx context.Context, body *UpdateServerAssignmentsRequest) (*http.Response, *core.ApiError) {
	path := "/permissions/server/assignments"
	req, err := s.client.NewRequest(http.MethodPost, path, body, nil)
	if err != nil {
		return nil, core.ApiErrorFromError(err)
	}

	httpResp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return httpResp, apiErr
	}

	return httpResp, nil
}

type UserOrRole struct {
	Type  UserOrRoleType
	Value string
}

func (sa *ServerAssignment) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Type ServerAssignmentType `json:"type"`
		Role *string              `json:"role,omitempty"`
		User *string              `json:"user,omitempty"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	sa.Assignment = aux.Type

	if aux.Role == nil && aux.User == nil {
		return fmt.Errorf("error reading server assignment, role or user must be provided")
	}

	if aux.Role != nil && aux.User != nil {
		return fmt.Errorf("error reading server assignment, role and user can't be both provided")
	}

	if aux.Role != nil {
		sa.Assignee = UserOrRole{
			RoleType,
			*aux.Role,
		}
		return nil
	}

	if aux.User != nil {
		sa.Assignee = UserOrRole{
			UserType,
			*aux.User,
		}
		return nil
	}
	return fmt.Errorf("incorrect server assignment")
}

func (sa ServerAssignment) MarshalJSON() ([]byte, error) {
	aux := make(map[string]string)

	switch sa.Assignee.Type {
	case RoleType:
		aux["role"] = sa.Assignee.Value
	case UserType:
		aux["user"] = sa.Assignee.Value
	}

	aux["type"] = string(sa.Assignment)

	return json.Marshal(aux)
}
