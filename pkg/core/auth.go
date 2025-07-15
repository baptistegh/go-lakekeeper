package core

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

// AuthSource is used to obtain access tokens.
type AuthSource interface {
	// Init is called once before making any requests.
	// If the token source needs access to client to initialize itself, it should do so here.
	Init(context.Context, Client) error

	// Header returns an authentication header. When no error is returned, the
	// key and value should never be empty.
	Header(ctx context.Context) (key, value string, err error)
}

// OAuthTokenSource wraps an oauth2.TokenSource to implement the AuthSource interface.
type OAuthTokenSource struct {
	TokenSource oauth2.TokenSource
}

func (*OAuthTokenSource) Init(context.Context, Client) error {
	return nil
}

func (as *OAuthTokenSource) Header(_ context.Context) (string, string, error) {
	t, err := as.TokenSource.Token()
	if err != nil {
		return "", "", err
	}

	return "Authorization", fmt.Sprintf("%s %s", t.TokenType, t.AccessToken), nil
}

// AccessTokenAuthSource is an AuthSource that uses a static access token.
// The token is added to the Authorization header using the Bearer scheme.
type AccessTokenAuthSource struct {
	Token string
}

func (*AccessTokenAuthSource) Init(context.Context, Client) error {
	return nil
}

func (s *AccessTokenAuthSource) Header(_ context.Context) (string, string, error) {
	return "Authorization", "Bearer " + s.Token, nil
}

// K8sServiceAccountAuthSource is an AuthSource that retrieves the service account token
// from the Kubernetes environment. This is typically used in Kubernetes pods where
// the service account token is mounted at a specific path.
type K8sServiceAccountAuthSource struct {
	// ServiceAccountTokenPath is the path to the service account token file.
	// Default is "/var/run/secrets/kubernetes.io/serviceaccount/token".
	ServiceAccountTokenPath string

	token string
}

func (s *K8sServiceAccountAuthSource) Init(context.Context, Client) error {
	// Get service account token
	// This is typically done by reading the token from a file mounted in the pod.
	// For example, the token is usually available at /var/run/secrets/kubernetes.io/serviceaccount/token.
	if s.ServiceAccountTokenPath == "" {
		s.ServiceAccountTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}

	token, err := os.ReadFile(s.ServiceAccountTokenPath)
	if err != nil {
		return fmt.Errorf("failed to read service account token: %w", err)
	}

	s.token = string(token)
	if s.token == "" {
		return fmt.Errorf("service account token is empty, please ensure the file at %s contains a valid token", s.ServiceAccountTokenPath)
	}

	return nil
}

func (s *K8sServiceAccountAuthSource) Header(_ context.Context) (string, string, error) {
	return "Authorization", "Bearer " + s.token, nil
}
