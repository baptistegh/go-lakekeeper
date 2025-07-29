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

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Interacts with projects",
}

// listProjectsCmd represents the project list command
var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the available projects for the current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, _, err := c.ProjectV1().List(cmd.Context())
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp.Projects)
	},
}

// GetProjectCmd represents the project list command
var GetProjectCmd = &cobra.Command{
	Use:   "get",
	Short: "Print the project info",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, _, err := c.ProjectV1().Get(cmd.Context(), project)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// addProjectCmd represents the project add command
var addProjectCmd = &cobra.Command{
	Use:   "add [flags] <name>",
	Short: "Add a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.CreateProjectOptions{
			Name: args[0],
		}

		resp, _, err := c.ProjectV1().Create(cmd.Context(), &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp.ID)
	},
}

// renameProjectCmd represents the project rename command
var renameProjectCmd = &cobra.Command{
	Use:   "rename [flags] <new-name>",
	Short: "Rename the project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.RenameProjectOptions{
			NewName: args[0],
		}

		if _, err := c.ProjectV1().Rename(cmd.Context(), project, &opt); err != nil {
			return err
		}

		return nil
	},
}

// accessProjectCmd represents the project access command
var accessProjectCmd = &cobra.Command{
	Use:   "access [flags] [--user <user> | --role <role>]",
	Short: "Get project access",
	Long:  "Get project access. By default, current user's access is returned",
	RunE: func(cmd *cobra.Command, args []string) error {
		user := viper.GetString("project_access_user")
		role := viper.GetString("project_access_role")

		if user != "" && role != "" {
			return errors.New("you only can filter by user OR role, both were supplied")
		}

		opt := permissionv1.GetProjectAccessOptions{}

		if user != "" {
			opt.PrincipalUser = core.Ptr(user)
		}

		if role != "" {
			opt.PrincipalRole = core.Ptr(role)
		}

		resp, _, err := c.PermissionV1().ProjectPermission().GetAccess(cmd.Context(), project, &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(listProjectsCmd)
	projectCmd.AddCommand(GetProjectCmd)
	projectCmd.AddCommand(addProjectCmd)
	projectCmd.AddCommand(renameProjectCmd)
	projectCmd.AddCommand(accessProjectCmd)

	accessProjectCmd.Flags().String("user", "", "filter by user")
	accessProjectCmd.Flags().String("role", "", "filter by role")

	_ = viper.BindPFlag("project_access_user", accessProjectCmd.Flags().Lookup("user"))
	_ = viper.BindPFlag("project_access_role", accessProjectCmd.Flags().Lookup("role"))
}
