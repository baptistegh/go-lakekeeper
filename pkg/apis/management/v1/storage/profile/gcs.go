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

package profile

import (
	"encoding/json"
)

type (
	// GCSStorageSettings represents the storage settings for a warehouse
	// where data are stored on Google Cloud Storage.
	GCSStorageSettings struct {
		// Name of the GCS bucket
		Bucket string `json:"bucket"`
		// Subpath in the bucket to use.
		KeyPrefix *string `json:"key-prefix,omitempty"`
	}

	GCSStorageSettingsOptions func(*GCSStorageSettings)
)

func (sp *GCSStorageSettings) GetStorageFamily() StorageFamily {
	return StorageFamilyGCS
}

// NewGCSStorageSettings creates a new GCS storage profile considering
// the options given.
func NewGCSStorageSettings(bucket string, opts ...GCSStorageSettingsOptions) *GCSStorageSettings {
	// Default configuration
	profile := GCSStorageSettings{
		Bucket: bucket,
	}

	// Apply options
	for _, v := range opts {
		v(&profile)
	}

	return &profile
}

func WithGCSKeyPrefix(prefix string) GCSStorageSettingsOptions {
	return func(sp *GCSStorageSettings) {
		sp.KeyPrefix = &prefix
	}
}

func (sp *GCSStorageSettings) AsProfile() StorageProfile {
	return StorageProfile{sp}
}

func (sp GCSStorageSettings) MarshalJSON() ([]byte, error) {
	type Alias GCSStorageSettings
	aux := struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  string(StorageFamilyGCS),
		Alias: Alias(sp),
	}
	return json.Marshal(aux)
}
