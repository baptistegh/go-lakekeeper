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

package main

import (
	"github.com/baptistegh/go-lakekeeper/pkg/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of lkctl",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("version_short") {
			cmd.Printf("%s\n", version.Version)
			return
		}

		cmd.Printf("version=%s, commit=%s, date=%s\n", version.Version, version.Commit, version.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().Bool("short", false, "print only the version number")

	_ = viper.BindPFlag("version_short", versionCmd.Flags().Lookup("short"))
}
