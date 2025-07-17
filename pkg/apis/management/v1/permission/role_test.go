package permission_test

import (
	"net/http"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
)

func TestRolePermissionService_GetAssignments(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/role/ed149356-70a0-4a9b-af80-b54b411dae33/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_role_get_assignments.json")
	})

	access, resp, err := client.PermissionV1().RolePermission().GetAssignments("ed149356-70a0-4a9b-af80-b54b411dae33", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetRoleAssignmentsResponse{
		Assignments: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.OwnershipRoleAssignment,
			},
		},
	}

	assert.Equal(t, want, access)
}

func TestRolePermissionService_Update(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opt := &permissionv1.UpdateRolePermissionsOptions{
		Deletes: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.OwnershipRoleAssignment,
			},
		},
		Writes: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-2",
				},
				Assignment: permissionv1.AssigneeRoleAssignment,
			},
		},
	}

	mux.HandleFunc("/management/v1/permissions/role/6068343f-7e97-4438-b5c1-866618e3619d/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Errorf("invalid request JSON body")
		}
	})

	resp, err := client.PermissionV1().RolePermission().Update("6068343f-7e97-4438-b5c1-866618e3619d", opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
