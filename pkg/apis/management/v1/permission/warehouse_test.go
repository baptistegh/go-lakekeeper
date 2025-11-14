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

func TestWarehousePermissionService_GetAuthzProperties(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/warehouse/6068343f-7e97-4438-b5c1-866618e3619d", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_warehouse_get_authz_properties.json")
	})

	resp, r, err := client.PermissionV1().WarehousePermission().GetAuthzProperties(t.Context(), "6068343f-7e97-4438-b5c1-866618e3619d")
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	want := &permissionv1.GetWarehouseAuthzPropertiesResponse{
		ManagedAccess: true,
	}

	assert.Equal(t, want, resp)
}

func TestWarehousePermissionService_GetAccess(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/warehouse/62709608-250c-41e0-9457-32bb4de3345c/access", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_warehouse_get_access.json")
	})

	access, resp, err := client.PermissionV1().WarehousePermission().GetAccess(t.Context(), "62709608-250c-41e0-9457-32bb4de3345c", nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetWarehouseAccessResponse{
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

func TestWarehousePermissionService_GetAssignments(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/permissions/warehouse/ed149356-70a0-4a9b-af80-b54b411dae33/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "../testdata/permissions_warehouse_get_assignments.json")
	})

	access, resp, err := client.PermissionV1().WarehousePermission().GetAssignments(t.Context(), "ed149356-70a0-4a9b-af80-b54b411dae33", nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &permissionv1.GetWarehouseAssignmentsResponse{
		Assignments: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.OwnershipWarehouseAssignment,
			},
		},
	}

	assert.Equal(t, want, access)
}

func TestWarehousePermissionService_Update(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opt := &permissionv1.UpdateWarehousePermissionsOptions{
		Deletes: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-1",
				},
				Assignment: permissionv1.ManageGrantsAdminWarehouseAssignment,
			},
		},
		Writes: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~test-user-2",
				},
				Assignment: permissionv1.PassGrantsAdminWarehouseAssignment,
			},
		},
	}

	mux.HandleFunc("/management/v1/permissions/warehouse/6068343f-7e97-4438-b5c1-866618e3619d/assignments", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusNoContent)
		if !testutil.TestBodyJSON(t, r, opt) {
			t.Errorf("invalid request JSON body")
		}
	})

	resp, err := client.PermissionV1().WarehousePermission().Update(t.Context(), "6068343f-7e97-4438-b5c1-866618e3619d", opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
