package profile

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTabularDeleteProfileHard(t *testing.T) {
	profile := NewTabularDeleteProfileHard()

	assert.Equal(t, HardDeleteProfileType, profile.GetDeteProfileType())

	b, err := json.Marshal(profile)
	assert.NoError(t, err)

	jsonStr := `{"type":"hard"}`

	assert.Equal(t, jsonStr, string(b))

	expected := &DeleteProfile{
		DeleteProfileSettings: &TabularDeleteProfileHard{},
	}

	assert.Equal(t, expected, profile.AsProfile())
}

func TestTabularDeleteProfileSoft(t *testing.T) {
	profile := NewTabularDeleteProfileSoft(7200)

	assert.Equal(t, SoftDeleteProfileType, profile.GetDeteProfileType())

	b, err := json.Marshal(profile)
	assert.NoError(t, err)

	jsonStr := `{"type":"soft","expiration-seconds":7200}`

	assert.Equal(t, jsonStr, string(b))

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
		assert.NoError(t, err)
		assert.Equal(t, test.expected, deleteProfile)
	}
}
