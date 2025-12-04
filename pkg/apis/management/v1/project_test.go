package v1_test

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectService_Get(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/project", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", "01f2fdfc-81fc-444d-8368-5b6701566e35")
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_project.json")
	})

	project, resp, err := client.ProjectV1().Get(t.Context(), "01f2fdfc-81fc-444d-8368-5b6701566e35")
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.Project{
		ID:   "01f2fdfc-81fc-444d-8368-5b6701566e35",
		Name: "test-project",
	}

	assert.Equal(t, want, project)
}

func TestProjectService_List(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/project-list", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "testdata/list_projects.json")
	})

	project, resp, err := client.ProjectV1().List(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.ListProjectsResponse{
		Projects: []*managementv1.Project{
			{
				ID:   "01f2fdfc-81fc-444d-8368-5b6701566e35",
				Name: "test-project-1",
			},
			{
				ID:   "f80ed5b3-2e5b-49df-a7a2-5f071f91e6dd",
				Name: "test-project-2",
			},
		},
	}

	assert.Equal(t, want, project)
}

func TestProjectService_Rename(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opts := &managementv1.RenameProjectOptions{
		NewName: "project-renamed",
	}

	mux.HandleFunc("/management/v1/project/rename", func(_ http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", "01f2fdfc-81fc-444d-8368-5b6701566e35")
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("wrong json body")
		}
	})

	resp, err := client.ProjectV1().Rename(t.Context(), "01f2fdfc-81fc-444d-8368-5b6701566e35", opts)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestProjectService_Delete(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/project", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodDelete)
		testutil.TestHeader(t, r, "x-project-id", "01f2fdfc-81fc-444d-8368-5b6701566e35")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.ProjectV1().Delete(t.Context(), "01f2fdfc-81fc-444d-8368-5b6701566e35")
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestProjectService_Create(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opts := managementv1.CreateProjectOptions{
		Name: "test-project",
	}

	mux.HandleFunc("/management/v1/project", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		if !testutil.TestBodyJSON(t, r, &opts) {
			t.Fatalf("wrong json body")
		}
		w.WriteHeader(http.StatusCreated)
		testutil.MustWriteHTTPResponse(t, w, "testdata/create_project.json")
	})
	project, resp, err := client.ProjectV1().Create(t.Context(), &opts)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.CreateProjectResponse{
		ID: "01f2fdfc-81fc-444d-8368-5b6701566e35",
	}

	assert.Equal(t, want, project)
}

func TestProjectService_GetAPIStatistics(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opts := managementv1.GetAPIStatisticsOptions{
		Warehouse: struct {
			Type string  "json:\"type\""
			ID   *string "json:\"id,omitempty\""
		}{
			Type: "all",
		},
	}

	mux.HandleFunc("/management/v1/endpoint-statistics", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", "01f2fdfc-81fc-444d-8368-5b6701566e35")
		if !testutil.TestBodyJSON(t, r, &opts) {
			t.Fatalf("wrong json body")
		}
		w.WriteHeader(http.StatusCreated)
		testutil.MustWriteHTTPResponse(t, w, "testdata/project_get_api_statistics.json")
	})
	project, resp, err := client.ProjectV1().GetAPIStatistics(t.Context(), "01f2fdfc-81fc-444d-8368-5b6701566e35", &opts)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.GetAPIStatisticsResponse{
		CalledEnpoints: [][]struct {
			Count         int64   `json:"count"`
			CreatedAt     string  `json:"created-at"`
			HTTPRoute     string  `json:"http-route"`
			StatusCode    int32   `json:"status-code"`
			UpdatedAt     *string `json:"updated-at,omitempty"`
			WarehouseID   *string `json:"warehouse-id,omitempty"`
			WarehouseName *string `json:"warehouse-name,omitempty"`
		}{
			{
				{
					Count:         0,
					CreatedAt:     "2019-08-24T14:15:22Z",
					HTTPRoute:     "string",
					StatusCode:    0,
					UpdatedAt:     core.Ptr("2019-08-24T14:15:22Z"),
					WarehouseID:   core.Ptr("019eee1f-0cac-41a0-9932-f7e58ee24619"),
					WarehouseName: core.Ptr("string"),
				},
			},
		},
		NextPageToken:     "string",
		PreviousPageToken: "string",
		Timestamps:        []string{"2019-08-24T14:15:22Z"},
	}

	assert.Equal(t, want, project)
}

func TestProjectService_GetAllowedActions(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"

	opt := &managementv1.GetProjectAllowedActionsOptions{
		PrincipalUser: core.Ptr("oidc~testuser"),
		PrincipalRole: core.Ptr("testrole"),
	}

	mux.HandleFunc("/management/v1/project/actions", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.TestParam(t, r, "principalUser", "oidc~testuser")
		testutil.TestParam(t, r, "principalRole", "testrole")
		testutil.MustWriteHTTPResponse(t, w, "./testdata/project_get_actions.json")
	})

	access, resp, err := client.ProjectV1().GetAllowedActions(t.Context(), projectID, opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.GetProjectAllowedActionsResponse{
		AllowedActions: []permissionv1.ProjectAction{
			permissionv1.CreateWarehouse,
			permissionv1.DeleteProject,
			permissionv1.RenameProject,
			permissionv1.ProjectGetMetadata,
			permissionv1.ListWarehouses,
			permissionv1.ProjectIncludeInList,
			permissionv1.CreateRole,
			permissionv1.ListRoles,
			permissionv1.SearchRoles,
			permissionv1.GetProjectEndpointStatistics,
		},
	}

	assert.Equal(t, want, access)
}
