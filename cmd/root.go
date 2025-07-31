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

package main

import (
	"os"
	"strings"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	c *client.Client

	project    string
	serverInfo *managementv1.ServerInfo
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "lkctl",
	Short:        "A CLI to interact with Lakekeeper's management - and Iceberg catalog APIs powered by go-iceberg.",
	SilenceUsage: true, // suppress usage on error
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		for _, arg := range os.Args[1:] {
			if arg == "--help" || arg == "-h" || arg == "help" || arg == "version" {
				return nil
			}
		}

		opts := []client.ClientOptionFunc{}

		oauthConfig := clientcredentials.Config{
			ClientID:     viper.GetString("client_id"),
			ClientSecret: viper.GetString("client_secret"),
			TokenURL:     viper.GetString("auth_url"),
			Scopes:       strings.Split(viper.GetString("scope"), " "),
		}

		if _, err := oauthConfig.Token(cmd.Context()); err != nil {
			return err
		}

		as := core.OAuthTokenSource{
			TokenSource: oauthConfig.TokenSource(cmd.Context()),
		}

		if viper.GetBool("bootstrap") {
			opts = append(opts, client.WithInitialBootstrapV1Enabled(true, true, core.Ptr(managementv1.ApplicationUserType)))
		}

		cli, err := client.NewAuthSourceClient(cmd.Context(), &as, viper.GetString("server"), opts...)
		if err != nil {
			return err
		}

		info, _, err := cli.ServerV1().Info(cmd.Context())
		if err != nil {
			return err
		}

		serverInfo = info
		c = cli

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", uuid.Nil.String(), "filter by project")
	rootCmd.PersistentFlags().String("server", "http://localhost:8181", "lakekeeper server base url. (can also be set with env LAKEKEEPER_SERVER)")
	rootCmd.PersistentFlags().String("auth-url", "http://localhost:30080/realms/iceberg/protocol/openid-connect/token", "oidc token endpoint (can also be set with env LAKEKEEPER_AUTH_URL)")
	rootCmd.PersistentFlags().String("client-id", "lakekeeper-admin", "oidc client_id  (can also be set with env LAKEKEEPER_CLIENT_ID)")
	rootCmd.PersistentFlags().String("client-secret", "", "oidc client_secret  (can also be set with env LAKEKEEPER_CLIENT_SECRET)")
	rootCmd.PersistentFlags().String("scope", "lakekeeper", "oidc scope, space separated (can also be set with env LAKEKEEPER_SCOPE)")
	rootCmd.PersistentFlags().Bool("bootstrap", false, "bootstrap the server on startup. The current user will have the operator role (can also be set with env LAKEKEEPER_BOOTSTRAP)")

	_ = viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	_ = viper.BindPFlag("auth_url", rootCmd.PersistentFlags().Lookup("auth-url"))
	_ = viper.BindPFlag("client_id", rootCmd.PersistentFlags().Lookup("client-id"))
	_ = viper.BindPFlag("client_secret", rootCmd.PersistentFlags().Lookup("client-secret"))
	_ = viper.BindPFlag("scope", rootCmd.PersistentFlags().Lookup("scope"))
	_ = viper.BindPFlag("bootstrap", rootCmd.PersistentFlags().Lookup("bootstrap"))

	// load .env file if exists
	_ = godotenv.Load()

	viper.SetEnvPrefix("lakekeeper")
	viper.AutomaticEnv()
}
