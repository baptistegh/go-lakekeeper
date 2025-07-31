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

	"github.com/baptistegh/go-lakekeeper/pkg/common"
	"github.com/joho/godotenv"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// load .env file if exists
	_ = godotenv.Load()
}

// NewCommand returns a new instance of an lkctl command
func NewCommand() *cobra.Command {
	var clientOpts clientOptions

	command := &cobra.Command{
		Use:   cliName,
		Short: "A CLI to interact with Lakekeeper's management - and Iceberg catalog APIs powered by go-iceberg.",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
		SilenceUsage:      true, // suppress usage on error
	}

	command.AddCommand(NewVersionCmd(&clientOpts))
	command.AddCommand(NewProjectCmd(&clientOpts))

	command.PersistentFlags().StringVar(&clientOpts.server, "server", common.GetEnvOr(common.EnvServer, common.DefaultServer), fmt.Sprintf("Lakekeeper base URL; set this or %s environment variable", common.EnvServer))
	command.PersistentFlags().StringVar(&clientOpts.authURL, "auth-url", common.GetEnvOr(common.EnvAuthURL, ""), fmt.Sprintf("OAuth2 token endpoint; set this or %s environment variable", common.EnvAuthURL))
	command.PersistentFlags().StringVar(&clientOpts.clientID, "client-id", common.GetEnvOr(common.EnvClientID, ""), fmt.Sprintf("OAuth2 client_id; set this or %s environment variable", common.EnvClientID))
	command.PersistentFlags().StringVar(&clientOpts.clientSecret, "client-secret", common.GetEnvOr(common.EnvClientSecret, ""), fmt.Sprintf("OAuth2 client_secret; set this or %s environment variable", common.EnvClientSecret))
	command.PersistentFlags().StringSliceVar(&clientOpts.scope, "scopes", common.GetEnvSlice(common.EnvScope, " ", common.DefaultScope), fmt.Sprintf("OAuth2 scopes; set this or %s environment variable", common.EnvScope))
	command.PersistentFlags().BoolVar(&clientOpts.boostrap, "bootstrap", common.GetBoolEnv(common.EnvBootstrap), fmt.Sprintf("If set to true, the CLI will try to bootstrap the server with the current user first; set this or %s environment variable", common.EnvBootstrap))

	return command
}
