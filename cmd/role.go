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
	"errors"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// roleCmd represents the role command
var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Interacts with roles",
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

// addRoleCmd represents the role add command
var addRoleCmd = &cobra.Command{
	Use:   "add [flags] <name> [--description <description>]",
	Short: "Add a new role",
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

// getRolesCmd represents the role get command
var getRoleCmd = &cobra.Command{
	Use:   "get [flags] <role-id>",
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

// deleteRoleCmd represents role delete command
var deleteRoleCmd = &cobra.Command{
	Use:   "delete [flags] <role-id>",
	Short: "Delete a role",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := c.RoleV1(project).Delete(cmd.Context(), args[0]); err != nil {
			return err
		}

		return nil
	},
}

// updateRoleCmd represents the role update command
var updateRoleCmd = &cobra.Command{
	Use:   "update [flags] <role-id> <new-name> [--description <new-description>]",
	Short: "Update role",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.UpdateRoleOptions{
			Name: args[1],
		}

		if len(viper.GetString("role_u_description")) > 0 {
			opt.Description = core.Ptr(viper.GetString("role_u_description"))
		}

		resp, _, err := c.RoleV1(project).Update(cmd.Context(), args[0], &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// accessRoleCmd represents role access command
var accessRoleCmd = &cobra.Command{
	Use:   "access [flags] [--user <user> | --role <role>]",
	Short: "Get role access. By default, current user's access is returned",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		user := viper.GetString("role_access_user")
		role := viper.GetString("role_access_role")

		if len(user) > 0 && len(role) > 0 {
			return errors.New("you only can filter by user OR role, both were supplied")
		}

		opt := permissionv1.GetRoleAccessOptions{}

		if user != "" {
			opt.PrincipalUser = core.Ptr(user)
		}

		if role != "" {
			opt.PrincipalRole = core.Ptr(role)
		}

		resp, _, err := c.PermissionV1().RolePermission().GetAccess(cmd.Context(), args[0], &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

func init() {
	rootCmd.AddCommand(roleCmd)
	roleCmd.AddCommand(listRolesCmd)
	roleCmd.AddCommand(addRoleCmd)
	roleCmd.AddCommand(getRoleCmd)
	roleCmd.AddCommand(updateRoleCmd)
	roleCmd.AddCommand(accessRoleCmd)
	roleCmd.AddCommand(deleteRoleCmd)

	listRolesCmd.Flags().String("name", "", "filter by name")
	_ = viper.BindPFlag("role_name", listRolesCmd.Flags().Lookup("name"))

	addRoleCmd.Flags().String("description", "", "set a description")
	_ = viper.BindPFlag("role_description", addRoleCmd.Flags().Lookup("description"))

	updateRoleCmd.Flags().String("description", "", "set a new description")
	_ = viper.BindPFlag("role_u_description", updateRoleCmd.Flags().Lookup("description"))

	accessRoleCmd.Flags().String("user", "", "filter by user")
	accessRoleCmd.Flags().String("role", "", "filter by role")
	_ = viper.BindPFlag("role_access_user", accessRoleCmd.Flags().Lookup("user"))
	_ = viper.BindPFlag("role_access_role", accessRoleCmd.Flags().Lookup("role"))
}
