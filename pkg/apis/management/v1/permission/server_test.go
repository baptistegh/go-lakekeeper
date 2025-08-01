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

package permission_test

import (
	"net/http"
	"testing"

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
