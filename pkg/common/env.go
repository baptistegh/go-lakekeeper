package common

import (
	"os"
	"strings"
)

const (
	EnvServer       = "LAKEKEEPER_SERVER"
	EnvAuthURL      = "LAKEKEEPER_AUTH_URL"
	EnvClientID     = "LAKEKEEPER_CLIENT_ID"
	EnvClientSecret = "LAKEKEEPER_CLIENT_SECRET"
	EnvScope        = "LAKEKEEPER_SCOPE"
	EnvBootstrap    = "LAKEKEEPER_BOOTSTRAP"
)

func GetEnvOr(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return v
}

func GetEnvSlice(key string, sep string, fallback []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return strings.Split(v, sep)
}

func GetBoolEnv(key string) bool {
	v := os.Getenv(key)
	return strings.ToLower(v) == "true"
}
