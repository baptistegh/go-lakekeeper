package permission_test

import (
	"net/http"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
)

func TestServerPermissionsService_GetAccess(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/server/access", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_server_get_access.json")
	})

	access, resp, err := client.PermissionV1().ServerPermissions().GetAccess(nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetServerAccessResponse{
		AllowedActions: []permissionv1.ProjectAction{
			permissionv1.CreateProject,
			permissionv1.UpdateUsers,
			permissionv1.DeleteUsers,
			permissionv1.ListUsers,
			permissionv1.GrantAdmin,
			permissionv1.ProvisionUsers,
			permissionv1.ReadAssignments,
		},
	}

	assert.Equal(t, want, access)
}

func TestServerPermissionsService_GetAssignments(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/server/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_server_get_assignments.json")
	})

	access, resp, err := client.PermissionV1().ServerPermissions().GetAssignments(nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetServerAssignmentsResponse{
		Assignments: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	}

	assert.Equal(t, want, access)
}

func TestServerPermissionsService_Update(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opt := &permissionv1.UpdateServerPermissionsOptions{
		Deletes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
		Writes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-2",
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	mux.HandleFunc("/management/v1/permissions/server/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Errorf("invalid request JSON body")
		}
	})

	resp, err := client.PermissionV1().ServerPermissions().Update(opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
