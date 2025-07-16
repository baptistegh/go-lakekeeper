//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	adminID = "oidc~6deeb417-cdf9-4320-8a30-ddecea77a4bd"

	defaultProjectID = new(uuid.UUID).String()
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

	c, err := client.NewAuthSourceClient(&as, os.Getenv("LAKEKEEPER_BASE_URL"), client.WithInitialBootstrapEnabled())
	if err != nil {
		t.Fatalf("could not create client, %v", err)
	}

	return c
}

func MustProvisionUser(t *testing.T, c *client.Client) *managementv1.User {
	id := uuid.New()
	rNb := rand.Int()

	u, _, err := c.UserV1().Provision(&managementv1.ProvisionUserOptions{
		ID:             core.Ptr(fmt.Sprintf("oidc~%s", id.String())),
		Name:           core.Ptr(fmt.Sprintf("test-user-%d", rNb)),
		Email:          core.Ptr(fmt.Sprintf("test-user-%d@exemple.com", rNb)),
		UpdateIfExists: core.Ptr(false),
		UserType:       core.Ptr(managementv1.HumanUserType),
	})
	if err != nil {
		t.Fatalf("could not create user, %v", err)
	}

	t.Cleanup(func() {
		if _, err := c.UserV1().Delete(u.ID); err != nil {
			t.Fatalf("could not delete user, %v", err)
		}

	})

	return u
}

func MustCreateRole(t *testing.T, c *client.Client, projectID string) *managementv1.Role {
	rNb := rand.Int()

	r, _, err := c.RoleV1(projectID).Create(&managementv1.CreateRoleOptions{
		Name: fmt.Sprintf("test-role-%d", rNb),
	})
	if err != nil {
		t.Fatalf("could not create role, %v", err)
	}

	t.Cleanup(func() {
		if _, err := c.RoleV1(projectID).Delete(r.ID); err != nil {
			t.Fatalf("could not delete role, %v", err)
		}

	})

	return r
}
