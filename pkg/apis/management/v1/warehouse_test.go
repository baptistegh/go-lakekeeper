package v1_test

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
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
	require.NoError(t, err)
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

	opt := &managementv1.CreateWarehouseOptions{
		Name:              "test-warehouse",
		StorageProfile:    sp,
		StorageCredential: sc,
	}

	mux.HandleFunc("/management/v1/warehouse", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Fatalf("error wrong body")
		}
		w.WriteHeader(http.StatusCreated)
		testutil.MustWriteHTTPResponse(t, w, "testdata/create_warehouse.json")
	})

	want := &managementv1.CreateWarehouseResponse{
		ID: warehouseID,
	}

	w, resp, err := client.WarehouseV1(projectID).Create(t.Context(), opt)
	require.NoError(t, err)
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
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestWarehouseService_Activate(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/activate", func(_ http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
	})

	resp, err := client.WarehouseV1(projectID).Activate(t.Context(), warehouseID)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_Deactivate(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/deactivate", func(_ http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
	})

	resp, err := client.WarehouseV1(projectID).Deactivate(t.Context(), warehouseID)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_Rename(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := &managementv1.RenameWarehouseOptions{
		NewName: "new-name",
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/rename", func(_ http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).Rename(t.Context(), warehouseID, opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_UpdateStorageProfile(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := &managementv1.UpdateStorageProfileOptions{
		StorageCredential: nil,
		StorageProfile:    profile.NewGCSStorageSettings("test-bucket").AsProfile(),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/storage", func(_ http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).UpdateStorageProfile(t.Context(), warehouseID, opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_UpdateDeleteProfile(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := managementv1.UpdateDeleteProfileOptions{
		DeleteProfile: *profile.NewTabularDeleteProfileSoft(3600).AsProfile(),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/delete-profile", func(_ http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, &opt) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).UpdateDeleteProfile(t.Context(), warehouseID, &opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_UpdateStorageCredential(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := managementv1.UpdateStorageCredentialOptions{
		StorageCredential: core.Ptr(credential.NewGCSCredentialSystemIdentity().AsCredential()),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/storage-credential", func(_ http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, &opt) {
			t.Fatalf("error wrong body")
		}
	})

	resp, err := client.WarehouseV1(projectID).UpdateStorageCredential(t.Context(), warehouseID, &opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWarehouseService_ListSoftDeletedTabulars(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := managementv1.ListSoftDeletedTabularsOptions{
		NamespaceID: core.Ptr("namespace_id"),
		ListOptions: managementv1.ListOptions{
			PageToken: core.Ptr("page_token"),
			PageSize:  core.Ptr(int64(250)),
		},
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/deleted-tabulars", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.TestParam(t, r, "namespaceId", "namespace_id")
		testutil.TestParam(t, r, "pageToken", "page_token")
		testutil.TestParam(t, r, "pageSize", "250")
		testutil.MustWriteHTTPResponse(t, w, "testdata/list_soft_deleted_tabulars.json")
	})

	want := &managementv1.ListSoftDeletedTabularsResponse{
		ListResponse: managementv1.ListResponse{
			NextPageToken: core.Ptr("string"),
		},
		Tabulars: []*managementv1.Tabular{
			{
				ID:             "497f6eca-6276-4993-bfeb-53cbbbba6f08",
				Name:           "string",
				Namespace:      []string{"string"},
				Type:           managementv1.TableTabularType,
				WarehouseID:    "019eee1f-0cac-41a0-9932-f7e58ee24619",
				CreatedAt:      "2019-08-24T14:15:22Z",
				DeletedAt:      "2019-08-24T14:15:22Z",
				ExpirationDate: "2019-08-24T14:15:22Z",
			},
		},
	}

	resp, r, err := client.WarehouseV1(projectID).ListSoftDeletedTabulars(t.Context(), warehouseID, &opt)
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, want, resp)
}

func TestWarehouseService_UndropTabular(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := managementv1.UndropTabularOptions{
		Targets: []struct {
			ID   string                   `json:"id"`
			Type managementv1.TabularType `json:"type"`
		}{
			{
				ID:   "test-id",
				Type: managementv1.ViewTabularType,
			},
			{
				ID:   "test-id-2",
				Type: managementv1.TableTabularType,
			},
		},
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/deleted-tabulars/undrop", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.TestBodyJSON(t, r, &opt)
		w.WriteHeader(http.StatusNoContent)
	})

	r, err := client.WarehouseV1(projectID).UndropTabular(t.Context(), warehouseID, &opt)
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)
}

func TestWarehouseService_GetEntityProtection(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"
	entityID := "74f558f9-1443-45f8-9856-fdfb10743d36"

	want := &managementv1.GetProtectionResponse{
		Protected: true,
		UpdatedAt: core.Ptr("2019-08-24T14:15:22Z"),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/namespace/"+entityID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_protection.json")
	})

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/table/"+entityID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_protection.json")
	})

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/view/"+entityID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_protection.json")
	})

	t.Run("Namespace protection", func(t *testing.T) {
		t.Parallel()
		resp, r, err := client.WarehouseV1(projectID).GetNamespaceProtection(t.Context(), warehouseID, entityID)
		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		assert.Equal(t, want, resp)
	})
	t.Run("Table protection", func(t *testing.T) {
		t.Parallel()
		resp, r, err := client.WarehouseV1(projectID).GetTableProtection(t.Context(), warehouseID, entityID)
		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		assert.Equal(t, want, resp)
	})
	t.Run("View protection", func(t *testing.T) {
		t.Parallel()
		resp, r, err := client.WarehouseV1(projectID).GetViewProtection(t.Context(), warehouseID, entityID)
		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		assert.Equal(t, want, resp)
	})
}

func TestWarehouseService_SetEntityProtection(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opt := &managementv1.SetProtectionOptions{
		Protected: true,
	}

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"
	entityID := "74f558f9-1443-45f8-9856-fdfb10743d36"

	want := &managementv1.GetProtectionResponse{
		Protected: true,
		UpdatedAt: core.Ptr("2019-08-24T14:15:22Z"),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Fatal("wrong JSON body")
		}
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_protection.json")
	})

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/namespace/"+entityID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Fatal("wrong JSON body")
		}
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_protection.json")
	})

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/table/"+entityID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Fatal("wrong JSON body")
		}
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_protection.json")
	})

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/view/"+entityID+"/protection", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Fatal("wrong JSON body")
		}
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_protection.json")
	})

	t.Run("Warehouse protection", func(t *testing.T) {
		t.Parallel()
		resp, r, err := client.WarehouseV1(projectID).SetWarehouseProtection(t.Context(), warehouseID, opt)
		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		assert.Equal(t, want, resp)
	})
	t.Run("Namespace protection", func(t *testing.T) {
		t.Parallel()
		resp, r, err := client.WarehouseV1(projectID).SetNamespaceProtection(t.Context(), warehouseID, entityID, opt)
		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		assert.Equal(t, want, resp)
	})
	t.Run("Table protection", func(t *testing.T) {
		t.Parallel()
		resp, r, err := client.WarehouseV1(projectID).SetTableProtection(t.Context(), warehouseID, entityID, opt)
		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		assert.Equal(t, want, resp)
	})
	t.Run("View protection", func(t *testing.T) {
		t.Parallel()
		resp, r, err := client.WarehouseV1(projectID).SetViewProtection(t.Context(), warehouseID, entityID, opt)
		require.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, http.StatusOK, r.StatusCode)

		assert.Equal(t, want, resp)
	})
}

func TestWarehouseService_GetStatistics(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := &managementv1.GetStatisticsOptions{
		PageToken: core.Ptr("page_token"),
		PageSize:  core.Ptr(int64(32)),
	}

	mux.HandleFunc("/management/v1/warehouse/"+warehouseID+"/statistics", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.TestParam(t, r, "page_token", "page_token")
		testutil.TestParam(t, r, "page_size", "32")
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_warehouse_statistics.json")
	})

	want := &managementv1.GetStatisticsResponse{
		ListResponse: managementv1.ListResponse{
			NextPageToken: core.Ptr("string"),
		},
		WarehouseID: "ffa0e747-387e-4f5a-a257-5f6bcf38297d",
		Stats: []struct {
			NumberOfTables int64  `json:"number-of-tables"`
			NumberOfView   int64  `json:"number-of-views"`
			Timestamp      string `json:"timestamp"`
			UpdatedAt      string `json:"updated-at"`
		}{{
			NumberOfTables: int64(12),
			NumberOfView:   int64(8),
			Timestamp:      "2019-08-24T14:15:22Z",
			UpdatedAt:      "2019-08-24T14:15:22Z",
		}},
	}

	resp, r, err := client.WarehouseV1(projectID).GetStatistics(t.Context(), warehouseID, opt)
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, want, resp)
}

func TestWarehouseService_GetAllowedActions(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	warehouseID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	opt := &managementv1.GetWarehouseAllowedActionsOptions{
		PrincipalUser: core.Ptr("oidc~testuser"),
	}

	mux.HandleFunc("/management/v1/warehouse/a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d/actions", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "./testdata/warehouse_get_actions.json")
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.TestParam(t, r, "principalUser", "oidc~testuser")
	})

	access, resp, err := client.WarehouseV1(projectID).GetAllowedActions(t.Context(), warehouseID, opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.GetWarehouseAllowedActionsResponse{
		AllowedActions: []permissionv1.WarehouseAction{
			permissionv1.CreateNamespace,
			permissionv1.DeleteWarehouse,
			permissionv1.ModifyStorage,
			permissionv1.ModifyStorageCredential,
			permissionv1.GetConfig,
			permissionv1.GetMetadata,
			permissionv1.ListNamespaces,
			permissionv1.IncludeInList,
			permissionv1.Deactivate,
			permissionv1.Activate,
			permissionv1.Rename,
			permissionv1.ListDeletedTabulars,
			permissionv1.ReadWarehouseAssignments,
			permissionv1.GrantCreate,
			permissionv1.GrantDescribe,
			permissionv1.GrantModify,
			permissionv1.GrantSelect,
			permissionv1.GrantPassGrants,
			permissionv1.GrantManageGrants,
			permissionv1.ChangeOwnership,
			permissionv1.SetWarehouseProtection,
			permissionv1.GetWarehouseEndpointStatistics,
		},
	}

	assert.Equal(t, want, access)
}
