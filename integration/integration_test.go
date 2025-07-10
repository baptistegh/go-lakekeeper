//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/clientcredentials"
)

func Setup(t *testing.T) *client.Client {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file, %v", err)
	}

	oauth := clientcredentials.Config{
		ClientID:     os.Getenv("LAKEKEEPER_CLIENT_ID"),
		ClientSecret: os.Getenv("LAKEKEEPER_CLIENT_SECRET"),
		TokenURL:     os.Getenv("LAKEKEEPER_TOKEN_URL"),
		Scopes:       []string{os.Getenv("LAKEKEEPER_SCOPE")},
	}

	as := core.OAuthTokenSource{
		TokenSource: oauth.TokenSource(context.Background()),
	}

	c, err := client.NewAuthSourceClient(as, os.Getenv("LAKEKEEPER_BASE_URL"), client.WithInitialBootstrapEnabled())
	if err != nil {
		t.Fatalf("could not create client, %v", err)
	}

	return c
}
