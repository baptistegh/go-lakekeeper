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
	"github.com/stretchr/testify/require"
)

func TestTabularDeleteProfileHard(t *testing.T) {
	profile := NewTabularDeleteProfileHard()

	assert.Equal(t, HardDeleteProfileType, profile.GetDeteProfileType())

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"hard"}`

	assert.JSONEq(t, jsonStr, string(b))

	expected := &DeleteProfile{
		DeleteProfileSettings: &TabularDeleteProfileHard{},
	}

	assert.Equal(t, expected, profile.AsProfile())
}

func TestTabularDeleteProfileSoft(t *testing.T) {
	profile := NewTabularDeleteProfileSoft(7200)

	assert.Equal(t, SoftDeleteProfileType, profile.GetDeteProfileType())

	b, err := json.Marshal(profile)
	require.NoError(t, err)

	jsonStr := `{"type":"soft","expiration-seconds":7200}`

	assert.JSONEq(t, jsonStr, string(b))

	expected := &DeleteProfile{
		DeleteProfileSettings: &TabularDeleteProfileSoft{7200},
	}

	assert.Equal(t, expected, profile.AsProfile())
}

func TestDeleteProfil_Unmarshal(t *testing.T) {
	tests := []struct {
		input    string
		expected DeleteProfile
	}{
		{
			`{"type":"hard"}`,
			*NewTabularDeleteProfileHard().AsProfile(),
		},
		{
			`{"type":"soft","expiration-seconds":7200}`,
			*NewTabularDeleteProfileSoft(7200).AsProfile(),
		},
	}

	for _, test := range tests {
		var deleteProfile DeleteProfile
		err := json.Unmarshal([]byte(test.input), &deleteProfile)
		require.NoError(t, err)
		assert.Equal(t, test.expected, deleteProfile)
	}
}
