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

package v1_test

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.Project{
		ID:   "01f2fdfc-81fc-444d-8368-5b6701566e35",
		Name: "test-project",
	}

	assert.Equal(t, want, project)
}

func TestProjectService_GetDefault(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/default-project", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", "")
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_project.json")
	})

	project, resp, err := client.ProjectV1().GetDefault(t.Context())
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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

func TestProjectService_RenameDefault(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opts := &managementv1.RenameProjectOptions{
		NewName: "project-renamed",
	}

	mux.HandleFunc("/management/v1/default-project/rename", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("wrong json body")
		}
	})

	resp, err := client.ProjectV1().RenameDefault(t.Context(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestProjectService_Rename(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opts := &managementv1.RenameProjectOptions{
		NewName: "project-renamed",
	}

	mux.HandleFunc("/management/v1/project/rename", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", "01f2fdfc-81fc-444d-8368-5b6701566e35")
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("wrong json body")
		}
	})

	resp, err := client.ProjectV1().Rename(t.Context(), "01f2fdfc-81fc-444d-8368-5b6701566e35", opts)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestProjectService_DeleteDefault(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/default-project", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodDelete)
		testutil.TestHeader(t, r, "x-project-id", "")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.ProjectV1().DeleteDefault(t.Context())
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
