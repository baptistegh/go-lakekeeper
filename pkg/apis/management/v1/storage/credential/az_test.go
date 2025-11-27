package credential

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	jsonStr := `{"type":"az","credential-type":"client-credentials","client-id":"client-id","client-secret":"client-secret","tenant-id":"tenant-id"}`

	assert.JSONEq(t, jsonStr, string(b))
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
	require.NoError(t, err)

	jsonStr := `{"type":"az","credential-type":"shared-access-key","key":"key"}`

	assert.JSONEq(t, jsonStr, string(b))
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
	require.NoError(t, err)

	jsonStr := `{"type":"az","credential-type":"azure-system-identity"}`

	assert.JSONEq(t, jsonStr, string(b))
}
