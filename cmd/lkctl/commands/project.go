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

func NewProjectCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:     "project",
		Aliases: []string{"proj"},
		Short:   "Manage projects",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewProjectListCmd(clientOpts))
	command.AddCommand(NewProjectGetCmd(clientOpts))
	command.AddCommand(NewProjectCreateCmd(clientOpts))
	command.AddCommand(NewProjectRenameCmd(clientOpts))
	command.AddCommand(NewProjectDeleteCmd(clientOpts))
	command.AddCommand(NewProjectAccessCmd(clientOpts))
	command.AddCommand(NewProjectAssignmentsCmd(clientOpts))
	command.AddCommand(NewProjectGrantCmd(clientOpts))

	return &command
}

func NewProjectListCmd(clientOpts *clientOptions) *cobra.Command {
	var output string

	command := cobra.Command{
		Use:   "list",
		Short: "List all the available projects for the current user",
		Example: `  # List all the available projects for the current user
  lkctl project list`,
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			resp, _, err := MustCreateClient(ctx, clientOpts).ProjectV1().List(ctx)
			errors.Check(err)

			switch output {
			case "text":
				if len(resp.Projects) == 0 {
					fmt.Println("No projects available")
					return
				}
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprint(w, "ID\tNAME\n")
				for _, p := range resp.Projects {
					fmt.Fprintf(w, "%s\t%s\n", p.ID, p.Name)
				}
				w.Flush()
			case "json":
				err := PrintResource(resp.Projects, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format %s\n", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewProjectGetCmd(clientOpts *clientOptions) *cobra.Command {
	var output string
	command := cobra.Command{
		Use:   "get PROJECT-ID",
		Short: "Get a project by id",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()
			resp, _, err := MustCreateClient(ctx, clientOpts).ProjectV1().Get(ctx, args[0])
			errors.Check(err)

			switch output {
			case "text":
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprint(w, "ID\tNAME\n")
				fmt.Fprintf(w, "%s\t%s\n", resp.ID, resp.Name)
				w.Flush()
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format %s\n", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewProjectCreateCmd(clientOpts *clientOptions) *cobra.Command {
	var output string
	command := cobra.Command{
		Use:   "create PROJECTNAME",
		Short: "Create a new project",
		Example: `  # Create a new project
  lkctl project create "New Project"`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			err := createProject(cmd.Context(), clientOpts, args[0], output)
			errors.Check(err)
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewProjectDeleteCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "delete PROJECT-ID",
		Short: "Delete a project",
		Example: `  # Delete project
  lkctl project delete 019861a0-6d4e-7bf3-96c6-9aef2d4a2749`,
		Aliases: []string{"rm"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			_, err := MustCreateClient(cmd.Context(), clientOpts).ProjectV1().Delete(cmd.Context(), args[0])
			errors.Check(err)

			fmt.Printf("Project %s deleted\n", args[0])
		},
	}

	return &command
}

func NewProjectRenameCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "rename PROJECT-ID NEW-NAME",
		Short: "Rename a project",
		Example: `  # Rename a project
  lkctl project rename 01986184-3cb1-7526-a98c-72fecfe97731 "New Project Name"`,
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

			fmt.Printf("Project %s renamed to %s\n", args[0], args[1])
		},
	}
	return &command
}

func NewProjectAccessCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		accessOpts accessOpts

		output string
	)

	command := cobra.Command{
		Use:   "access PROJECT-ID",
		Short: "Get project access",
		Long:  "Get project access. By default, current user's access is returned",
		Example: `  # Get default project access
  lkctl project access

  # Get specific project access
  lkctl project access 01986184-3cb1-7526-a98c-72fecfe97731

  # Get project access for a specific user
  lkctl project access 01986184-3cb1-7526-a98c-72fecfe97731 --user oidc~0198618c-5be8-7a82-a0b9-1076c9dd12f0

  # Get project access for a specific role
  lkctl project access 01986184-3cb1-7526-a98c-72fecfe97731 --role oidc~0198618c-5be8-7a82-a0b9-1076c9dd12f0`,
		Run: func(cmd *cobra.Command, args []string) {
			var project string
			if len(args) != 1 {
				project = uuid.Nil.String()
			} else {
				project = args[0]
			}

			ctx := cmd.Context()

			if accessOpts.role != "" && accessOpts.user != "" {
				log.Fatal("you only can filter by user OR role, both were supplied")
			}

			opt := permissionv1.GetProjectAccessOptions{}

			if accessOpts.user != "" {
				opt.PrincipalUser = core.Ptr(accessOpts.user)
			}

			if accessOpts.role != "" {
				opt.PrincipalRole = core.Ptr(accessOpts.role)
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).PermissionV1().ProjectPermission().GetAccess(ctx, project, &opt)
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

func NewProjectAssignmentsCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		assignmentsOpts assignmentsOpts

		output string
	)

	command := cobra.Command{
		Use:   "assignments PROJECT-ID",
		Short: "Get project assignments",
		Example: `  # Get default project assignments
  lkctl project assignments

  # Filter by assignment type
  lkctl project assignments 01986184-3cb1-7526-a98c-72fecfe97731 --relations project_admin

  # Filter by multiple assignment types
  lkctl project assignments 01986184-3cb1-7526-a98c-72fecfe97731 --relations project_admin --relations select`,
		Run: func(cmd *cobra.Command, args []string) {
			var project string
			if len(args) != 1 {
				project = uuid.Nil.String()
			} else {
				project = args[0]
			}

			ctx := cmd.Context()

			var relations []permissionv1.ProjectAssignmentType
			for _, v := range assignmentsOpts.relations {
				relations = append(relations, permissionv1.ProjectAssignmentType(v))
			}

			opt := permissionv1.GetProjectAssignmentsOptions{
				Relations: relations,
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).PermissionV1().ProjectPermission().GetAssignments(ctx, project, &opt)
			errors.Check(err)

			switch output {
			case "text":
				PrintAssignments(resp.Assignments...)
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format %s\n", output)
			}
		},
	}

	AddAssignmentsFlags(&command, &assignmentsOpts)
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewProjectGrantCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		users []string
		roles []string

		assignments []string
	)

	command := cobra.Command{
		Use:     "grant PROJECT-ID",
		Short:   "add project assignments",
		Aliases: []string{"assign"},
		Run: func(cmd *cobra.Command, args []string) {
			var project string
			if len(args) != 1 {
				project = uuid.Nil.String()
			} else {
				project = args[0]
			}

			opt := permissionv1.UpdateProjectPermissionsOptions{}
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
					opt.Writes = append(opt.Writes, &permissionv1.ProjectAssignment{
						Assignee:   assignee,
						Assignment: permissionv1.ProjectAssignmentType(assignment),
					})
				}
			}

			ctx := cmd.Context()
			c := MustCreateClient(ctx, clientOpts).PermissionV1().ProjectPermission()

			_, err := c.Update(cmd.Context(), project, &opt)
			errors.Check(err)

			fmt.Println("Project permissions updated")
		},
	}

	command.Flags().StringSliceVar(&users, "users", []string{}, "Grant access to users; can be repeated multiple times to add multiple users")
	command.Flags().StringSliceVar(&roles, "roles", []string{}, "Grant access to roles; can be repeated multiple times to add multiple roles")
	command.Flags().StringSliceVar(&assignments, "assignments", []string{}, "Assignments to use; can be repeated multiple times to add multiple assignments")

	err := command.MarkFlagRequired("assignments")
	errors.Check(err)

	return &command
}

func createProject(ctx context.Context, clientOpts *clientOptions, name, output string) error {
	opt := managementv1.CreateProjectOptions{
		Name: name,
	}

	c := MustCreateClient(ctx, clientOpts).ProjectV1()

	switch output {
	case "text":
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
