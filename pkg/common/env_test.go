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
