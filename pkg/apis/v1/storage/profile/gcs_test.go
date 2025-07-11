package profile

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCSStorageSettings_NoOpts(t *testing.T) {
	profile := NewGCSStorageSettings("bucket")

	assert.Equal(t, StorageFamilyGCS, profile.GetStorageFamily())

	b, err := json.Marshal(profile)
	assert.NoError(t, err)

	jsonStr := `{"type":"gcs","bucket":"bucket"}`

	assert.Equal(t, jsonStr, string(b))
}

func TestGCSStorageSettings_KeyPrefix(t *testing.T) {
	profile := NewGCSStorageSettings("bucket", WithGCSKeyPrefix("prefix"))

	assert.Equal(t, StorageFamilyGCS, profile.GetStorageFamily())
	assert.Equal(t, "prefix", *profile.KeyPrefix)

	b, err := json.Marshal(profile)
	assert.NoError(t, err)

	jsonStr := `{"type":"gcs","bucket":"bucket","key-prefix":"prefix"}`

	assert.Equal(t, jsonStr, string(b))
}
