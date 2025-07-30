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
	"fmt"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// assignmentsRoleCmd represents the role assignments command
var assignmentsRoleCmd = &cobra.Command{
	Use:   "assignments",
	Short: "Interacts with role assignments",
}

// getAssignmentsRoleCmd represents the role assignments get command
var getAssignmentsRoleCmd = &cobra.Command{
	Use:   "get [flags] <role-id> [--relation <assignment> ...]",
	Short: "Get user and role assignments of a role. By default all assignments are returned",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := permissionv1.GetRoleAssignmentsOptions{}

		for _, v := range viper.GetStringSlice("role_relations") {
			opt.Relations = append(opt.Relations, permissionv1.RoleAssignmentType(v))
		}

		resp, _, err := c.PermissionV1().RolePermission().GetAssignments(cmd.Context(), args[0], &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// wRoleAssignmentsCmd represents the project assignments add command
var wRoleAssignmentsCmd = &cobra.Command{
	Use:   "add [flags] <role-id> --assignment <assignment> [--user <user> --role <role>]",
	Short: "add role assignments",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := permissionv1.UpdateRolePermissionsOptions{}
		assignees := []permissionv1.UserOrRole{}

		if len(viper.GetStringSlice("role_ass_assignments")) < 1 {
			return errors.New("you must set at lest one assignment")
		}

		if len(viper.GetStringSlice("role_ass_users")) < 1 && len(viper.GetStringSlice("role_ass_roles")) < 1 {
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
				opt.Writes = append(opt.Writes, &permissionv1.RoleAssignment{
					Assignee:   assignee,
					Assignment: permissionv1.RoleAssignmentType(assignment),
				})
			}
		}

		_, err := c.PermissionV1().RolePermission().Update(cmd.Context(), args[0], &opt)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	roleCmd.AddCommand(assignmentsRoleCmd)
	assignmentsRoleCmd.AddCommand(getAssignmentsRoleCmd)
	assignmentsRoleCmd.AddCommand(wRoleAssignmentsCmd)

	getAssignmentsRoleCmd.Flags().StringSlice("relation", []string{}, fmt.Sprintf("relations to be loaded. If not specified, all relations are returned. Values can be: %v", permissionv1.ValidRoleAssignmentTypes))
	_ = viper.BindPFlag("role_relations", getAssignmentsRoleCmd.Flags().Lookup("relation"))

	wRoleAssignmentsCmd.Flags().StringSlice("user", []string{}, "Add user as an assignee")
	_ = viper.BindPFlag("role_ass_users", wRoleAssignmentsCmd.Flags().Lookup("user"))

	wRoleAssignmentsCmd.Flags().StringSlice("role", []string{}, "Add role as an assignee")
	_ = viper.BindPFlag("role_ass_roles", wRoleAssignmentsCmd.Flags().Lookup("role"))

	wRoleAssignmentsCmd.Flags().StringSlice("assignment", []string{}, fmt.Sprintf("Add assignment, values can be %s", permissionv1.ValidRoleAssignmentTypes))
	_ = viper.BindPFlag("role_ass_assignments", wRoleAssignmentsCmd.Flags().Lookup("assignment"))
}
