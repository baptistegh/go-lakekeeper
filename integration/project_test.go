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
	"context"
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/stretchr/testify/assert"
)

func TestProject_Create(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.ProjectV1().Create(t.Context(), &managementv1.CreateProjectOptions{
		Name: "test-project",
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
	assert.NotEmpty(t, resp.ID)

	t.Cleanup(func() {
		r, err = client.ProjectV1().Delete(context.Background(), resp.ID)
		if err != nil {
			t.Fatalf("could not delete project, %v", err)
		}
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})
}

func TestProject_Rename(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.ProjectV1().Create(t.Context(), &managementv1.CreateProjectOptions{
		Name: "test-project-2",
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
	assert.NotEmpty(t, resp.ID)

	t.Cleanup(func() {
		r, err = client.ProjectV1().Delete(context.Background(), resp.ID)
		if err != nil {
			t.Fatalf("could not delete project, %v", err)
		}
		assert.Equal(t, http.StatusNoContent, r.StatusCode)
	})

	r, err = client.ProjectV1().Rename(t.Context(), resp.ID, &managementv1.RenameProjectOptions{
		NewName: "test-project-renamed",
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	project, r, err := client.ProjectV1().Get(t.Context(), resp.ID)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "test-project-renamed", project.Name)
}

func TestProject_RenameDefault(t *testing.T) {
	client := Setup(t)

	t.Cleanup(func() {
		r, err := client.ProjectV1().RenameDefault(context.Background(), &managementv1.RenameProjectOptions{
			NewName: "Default Project",
		})
		if err != nil {
			t.Fatalf("could not rename default project, %v", err)
		}
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

	r, err := client.ProjectV1().RenameDefault(t.Context(), &managementv1.RenameProjectOptions{
		NewName: "test-project-renamed",
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	project, r, err := client.ProjectV1().GetDefault(t.Context())

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "test-project-renamed", project.Name)
}

func TestProject_Default(t *testing.T) {
	client := Setup(t)

	project, r, err := client.ProjectV1().GetDefault(t.Context())

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, defaultProjectID, project.ID)
	assert.Equal(t, "Default Project", project.Name)
}

func TestProject_Delete(t *testing.T) {
	client := Setup(t)

	project, r, err := client.ProjectV1().Create(t.Context(), &managementv1.CreateProjectOptions{
		Name: "test-project-3",
	})

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusCreated, r.StatusCode)
	assert.NotEmpty(t, project.ID)

	r, err = client.ProjectV1().Delete(t.Context(), project.ID)

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	p, r, err := client.ProjectV1().Get(t.Context(), project.ID)

	// Lakekeeper API sends 403 when trying to read a non existent object
	assert.ErrorContains(t, err, "Forbidden")
	assert.NotNil(t, r)
	assert.Nil(t, p)
}

func TestProject_List(t *testing.T) {
	client := Setup(t)

	resp, r, err := client.ProjectV1().List(t.Context())

	want := &managementv1.ListProjectsResponse{
		Projects: []*managementv1.Project{
			{
				ID:   defaultProjectID,
				Name: "Default Project",
			},
		},
	}

	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusOK, r.StatusCode)

	assert.Equal(t, want, resp)
}

// APIStatistics gives 0 called endpoints because when tests are run,
// no endpoints are being called before this test (or the call comes too fast).
//
// TODO: this integration test needs to be fixed
//
// func TestProject_GetAPIStatistics(t *testing.T) {
// 	client := Setup(t)
//
// 	resp, r, err := client.ProjectV1().GetAPIStatistics(defaultProjectID, &v1.GetAPIStatisticsOptions{
// 		Warehouse: struct {
// 			Type string  `json:"type"`
// 			ID   *string `json:"id,omitempty"`
// 		}{
// 			Type: "all",
// 		},
// 	})
//
// 	assert.NoError(t, err)
// 	assert.NotNil(t, r)
// 	assert.Equal(t, http.StatusOK, r.StatusCode)
//
// 	assert.IsType(t, &v1.GetAPIStatisticsResponse{}, resp)
// 	assert.NotEmpty(t, resp.CalledEnpoints)
// 	assert.NotEmpty(t, resp.NextPageToken)
// 	assert.NotEmpty(t, resp.PreviousPageToken)
// 	assert.NotEmpty(t, resp.Timestamps)
// }
