package permission_test

import (
	"net/http"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
)

func TestProjectPermissionService_GetAccess(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/project/62709608-250c-41e0-9457-32bb4de3345c/access", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_project_get_access.json")
	})

	access, resp, err := client.PermissionV1().ProjectPermission().GetAccess(t.Context(), "62709608-250c-41e0-9457-32bb4de3345c", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetProjectAccessResponse{
		AllowedActions: []permissionv1.ProjectAction{
			permissionv1.CreateWarehouse,
			permissionv1.DeleteProject,
			permissionv1.RenameProject,
			permissionv1.ListWarehouses,
			permissionv1.CreateRole,
			permissionv1.ListRoles,
			permissionv1.SearchRoles,
			permissionv1.ReadProjectAssignments,
			permissionv1.GrantProjectRoleCreator,
			permissionv1.GrantProjectCreate,
			permissionv1.GrantProjectDescribe,
			permissionv1.GrantProjectModify,
			permissionv1.GrantProjectSelet,
			permissionv1.GrantProjectAdmin,
			permissionv1.GrantSecurityAdmin,
			permissionv1.GrantDataAdmin,
		},
	}

	assert.Equal(t, want, access)
}

func TestProjectPermissionService_GetAssignments(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/project/ed149356-70a0-4a9b-af80-b54b411dae33/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_project_get_assignments.json")
	})

	access, resp, err := client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), "ed149356-70a0-4a9b-af80-b54b411dae33", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
		},
	}

	assert.Equal(t, want, access)
}

func TestProjectPermissionService_Update(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opt := &permissionv1.UpdateProjectPermissionsOptions{
		Deletes: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.SecurityAdminProjectAssignment,
			},
		},
		Writes: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-2",
				},
				Assignment: permissionv1.ModifyProjectAssignment,
			},
		},
	}

	mux.HandleFunc("/management/v1/permissions/project/6068343f-7e97-4438-b5c1-866618e3619d/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Errorf("invalid request JSON body")
		}
	})

	resp, err := client.PermissionV1().ProjectPermission().Update(t.Context(), "6068343f-7e97-4438-b5c1-866618e3619d", opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
