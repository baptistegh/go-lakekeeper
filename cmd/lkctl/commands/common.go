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

package commands

import (
	"encoding/json"
	"fmt"
)

const (
	cliName = "lkctl"
)

// PrintResource prints a single resource in YAML or JSON format to stdout according to the output format
func PrintResource(resource any, output string) error {
	switch output {
	case "json":
		jsonBytes, err := json.MarshalIndent(resource, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal resource to json: %w", err)
		}
		fmt.Println(string(jsonBytes))
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}
	return nil
}
