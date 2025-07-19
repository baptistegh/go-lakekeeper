//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/stretchr/testify/assert"
)

func TestPermissions_Warehouse_GetAuthzProps(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAuthzProperties(wh, nil)
	assert.NoError(t, err)
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
	wh := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAuthzProperties(wh, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	want := &permissionv1.GetWarehouseAuthzPropertiesResponse{
		ManagedAccess: false,
	}

	// set the managed access to true
	r, err = client.PermissionV1().WarehousePermission().SetManagedAccess(wh, &permissionv1.SetWarehouseManagedAccessOptions{
		ManagedAccess: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	want = &permissionv1.GetWarehouseAuthzPropertiesResponse{
		ManagedAccess: true,
	}

	resp, r, err = client.PermissionV1().WarehousePermission().GetAuthzProperties(wh, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, want, resp)
}

func TestPermissions_Warehouse_GetAccess(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAccess(wh, nil)
	assert.NoError(t, err)
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
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Warehouse_GetAssignments(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh := MustCreateWarehouse(t, client, project)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAssignments(wh, nil)
	assert.NoError(t, err)
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
	wh := MustCreateWarehouse(t, client, project)

	user := MustProvisionUser(t, client)

	resp, _, err := client.PermissionV1().WarehousePermission().GetAssignments(wh, nil)
	assert.NoError(t, err)
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
	r, err := client.PermissionV1().WarehousePermission().Update(wh, &permissionv1.UpdateWarehousePermissionsOptions{
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

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().WarehousePermission().GetAssignments(wh, nil)
	assert.NoError(t, err)
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
	r, err = client.PermissionV1().WarehousePermission().Update(wh, &permissionv1.UpdateWarehousePermissionsOptions{
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

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().WarehousePermission().GetAssignments(wh, nil)
	assert.NoError(t, err)
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
	wh := MustCreateWarehouse(t, client, project)

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
	r, err := client.PermissionV1().WarehousePermission().Update(wh, opt)

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// adding same permission
	r, err = client.PermissionV1().WarehousePermission().Update(wh, opt)

	assert.ErrorContains(t, err, "TupleAlreadyExistsError")
}

func TestPermissions_Warehouse_Add_Role(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	wh := MustCreateWarehouse(t, client, project)
	role := MustCreateRole(t, client, project)

	r, err := client.PermissionV1().WarehousePermission().Update(wh, &permissionv1.UpdateWarehousePermissionsOptions{
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
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, r, err := client.PermissionV1().WarehousePermission().GetAssignments(wh, nil)
	assert.NoError(t, err)
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
