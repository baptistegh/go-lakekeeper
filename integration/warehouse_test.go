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
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarehouse_Create_Default(t *testing.T) {
	t.Parallel()
	client := Setup(t)

	sp := profile.NewS3StorageSettings(
		"testacc",
		"eu-local-1",
		profile.WithPathStyleAccess(),
		profile.WithEndpoint("http://minio:9000/"),
		profile.WithRemoteSigningURLStyle(profile.PathSigningURLStyle),
	).AsProfile()

	sc := credential.NewS3CredentialAccessKey("minio-root-user", "minio-root-password").AsCredential()

	resp, r, err := client.WarehouseV1(defaultProjectID).Create(
		t.Context(),
		&managementv1.CreateWarehouseOptions{
			Name:              "test",
			StorageProfile:    sp,
			StorageCredential: sc,
		},
	)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, r.StatusCode)

	t.Cleanup(func() {
		r, err = client.WarehouseV1(defaultProjectID).Delete(context.Background(), resp.ID, nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	w, r, err := client.WarehouseV1(defaultProjectID).Get(t.Context(), resp.ID)
	require.NoError(t, err)
	assert.NotNil(t, w)

	want := &managementv1.Warehouse{
		ID:             resp.ID,
		Name:           "test",
		ProjectID:      defaultProjectID,
		StorageProfile: sp,
		Status:         managementv1.WarehouseStatusActive,
		DeleteProfile:  profile.NewTabularDeleteProfileHard().AsProfile(),
		Protected:      false,
	}

	assert.Equal(t, want, w)
}

func TestWarehouse_Create_NewProject(t *testing.T) {
	t.Parallel()
	client := Setup(t)

	p, r, err := client.ProjectV1().Create(t.Context(), &managementv1.CreateProjectOptions{
		Name: "test-project",
	})
	require.NoError(t, err)
	assert.NotNil(t, p)

	sp := profile.NewS3StorageSettings(
		"testacc", "eu-local-1",
		profile.WithPathStyleAccess(), profile.WithEndpoint("http://minio:9000/")).AsProfile()

	resp, r, err := client.WarehouseV1(p.ID).Create(
		t.Context(),
		&managementv1.CreateWarehouseOptions{
			Name:              "test",
			StorageProfile:    sp,
			StorageCredential: credential.NewS3CredentialAccessKey("minio-root-user", "minio-root-password").AsCredential(),
		},
	)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, r.StatusCode)

	t.Cleanup(func() {
		r, err = client.WarehouseV1(p.ID).Delete(context.Background(), resp.ID, nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)

		r, err = client.ProjectV1().Delete(context.Background(), p.ID)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})
}

func TestWarehouse_ListSoftDeletedTabulars(t *testing.T) {
	t.Parallel()
	client := Setup(t)

	project := MustCreateProject(t, client)
	warehouseID, _ := MustCreateWarehouse(t, client, project)

	resp, r, err := client.WarehouseV1(project).ListSoftDeletedTabulars(t.Context(), warehouseID, nil)
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.NotNil(t, resp)

	// TODO: add better test
	// 1. Create Table (or view)
	// 2. Enable Soft delete
	// 3. Delete the table
	// 4. List the soft deleted tabulars, we should see an non empty answer
	want := &managementv1.ListSoftDeletedTabularsResponse{
		Tabulars: []*managementv1.Tabular{},
	}

	assert.Equal(t, want, resp)
}

func TestWarehouse_Statistics(t *testing.T) {
	t.Parallel()
	client := Setup(t)

	project := MustCreateProject(t, client)
	warehouseID, _ := MustCreateWarehouse(t, client, project)

	resp, r, err := client.WarehouseV1(project).GetStatistics(t.Context(), warehouseID, nil)
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// It's hard to test against computed values.
	// we can't determine correctly the timestamps.
	// But maybe we can create tables/views and test the correct numbers.
	// TODO: see above
	assert.NotEmpty(t, resp.Stats, resp)
	assert.Equal(t, resp.WarehouseID, warehouseID)
}

func TestWarehouse_SetProtection(t *testing.T) {
	t.Parallel()
	client := Setup(t)

	project := MustCreateProject(t, client)
	warehouseID, _ := MustCreateWarehouse(t, client, project)

	resp, r, err := client.WarehouseV1(project).SetWarehouseProtection(t.Context(), warehouseID, &managementv1.SetProtectionOptions{
		Protected: true,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, true, resp.Protected)
	assert.NotNil(t, resp.UpdatedAt)
}

// TODO: add missing tests
