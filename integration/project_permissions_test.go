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

func TestPermissions_Project_GetAccess(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.PermissionV1().ProjectPermission().GetAccess(t.Context(), defaultProjectID, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the project
	want := []permissionv1.ProjectAction{
		permissionv1.CreateWarehouse,
		permissionv1.DeleteProject,
		permissionv1.RenameProject,
		permissionv1.ListWarehouses,
		permissionv1.CreateRole,
		permissionv1.ListRoles,
		permissionv1.SearchRoles,
		permissionv1.ReadProjectAssignments,
		permissionv1.GrantProjectRoleCreator,
		permissionv1.GrantProjectCreate,
		permissionv1.GrantProjectDescribe,
		permissionv1.GrantProjectModify,
		permissionv1.GrantProjectSelet,
		permissionv1.GrantProjectAdmin,
		permissionv1.GrantSecurityAdmin,
		permissionv1.GrantDataAdmin,
		permissionv1.GetProjectEndpointStatistics,
	}

	assert.Subset(t, want, resp.AllowedActions)
}

func TestPermissions_Project_GetAssignments(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), defaultProjectID, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the project
	want := &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
		},
		ProjectID: defaultProjectID,
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Project_Update(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)

	resp, _, err := client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), defaultProjectID, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// initial permissions
	want := &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
		},
		ProjectID: defaultProjectID,
	}

	assert.Equal(t, want, resp)

	// adding permission
	r, err := client.PermissionV1().ProjectPermission().Update(t.Context(), defaultProjectID, &permissionv1.UpdateProjectPermissionsOptions{
		Writes: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.SelectProjectAssignment,
			},
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), defaultProjectID, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// permission added
	want = &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.SelectProjectAssignment,
			},
		},
		ProjectID: defaultProjectID,
	}

	assert.Equal(t, want, resp)

	// removing permission
	r, err = client.PermissionV1().ProjectPermission().Update(t.Context(), defaultProjectID, &permissionv1.UpdateProjectPermissionsOptions{
		Deletes: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.SelectProjectAssignment,
			},
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), defaultProjectID, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// permission deleted
	want = &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
		},
		ProjectID: defaultProjectID,
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Project_SameAdd(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)

	opt := &permissionv1.UpdateProjectPermissionsOptions{
		Writes: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.ModifyProjectAssignment,
			},
		},
	}

	// adding permission
	r, err := client.PermissionV1().ProjectPermission().Update(t.Context(), defaultProjectID, opt)

	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// adding same permission
	r, err = client.PermissionV1().ProjectPermission().Update(t.Context(), defaultProjectID, opt)

	require.ErrorContains(t, err, "TupleAlreadyExistsError")
}

func TestPermissions_Project_Add_NewProject(t *testing.T) {
	client := Setup(t)

	user := MustProvisionUser(t, client)

	projectID := MustCreateProject(t, client)

	resp, r, err := client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), projectID, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// only creator should have assignments on new project
	want := &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
		},
		ProjectID: projectID,
	}

	assert.Equal(t, want, resp)

	opt := &permissionv1.UpdateProjectPermissionsOptions{
		Writes: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.ModifyProjectAssignment,
			},
		},
	}

	// adding permission
	r, err = client.PermissionV1().ProjectPermission().Update(t.Context(), projectID, opt)
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, r, err = client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), projectID, nil)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// we should have the created assignments for the new user
	want = &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: user.ID,
				},
				Assignment: permissionv1.ModifyProjectAssignment,
			},
		},
		ProjectID: projectID,
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Project_Add_Role(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)
	role := MustCreateRole(t, client, project)

	r, err := client.PermissionV1().ProjectPermission().Update(t.Context(), project, &permissionv1.UpdateProjectPermissionsOptions{
		Writes: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.RoleType,
					Value: role.ID,
				},
				Assignment: permissionv1.DescribeProjectAssignment,
			},
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, r, err := client.PermissionV1().ProjectPermission().GetAssignments(t.Context(), project, nil)
	require.NoError(t, err)
	assert.NotNil(t, r)

	want := &permissionv1.GetProjectAssignmentsResponse{
		Assignments: []*permissionv1.ProjectAssignment{
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: adminID,
				},
				Assignment: permissionv1.AdminProjectAssignment,
			},
			{
				Assignee: permissionv1.UserOrRole{
					Type:  permissionv1.RoleType,
					Value: role.ID,
				},
				Assignment: permissionv1.DescribeProjectAssignment,
			},
		},
		ProjectID: project,
	}

	assert.Equal(t, want, resp)
}
