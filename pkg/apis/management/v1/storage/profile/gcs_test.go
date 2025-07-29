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
