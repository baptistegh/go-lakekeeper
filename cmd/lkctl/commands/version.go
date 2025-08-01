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
	"log"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	"github.com/baptistegh/go-lakekeeper/pkg/version"

	"github.com/spf13/cobra"
)

func NewVersionCmd(clientOpts *clientOptions) *cobra.Command {
	var (
		short  bool
		client bool
		output string
	)

	command := cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Example: `  # Print the full version of client and server to stdout
  lkctl version

  # Print only full version of the client - no connection to server will be made
  lkctl version --client

  # Print the full version of client and server in JSON format
  lkctl version

  # Print only client and server core version strings in YAML format
  lkctl version --short`,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			cv := version.GetVersion()
			switch output {
			case "json":
				v := make(map[string]any)

				if short {
					v["client"] = map[string]string{cliName: cv.Version}
				} else {
					v["client"] = cv
				}

				if !client {
					sv := getServerVersion(ctx, clientOpts)

					if short {
						v["server"] = map[string]string{"lakekeeper": sv}
					} else {
						v["server"] = sv
					}
				}

				err := PrintResource(v, output)
				errors.Check(err)
			case "text", "short", "":
				fmt.Fprint(cmd.OutOrStdout(), printClientVersion(&cv, short || (output == "short")))
				if !client {
					sv := getServerVersion(ctx, clientOpts)
					fmt.Fprint(cmd.OutOrStdout(), printServerVersion(sv))
				}
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text|short")
	command.Flags().BoolVar(&short, "short", false, "print just the version number")
	command.Flags().BoolVar(&client, "client", false, "client version only (no server required)")
	return &command
}

func getServerVersion(ctx context.Context, opts *clientOptions) string {
	info, _, err := MustCreateClient(ctx, opts).ServerV1().Info(ctx)
	errors.Check(err)

	return info.Version
}

func printClientVersion(v *version.Version, short bool) string {
	if short {
		return fmt.Sprintf("lkctl: %s\n", v.Version)
	}

	output := fmt.Sprintf("%s: %s\n", cliName, v)

	output += fmt.Sprintf("  BuildDate: %s\n", v.BuildDate)
	output += fmt.Sprintf("  GitCommit: %s\n", v.GitCommit)
	output += fmt.Sprintf("  GitTreeState: %s\n", v.GitTreeState)
	if v.GitTag != "" {
		output += fmt.Sprintf("  GitTag: %s\n", v.GitTag)
	}
	output += fmt.Sprintf("  GoVersion: %s\n", v.GoVersion)
	output += fmt.Sprintf("  Compiler: %s\n", v.Compiler)
	output += fmt.Sprintf("  Platform: %s\n", v.Platform)

	return output
}

func printServerVersion(v string) string {
	return fmt.Sprintf("%s: %s\n", "lakekeeper", v)
}
