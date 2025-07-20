package v1_test

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseService_Get(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID, func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_warehouse.json")
	})

	wh, resp, err := client.WarehouseV1(projectID).Get(t.Context(), warehouseID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	want := &managementv1.Warehouse{
		ID:             warehouseID,
		ProjectID:      projectID,
		Name:           "test-warehouse",
		Protected:      false,
		Status:         managementv1.WarehouseStatusActive,
		StorageProfile: profile.NewS3StorageSettings("test-bucket", "eu-west-1").AsProfile(),
		DeleteProfile:  profile.NewTabularDeleteProfileHard().AsProfile(),
	}

	assert.Equal(t, want, wh)
}

func TestWarehouseService_List(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"

	mux.HandleFunc("/management/v1/warehouse", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/list_warehouses.json")
	})

	warehouses, resp, err := client.WarehouseV1(projectID).List(t.Context(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	want := &managementv1.ListWarehouseResponse{
		Warehouses: []*managementv1.Warehouse{
			{
				ID:             "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d",
				ProjectID:      projectID,
				Name:           "test-warehouse-1",
				Protected:      false,
				Status:         managementv1.WarehouseStatusActive,
				StorageProfile: profile.NewS3StorageSettings("test-bucket-1", "eu-west-1").AsProfile(),
			},
			{
				ID:             "b5c3d2e1-f4a5-6b7c-8d9e-0f1a2b3c4d5e",
				ProjectID:      projectID,
				Name:           "test-warehouse-2",
				Protected:      true,
				Status:         managementv1.WarehouseStatusInactive,
				StorageProfile: profile.NewS3StorageSettings("test-bucket-2", "eu-west-1").AsProfile(),
			},
		},
	}

	assert.Equal(t, want, warehouses)
}

func TestWarehouseService_Create(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	sp := profile.NewS3StorageSettings("test-bucket", "eu-west-1").AsProfile()

	sc := credential.NewS3CredentialAccessKey("test-access-key", "test-secret-key").AsCredential()

	opts := &managementv1.CreateWarehouseOptions{
		Name:              "test-warehouse",
		StorageProfile:    sp,
		StorageCredential: sc,
	}

	mux.HandleFunc("/management/v1/warehouse", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("error wrong body")
		}
		w.WriteHeader(http.StatusCreated)
		testutil.MustWriteHTTPResponse(t, w, "testdata/create_warehouse.json")
	})

	want := &managementv1.CreateWarehouseResponse{
		ID: warehouseID,
	}

	w, resp, err := client.WarehouseV1(projectID).Create(t.Context(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	assert.Equal(t, want, w)
}

func TestWarehouseService_Delete(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID, func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodDelete)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.WarehouseV1(projectID).Delete(t.Context(), warehouseID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestWarehouseService_SetProtection(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestBodyJSON(t, r, &struct {
			Protected bool `json:"protected"`
		}{Protected: true})
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/set_protected.json")
	})

	resp, r, err := client.WarehouseV1(projectID).SetProtection(t.Context(), warehouseID, true)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, true, resp.Protected)
}

func TestWarehouseService_Activate(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/activate", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
	})

	resp, err := client.WarehouseV1(projectID).Activate(t.Context(), warehouseID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_Deactivate(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/deactivate", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
	})

	resp, err := client.WarehouseV1(projectID).Deactivate(t.Context(), warehouseID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_Rename(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opts := &managementv1.RenameWarehouseOptions{
		NewName: "new-name",
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/rename", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).Rename(t.Context(), warehouseID, opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_UpdateStorageProfile(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opts := &managementv1.UpdateStorageProfileOptions{
		StorageCredential: nil,
		StorageProfile:    profile.NewGCSStorageSettings("test-bucket").AsProfile(),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/storage", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).UpdateStorageProfile(t.Context(), warehouseID, opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_UpdateDeleteProfile(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opts := managementv1.UpdateDeleteProfileOptions{
		DeleteProfile: *profile.NewTabularDeleteProfileSoft(3600).AsProfile(),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/delete-profile", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, &opts) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).UpdateDeleteProfile(t.Context(), warehouseID, &opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_UpdateStorageCredential(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opts := managementv1.UpdateStorageCredentialOptions{
		StorageCredential: core.Ptr(credential.NewGCSCredentialSystemIdentity().AsCredential()),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/storage-credential", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, &opts) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).UpdateStorageCredential(t.Context(), warehouseID, &opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
