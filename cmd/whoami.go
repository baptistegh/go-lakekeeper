// Copyright © 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

// whoamiCmd represents the user command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Print the current user's informations",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, _, err := c.UserV1().Whoami(cmd.Context())
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
