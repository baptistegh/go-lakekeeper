package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewServerCmd(clientOptions *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:     "server",
		Aliases: []string{"srv"},
		Short:   "Manage server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewInfoCmd(clientOptions))
	command.AddCommand(NewBootstrapCmd(clientOptions))
	command.AddCommand(NewServerAccessCmd(clientOptions))
	command.AddCommand(NewServerAssignmentsCmd(clientOptions))
	command.AddCommand(NewServerGrantCmd(clientOptions))

	return &command
}

func NewBootstrapCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		asOperator       bool
		acceptTermsOfUse bool

		output string
	)

	command := cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstraps the server with the current user",
		Example: `  # Bootstrap the server and get the server admin role
  lkctl bootstrap --accept-terms-of-use

  # Bootstrap the server as an operator
  lkctl bootstrap --accept-terms-of-use --as-operator`,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			opt := managementv1.BootstrapServerOptions{
				AcceptTermsOfUse: acceptTermsOfUse,
				IsOperator:       &asOperator,
			}

			client := MustCreateClient(ctx, clientOpts).ServerV1()

			_, err := client.Bootstrap(ctx, &opt)
			errors.Check(err)

			switch output {
			case "json":
				info, _, err := client.Info(ctx)
				errors.Check(err)

				err = PrintResource(info, output)
				errors.Check(err)
			case "text":
				fmt.Println("Server bootstrapped successfully")
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	command.Flags().BoolVar(&asOperator, "as-operator", false, "Bootstrap the server as an operator")
	command.Flags().BoolVar(&acceptTermsOfUse, "accept-terms-of-use", false, "Accept the terms of use")

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewServerAssignmentsCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		assignmentsOpts assignmentsOpts
		output          string
	)

	command := cobra.Command{
		Use:   "assignments",
		Short: "Get server assignments",
		Example: `  # Get server assignments
  lkctl server assignments

  # Filter by assignment type
  lkctl server assignments --relations admin`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()

			var relations []permissionv1.ServerAssignmentType
			for _, v := range assignmentsOpts.relations {
				relations = append(relations, permissionv1.ServerAssignmentType(v))
			}

			opt := permissionv1.GetServerAssignmentsOptions{
				Relations: relations,
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).PermissionV1().ServerPermission().GetAssignments(ctx, &opt)
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

func NewServerGrantCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		users       []string
		roles       []string
		assignments []string
	)

	command := cobra.Command{
		Use:     "grant",
		Short:   "add server assignments",
		Aliases: []string{"assign"},
		Example: `  # Grant admin assignment to a user
  lkctl server grant --users 0198618c-5be8-7a82-a0b9-1076c9dd12f0 --assignments admin

  # Grant operator assignment to a role
  lkctl server grant --roles 0198618c-5be8-7a82-a0b9-1076c9dd12f0 --assignments operator`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			opt := permissionv1.UpdateServerPermissionsOptions{}
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
					opt.Writes = append(opt.Writes, &permissionv1.ServerAssignment{
						Assignee:   assignee,
						Assignment: permissionv1.ServerAssignmentType(assignment),
					})
				}
			}

			ctx := cmd.Context()
			c := MustCreateClient(ctx, clientOpts).PermissionV1().ServerPermission()

			_, err := c.Update(ctx, &opt)
			errors.Check(err)

			fmt.Println("Server permissions updated")
		},
	}

	command.Flags().StringSliceVar(&users, "users", []string{}, "Grant access to users; can be repeated multiple times to add multiple users")
	command.Flags().StringSliceVar(&roles, "roles", []string{}, "Grant access to roles; can be repeated multiple times to add multiple roles")
	command.Flags().StringSliceVar(&assignments, "assignments", []string{}, "Assignments to use; can be repeated multiple times to add multiple assignments")

	err := command.MarkFlagRequired("assignments")
	errors.Check(err)

	return &command
}

func NewInfoCmd(clientOptions *clientOptions) *cobra.Command {
	var output string

	command := cobra.Command{
		Use:   "info",
		Short: "Print server informations",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			resp, _, err := MustCreateClient(ctx, clientOptions).ServerV1().Info(ctx)
			errors.Check(err)

			switch output {
			case "text":
				fmt.Printf("ID: %s\n", resp.ServerID)
				fmt.Printf("Version: %s\n", resp.Version)
				fmt.Printf("Default Project ID: %s\n", resp.DefaultProjectID)
				fmt.Printf("Bootstraped: %t\n", resp.Bootstrapped)
				fmt.Printf("Authorization Backend: %s\n", resp.AuthzBackend)
				fmt.Printf("AWS System Identities Enabled: %t\n", resp.AWSSystemIdentitiesEnabled)
				fmt.Printf("Azure System Identities Enabled: %t\n", resp.AzureSystemIdentitiesEnabled)
				fmt.Printf("GCP System Identities Enableds: %t\n", resp.GCPSystemIdentitiesEnabled)
				fmt.Println("Queues:")
				for _, q := range resp.Queues {
					fmt.Printf("  %s\n", q)
				}
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Printf("unknown output format: %s\n", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewServerAccessCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		accessOpts accessOpts
		output     string
	)

	command := cobra.Command{
		Use:   "access",
		Short: "Get server access",
		Long:  "Get server-level access. By default, the current user's access is returned.",
		Example: `  # Get current user's server access
  lkctl server access

  # Get server access for a specific user
  lkctl server access --user 0198618c-5be8-7a82-a0b9-1076c9dd12f0

  # Get server access for a specific role
  lkctl server access --role 0198618c-5be8-7a82-a0b9-1076c9dd12f0`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()

			if accessOpts.role != "" && accessOpts.user != "" {
				log.Fatal("you only can filter by user OR role, both were supplied")
			}

			opt := permissionv1.GetServerAccessOptions{}
			if accessOpts.user != "" {
				opt.PrincipalUser = core.Ptr(accessOpts.user)
			}
			if accessOpts.role != "" {
				opt.PrincipalRole = core.Ptr(accessOpts.role)
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).PermissionV1().ServerPermission().GetAccess(ctx, &opt)
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
