//go:build integration
// +build integration

package integration

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestPermissions_Project_GetAccess(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.PermissionV1().ProjectPermission().GetAccess("00000000-0000-0000-0000-000000000000", nil)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	// User should have all permissions on the project
	want := &permissionv1.GetProjectAccessResponse{
		AllowedActions: []permissionv1.ProjectAction{
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
		},
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Project_GetAssignments(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.PermissionV1().ProjectPermission().GetAssignments("00000000-0000-0000-0000-000000000000", nil)
	assert.NoError(t, err)
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
		ProjectID: "00000000-0000-0000-0000-000000000000",
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Project_Update(t *testing.T) {
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

	resp, _, err := client.PermissionV1().ProjectPermission().GetAssignments("00000000-0000-0000-0000-000000000000", nil)
	assert.NoError(t, err)
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
		ProjectID: "00000000-0000-0000-0000-000000000000",
	}

	assert.Equal(t, want, resp)

	// adding permission
	r, err := client.PermissionV1().ProjectPermission().Update("00000000-0000-0000-0000-000000000000", &permissionv1.UpdateProjectPermissionsOptions{
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

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ProjectPermission().GetAssignments("00000000-0000-0000-0000-000000000000", nil)
	assert.NoError(t, err)
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
		ProjectID: "00000000-0000-0000-0000-000000000000",
	}

	assert.Equal(t, want, resp)

	// removing permission
	r, err = client.PermissionV1().ProjectPermission().Update("00000000-0000-0000-0000-000000000000", &permissionv1.UpdateProjectPermissionsOptions{
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

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, _, err = client.PermissionV1().ProjectPermission().GetAssignments("00000000-0000-0000-0000-000000000000", nil)
	assert.NoError(t, err)
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
		ProjectID: "00000000-0000-0000-0000-000000000000",
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Project_SameAdd(t *testing.T) {
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
	r, err := client.PermissionV1().ProjectPermission().Update("00000000-0000-0000-0000-000000000000", opt)

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// adding same permission
	r, err = client.PermissionV1().ProjectPermission().Update("00000000-0000-0000-0000-000000000000", opt)

	assert.ErrorContains(t, err, "TupleAlreadyExistsError")
}

func TestPermissions_Project_Add_NewProject(t *testing.T) {
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

	project, _, err := client.ProjectV1().Create(&managementv1.CreateProjectOptions{
		Name: fmt.Sprintf("test-project-%d", rand.Int()),
	})
	assert.NoError(t, err)
	assert.NotNil(t, project)

	t.Cleanup(func() {
		r, err := client.UserV1().Delete(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)

		r, err = client.ProjectV1().Delete(project.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	resp, r, err := client.PermissionV1().ProjectPermission().GetAssignments(project.ID, nil)
	assert.NoError(t, err)
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
		ProjectID: project.ID,
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
	r, err = client.PermissionV1().ProjectPermission().Update(project.ID, opt)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, r, err = client.PermissionV1().ProjectPermission().GetAssignments(project.ID, nil)
	assert.NoError(t, err)
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
		ProjectID: project.ID,
	}

	assert.Equal(t, want, resp)
}

func TestPermissions_Project_Add_Role(t *testing.T) {
	client := Setup(t)

	role, _, err := client.RoleV1("00000000-0000-0000-0000-000000000000").Create(&managementv1.CreateRoleOptions{
		Name: fmt.Sprintf("test-role-%d", rand.Int()),
	})
	assert.NoError(t, err)
	assert.NotNil(t, role)

	t.Cleanup(func() {
		r, err := client.RoleV1("00000000-0000-0000-0000-000000000000").Delete(role.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	r, err := client.PermissionV1().ProjectPermission().Update("00000000-0000-0000-0000-000000000000", &permissionv1.UpdateProjectPermissionsOptions{
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
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	resp, r, err := client.PermissionV1().ProjectPermission().GetAssignments("00000000-0000-0000-0000-000000000000", nil)
	assert.NoError(t, err)
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
		ProjectID: "00000000-0000-0000-0000-000000000000",
	}

	assert.Equal(t, want, resp)
}
