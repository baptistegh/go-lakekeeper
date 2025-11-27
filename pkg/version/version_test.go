package version

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	// Save original values to restore them after the test.
	// This is important to not interfere with other tests.
	originalVersion := version
	originalBuildDate := buildDate
	originalGitCommit := gitCommit
	originalGitTag := gitTag
	originalGitTreeState := gitTreeState
	t.Cleanup(func() {
		version = originalVersion
		buildDate = originalBuildDate
		gitCommit = originalGitCommit
		gitTag = originalGitTag
		gitTreeState = originalGitTreeState
	})

	// Common runtime values
	goVersion := runtime.Version()
	compiler := runtime.Compiler
	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	testCases := []struct {
		name           string
		version        string
		buildDate      string
		gitCommit      string
		gitTag         string
		gitTreeState   string
		expectedString string
	}{
		{
			name:           "official release",
			version:        "1.2.3",
			buildDate:      "2025-01-01T12:00:00Z",
			gitCommit:      "abcdef1234567890",
			gitTag:         "v1.2.3",
			gitTreeState:   "clean",
			expectedString: "v1.2.3",
		},
		{
			name:           "dev build with clean tree",
			version:        "1.2.3-dev",
			buildDate:      "2025-01-01T12:00:00Z",
			gitCommit:      "abcdef1234567890",
			gitTag:         "",
			gitTreeState:   "clean",
			expectedString: "v1.2.3-dev+abcdef1",
		},
		{
			name:           "dev build with dirty tree",
			version:        "1.2.3-dev",
			buildDate:      "2025-01-01T12:00:00Z",
			gitCommit:      "abcdef1234567890",
			gitTag:         "",
			gitTreeState:   "dirty",
			expectedString: "v1.2.3-dev+abcdef1.dirty",
		},
		{
			name:           "dev build with short git commit",
			version:        "1.2.3-dev",
			buildDate:      "2025-01-01T12:00:00Z",
			gitCommit:      "abc",
			gitTag:         "",
			gitTreeState:   "clean",
			expectedString: "v1.2.3-dev+unknown",
		},
		{
			name:           "dev build with no git commit",
			version:        "1.2.3-dev",
			buildDate:      "2025-01-01T12:00:00Z",
			gitCommit:      "",
			gitTag:         "",
			gitTreeState:   "clean",
			expectedString: "v1.2.3-dev+unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set package-level variables for this test case
			version, buildDate, gitCommit, gitTag, gitTreeState = tc.version, tc.buildDate, tc.gitCommit, tc.gitTag, tc.gitTreeState

			expected := Version{
				Version:      tc.expectedString,
				BuildDate:    tc.buildDate,
				GitCommit:    tc.gitCommit,
				GitTag:       tc.gitTag,
				GitTreeState: tc.gitTreeState,
				GoVersion:    goVersion,
				Compiler:     compiler,
				Platform:     platform,
			}

			assert.Equal(t, expected, GetVersion())
		})
	}
}

func TestVersion_String(t *testing.T) {
	v := Version{Version: "v1.2.3"}
	assert.Equal(t, "v1.2.3", v.String())
}
