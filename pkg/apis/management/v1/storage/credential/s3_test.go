package credential

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3CredentialSystemIdentity(t *testing.T) {
	t.Parallel()

	creds := NewS3CredentialSystemIdentity("external-id")

	assert.Equal(t, S3CredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, AWSSystemIdentity, creds.GetS3CredentialType())

	want := StorageCredential{
		Settings: creds,
	}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	require.NoError(t, err)

	jsonStr := `{"type":"s3","credential-type":"aws-system-identity","external-id":"external-id"}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestS3CredentialAccessKey(t *testing.T) {
	t.Parallel()

	creds := NewS3CredentialAccessKey("access-key", "secret-key", WithExternalID("external-id"))

	assert.Equal(t, S3CredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, AccessKey, creds.GetS3CredentialType())

	want := StorageCredential{
		Settings: creds,
	}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	require.NoError(t, err)

	jsonStr := `{"type":"s3","credential-type":"access-key","aws-access-key-id":"access-key","aws-secret-access-key":"secret-key","external-id":"external-id"}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestCloudflareR2Credential(t *testing.T) {
	t.Parallel()

	creds := NewCloudflareR2Credential("access-key", "secret-key", "account-id", "token")

	assert.Equal(t, S3CredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, CloudflareR2, creds.GetS3CredentialType())

	want := StorageCredential{
		Settings: creds,
	}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	require.NoError(t, err)

	jsonStr := `{"type":"s3","credential-type":"cloudflare-r2","access-key-id":"access-key","secret-access-key":"secret-key","account=id":"account-id","token":"token"}`

	assert.JSONEq(t, jsonStr, string(b))
}
