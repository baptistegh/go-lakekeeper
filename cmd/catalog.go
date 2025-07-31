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
	"errors"

	"github.com/spf13/cobra"
)

// catalogCmd represents the catalog command
var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Invoke go-iceberg command line",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("catalog command is not implemented")
	},
}

func init() {
	rootCmd.AddCommand(catalogCmd)
}
