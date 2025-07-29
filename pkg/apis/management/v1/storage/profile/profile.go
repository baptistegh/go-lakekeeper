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
	"fmt"
)

type (
	StorageProfile struct {
		StorageSettings StorageSettings
	}

	StorageFamily string

	StorageSettings interface {
		GetStorageFamily() StorageFamily
		AsProfile() StorageProfile

		json.Marshaler
	}
)

const (
	StorageFamilyADLS StorageFamily = "adls"
	StorageFamilyGCS  StorageFamily = "gcs"
	StorageFamilyS3   StorageFamily = "s3"
)

// Check the implementation
var (
	_ StorageSettings = (*ADLSStorageSettings)(nil)
	_ StorageSettings = (*GCSStorageSettings)(nil)
	_ StorageSettings = (*S3StorageSettings)(nil)
)

func (sc *StorageProfile) UnmarshalJSON(data []byte) error {
	var peek struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	switch peek.Type {
	case "s3":
		var cfg S3StorageSettings
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageSettings = &cfg
	case "adls":
		var cfg ADLSStorageSettings
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageSettings = &cfg
	case "gcs":
		var cfg GCSStorageSettings
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		sc.StorageSettings = &cfg
	default:
		return fmt.Errorf("unsupported storage type: %s", peek.Type)
	}
	return nil
}

func (sc StorageProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(sc.StorageSettings)
}

// Type-safe helpers
func (sc StorageProfile) AsS3() (*S3StorageSettings, bool) {
	cfg, ok := sc.StorageSettings.(*S3StorageSettings)
	return cfg, ok
}

func (sc StorageProfile) AsADLS() (*ADLSStorageSettings, bool) {
	cfg, ok := sc.StorageSettings.(*ADLSStorageSettings)
	return cfg, ok
}

func (sc StorageProfile) AsGCS() (*GCSStorageSettings, bool) {
	cfg, ok := sc.StorageSettings.(*GCSStorageSettings)
	return cfg, ok
}
