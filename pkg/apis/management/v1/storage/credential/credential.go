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
	"fmt"
)

type (
	StorageCredential struct {
		Settings CredentialSettings
	}

	CredentialFamily string

	CredentialSettings interface {
		GetCredentialFamily() CredentialFamily
		AsCredential() StorageCredential

		json.Marshaler
	}
)

const (
	S3CredentialFamily  CredentialFamily = "s3"
	GCSCredentialFamily CredentialFamily = "gcs"
	AZCredentialFamily  CredentialFamily = "az"
)

func (sc *StorageCredential) UnmarshalJSON(data []byte) error {
	var peek struct {
		Type           string `json:"type"`
		CredentialType string `json:"credential-type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	switch fmt.Sprintf("%s:%s", peek.Type, peek.CredentialType) {
	case fmt.Sprintf("%s:%s", S3CredentialFamily, AccessKey):
		var cfg S3CredentialAccessKey
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", S3CredentialFamily, AWSSystemIdentity):
		var cfg S3CredentialSystemIdentity
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", S3CredentialFamily, CloudflareR2):
		var cfg CloudflareR2Credential
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", GCSCredentialFamily, ServiceAccountKey):
		var cfg GCSCredentialServiceAccountKey
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", GCSCredentialFamily, GCPSystemIdentity):
		var cfg GCSCredentialSystemIdentity
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", AZCredentialFamily, ClientCredentials):
		var cfg AZCredentialClientCredentials
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", AZCredentialFamily, SharedAccessKey):
		var cfg AZCredentialSharedAccessKey
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	case fmt.Sprintf("%s:%s", AZCredentialFamily, AzureSystemIdentity):
		var cfg AZCredentialManagedIdentity
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.Settings = &cfg
	default:
		return fmt.Errorf("unsupported storage credential type: %s / %s", peek.Type, peek.CredentialType)
	}
	return nil
}

func (sc StorageCredential) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.Settings)
}

// Type-safe helpers

func (sc StorageCredential) AsS3() (S3SCredentialSettings, bool) {
	cfg, ok := sc.Settings.(S3SCredentialSettings)
	return cfg, ok
}

func (sc StorageCredential) AsAZ() (AZCredentialSettings, bool) {
	cfg, ok := sc.Settings.(AZCredentialSettings)
	return cfg, ok
}

func (sc StorageCredential) AsGCS() (GCSSCredentialSettings, bool) {
	cfg, ok := sc.Settings.(GCSSCredentialSettings)
	return cfg, ok
}
