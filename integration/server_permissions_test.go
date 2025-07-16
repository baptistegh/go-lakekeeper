//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestPermissions_Server_GetAccess(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.PermissionV1().ServerPermission().GetAccess(nil)
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

	resp, r, err := client.PermissionV1().ServerPermission().GetAssignments(nil)
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
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_Update(t *testing.T) {
	client := Setup(t)

	user, _, err := client.UserV1().Provision(&managementv1.ProvisionUserOptions{
		ID:             core.Ptr("oidc~7b98af91-a814-4498-98cb-2730064db4bc"),
		Email:          core.Ptr("test-user@lakekeeper.io"),
		Name:           core.Ptr("Test User"),
		UpdateIfExists: core.Ptr(true),
		UserType:       core.Ptr(managementv1.HumanUserType),
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)

	t.Cleanup(func() {
		r, err := client.UserV1().Delete(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	resp, _, err := client.PermissionV1().ServerPermission().GetAssignments(nil)
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
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// adding permission
	r, err := client.PermissionV1().ServerPermission().Update(&permissionv1.UpdateServerPermissionsOptions{
		Writes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~7b98af91-a814-4498-98cb-2730064db4bc",
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ServerPermission().GetAssignments(nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// permission added
	want = &permissionv1.GetServerAssignmentsResponse{
		Assignments: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~7b98af91-a814-4498-98cb-2730064db4bc",
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)

	// removing permission
	r, err = client.PermissionV1().ServerPermission().Update(&permissionv1.UpdateServerPermissionsOptions{
		Deletes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~7b98af91-a814-4498-98cb-2730064db4bc",
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ServerPermission().GetAssignments(nil)
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
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_SameAdd(t *testing.T) {
	client := Setup(t)

	user, _, err := client.UserV1().Provision(&managementv1.ProvisionUserOptions{
		ID:             core.Ptr("oidc~fe7b7575-b390-4404-90ce-375421f936bd"),
		Email:          core.Ptr("test-user@exemple.com"),
		Name:           core.Ptr("Test User"),
		UpdateIfExists: core.Ptr(true),
		UserType:       core.Ptr(managementv1.HumanUserType),
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)

	t.Cleanup(func() {
		r, err := client.UserV1().Delete(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	opt := &permissionv1.UpdateServerPermissionsOptions{
		Writes: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~fe7b7575-b390-4404-90ce-375421f936bd",
				},
				Assignment: permissionv1.OperatorServerAssignment,
			},
		},
	}

	// adding permission
	r, err := client.PermissionV1().ServerPermission().Update(opt)

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// adding same permission
	r, err = client.PermissionV1().ServerPermission().Update(opt)

	assert.ErrorContains(t, err, "TupleAlreadyExistsError")
}
