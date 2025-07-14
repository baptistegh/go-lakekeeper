package v1_test

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRoleService_Get(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	roleID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/role/"+roleID, func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_role.json")
	})

	role, resp, err := client.RoleV1(projectID).Get(roleID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.Role{
		ID:          roleID,
		ProjectID:   projectID,
		Name:        "test-role",
		Description: core.Ptr("description of the role"),
		CreatedAt:   "2019-08-24T14:15:22Z",
		UpdatedAt:   core.Ptr("2019-08-24T14:15:22Z"),
	}

	assert.Equal(t, want, role)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRoleService_Create(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	roleID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	want := &managementv1.Role{
		ID:          roleID,
		ProjectID:   projectID,
		Name:        "test-role",
		Description: core.Ptr("description of the role"),
		CreatedAt:   "2019-08-24T14:15:22Z",
		UpdatedAt:   core.Ptr("2019-08-24T14:15:22Z"),
	}

	opts := managementv1.CreateRoleOptions{
		Name:        "test-role",
		Description: core.Ptr("description of the role"),
	}

	mux.HandleFunc("/management/v1/role", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, &opts) {
			t.Fatalf("wrong json body")
		}
		w.WriteHeader(http.StatusCreated)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_role.json")
	})

	role, resp, err := client.RoleV1(projectID).Create(&opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	assert.Equal(t, want, role)
}

func TestRoleService_Update(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	roleID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	want := &managementv1.Role{
		ID:          roleID,
		ProjectID:   projectID,
		Name:        "test-role",
		Description: core.Ptr("description of the role"),
		CreatedAt:   "2019-08-24T14:15:22Z",
		UpdatedAt:   core.Ptr("2019-08-24T14:15:22Z"),
	}

	opts := managementv1.UpdateRoleOptions{
		Name:        "test-role",
		Description: core.Ptr("description of the role"),
	}

	mux.HandleFunc("/management/v1/role/"+roleID, func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, &opts) {
			t.Fatalf("wrong json body")
		}
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_role.json")
	})

	role, resp, err := client.RoleV1(projectID).Update(roleID, &opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, want, role)
}

func TestRoleService_Delete(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"
	roleID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/role/"+roleID, func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodDelete)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.RoleV1(projectID).Delete(roleID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestRoleService_List(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"

	mux.HandleFunc("/management/v1/role", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		testutil.MustWriteHTTPResponse(t, w, "testdata/list_roles.json")
	})

	nextPage := "8bd02c7f-1d9a-4c5c-afbb-eba7f174da09"
	roleID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	r := &managementv1.Role{
		ID:          roleID,
		ProjectID:   projectID,
		Name:        "test-role",
		Description: core.Ptr("description of the role"),
		CreatedAt:   "2019-08-24T14:15:22Z",
		UpdatedAt:   core.Ptr("2019-08-24T14:15:22Z"),
	}

	want := managementv1.ListRolesResponse{
		NextPageToken: &nextPage,
		Roles:         []*managementv1.Role{r, r},
	}

	roles, resp, err := client.RoleV1(projectID).List(nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, &want, roles)
}

func TestRoleService_Search(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	projectID := "01f2fdfc-81fc-444d-8368-5b6701566e35"

	opts := &managementv1.SearchRoleOptions{
		Search: "test-role",
	}

	mux.HandleFunc("/management/v1/search/role", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		testutil.TestHeader(t, r, "x-project-id", projectID)
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("wrong json body")
		}
		testutil.MustWriteHTTPResponse(t, w, "testdata/search_roles.json")
	})

	roleID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	r := &managementv1.Role{
		ID:          roleID,
		ProjectID:   projectID,
		Name:        "test-role",
		Description: core.Ptr("description of the role"),
		CreatedAt:   "2019-08-24T14:15:22Z",
		UpdatedAt:   core.Ptr("2019-08-24T14:15:22Z"),
	}

	want := managementv1.SearchRoleResponse{
		Roles: []*managementv1.Role{r, r},
	}

	roles, resp, err := client.RoleV1(projectID).Search(opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, &want, roles)
}
