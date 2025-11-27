package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewUserCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "user",
		Short: "Manage users",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
			os.Exit(1)
		},
	}

	command.AddCommand(NewUserListCmd(clientOpts))
	command.AddCommand(NewUserGetCmd(clientOpts))
	command.AddCommand(NewUserDeleteCmd(clientOpts))
	command.AddCommand(NewUserCreateCmd(clientOpts))

	return &command
}

func NewUserListCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		listOpts listOpts

		output string
	)

	command := cobra.Command{
		Use:     "list",
		Short:   "List users",
		Aliases: []string{"ls"},
		Example: `  # List users
  lkctl user ls`,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			opt := managementv1.ListUsersOptions{
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

			resp, _, err := MustCreateClient(ctx, clientOpts).UserV1().List(ctx, &opt)
			errors.Check(err)

			switch output {
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			case "text", "wide":
				printUsers(output, resp.NextPageToken, resp.Users...)
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	AddListFlags(&command, &listOpts)
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")

	return &command
}

func NewUserGetCmd(clientOpts *clientOptions) *cobra.Command {
	var output string

	command := cobra.Command{
		Use:   "get USERID",
		Short: "Get a user by id",
		Example: `  # Get a user by its id
  lkctl user get oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()
			resp, _, err := MustCreateClient(ctx, clientOpts).UserV1().Get(ctx, args[0])
			errors.Check(err)

			switch output {
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			case "text", "wide":
				printUsers(output, nil, resp)
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text|wide")

	return &command
}

func NewUserDeleteCmd(clientOpts *clientOptions) *cobra.Command {
	command := cobra.Command{
		Use:   "delete USERID",
		Short: "Delete a user by id",
		Example: `  # Delete a user
  lkctl user delete oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()
			_, err := MustCreateClient(ctx, clientOpts).UserV1().Delete(ctx, args[0])
			errors.Check(err)

			fmt.Printf("User %s deleted\n", args[0])
		},
	}

	return &command
}

func NewUserCreateCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		output string

		email  string
		update bool
	)

	command := cobra.Command{
		Use:     "create USERID NAME USERTYPE",
		Short:   "Create a new user",
		Aliases: []string{"add"},
		Example: `  # Create a new human user authenticated from OIDC
  lkctl user create oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6 "Peter Cold" human
  
  # Create an application user from kubernetes
  lkctl user create kubernetes~d223d88c-85b6-4859-b5c5-27f3825e47f6 "Service Account" application

  # Create a user with an email
  lkctl user create oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6 "Peter Cold" human --email peter.cold@example.com`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx := cmd.Context()
			opt := managementv1.ProvisionUserOptions{
				ID:       core.Ptr(args[0]),
				Name:     core.Ptr(args[1]),
				UserType: core.Ptr(managementv1.UserType(args[2])),
			}

			if email != "" {
				opt.Email = core.Ptr(email)
			}

			if update {
				opt.UpdateIfExists = core.Ptr(update)
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).UserV1().Provision(ctx, &opt)
			errors.Check(err)

			switch output {
			case "text":
				fmt.Printf("User %s registered with id %s\n", args[1], resp.ID)
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text")
	command.Flags().StringVar(&email, "email", "", "Add an email to the user")
	command.Flags().BoolVar(&update, "update", false, "Update the user if exists")

	return &command
}

func printUsers(output string, nextPageToken *string, users ...*managementv1.User) {
	switch output {
	case "text":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tNAME\tEMAIL\tUSER TYPE\n")
		for _, u := range users {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", u.ID, u.Name, FormatPString(u.Email), u.UserType)
		}
		w.Flush()
		if nextPageToken != nil {
			fmt.Printf("\nNext page token: %s\n", *nextPageToken)
		}
	case "wide":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID\tNAME\tEMAIL\tUSER TYPE\tCREATED AT\tUPDATED AT\tLAST UPDATED WITH\n")
		for _, u := range users {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", u.ID, u.Name, FormatPString(u.Email), u.UserType, u.CreatedAt, FormatPString(u.UpdatedAt), u.LastUpdatedWith)
		}
		w.Flush()
		if nextPageToken != nil {
			fmt.Printf("\nNext page token: %s\n", *nextPageToken)
		}
	}
}
