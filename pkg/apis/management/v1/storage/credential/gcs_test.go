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

func TestGCSCredentialServiceAccountKey(t *testing.T) {
	t.Parallel()

	creds := NewGCSCredentialServiceAccountKey(GCSServiceKey{
		AuthProviderX509CertURL: "auth-provider-x509-cert-url",
		AuthURI:                 "auth-uri",
		ClientEmail:             "client-email",
		ClientID:                "client-id",
		ClientX509CertURL:       "client-x509-cert-url",
		PrivateKey:              "private-key",
		PrivateKeyID:            "private-key-id",
		ProjectID:               "project-id",
		TokenURI:                "token-uri",
		Type:                    "type",
		UniverseDomain:          "universe-domain",
	})

	assert.Equal(t, GCSCredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, ServiceAccountKey, creds.GetGCSCredentialType())

	want := StorageCredential{Settings: creds}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	assert.NoError(t, err)

	jsonStr := `{"type":"gcs","credential-type":"service-account-key","key":{"auth_provider_x509_cert_url":"auth-provider-x509-cert-url","auth_uri":"auth-uri","client_email":"client-email","client_id":"client-id","client_x509_cert_url":"client-x509-cert-url","private_key":"private-key","private_key_id":"private-key-id","project_id":"project-id","token_uri":"token-uri","type":"type","universe_domain":"universe-domain"}}`

	assert.Equal(t, jsonStr, string(b))
}

func TestGCSCredentialSystemIdentity(t *testing.T) {
	t.Parallel()

	creds := NewGCSCredentialSystemIdentity()

	assert.Equal(t, GCSCredentialFamily, creds.GetCredentialFamily())
	assert.Equal(t, GCPSystemIdentity, creds.GetGCSCredentialType())

	want := StorageCredential{Settings: creds}

	assert.Equal(t, want, creds.AsCredential())

	b, err := json.Marshal(creds)
	assert.NoError(t, err)

	jsonStr := `{"type":"gcs","credential-type":"gcp-system-identity"}`

	assert.Equal(t, jsonStr, string(b))
}
