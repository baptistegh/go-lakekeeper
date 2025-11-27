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
