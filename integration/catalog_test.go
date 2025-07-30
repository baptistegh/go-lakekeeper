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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatalog_Basic(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)

	_, warehouseDefault := MustCreateWarehouse(t, client, defaultProjectID)
	_, err := client.CatalogV1(context.Background(), defaultProjectID, warehouseDefault)
	assert.NoError(t, err)

	_, warehouseProject := MustCreateWarehouse(t, client, project)
	_, err = client.CatalogV1(context.Background(), project, warehouseProject)
	assert.NoError(t, err)
}
