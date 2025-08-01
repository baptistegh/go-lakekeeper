// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRoleCmd(clientOptions *clientOptions) *cobra.Command {
	var project string

	command := cobra.Command{
		Use:   "role",
		Short: "Manage roles",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.PersistentFlags().StringVarP(&project, "project", "p", uuid.Nil.String(), "Select a project")

	command.AddCommand(NewRoleListCmd(clientOptions, &project))
	command.AddCommand(NewRoleGetCmd(clientOptions, &project))
	command.AddCommand(NewRoleCreateCmd(clientOptions, &project))
	command.AddCommand(NewRoleDeleteCmd(clientOptions, &project))
	command.AddCommand(NewRoleUpdateCmd(clientOptions, &project))
	command.AddCommand(NewRoleAccessCmd(clientOptions, &project))
	command.AddCommand(NewRoleAssignmentsCmd(clientOptions, &project))
	command.AddCommand(NewRoleGrantCmd(clientOptions, &project))

	return &command
}

func NewRoleListCmd(clientOptions *clientOptions, project *string) *cobra.Command {
	var (
		listOpts listOpts

		output string
	)

	command := cobra.Command{
		Use:     "list",
		Short:   "List available roles",
		Aliases: []string{"ls"},
		Example: `  # List available roles
  lkctl role ls`,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			opt := managementv1.ListRolesOptions{
				ProjectID: project,
				ListOptions: managementv1.ListOptions{
					PageSize: core.Ptr(listOpts.limit),
				},
			}

			if listOpts.token != "" {
				opt.PageToken = core.Ptr(listOpts.token)
			}

			if listOpts.name != "" {
				opt.Name = core.Ptr(listOpts.name)
			}

			resp, _, err := MustCreateClient(ctx, clientOptions).RoleV1(*project).List(ctx, &opt)
			errors.Check(err)

			switch output {
			case "text", "wide":
				if len(resp.Roles) == 0 {
					fmt.Println("No roles available")
					return
				}
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprint(w, "ID\tNAME\tPROJECT ID\tCREATED AT\tUPDATED AT")
				if output == "wide" {
					fmt.Fprint(w, "\tDESCRIPTION")
				}
				fmt.Fprint(w, "\n")
				for _, r := range resp.Roles {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s", r.ID, r.Name, r.ProjectID, r.CreatedAt, FormatPString(r.UpdatedAt))
					if output == "wide" {
						fmt.Fprintf(w, "\t%s", FormatPString(r.Description))
					}
					fmt.Fprint(w, "\n")
				}
				w.Flush()

				fmt.Println()

				if resp.NextPageToken != nil {
					fmt.Printf("Next page token: %s\n", *resp.NextPageToken)
				}
			case "json":
				err := PrintResource(resp.Roles, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format %s\n", output)
			}
		},
	}

	AddListFlags(&command, &listOpts)
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text|wide")

	return &command
}

func NewRoleGetCmd(clientOptions *clientOptions, project *string) *cobra.Command {
	var output string

	command := cobra.Command{
		Use:   "get",
		Short: "Get a role by id",
		Example: `  # Get a role by id
  lkctl role get 01986184-3cb1-7526-a98c-72fecfe97731`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			resp, _, err := MustCreateClient(ctx, clientOptions).RoleV1(*project).Get(ctx, args[0])
			errors.Check(err)

			switch output {
			case "text", "wide":
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprint(w, "ID\tNAME\tPROJECT ID\tCREATED AT\tUPDATED AT")
				if output == "wide" {
					fmt.Fprint(w, "\tDESCRIPTION")
				}
				fmt.Fprintf(w, "\n")
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s", resp.ID, resp.Name, resp.ProjectID, resp.CreatedAt, FormatPString(resp.UpdatedAt))
				if output == "wide" {
					fmt.Fprintf(w, "\t%s", PrintNil(resp.Description))
				}
				fmt.Fprintf(w, "\n")
				w.Flush()
			case "json":
				err = PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format %s\n", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text|wide")

	return &command
}

func NewRoleCreateCmd(clientOpts *clientOptions, project *string) *cobra.Command {
	var (
		output      string
		description string
	)
	command := cobra.Command{
		Use:     "create ROLENAME",
		Short:   "Create a new role",
		Aliases: []string{"add"},
		Example: `  # Create a new role
  lkctl role create "New Role"

  # Create a new role with a description
  lkctl role create "New Role" --description "With a description"`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			err := createRole(cmd.Context(), clientOpts, args[0], *project, description, output)
			errors.Check(err)
		},
	}

	command.Flags().StringVar(&description, "description", "", "Add a description to the role")
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewRoleDeleteCmd(clientOpts *clientOptions, project *string) *cobra.Command {
	command := cobra.Command{
		Use:     "delete ROLEID",
		Short:   "Delete a role by id",
		Aliases: []string{"rm"},
		Example: `  # Delete role
  lkctl role rm 01986184-3cb1-7526-a98c-72fecfe97731`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()
			_, err := MustCreateClient(ctx, clientOpts).RoleV1(*project).Delete(ctx, args[0])
			errors.Check(err)

			fmt.Printf("Role %s deleted\n", args[0])
		},
	}

	return &command
}

func NewRoleUpdateCmd(clientOpts *clientOptions, project *string) *cobra.Command {
	var (
		output      string
		description string
	)
	command := cobra.Command{
		Use:   "update ROLEID ROLENAME",
		Short: "Update role",
		Example: `  # Update role
  lkctl role update 01986184-3cb1-7526-a98c-72fecfe97731 "Updated Name" --description "Updated Description"

  # Update role and delete its description
  lkctl role update 01986184-3cb1-7526-a98c-72fecfe97731 "Updated Name"`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			err := updateRole(cmd.Context(), clientOpts, args[0], *project, args[1], description, output)
			errors.Check(err)
		},
	}

	command.Flags().StringVar(&description, "description", "", "Add a description to the role")
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewRoleAccessCmd(clientOpts *clientOptions, _ *string) *cobra.Command {
	var (
		accessOpts accessOpts

		output string
	)

	command := cobra.Command{
		Use:   "access ROLEID",
		Short: "Get role access",
		Long:  "Get role access. By default, current user's access is returned",
		Example: `  # Get role access
  lkctl role access 01986184-3cb1-7526-a98c-72fecfe97731

  # Get role access for a specific user
  lkctl role access 01986184-3cb1-7526-a98c-72fecfe97731 --user oidc~0198618c-5be8-7a82-a0b9-1076c9dd12`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()
			if accessOpts.role != "" && accessOpts.user != "" {
				log.Fatal("you only can filter by user OR role, both were supplied")
			}

			opt := permissionv1.GetRoleAccessOptions{}

			if accessOpts.user != "" {
				opt.PrincipalUser = core.Ptr(accessOpts.user)
			}

			if accessOpts.role != "" {
				opt.PrincipalRole = core.Ptr(accessOpts.role)
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).PermissionV1().RolePermission().GetAccess(ctx, args[0], &opt)
			errors.Check(err)

			switch output {
			case "text":
				if len(resp.AllowedActions) == 0 {
					fmt.Println("No access")
					return
				}
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintf(w, "ALLOWED ACTIONS\n")
				for _, a := range resp.AllowedActions {
					fmt.Fprintf(w, "%s\n", a)
				}
				w.Flush()
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format %s\n", output)
			}
		},
	}

	AddAccessFlags(&command, &accessOpts)
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewRoleAssignmentsCmd(clientOpts *clientOptions, _ *string) *cobra.Command {
	var (
		assignmentsOpts assignmentsOpts

		output string
	)

	command := cobra.Command{
		Use:   "assignments ROLEID",
		Short: "Get role assignments",
		Example: `  # Get default role assignments
  lkctl role assignments 01986184-3cb1-7526-a98c-72fecfe97731

  # Filter by assignment type
  lkctl role assignments 01986184-3cb1-7526-a98c-72fecfe97731 --relations ownership

  # Filter by multiple assignment types
  lkctl role assignments 01986184-3cb1-7526-a98c-72fecfe97731 --relations ownership --relations assignee`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()

			var relations []permissionv1.RoleAssignmentType
			for _, v := range assignmentsOpts.relations {
				relations = append(relations, permissionv1.RoleAssignmentType(v))
			}

			opt := permissionv1.GetRoleAssignmentsOptions{
				Relations: relations,
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).PermissionV1().RolePermission().GetAssignments(ctx, args[0], &opt)
			errors.Check(err)

			switch output {
			case "text":
				PrintAssignments(resp.Assignments...)
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			}
		},
	}

	AddAssignmentsFlags(&command, &assignmentsOpts)
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewRoleGrantCmd(clientOpts *clientOptions, _ *string) *cobra.Command {
	var (
		users []string
		roles []string

		assignments []string
	)

	command := cobra.Command{
		Use:     "grant ROLEID",
		Short:   "add role assignments",
		Aliases: []string{"assign"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			opt := permissionv1.UpdateRolePermissionsOptions{}
			assignees := []permissionv1.UserOrRole{}
			if len(assignments) < 1 {
				log.Fatal("you must set at lest one assignment")
			}
			if len(users) < 1 && len(roles) < 1 {
				log.Fatal("you must set at least one user or role")
			}
			for _, v := range users {
				assignees = append(assignees, permissionv1.UserOrRole{
					Type:  permissionv1.UserType,
					Value: v,
				})
			}

			for _, v := range roles {
				assignees = append(assignees, permissionv1.UserOrRole{
					Type:  permissionv1.RoleType,
					Value: v,
				})
			}

			for _, assignee := range assignees {
				for _, assignment := range assignments {
					opt.Writes = append(opt.Writes, &permissionv1.RoleAssignment{
						Assignee:   assignee,
						Assignment: permissionv1.RoleAssignmentType(assignment),
					})
				}
			}

			ctx := cmd.Context()
			c := MustCreateClient(ctx, clientOpts).PermissionV1().RolePermission()

			_, err := c.Update(cmd.Context(), args[0], &opt)
			errors.Check(err)

			fmt.Println("Role permissions updated")
		},
	}

	command.Flags().StringSliceVar(&users, "users", []string{}, "Grant access to users; can be repeated multiple times to add multiple users")
	command.Flags().StringSliceVar(&roles, "roles", []string{}, "Grant access to roles; can be repeated multiple times to add multiple roles")
	command.Flags().StringSliceVar(&assignments, "assignments", []string{}, "Assignments to use; can be repeated multiple times to add multiple assignments")

	err := command.MarkFlagRequired("assignments")
	errors.Check(err)

	return &command
}

func createRole(ctx context.Context, clientOpts *clientOptions, name, project, description, output string) error {
	opt := managementv1.CreateRoleOptions{
		Name:      name,
		ProjectID: core.Ptr(project),
	}

	if description != "" {
		opt.Description = core.Ptr(description)
	}

	resp, _, err := MustCreateClient(ctx, clientOpts).RoleV1(project).Create(ctx, &opt)
	if err != nil {
		return err
	}

	switch output {
	case "text":
		fmt.Printf("Role %s created with id %s\n", name, resp.ID)
		return nil
	case "json":
		return PrintResource(resp, output)
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}
}

func updateRole(ctx context.Context, clientOpts *clientOptions, id, project, name, description, output string) error {
	opt := managementv1.UpdateRoleOptions{
		Name: name,
	}

	if description != "" {
		opt.Description = core.Ptr(description)
	}

	resp, _, err := MustCreateClient(ctx, clientOpts).RoleV1(project).Update(ctx, id, &opt)
	if err != nil {
		return err
	}

	switch output {
	case "text":
		fmt.Printf("Role %s updated\n", id)
		return nil
	case "json":
		return PrintResource(resp, output)
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}
}
