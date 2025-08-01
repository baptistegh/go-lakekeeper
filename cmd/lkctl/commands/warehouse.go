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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/baptistegh/go-lakekeeper/cmd/lkctl/errors"
	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	profilev1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewWarehouseCmd(clientOpts *clientOptions) *cobra.Command {
	var project string

	command := cobra.Command{
		Use:     "warehouse",
		Aliases: []string{"wh"},
		Short:   "Manage warehouses",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.PersistentFlags().StringVarP(&project, "project", "p", uuid.Nil.String(), "Select a project")

	command.AddCommand(NewWarehouseListCmd(clientOpts, &project))
	command.AddCommand(NewWarehouseGetCmd(clientOpts, &project))
	command.AddCommand(NewWarehouseCreateCmd(clientOpts, &project))
	command.AddCommand(NewWarehouseDeleteCmd(clientOpts, &project))

	return &command
}

func NewWarehouseCreateCmd(clientOpts *clientOptions, project *string) *cobra.Command {
	var config string

	command := cobra.Command{
		Use:   "create WAREHOUSENAME -f JSONCONFIGFILE",
		Short: "Create a new warehouse",
		Example: `  # Create a warehouse from file
  lkctl warehouse create "New Warehouse" -f warehouse-config.json
  
  # Create a warehouse from stdin
  cat warehouse-config.json | lkctl warehouse create "New Warehouse" -f -`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			if len(args) != 1 || config == "" {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			var reader io.Reader

			if config == "-" {
				reader = cmd.InOrStdin()
			} else {
				file, err := os.Open(config)
				errors.Check(err)
				defer file.Close()

				reader = file
			}

			var opt managementv1.CreateWarehouseOptions

			err := json.NewDecoder(reader).Decode(&opt)
			errors.Check(err)

			if opt.Name != args[0] {
				log.Fatal("Warehouse name provided in config does not match the name supplied as argument")
			}

			//nolint:staticcheck // project id needs to be remove from the API first
			if opt.ProjectID == nil {
				opt.ProjectID = project //nolint:staticcheck // project id needs to be remove from the API first
			}

			//nolint:staticcheck // project id needs to be remove from the API first
			if *project != *opt.ProjectID {
				log.Fatal("Project ID provided in config does not match the project ID supplied as argument")
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).WarehouseV1(*project).Create(cmd.Context(), &opt)
			errors.Check(err)

			fmt.Printf("Warehouse %s created with id %s\n", opt.Name, resp.ID)
		},
	}

	command.Flags().StringVarP(&config, "file", "f", "", "Warehouse config file. JSON file or '-' for stdin")

	return &command
}

func NewWarehouseListCmd(clientOpts *clientOptions, project *string) *cobra.Command {
	var (
		status []string

		output string
	)

	command := cobra.Command{
		Use:     "list",
		Short:   "List warehouses",
		Aliases: []string{"ls"},
		Example: `  # List warehouses
  lkctl warehouse ls
  
  # Filter by inactive status
  lkctl warehouse ls --status inactive`,
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			opt := managementv1.ListWarehouseOptions{
				ProjectID:       project,
				WarehouseStatus: []managementv1.WarehouseStatus{},
			}

			if len(status) > 0 {
				for _, s := range status {
					opt.WarehouseStatus = append(opt.WarehouseStatus, managementv1.WarehouseStatus(s))
				}
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).WarehouseV1(*project).List(ctx, &opt)
			errors.Check(err)

			switch output {
			case "text", "wide":
				printWarehouses(output, resp.Warehouses...)
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	command.Flags().StringSliceVar(&status, "status", []string{}, "Filter by status. Can be repeated multiple times to filter by multiple statuses. One of: active|inactice")
	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text|wide")

	return &command
}

func NewWarehouseGetCmd(clientOpts *clientOptions, project *string) *cobra.Command {
	var output string

	command := cobra.Command{
		Use:   "get WAREHOUSEIDs",
		Short: "get a warehouse by id",
		Example: `  # get a warehouse by id
  lkctl warehouse get 019861a0-6d4e-7bf3-96c6-9aef2d4a2749`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			resp, _, err := MustCreateClient(ctx, clientOpts).WarehouseV1(*project).Get(ctx, args[0])
			errors.Check(err)

			switch output {
			case "text":
				printWarehouses(output, resp)
			case "wide":
				printWarehouse(resp)
			case "json":
				err := PrintResource(resp, output)
				errors.Check(err)
			default:
				log.Fatalf("unknown output format: %s", output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: json|text|wide")

	return &command
}

func NewWarehouseDeleteCmd(clientOpts *clientOptions, project *string) *cobra.Command {
	var force bool

	command := cobra.Command{
		Use:     "delete WAREHOUSEID",
		Aliases: []string{"rm"},
		Short:   "delete a warehouse by id",
		Example: `  # delete a warehouse by id
  lkctl warehouse delete 019861a0-6d4e-7bf3-96c6-9aef2d4a2749
  
  # force delete a warehouse
  lkctl warehouse rm 019861a0-6d4e-7bf3-96c6-9aef2d4a2749 --force`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			opt := managementv1.DeleteWarehouseOptions{
				Force: core.Ptr(force),
			}

			_, err := MustCreateClient(ctx, clientOpts).WarehouseV1(*project).Delete(ctx, args[0], &opt)
			errors.Check(err)

			fmt.Printf("Warehouse %s deleted\n", args[0])
		},
	}

	command.Flags().BoolVar(&force, "force", false, "Force delete the warehouse")

	return &command
}

func printWarehouses(output string, warehouses ...*managementv1.Warehouse) {
	if len(warehouses) < 1 {
		fmt.Println("No warehouses available")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if output == "text" {
		fmt.Fprintf(w, "ID\tNAME\tSTORAGE PROFILE\tDELETE PROFILE\tSTATUS\tPROJECT ID\n")
		for _, wh := range warehouses {
			dp := "hard"
			if wh.DeleteProfile.DeleteProfileSettings.GetDeteProfileType() == profilev1.SoftDeleteProfileType {
				exp := wh.DeleteProfile.DeleteProfileSettings.(*profilev1.TabularDeleteProfileSoft).ExpirationSeconds
				dp = fmt.Sprintf("soft (%d)", exp)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", wh.ID, wh.Name, wh.StorageProfile.StorageSettings.GetStorageFamily(), dp, wh.Status, wh.ProjectID)
		}
		w.Flush()
	}
}

func printWarehouse(warehouse *managementv1.Warehouse) {
	if warehouse == nil {
		fmt.Println("Warehouse not found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	switch warehouse.StorageProfile.StorageSettings.GetStorageFamily() {
	case profilev1.StorageFamilyS3:
		sp, ok := warehouse.StorageProfile.AsS3()
		if !ok {
			log.Fatal("could not unmarshal warehouse storage profile")
		}
		fmt.Fprintf(w, "ID\tNAME\tSTORAGE PROFILE\tBUCKET\tREGION\tSTS ENABLED\tALTERNATIVE PROTOCOLS\tASSUME ROLE\tKMS KEY\tENPOINT\tPREFIX\tPATH STYLE\tPUSH S3 DISABLED\tREMOTE SIGNING\tSTS ROLE\tSTS TOKEN VALIDITY\tSTATUS\tPROJECT ID\n")
		fmt.Fprintf(w, "%s", warehouse.ID)
		fmt.Fprintf(w, "\t%s", warehouse.Name)
		fmt.Fprintf(w, "\t%s", warehouse.StorageProfile.StorageSettings.GetStorageFamily())
		fmt.Fprintf(w, "\t%s", sp.Bucket)
		fmt.Fprintf(w, "\t%s", sp.Region)
		fmt.Fprintf(w, "\t%t", sp.STSEnabled)
		if sp.AllowAlternativeProtocols != nil {
			fmt.Fprintf(w, "\t%t", *sp.AllowAlternativeProtocols)
		} else {
			fmt.Fprintf(w, "\t%t", false)
		}
		fmt.Fprintf(w, "\t%s", FormatPString(sp.AssumeRoleARN))
		fmt.Fprintf(w, "\t%s", FormatPString(sp.AWSKMSKeyARN))
		fmt.Fprintf(w, "\t%s", FormatPString(sp.Endpoint))
		fmt.Fprintf(w, "\t%s", FormatPString(sp.KeyPrefix))
		if sp.PathStyleAccess != nil {
			fmt.Fprintf(w, "\t%t", *sp.PathStyleAccess)
		} else {
			fmt.Fprintf(w, "\t%s", "")
		}
		if sp.PushS3DeleteDisabled != nil {
			fmt.Fprintf(w, "\t%t", *sp.PushS3DeleteDisabled)
		} else {
			fmt.Fprintf(w, "\t%s", "")
		}
		if sp.RemoteSigningURLStyle != nil {
			fmt.Fprintf(w, "\t%s", *sp.RemoteSigningURLStyle)
		} else {
			fmt.Fprintf(w, "\t%s", "")
		}
		fmt.Fprintf(w, "\t%s", FormatPString(sp.STSRoleARN))
		fmt.Fprintf(w, "\t%d", sp.STSTokenValiditySeconds)
		fmt.Fprintf(w, "\t%s", warehouse.Status)
		fmt.Fprintf(w, "\t%s", warehouse.ProjectID)
		fmt.Fprintf(w, "\n")
	case profilev1.StorageFamilyGCS:
		sp, ok := warehouse.StorageProfile.AsGCS()
		if !ok {
			log.Fatal("could not unmarshal warehouse storage profile")
		}
		fmt.Fprintf(w, "ID\tNAME\tSTORAGE PROFILE\tBUCKET\tPREFIX\tSTATUS\tPROJECT ID\n")
		fmt.Fprintf(w, "%s", warehouse.ID)
		fmt.Fprintf(w, "\t%s", warehouse.Name)
		fmt.Fprintf(w, "\t%s", sp.GetStorageFamily())
		fmt.Fprintf(w, "\t%s", sp.Bucket)
		fmt.Fprintf(w, "\t%s", FormatPString(sp.KeyPrefix))
		fmt.Fprintf(w, "\t%s", warehouse.Status)
		fmt.Fprintf(w, "\t%s", warehouse.ProjectID)
		fmt.Fprintf(w, "\n")
	case profilev1.StorageFamilyADLS:
		sp, ok := warehouse.StorageProfile.AsADLS()
		if !ok {
			log.Fatal("could not unmarshal warehouse storage profile")
		}
		fmt.Fprintf(w, "ID\tNAME\tSTORAGE PROFILE\tACCOUNT NAME\tFILESYSTEM\tALLOW ALTERNATIVE PROTOCOLS\tAUTHORITY HOST\tHOST\tKEY PREFIX\tSAS TOKEN VALIDITY\tSTATUS\tPROJECT ID\n")
		fmt.Fprintf(w, "%s", warehouse.ID)
		fmt.Fprintf(w, "\t%s", warehouse.Name)
		fmt.Fprintf(w, "\t%s", sp.GetStorageFamily())
		fmt.Fprintf(w, "\t%s", sp.AccountName)
		fmt.Fprintf(w, "\t%s", sp.Filesystem)
		if sp.AllowAlternativeProtocols != nil {
			fmt.Fprintf(w, "\t%t", *sp.AllowAlternativeProtocols)
		} else {
			fmt.Fprintf(w, "\t%t", false)
		}
		fmt.Fprintf(w, "\t%s", FormatPString(sp.AuthorityHost))
		fmt.Fprintf(w, "\t%s", FormatPString(sp.Host))
		fmt.Fprintf(w, "\t%s", FormatPString(sp.KeyPrefix))
		fmt.Fprintf(w, "\t%d", sp.SASTokenValiditySeconds)
		fmt.Fprintf(w, "\t%s", warehouse.Status)
		fmt.Fprintf(w, "\t%s", warehouse.ProjectID)
		fmt.Fprintf(w, "\n")
	default:
		log.Fatalf("unknown storage profile %s", warehouse.StorageProfile.StorageSettings.GetStorageFamily())
	}
	w.Flush()
}
