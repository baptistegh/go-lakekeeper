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
	credentialv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	profilev1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"

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

	c, err := client.NewAuthSourceClient(t.Context(), &as, os.Getenv("LAKEKEEPER_BASE_URL"), client.WithInitialBootstrapV1Enabled(true, true, core.Ptr(managementv1.ApplicationUserType)))
	if err != nil {
		t.Fatalf("could not create client, %v", err)
	}

	return c
}

func MustProvisionUser(t *testing.T, c *client.Client) *managementv1.User {
	id := uuid.New()
	rNb := rand.Int()

	u, _, err := c.UserV1().Provision(t.Context(), &managementv1.ProvisionUserOptions{
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
		if _, err := c.UserV1().Delete(context.Background(), u.ID); err != nil {
			t.Fatalf("could not delete user, %v", err)
		}

	})

	return u
}

func MustCreateRole(t *testing.T, c *client.Client, projectID string) *managementv1.Role {
	rNb := rand.Int()

	r, _, err := c.RoleV1(projectID).Create(t.Context(), &managementv1.CreateRoleOptions{
		Name: fmt.Sprintf("test-role-%d", rNb),
	})
	if err != nil {
		t.Fatalf("could not create role, %v", err)
	}

	t.Cleanup(func() {
		if _, err := c.RoleV1(projectID).Delete(context.Background(), r.ID); err != nil {
			t.Fatalf("could not delete role, %v", err)
		}

	})

	return r
}

func MustCreateProject(t *testing.T, c *client.Client) string {
	p, _, err := c.ProjectV1().Create(t.Context(), &managementv1.CreateProjectOptions{
		Name: fmt.Sprintf("test-project-%d", rand.Int()),
	})
	if err != nil {
		t.Fatalf("could not create user, %v", err)
	}

	t.Cleanup(func() {
		if _, err := c.ProjectV1().Delete(context.Background(), p.ID); err != nil {
			t.Fatalf("could not delete user, %v", err)
		}

	})

	return p.ID
}

func MustCreateWarehouse(t *testing.T, c *client.Client, projectID string) (string, string) {
	rNb := rand.Int()
	name := fmt.Sprintf("test-role-%d", rNb)

	w, _, err := c.WarehouseV1(projectID).Create(t.Context(), &managementv1.CreateWarehouseOptions{
		Name:              name,
		StorageProfile:    profilev1.NewS3StorageSettings("testacc", "eu-local-1").AsProfile(),
		StorageCredential: credentialv1.NewS3CredentialAccessKey("access-key", "secret-key").AsCredential(),
	})
	if err != nil {
		t.Fatalf("could not create warehouse, %v", err)
	}

	t.Cleanup(func() {
		if _, err := c.WarehouseV1(projectID).Delete(context.Background(), w.ID, &managementv1.DeleteWarehouseOptions{Force: core.Ptr(true)}); err != nil {
			t.Fatalf("could not delete warehouse, %v", err)
		}

	})

	return w.ID, name
}
