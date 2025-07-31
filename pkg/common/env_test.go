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

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOr(t *testing.T) {
	t.Setenv("TEST_ENV", "test")

	assert.Equal(t, "test", GetEnvOr("TEST_ENV", "fallback"))
	assert.Equal(t, "fallback", GetEnvOr("TEST_ENV2", "fallback"))
}

func TestGetEnvSlice(t *testing.T) {
	t.Setenv("TEST_ENV", "test1,test2")

	assert.Equal(t, []string{"test1", "test2"}, GetEnvSlice("TEST_ENV", ",", []string{"test"}))
	assert.Equal(t, []string{"test", "test2"}, GetEnvSlice("TEST_ENV2", ",", []string{"test", "test2"}))
}

func TestGetBoolEnv(t *testing.T) {
	t.Setenv("TEST_ENV", "true")

	assert.True(t, GetBoolEnv("TEST_ENV"))
	assert.False(t, GetBoolEnv("TEST_ENV2"))
}
