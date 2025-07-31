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

import "encoding/json"

type (
	AZSCredentialType string

	AZCredentialSettings interface {
		GetAZCredentialType() AZSCredentialType

		CredentialSettings
	}

	AZCredentialClientCredentials struct {
		ClientID     string `json:"client-id"`
		ClientSecret string `json:"client-secret"`
		TenantID     string `json:"tenant-id"`
	}

	AZCredentialSharedAccessKey struct {
		Key string `json:"key"`
	}

	AZCredentialManagedIdentity struct{}
)

const (
	ClientCredentials   AZSCredentialType = "client-credentials"
	SharedAccessKey     AZSCredentialType = "shared-access-key"
	AzureSystemIdentity AZSCredentialType = "azure-system-identity"
)

// verify implementations
var (
	_ AZCredentialSettings = (*AZCredentialClientCredentials)(nil)
	_ AZCredentialSettings = (*AZCredentialSharedAccessKey)(nil)
	_ AZCredentialSettings = (*AZCredentialManagedIdentity)(nil)

	_ CredentialSettings = (*AZCredentialClientCredentials)(nil)
	_ CredentialSettings = (*AZCredentialSharedAccessKey)(nil)
	_ CredentialSettings = (*AZCredentialManagedIdentity)(nil)
)

func NewAZCredentialClientCredentials(clientID, clientSecret, tenantID string) *AZCredentialClientCredentials {
	return &AZCredentialClientCredentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TenantID:     tenantID,
	}
}

func NewAZCredentialSharedAccessKey(key string) *AZCredentialSharedAccessKey {
	return &AZCredentialSharedAccessKey{
		Key: key,
	}
}

func NewAZCredentialManagedIdentity() *AZCredentialManagedIdentity {
	return &AZCredentialManagedIdentity{}
}

func (*AZCredentialClientCredentials) GetCredentialFamily() CredentialFamily {
	return AZCredentialFamily
}

func (*AZCredentialClientCredentials) GetAZCredentialType() AZSCredentialType {
	return ClientCredentials
}

func (c *AZCredentialClientCredentials) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}

func (c AZCredentialClientCredentials) MarshalJSON() ([]byte, error) {
	type Alias AZCredentialClientCredentials
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(c.GetCredentialFamily()),
		CredentialType: string(c.GetAZCredentialType()),
		Alias:          Alias(c),
	}
	return json.Marshal(aux)
}

func (*AZCredentialSharedAccessKey) GetCredentialFamily() CredentialFamily {
	return AZCredentialFamily
}

func (*AZCredentialSharedAccessKey) GetAZCredentialType() AZSCredentialType {
	return SharedAccessKey
}

func (c *AZCredentialSharedAccessKey) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}

func (c AZCredentialSharedAccessKey) MarshalJSON() ([]byte, error) {
	type Alias AZCredentialSharedAccessKey
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
		Alias
	}{
		Type:           string(c.GetCredentialFamily()),
		CredentialType: string(c.GetAZCredentialType()),
		Alias:          Alias(c),
	}
	return json.Marshal(aux)
}

func (*AZCredentialManagedIdentity) GetCredentialFamily() CredentialFamily {
	return AZCredentialFamily
}

func (*AZCredentialManagedIdentity) GetAZCredentialType() AZSCredentialType {
	return AzureSystemIdentity
}

func (c *AZCredentialManagedIdentity) AsCredential() StorageCredential {
	return StorageCredential{Settings: c}
}

func (c AZCredentialManagedIdentity) MarshalJSON() ([]byte, error) {
	aux := struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
	}{
		Type:           string(c.GetCredentialFamily()),
		CredentialType: string(c.GetAZCredentialType()),
	}
	return json.Marshal(aux)
}
