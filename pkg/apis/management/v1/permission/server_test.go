package permission_test

import (
	"net/http"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
)

func TestServerPermissionService_GetAccess(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/server/access", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_server_get_access.json")
	})

	access, resp, err := client.PermissionV1().ServerPermission().GetAccess(t.Context(), nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetServerAccessResponse{
		AllowedActions: []permissionv1.ServerAction{
			permissionv1.CreateProject,
			permissionv1.UpdateUsers,
			permissionv1.DeleteUsers,
			permissionv1.ListUsers,
			permissionv1.GrantServerAdmin,
			permissionv1.ProvisionUsers,
			permissionv1.ReadAssignments,
		},
	}

	assert.Equal(t, want, access)
}

func TestServerPermissionService_GetAssignments(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/server/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_server_get_assignments.json")
	})

	access, resp, err := client.PermissionV1().ServerPermission().GetAssignments(t.Context(), nil)
	require.NoError(t, err)
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

func TestServerPermissionService_Update(t *testing.T) {
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

	resp, err := client.PermissionV1().ServerPermission().Update(t.Context(), opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestServerPermissionService_GetAllowedAuthorizerActions(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opt := &permissionv1.GetServerAllowedAuthorizerActionsOptions{
		PrincipalUser: core.Ptr("oidc~testuser"),
		PrincipalRole: core.Ptr("testrole"),
	}

	mux.HandleFunc("/management/v1/permissions/server/authorizer-actions", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.TestParam(t, r, "principalUser", "oidc~testuser")
		testutil.TestParam(t, r, "principalRole", "testrole")
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_server_get_authorizer_actions.json")
	})

	access, resp, err := client.PermissionV1().ServerPermission().GetAllowedAuthorizerActions(t.Context(), opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetServerAllowedAuthorizerActionsResponse{
		AllowedActions: []permissionv1.OpenFGAServerAction{
			permissionv1.ServerGrantAdmin,
			permissionv1.ServerReadAssignments,
		},
	}

	assert.Equal(t, want, access)
}
