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

package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"

	"github.com/spf13/cobra"
)

func NewProjectCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "project",
		Short: "Interacts with projects",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewProjectListCmd(clientOpts))
	command.AddCommand(NewProjectGetCmd(clientOpts))
	command.AddCommand(NewProjectAddCmd(clientOpts))
	command.AddCommand(NewProjectRenameCmd(clientOpts))

	return &command
}

func NewProjectListCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "list",
		Short: "List all the available projects for the current user",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			resp, _, err := MustCreateClient(ctx, clientOpts).ProjectV1().List(ctx)
			errors.Check(err)

			err = json.NewEncoder(cmd.OutOrStdout()).Encode(resp.Projects)
			errors.Check(err)
		},
	}

	return &command
}

func NewProjectGetCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "get PROJECT_ID",
		Short: "Get project information",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()
			resp, _, err := MustCreateClient(ctx, clientOpts).ProjectV1().Get(ctx, args[0])
			errors.Check(err)

			err = json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
			errors.Check(err)
		},
	}

	return &command
}

func NewProjectAddCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		output string
	)
	command := cobra.Command{
		Use:   "create PROJECT",
		Short: "Create a new project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
			}

			err := createProject(cmd.Context(), clientOpts, args[0], output)
			errors.Check(err)
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewProjectRenameCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "rename PROJECT_ID NEW_NAME",
		Short: "Rename a project",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			if len(args) != 2 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			opt := managementv1.RenameProjectOptions{
				NewName: args[1],
			}

			_, err := MustCreateClient(ctx, clientOpts).ProjectV1().Rename(cmd.Context(), args[0], &opt)
			errors.Check(err)

			fmt.Printf("Project %s renamed\n", args[0])
		},
	}
	return &command
}

// accessProjectCmd represents the project access command
/*
var accessProjectCmd = &cobra.Command{
	Use:   "access [flags] [--user <user> | --role <role>]",
	Short: "Get project access",
	Long:  "Get project access. By default, current user's access is returned",
	RunE: func(cmd *cobra.Command, _ []string) error {
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
}*/

func createProject(ctx context.Context, clientOpts *clientOptions, name, output string) error {
	opt := managementv1.CreateProjectOptions{
		Name: name,
	}

	c := MustCreateClient(ctx, clientOpts).ProjectV1()

	switch output {
	case "wide":
		resp, _, err := c.Create(ctx, &opt)
		if err != nil {
			return err
		}

		fmt.Printf("Project %s created with id %s\n", name, resp.ID)
	case "json", "yaml":
		resp, _, err := c.Create(ctx, &opt)
		if err != nil {
			return err
		}

		project, _, err := c.Get(ctx, resp.ID)
		if err != nil {
			return err
		}

		return PrintResource(project, output)
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}

	return nil
}
