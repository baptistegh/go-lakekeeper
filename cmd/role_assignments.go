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

func init() {
	roleCmd.AddCommand(assignmentsRoleCmd)
	assignmentsRoleCmd.AddCommand(getAssignmentsRoleCmd)

	getAssignmentsRoleCmd.Flags().StringSlice("relation", []string{}, fmt.Sprintf("relations to be loaded. If not specified, all relations are returned. Values can be: %v", permissionv1.ValidRoleAssignmentTypes))
	_ = viper.BindPFlag("role_relations", getAssignmentsRoleCmd.Flags().Lookup("relation"))
}
