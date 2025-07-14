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
	t.Parallel()
	client := Setup(t)

	resp, r, err := client.PermissionV1().ServerPermissions().GetAccess(nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the default project
	want := &permissionv1.GetServerAccessResponse{
		AllowedActions: []permissionv1.ProjectAction{
			permissionv1.CreateProject,
			permissionv1.UpdateUsers,
			permissionv1.DeleteUsers,
			permissionv1.ListUsers,
			permissionv1.GrantAdmin,
			permissionv1.ProvisionUsers,
			permissionv1.ReadAssignments,
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_GetAssignments(t *testing.T) {
	t.Parallel()
	client := Setup(t)

	resp, r, err := client.PermissionV1().ServerPermissions().GetAssignments(nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the default project
	want := &permissionv1.GetServerAssignmentsResponse{
		Assignments: []*permissionv1.ServerAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: "oidc~6deeb417-cdf9-4320-8a30-ddecea77a4bd",
				},
				Assignment: permissionv1.AdminServerAssignment,
			},
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Server_Update(t *testing.T) {
	t.Parallel()
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

	resp, r, err := client.PermissionV1().ServerPermissions().GetAccess(&permissionv1.GetServerAccessOptions{
		Principal: permissionv1.UserOrRole{
			Type:  permissionv1.UserType,
			Value: user.ID,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	want := &permissionv1.GetServerAccessResponse{
		AllowedActions: []permissionv1.ProjectAction{
			permissionv1.CreateProject,
		},
	}

	assert.Equal(t, want, resp)
}
