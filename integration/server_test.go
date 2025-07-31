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
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerService_Info(t *testing.T) {
	client := Setup(t)

	info, _, err := client.ServerV1().Info(t.Context())

	require.NoError(t, err)
	assert.Equal(t, true, info.Bootstrapped)
	assert.NotEmpty(t, true, info.ServerID)
	assert.NotEmpty(t, true, info.DefaultProjectID)
	assert.Equal(t, info.AuthzBackend, "openfga")
}

func TestServerService_Bootstrap(t *testing.T) {
	client := Setup(t)

	resp, err := client.ServerV1().Bootstrap(t.Context(), &managementv1.BootstrapServerOptions{
		AcceptTermsOfUse: true,
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)

	info, _, err := client.ServerV1().Info(t.Context())

	require.NoError(t, err)
	assert.Equal(t, true, info.Bootstrapped)
}
