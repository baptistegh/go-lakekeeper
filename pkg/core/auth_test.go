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

package core

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestK8sAuthSource(t *testing.T) {
	t.Parallel()

	t.Run("Provided ServiceAccountTokenPath", func(t *testing.T) {
		t.Parallel()

		path := "/tmp/service-account-go-lakekeeper-test"
		if err := os.WriteFile(path, []byte("service-account-token"), 0o600); err != nil {
			t.Fatalf("Failed to create service account token file: %v", err)
		}

		t.Cleanup(func() {
			_ = os.Remove(path)
		})

		as := K8sServiceAccountAuthSource{
			ServiceAccountTokenPath: Ptr(path),
		}

		err := as.Init(context.Background())
		require.NoError(t, err)

		key, value, err := as.Header(context.Background())
		require.NoError(t, err)

		assert.Equal(t, "Authorization", key)
		assert.Contains(t, "Bearer service-account-token", value)
	})

	t.Run("Default ServiceAccountTokenPath", func(t *testing.T) {
		t.Parallel()

		as := K8sServiceAccountAuthSource{}

		err := as.Init(context.Background())
		require.Error(t, err, "failed to read service account token")

		assert.Equal(t, Ptr("/var/run/secrets/kubernetes.io/serviceaccount/token"), as.ServiceAccountTokenPath)
	})
}
