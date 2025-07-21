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

func (as *OAuthTokenSource) GetToken(ctx context.Context) (string, error) {
	t, err := as.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	return t.AccessToken, nil
}

func (*AccessTokenAuthSource) Init(context.Context) error {
	return nil
}

func (s *AccessTokenAuthSource) Header(context.Context) (string, string, error) {
	return "Authorization", "Bearer " + s.Token, nil
}

func (as *AccessTokenAuthSource) GetToken(context.Context) (string, error) {
	return as.Token, nil
}

func (s *K8sServiceAccountAuthSource) Init(context.Context) error {
	// Get service account token
	// This is typically done by reading the token from a file mounted in the pod.
	// For example, the token is usually available at /var/run/secrets/kubernetes.io/serviceaccount/token.
	var err error
	s.doOnce.Do(func() {
		if s.ServiceAccountTokenPath == nil {
			s.ServiceAccountTokenPath = Ptr("/var/run/secrets/kubernetes.io/serviceaccount/token")
		}

		token, e := os.ReadFile(*s.ServiceAccountTokenPath)
		if e != nil {
			err = fmt.Errorf("failed to read service account token: %w", e)
		}

		s.token = string(token)
		if s.token == "" {
			err = fmt.Errorf("service account token is empty, please ensure the file at %s contains a valid token", *s.ServiceAccountTokenPath)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *K8sServiceAccountAuthSource) Header(context.Context) (string, string, error) {
	return "Authorization", "Bearer " + s.token, nil
}

func (s *K8sServiceAccountAuthSource) GetToken(ctx context.Context) (string, error) {
	if err := s.Init(ctx); err != nil {
		return "", err
	}
	return s.token, nil
}
