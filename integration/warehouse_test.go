//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	"github.com/stretchr/testify/assert"
)

func TestWarehouse_Create_Default(t *testing.T) {
	t.Parallel()
	client := Setup(t)

	defaultPrj := "00000000-0000-0000-0000-000000000000"

	sp := profile.NewS3StorageSettings(
		"testacc",
		"eu-local-1",
		profile.WithPathStyleAccess(),
		profile.WithEndpoint("http://minio:9000/"),
	).AsProfile()

	sc := credential.NewS3CredentialAccessKey("minio-root-user", "minio-root-password").AsCredential()

	resp, r, err := client.WarehouseV1(defaultPrj).Create(
		&managementv1.CreateWarehouseOptions{
			Name:              "test",
			StorageProfile:    sp,
			StorageCredential: sc,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, r.StatusCode)

	t.Cleanup(func() {
		r, err = client.WarehouseV1("00000000-0000-0000-0000-000000000000").Delete(resp.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	w, r, err := client.WarehouseV1(defaultPrj).Get(resp.ID)
	assert.NoError(t, err)
	assert.NotNil(t, w)

	want := &managementv1.Warehouse{
		ID:             resp.ID,
		Name:           "test",
		ProjectID:      "00000000-0000-0000-0000-000000000000",
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

	p, r, err := client.ProjectV1().Create(&managementv1.CreateProjectOptions{
		Name: "test-project",
	})
	assert.NoError(t, err)
	assert.NotNil(t, p)

	sp := profile.NewS3StorageSettings(
		"testacc", "eu-local-1",
		profile.WithPathStyleAccess(), profile.WithEndpoint("http://minio:9000/")).AsProfile()

	resp, r, err := client.WarehouseV1(p.ID).Create(
		&managementv1.CreateWarehouseOptions{
			Name:              "test",
			StorageProfile:    sp,
			StorageCredential: credential.NewS3CredentialAccessKey("minio-root-user", "minio-root-password").AsCredential(),
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, r.StatusCode)

	t.Cleanup(func() {
		r, err = client.WarehouseV1(p.ID).Delete(resp.ID, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)

		r, err = client.ProjectV1().Delete(p.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})
}
