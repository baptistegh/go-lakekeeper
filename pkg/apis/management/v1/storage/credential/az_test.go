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

package credential

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAZCredentialClientCredentials(t *testing.T) {
	t.Parallel()

	creds := NewAZCredentialClientCredentials("client-id", "client-secret", "tenant-id")

	assert.Equal(t, AZCredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, ClientCredentials, creds.GetAZCredentialType())

	want := StorageCredential{
		Settings: creds,
	}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	assert.NoError(t, err)

	jsonStr := `{"type":"az","credential-type":"client-credentials","client-id":"client-id","client-secret":"client-secret","tenant-id":"tenant-id"}`

	assert.Equal(t, jsonStr, string(b))
}

func TestAZCredentialSharedAccessKey(t *testing.T) {
	t.Parallel()

	creds := NewAZCredentialSharedAccessKey("key")

	assert.Equal(t, AZCredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, SharedAccessKey, creds.GetAZCredentialType())

	want := StorageCredential{
		Settings: creds,
	}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	assert.NoError(t, err)

	jsonStr := `{"type":"az","credential-type":"shared-access-key","key":"key"}`

	assert.Equal(t, jsonStr, string(b))
}

func TestAZCredentialManagedIdentity(t *testing.T) {
	t.Parallel()

	creds := NewAZCredentialManagedIdentity()

	assert.Equal(t, AZCredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, AzureSystemIdentity, creds.GetAZCredentialType())

	want := StorageCredential{
		Settings: creds,
	}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	assert.NoError(t, err)

	jsonStr := `{"type":"az","credential-type":"azure-system-identity"}`

	assert.Equal(t, jsonStr, string(b))
}
