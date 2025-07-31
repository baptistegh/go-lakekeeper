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
	"encoding/json"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Interacts with users",
}

// listUsersCmd represents the user list command
var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.ListUsersOptions{}

		if viper.GetString("user_name") != "" {
			opt.Name = core.Ptr(viper.GetString("role_name"))
		}

		resp, _, err := c.UserV1().List(cmd.Context(), &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// addUserCmd represents the user add command
var addUserCmd = &cobra.Command{
	Use:   "add [flags] <user-id> <name> <user-type> [--email <email>] [--update]",
	Short: "Add a new user",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.ProvisionUserOptions{
			ID:       core.Ptr(args[0]),
			Name:     core.Ptr(args[1]),
			UserType: core.Ptr(managementv1.UserType(args[2])),
		}

		if len(viper.GetString("email")) > 0 {
			opt.Email = core.Ptr(viper.GetString("email"))
		}

		if viper.GetBool("user_update") {
			opt.UpdateIfExists = core.Ptr(true)
		}

		resp, _, err := c.UserV1().Provision(cmd.Context(), &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// getUserCmd represents the user get command
var getUserCmd = &cobra.Command{
	Use:   "get [flags] <user-id>",
	Short: "Get user by id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, _, err := c.UserV1().Get(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// deleteUserCmd represents the user delete command
var deleteUserCmd = &cobra.Command{
	Use:   "delete [flags] <user-id>",
	Short: "Delete a user by id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := c.UserV1().Delete(cmd.Context(), args[0]); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(listUsersCmd)
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(getUserCmd)
	userCmd.AddCommand(deleteUserCmd)

	listUsersCmd.Flags().String("name", "", "filter by name")
	_ = viper.BindPFlag("user_name", listRolesCmd.Flags().Lookup("name"))

	addUserCmd.Flags().String("email", "", "add an email to the user")
	_ = viper.BindPFlag("email", addUserCmd.Flags().Lookup("email"))

	addUserCmd.Flags().Bool("update", false, "update if user exists")
	_ = viper.BindPFlag("user_update", addUserCmd.Flags().Lookup("update"))
}
