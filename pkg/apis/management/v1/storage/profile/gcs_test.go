package profile

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGCSStorageSettings_NoOpts(t *testing.T) {
	profile := NewGCSStorageSettings("bucket")

	assert.Equal(t, StorageFamilyGCS, profile.GetStorageFamily())

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"gcs","bucket":"bucket"}`

	assert.JSONEq(t, jsonStr, string(b))
}

func TestGCSStorageSettings_KeyPrefix(t *testing.T) {
	profile := NewGCSStorageSettings("bucket", WithGCSKeyPrefix("prefix"))

	assert.Equal(t, StorageFamilyGCS, profile.GetStorageFamily())
	assert.Equal(t, "prefix", *profile.KeyPrefix)

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"gcs","bucket":"bucket","key-prefix":"prefix"}`

	assert.JSONEq(t, jsonStr, string(b))
}
