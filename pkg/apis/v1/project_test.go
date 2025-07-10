package v1_test

import (
	"net/http"
	"testing"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
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

	project, resp, err := client.ProjectV1().Get("01f2fdfc-81fc-444d-8368-5b6701566e35")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &v1.Project{
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

	project, resp, err := client.ProjectV1().List()
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &v1.ListProjectsResponse{
		Projects: []*v1.Project{
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

func TestProjectService_Delete(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/project", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodDelete)
		testutil.TestHeader(t, r, "x-project-id", "01f2fdfc-81fc-444d-8368-5b6701566e35")
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.ProjectV1().Delete("01f2fdfc-81fc-444d-8368-5b6701566e35")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestProjectService_Create(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opts := v1.CreateProjectOptions{
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
	project, resp, err := client.ProjectV1().Create(&opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &v1.CreateProjectResponse{
		ID: "01f2fdfc-81fc-444d-8368-5b6701566e35",
	}

	assert.Equal(t, want, project)
}
