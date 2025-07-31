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
	"errors"
	"fmt"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// projectAssignmentsCmd represents the project assignments command
var projectAssignmentsCmd = &cobra.Command{
	Use:   "assignments",
	Short: "Interacts with project assignments",
}

// getProjectAssignmentsCmd represents the project assignments get command
var getProjectAssignmentsCmd = &cobra.Command{
	Use:   "get [flags] [--relation <assignment> ...]",
	Short: "Get project assignments",
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := permissionv1.GetProjectAssignmentsOptions{}

		for _, v := range viper.GetStringSlice("project_relations") {
			opt.Relations = append(opt.Relations, permissionv1.ProjectAssignmentType(v))
		}

		resp, _, err := c.PermissionV1().ProjectPermission().GetAssignments(cmd.Context(), project, &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// wProjectAssignmentsCmd represents the project assignments add command
var wProjectAssignmentsCmd = &cobra.Command{
	Use:   "add [flags] --assignment <assignment> [--user <user> --role <role>]",
	Short: "add project assignments",
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := permissionv1.UpdateProjectPermissionsOptions{}
		assignees := []permissionv1.UserOrRole{}

		if len(viper.GetStringSlice("project_ass_assignments")) < 1 {
			return errors.New("you must set at lest one assignment")
		}

		if len(viper.GetStringSlice("project_ass_users")) < 1 && len(viper.GetStringSlice("project_ass_roles")) < 1 {
			return errors.New("you must set at least one user or role")
		}

		for _, v := range viper.GetStringSlice("project_ass_users") {
			assignees = append(assignees, permissionv1.UserOrRole{
				Type:  permissionv1.UserType,
				Value: v,
			})
		}

		for _, v := range viper.GetStringSlice("project_ass_roles") {
			assignees = append(assignees, permissionv1.UserOrRole{
				Type:  permissionv1.RoleType,
				Value: v,
			})
		}

		for _, assignee := range assignees {
			for _, assignment := range viper.GetStringSlice("project_ass_assignments") {
				opt.Writes = append(opt.Writes, &permissionv1.ProjectAssignment{
					Assignee:   assignee,
					Assignment: permissionv1.ProjectAssignmentType(assignment),
				})
			}
		}

		_, err := c.PermissionV1().ProjectPermission().Update(cmd.Context(), project, &opt)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectAssignmentsCmd)

	projectAssignmentsCmd.AddCommand(getProjectAssignmentsCmd)
	projectAssignmentsCmd.AddCommand(wProjectAssignmentsCmd)

	getProjectAssignmentsCmd.Flags().StringSlice("relation", []string{}, fmt.Sprintf("filter by assignments, values can be: %v", permissionv1.ValidProjectAssignmentTypes))
	_ = viper.BindPFlag("project_relations", getProjectAssignmentsCmd.Flags().Lookup("relation"))

	wProjectAssignmentsCmd.Flags().StringSlice("user", []string{}, "Add user as an assignee")
	_ = viper.BindPFlag("project_ass_users", wProjectAssignmentsCmd.Flags().Lookup("user"))

	wProjectAssignmentsCmd.Flags().StringSlice("role", []string{}, "Add role as an assignee")
	_ = viper.BindPFlag("project_ass_roles", wProjectAssignmentsCmd.Flags().Lookup("role"))

	wProjectAssignmentsCmd.Flags().StringSlice("assignment", []string{}, fmt.Sprintf("Add assignment, values can be %s", permissionv1.ValidProjectAssignmentTypes))
	_ = viper.BindPFlag("project_ass_assignments", wProjectAssignmentsCmd.Flags().Lookup("assignment"))
}
