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

//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissions_Warehouse_GetAuthzProps(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh, _ := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAuthzProperties(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	want := &permissionv1.GetWarehouseAuthzPropertiesResponse{
		ManagedAccess: false,
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Warehouse_SetManagedAccess(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh, _ := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAuthzProperties(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	want := &permissionv1.GetWarehouseAuthzPropertiesResponse{
		ManagedAccess: false,
	}

	// set the managed access to true
	r, err = client.PermissionV1().WarehousePermission().SetManagedAccess(t.Context(), wh, &permissionv1.SetWarehouseManagedAccessOptions{
		ManagedAccess: true,
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	want = &permissionv1.GetWarehouseAuthzPropertiesResponse{
		ManagedAccess: true,
	}

	resp, r, err = client.PermissionV1().WarehousePermission().GetAuthzProperties(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, want, resp)
}

func TestPermissions_Warehouse_GetAccess(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh, _ := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAccess(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the project
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
			permissionv1.GetAllTasks,
			permissionv1.ControlAllTasks,
			permissionv1.SetWarehouseProtection,
			permissionv1.GetWarehouseEndpointStatistics,
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Warehouse_GetAssignments(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh, _ := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAssignments(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the project
	want := &permissionv1.GetWarehouseAssignmentsResponse{
		Assignments: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipWarehouseAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Warehouse_Update(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh, _ := MustCreateWarehouse(t, client, project)

	user := MustProvisionUser(t, client)

	resp, _, err := client.PermissionV1().WarehousePermission().GetAssignments(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// initial permissions
	want := &permissionv1.GetWarehouseAssignmentsResponse{
		Assignments: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipWarehouseAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// adding permission
	r, err := client.PermissionV1().WarehousePermission().Update(t.Context(), wh, &permissionv1.UpdateWarehousePermissionsOptions{
		Writes: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.DescribeWarehouseAssignment,
			},
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().WarehousePermission().GetAssignments(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// permission added
	want = &permissionv1.GetWarehouseAssignmentsResponse{
		Assignments: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipWarehouseAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.DescribeWarehouseAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// removing permission
	r, err = client.PermissionV1().WarehousePermission().Update(t.Context(), wh, &permissionv1.UpdateWarehousePermissionsOptions{
		Deletes: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.DescribeWarehouseAssignment,
			},
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().WarehousePermission().GetAssignments(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// permission deleted
	want = &permissionv1.GetWarehouseAssignmentsResponse{
		Assignments: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipWarehouseAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Warehouse_SameAdd(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)

	project := MustCreateProject(t, client)
	wh, _ := MustCreateWarehouse(t, client, project)

	opt := &permissionv1.UpdateWarehousePermissionsOptions{
		Writes: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.ModifyWarehouseAssignment,
			},
		},
	}

	// adding permission
	r, err := client.PermissionV1().WarehousePermission().Update(t.Context(), wh, opt)

	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// adding same permission
	r, err = client.PermissionV1().WarehousePermission().Update(t.Context(), wh, opt)

	require.ErrorContains(t, err, "TupleAlreadyExistsError")
}

func TestPermissions_Warehouse_Add_Role(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh, _ := MustCreateWarehouse(t, client, project)
	role := MustCreateRole(t, client, project)

	r, err := client.PermissionV1().WarehousePermission().Update(t.Context(), wh, &permissionv1.UpdateWarehousePermissionsOptions{
		Writes: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.RoleType,
					Value: role.ID,
				},
				Assignment: permissionv1.DescribeWarehouseAssignment,
			},
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAssignments(t.Context(), wh, nil)
	require.NoError(t, err)
	assert.NotNil(t, r)

	want := &permissionv1.GetWarehouseAssignmentsResponse{
		Assignments: []*permissionv1.WarehouseAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipWarehouseAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.RoleType,
					Value: role.ID,
				},
				Assignment: permissionv1.DescribeWarehouseAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}
