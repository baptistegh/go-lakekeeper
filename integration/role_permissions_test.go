//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/stretchr/testify/assert"
)

func TestPermissions_Role_GetAccess(t *testing.T) {
	client := Setup(t)

	role := MustCreateRole(t, client, defaultProjectID)

	resp, r, err := client.PermissionV1().RolePermission().GetAccess(role.ID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should be owner on this role
	want := &permissionv1.GetRoleAccessResponse{
		AllowedActions: []permissionv1.RoleAction{
			permissionv1.Assume,
			permissionv1.CanGrantAssignee,
			permissionv1.CanChangeOwnership,
			permissionv1.DeleteRole,
			permissionv1.UpdateRole,
			permissionv1.ReadRole,
			permissionv1.ReadRoleAssignments,
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Role_GetAssignments(t *testing.T) {
	client := Setup(t)

	role := MustCreateRole(t, client, defaultProjectID)

	resp, r, err := client.PermissionV1().RolePermission().GetAssignments(role.ID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should be owner on this role
	want := &permissionv1.GetRoleAssignmentsResponse{
		Assignments: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipRoleAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Role_Update(t *testing.T) {
	client := Setup(t)

	projectID := MustCreateProject(t, client)
	role := MustCreateRole(t, client, projectID)
	user := MustProvisionUser(t, client)

	resp, _, err := client.PermissionV1().RolePermission().GetAssignments(role.ID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// initial permissions
	want := &permissionv1.GetRoleAssignmentsResponse{
		Assignments: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipRoleAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// adding permission
	r, err := client.PermissionV1().RolePermission().Update(role.ID, &permissionv1.UpdateRolePermissionsOptions{
		Writes: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AssigneeRoleAssignment,
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().RolePermission().GetAssignments(role.ID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// permission added
	want = &permissionv1.GetRoleAssignmentsResponse{
		Assignments: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AssigneeRoleAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipRoleAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// removing permission
	r, err = client.PermissionV1().RolePermission().Update(role.ID, &permissionv1.UpdateRolePermissionsOptions{
		Deletes: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AssigneeRoleAssignment,
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().RolePermission().GetAssignments(role.ID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// permission deleted
	want = &permissionv1.GetRoleAssignmentsResponse{
		Assignments: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OwnershipRoleAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Role_SameAdd(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)
	role := MustCreateRole(t, client, defaultProjectID)

	opt := &permissionv1.UpdateRolePermissionsOptions{
		Writes: []*permissionv1.RoleAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AssigneeRoleAssignment,
			},
		},
	}

	// adding permission
	r, err := client.PermissionV1().RolePermission().Update(role.ID, opt)

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// adding same permission
	r, err = client.PermissionV1().RolePermission().Update(role.ID, opt)

	assert.ErrorContains(t, err, "TupleAlreadyExistsError")
}
