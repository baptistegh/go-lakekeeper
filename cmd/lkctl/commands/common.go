package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/spf13/cobra"
)

const (
	cliName = "lkctl"
)

type clientOptions struct {
	server       string
	authURL      string
	clientID     string
	clientSecret string
	scope        []string
	boostrap     bool
	debug        bool
}

type accessOpts struct {
	user string
	role string
}

type assignmentsOpts struct {
	relations []string
}

type listOpts struct {
	limit int64
	token string
	name  string
}

// PrintResource prints a single resource in YAML or JSON format to stdout according to the output format
func PrintResource(resource any, output string) error {
	switch output {
	case "json":
		jsonBytes, err := json.MarshalIndent(resource, "", "  ")
		if err != nil {
			return fmt.Errorf("unable to marshal resource to json: %w", err)
		}
		fmt.Println(string(jsonBytes))
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}
	return nil
}

func PrintAssignments[T permissionv1.Assignment](assignments ...T) {
	if len(assignments) == 0 {
		fmt.Println("No assignments")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "PRINCIPAL TYPE\tPRINCIPAL ID\tASSIGNMENT\n")
	for _, a := range assignments {
		fmt.Fprintf(w, "%s\t%s\t%s\n", a.GetPrincipalType(), a.GetPrincipalID(), a.GetAssignment())
	}
	w.Flush()
}

func FormatPString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func AddAccessFlags(cmd *cobra.Command, opts *accessOpts) {
	cmd.Flags().StringVar(&opts.user, "user", "", "Filter by user")
	cmd.Flags().StringVar(&opts.role, "role", "", "Filter by role")
}

func AddAssignmentsFlags(cmd *cobra.Command, opts *assignmentsOpts) {
	cmd.Flags().StringSliceVar(&opts.relations, "relations", []string{}, "Filter by relations. (Can be repeated multiple times to add multiple relations, also supports comma separated relations)")
}

func AddListFlags(cmd *cobra.Command, opts *listOpts) {
	cmd.Flags().Int64Var(&opts.limit, "limit", int64(100), "Signals an upper bound of the number of results that the client will receive")
	cmd.Flags().StringVar(&opts.token, "token", "", "Pagination token")
	cmd.Flags().StringVar(&opts.name, "name", "", "Filter by name")
}

func PrintNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
