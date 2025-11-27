package commands

import (
	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewWhoamiCmd(clientOptions *clientOptions) *cobra.Command {
	var output string

	command := cobra.Command{
		Use:   "whoami",
		Short: "Print the current user",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()
			resp, _, err := MustCreateClient(ctx, clientOptions).UserV1().Whoami(ctx)
			errors.Check(err)

			switch output {
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			case "test", "wide":
				printUsers(output, nil, resp)
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text|wide")

	return &command
}
