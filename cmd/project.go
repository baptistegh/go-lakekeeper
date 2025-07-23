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

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Interacts with the projects",
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

var createProjectCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
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

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(listProjectsCmd)
	projectCmd.AddCommand(GetProjectCmd)
	projectCmd.AddCommand(createProjectCmd)
}
