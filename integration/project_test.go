//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestProject_Create(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.ProjectV1().Create(&v1.CreateProjectOptions{
		Name: "test-project",
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
	assert.NotEmpty(t, resp.ID)

	t.Cleanup(func() {
		r, err = client.ProjectV1().Delete(resp.ID)
		if err != nil {
			t.Fatalf("could not delete project, %v", err)
		}
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})
}

func TestProject_Rename(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.ProjectV1().Create(&v1.CreateProjectOptions{
		Name: "test-project-2",
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
	assert.NotEmpty(t, resp.ID)

	t.Cleanup(func() {
		r, err = client.ProjectV1().Delete(resp.ID)
		if err != nil {
			t.Fatalf("could not delete project, %v", err)
		}
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	r, err = client.ProjectV1().Rename(resp.ID, &v1.RenameProjectOptions{
		NewName: "test-project-renamed",
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	project, r, err := client.ProjectV1().Get(resp.ID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "test-project-renamed", project.Name)
}

func TestProject_Default(t *testing.T) {
	client := Setup(t)

	project, r, err := client.ProjectV1().Default()

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, "00000000-0000-0000-0000-000000000000", project.ID)
	assert.Equal(t, "Default Project", project.Name)
}

func TestProject_Delete(t *testing.T) {
	client := Setup(t)

	project, r, err := client.ProjectV1().Create(&v1.CreateProjectOptions{
		Name: "test-project-3",
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
	assert.NotEmpty(t, project.ID)

	r, err = client.ProjectV1().Delete(project.ID)

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	p, r, err := client.ProjectV1().Get(project.ID)

	// Lakekeeper API sends 403 when trying to read a non existent object
	assert.ErrorContains(t, err, "Forbidden")
	assert.NotNil(t, r)
	assert.Nil(t, p)
}

func TestProject_List(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.ProjectV1().List()

	want := &v1.ListProjectsResponse{
		Projects: []*v1.Project{
			{
				ID:   "00000000-0000-0000-0000-000000000000",
				Name: "Default Project",
			},
		},
	}

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, want, resp)
}
