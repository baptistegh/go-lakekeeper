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

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// roleCmd represents the role command
var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Interacts with the roles",
}

// listRolesCmd represents the role command
var listRolesCmd = &cobra.Command{
	Use:   "list",
	Short: "List available roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.ListRolesOptions{}

		if viper.GetString("role_name") != "" {
			opt.Name = core.Ptr(viper.GetString("role_name"))
		}

		resp, _, err := c.RoleV1(project).List(cmd.Context(), &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp.Roles)
	},
}

var createRoleCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.CreateRoleOptions{
			Name: args[0],
		}

		if viper.GetString("role_description") != "" {
			opt.Description = core.Ptr(viper.GetString("role_description"))
		}

		resp, _, err := c.RoleV1(project).Create(cmd.Context(), &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp.ID)
	},
}

var getRoleCmd = &cobra.Command{
	Use:   "get [role_id]",
	Short: "Get role informations by its id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, _, err := c.RoleV1(project).Get(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

func init() {
	rootCmd.AddCommand(roleCmd)
	roleCmd.AddCommand(listRolesCmd)
	roleCmd.AddCommand(createRoleCmd)
	roleCmd.AddCommand(getRoleCmd)

	listRolesCmd.Flags().String("name", "", "filter by name")
	_ = viper.BindPFlag("role_name", listRolesCmd.Flags().Lookup("name"))

	createRoleCmd.Flags().String("description", "", "set a description")
	_ = viper.BindPFlag("role_description", createRoleCmd.Flags().Lookup("description"))
}
