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
	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"

	"github.com/spf13/cobra"
)

var (
	asOperator       bool
	acceptTermsOfUse bool
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstraps the server with the current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.BootstrapServerOptions{
			AcceptTermsOfUse: acceptTermsOfUse,
			IsOperator:       &asOperator,
		}

		_, err := c.ServerV1().Bootstrap(cmd.Context(), &opt)
		return err
	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)

	bootstrapCmd.Flags().BoolVar(&asOperator, "as-operator", false, "Bootstrap the server as an operator")
	bootstrapCmd.Flags().BoolVar(&acceptTermsOfUse, "accept-terms-of-use", false, "Accept the terms of use")
}
