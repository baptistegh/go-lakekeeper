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
	"fmt"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	managmentv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewServerCmd(clientOptions *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:     "server",
		Aliases: []string{"srv"},
		Short:   "Manage server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewInfoCmd(clientOptions))
	command.AddCommand(NewBootstrapCmd(clientOptions))

	return &command
}

func NewBootstrapCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		asOperator       bool
		acceptTermsOfUse bool

		output string
	)

	command := cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstraps the server with the current user",
		Example: `  # Bootstrap the server and get the server admin role
  lkctl bootstrap --accept-terms-of-use

  # Bootstrap the server as an operator
  lkctl bootstrap --accept-terms-of-use --as-operator`,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			opt := managmentv1.BootstrapServerOptions{
				AcceptTermsOfUse: acceptTermsOfUse,
				IsOperator:       &asOperator,
			}

			client := MustCreateClient(ctx, clientOpts).ServerV1()

			_, err := client.Bootstrap(ctx, &opt)
			errors.Check(err)

			switch output {
			case "json":
				info, _, err := client.Info(ctx)
				errors.Check(err)

				err = PrintResource(info, output)
				errors.Check(err)
			case "text":
				fmt.Println("Server bootstrapped successfully")
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	command.Flags().BoolVar(&asOperator, "as-operator", false, "Bootstrap the server as an operator")
	command.Flags().BoolVar(&acceptTermsOfUse, "accept-terms-of-use", false, "Accept the terms of use")

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewInfoCmd(clientOptions *clientOptions) *cobra.Command {
	var output string

	command := cobra.Command{
		Use:   "info",
		Short: "Print server informations",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			resp, _, err := MustCreateClient(ctx, clientOptions).ServerV1().Info(ctx)
			errors.Check(err)

			switch output {
			case "text":
				fmt.Printf("ID: %s\n", resp.ServerID)
				fmt.Printf("Version: %s\n", resp.Version)
				fmt.Printf("Default Project ID: %s\n", resp.DefaultProjectID)
				fmt.Printf("Bootstraped: %t\n", resp.Bootstrapped)
				fmt.Printf("Authorization Backend: %s\n", resp.AuthzBackend)
				fmt.Printf("AWS System Identities Enabled: %t\n", resp.AWSSystemIdentitiesEnabled)
				fmt.Printf("Azure System Identities Enabled: %t\n", resp.AzureSystemIdentitiesEnabled)
				fmt.Printf("GCP System Identities Enableds: %t\n", resp.GCPSystemIdentitiesEnabled)
				fmt.Println("Queues:")
				for _, q := range resp.Queues {
					fmt.Printf("  %s\n", q)
				}
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Printf("unknown output format: %s\n", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}
