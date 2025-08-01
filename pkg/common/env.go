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

func GetEnvOr(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return v
}

func GetEnvSlice(key, sep string, fallback []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return strings.Split(v, sep)
}

func GetBoolEnv(key string) bool {
	v := os.Getenv(key)
	return strings.EqualFold(v, "true")
}
