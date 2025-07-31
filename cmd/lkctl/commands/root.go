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
	var (
		clientOpts clientOptions
	)

	command := &cobra.Command{
		Use:   cliName,
		Short: "A CLI to interact with Lakekeeper's management - and Iceberg catalog APIs powered by go-iceberg.",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
		SilenceUsage:      true, // suppress usage on error
	}

	command.PersistentFlags().StringVar(&clientOpts.server, "server", common.GetEnvOr(common.EnvServer, common.DefaultServer), fmt.Sprintf("Lakekeeper base URL; set this or %s environment variable", common.EnvServer))
	command.PersistentFlags().StringVar(&clientOpts.authURL, "auth-url", common.GetEnvOr(common.EnvServer, ""), fmt.Sprintf("OAuth2 token endpoint; set this or %s environment variable", common.EnvAuthURL))
	command.PersistentFlags().StringVar(&clientOpts.clientID, "client-id", common.GetEnvOr(common.EnvClientID, ""), fmt.Sprintf("OAuth2 client_id; set this or %s environment variable", common.EnvClientID))
	command.PersistentFlags().StringVar(&clientOpts.clientSecret, "client-secret", common.GetEnvOr(common.EnvClientSecret, ""), fmt.Sprintf("OAuth2 client_secret; set this or %s environment variable", common.EnvClientSecret))
	command.PersistentFlags().StringSliceVar(&clientOpts.scope, "scopes", common.GetEnvSlice(common.EnvScope, " ", common.DefaultScope), fmt.Sprintf("OAuth2 scopes; set this or %s environment variable", common.EnvScope))
	command.PersistentFlags().BoolVar(&clientOpts.boostrap, "bootstrap", common.GetBoolEnv(common.EnvBootstrap), fmt.Sprintf("If set to true, the CLI will try to bootstrap the server with the current user first; set this or %s environment variable", common.EnvBootstrap))

	return command
}
