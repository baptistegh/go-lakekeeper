//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/stretchr/testify/assert"
)

func TestPermissions_Server_GetAccess(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.PermissionV1().ServerPermission().GetAccess(t.Context(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the server
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

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_GetAssignments(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.PermissionV1().ServerPermission().GetAssignments(t.Context(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the server
	want := &permissionv1.GetServerAssignmentsResponse{
		Assignments: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_Update(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)

	resp, _, err := client.PermissionV1().ServerPermission().GetAssignments(t.Context(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// initial permissions
	want := &permissionv1.GetServerAssignmentsResponse{
		Assignments: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// adding permission
	r, err := client.PermissionV1().ServerPermission().Update(t.Context(), &permissionv1.UpdateServerPermissionsOptions{
		Writes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ServerPermission().GetAssignments(t.Context(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// permission added
	want = &permissionv1.GetServerAssignmentsResponse{
		Assignments: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// removing permission
	r, err = client.PermissionV1().ServerPermission().Update(t.Context(), &permissionv1.UpdateServerPermissionsOptions{
		Deletes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ServerPermission().GetAssignments(t.Context(), nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// permission deleted
	want = &permissionv1.GetServerAssignmentsResponse{
		Assignments: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_SameAdd(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)

	opt := &permissionv1.UpdateServerPermissionsOptions{
		Writes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	// adding permission
	r, err := client.PermissionV1().ServerPermission().Update(t.Context(), opt)

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// adding same permission
	r, err = client.PermissionV1().ServerPermission().Update(t.Context(), opt)

	assert.ErrorContains(t, err, "TupleAlreadyExistsError")
}

func TestPermissions_Server_GetAccess_UserFilter(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)

	opt := permissionv1.GetServerAccessOptions{
		PrincipalUser: &user.ID,
	}

	resp, _, err := client.PermissionV1().ServerPermission().GetAccess(t.Context(), &opt)
	assert.NoError(t, err)

	// initial permissions
	want := &permissionv1.GetServerAccessResponse{
		AllowedActions: []permissionv1.ServerAction{},
	}

	assert.Equal(t, want, resp)

	// add user permission
	_, err = client.PermissionV1().ServerPermission().Update(t.Context(), &permissionv1.UpdateServerPermissionsOptions{
		Writes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	})
	assert.NoError(t, err)

	resp, _, err = client.PermissionV1().ServerPermission().GetAccess(t.Context(), &opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// updated permissions
	want = &permissionv1.GetServerAccessResponse{
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

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_GetAccess_RoleFilter(t *testing.T) {
	client := Setup(t)

	role := MustCreateRole(t, client, defaultProjectID)

	opt := permissionv1.GetServerAccessOptions{
		PrincipalRole: &role.ID,
	}

	resp, _, err := client.PermissionV1().ServerPermission().GetAccess(t.Context(), &opt)
	assert.NoError(t, err)

	// initial permissions
	want := &permissionv1.GetServerAccessResponse{
		AllowedActions: []permissionv1.ServerAction{},
	}

	assert.Equal(t, want, resp)

	// add user permission
	_, err = client.PermissionV1().ServerPermission().Update(t.Context(), &permissionv1.UpdateServerPermissionsOptions{
		Writes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.RoleType,
					Value: role.ID,
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	})
	assert.NoError(t, err)

	resp, _, err = client.PermissionV1().ServerPermission().GetAccess(t.Context(), &opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// updated permissions
	want = &permissionv1.GetServerAccessResponse{
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

	assert.Equal(t, want, resp)
}
