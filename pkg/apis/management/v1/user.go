package v1

import (
	"net/http"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	UserServiceInterface interface {
		// TODO: implement missing endpoints
		// Replace()

		// Creates a new user or updates an existing user's metadata from the provided token.
		// The token should include "profile" and "email" scopes for complete user information.
		Provision(opts *ProvisionUserOptions, options ...core.RequestOptionFunc) (*User, *http.Response, error)
		// Retrieves detailed information about a specific user.
		Get(id string, options ...core.RequestOptionFunc) (*User, *http.Response, error)
		// Returns information about the user associated with the current authentication token.
		Whoami(options ...core.RequestOptionFunc) (*User, *http.Response, error)
		// Permanently removes a user and all their associated permissions.
		// If the user is re-registered later, their permissions will need to be re-added.
		Delete(id string, options ...core.RequestOptionFunc) (*http.Response, error)
		// Performs a fuzzy search for users based on the provided criteria.
		Search(opt *SearchUserOptions, options ...core.RequestOptionFunc) (*SearchUserResponse, *http.Response, error)
		// Returns a paginated list of users based on the provided query parameters.
		List(opt *ListUsersOptions, options ...core.RequestOptionFunc) (*ListUsersResponse, *http.Response, error)
	}

	// UserService handles communication with user endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user
	UserService struct {
		client core.Client
	}
)

var _ UserServiceInterface = (*UserService)(nil)

func NewUserService(client core.Client) UserServiceInterface {
	return &UserService{
		client: client,
	}
}

// User represents a lakekeeper user
type User struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Email           *string  `json:"email,omitempty"`
	UserType        UserType `json:"user-type"`
	CreatedAt       string   `json:"created-at"`
	UpdatedAt       *string  `json:"updated-at,omitempty"`
	LastUpdatedWith string   `json:"last-updated-with"`
}

type UserType string

const (
	HumanUserType       UserType = "human"
	ApplicationUserType UserType = "application"
)

// Get retrieves detailed information about a specific user.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/get_user
func (s *UserService) Get(id string, options ...core.RequestOptionFunc) (*User, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/user/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var user User

	resp, apiErr := s.client.Do(req, &user)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &user, resp, nil
}

// Whoami returns information about the user associated with the current authentication token.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/whoami
func (s *UserService) Whoami(options ...core.RequestOptionFunc) (*User, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/whoami", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var user User

	resp, apiErr := s.client.Do(req, &user)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &user, resp, nil
}

// ProvisionUserOptions represents Provision() options.
//
// The id must be identical to the subject in JWT tokens, prefixed with <idp-identifier>~.
// For example: oidc~1234567890 for OIDC users or kubernetes~1234567890 for Kubernetes users.
// To create users in self-service manner, do not set the id.
// The id is then extracted from the passed JWT token.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/create_user
type ProvisionUserOptions struct {
	ID             *string   `json:"id,omitempty"`
	Email          *string   `json:"email,omitempty"`
	Name           *string   `json:"name,omitempty"`
	UpdateIfExists *bool     `json:"update-if-exists,omitempty"`
	UserType       *UserType `json:"user-type,omitempty"`
}

// Provision creates a new user or updates an existing user's metadata from the provided token.
// The token should include "profile" and "email" scopes for complete user information.
// If opts is provided, the associated user will be created
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/create_user
func (s *UserService) Provision(opts *ProvisionUserOptions, options ...core.RequestOptionFunc) (*User, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/user", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var user User

	resp, apiErr := s.client.Do(req, &user)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &user, resp, nil
}

// Delete permanently removes a user and all their associated permissions.
// If the user is re-registered later, their permissions will need to be re-added.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/delete_user
func (s *UserService) Delete(id string, options ...core.RequestOptionFunc) (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodDelete, "/user/"+id, nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// ListUsersOptions represents options for the List() method.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/list_user
type ListUsersOptions struct {
	Name *string `url:"name,omitempty"`

	ListOptions `url:",inline"` // Embed ListOptions for pagination support
}

// ListUsersResponse represents the response from the List() method.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/list_user
type ListUsersResponse struct {
	Users []*User `json:"users"`

	ListResponse `json:",inline"` // Embed ListResponse for pagination support
}

// List returns a paginated list of users based on the provided query parameters.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/list_user
func (s *UserService) List(opt *ListUsersOptions, options ...core.RequestOptionFunc) (*ListUsersResponse, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/user", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp ListUsersResponse
	httpResp, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, httpResp, apiErr
	}

	return &resp, httpResp, nil
}

// SearchUserOptions represents options for the Search() method.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/search_user
type SearchUserOptions struct {
	Search string `url:"search"`
}

// SearchUserResponse represents the response from the Search() method.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/search_user
type SearchUserResponse struct {
	Users []*User `json:"users"`
}

// Search performs a fuzzy search for users based on the provided criteria.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/search_user
func (s *UserService) Search(opt *SearchUserOptions, options ...core.RequestOptionFunc) (*SearchUserResponse, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/search/user", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var resp SearchUserResponse
	httpResp, apiErr := s.client.Do(req, &resp)
	if apiErr != nil {
		return nil, httpResp, apiErr
	}

	return &resp, httpResp, nil
}
