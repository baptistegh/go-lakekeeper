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
	"io"
	"os"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// warehouseCmd represents the warehouse command
var warehouseCmd = &cobra.Command{
	Use:   "warehouse",
	Short: "Interact with warehouses",
}

// listWarehouseCmd represents the warehouse list command
var listWarehouseCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List available warehouses",
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.ListWarehouseOptions{}

		for _, v := range viper.GetStringSlice("warehouse_list_status") {
			opt.WarehouseStatus = append(opt.WarehouseStatus, managementv1.WarehouseStatus(v))
		}

		resp, _, err := c.WarehouseV1(project).List(cmd.Context(), &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// addWarehouseCmd represents the warehouse add command
var addWarehouseCmd = &cobra.Command{
	Use:   "add [flags] <json-config-file> ",
	Short: "Add a new warehouse. if '-' is supplied, config file is read from stdin.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var reader io.Reader

		if args[0] == "-" {
			reader = cmd.InOrStdin()
		} else {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()

			reader = file
		}

		var opt managementv1.CreateWarehouseOptions

		if err := json.NewDecoder(reader).Decode(&opt); err != nil {
			return err
		}

		resp, _, err := c.WarehouseV1(project).Create(cmd.Context(), &opt)
		if err != nil {
			return err
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
	},
}

// deleteWarehouseCmd represents the warehouse delete command
var deleteWarehouseCmd = &cobra.Command{
	Use:   "delete [flags] <warehouse-id> [--force]",
	Short: "delete warehouse by id.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		opt := managementv1.DeleteWarehouseOptions{}

		if viper.GetBool("warehouse_delete_force") {
			opt.Force = core.Ptr(true)
		}

		if _, err := c.WarehouseV1(project).Delete(cmd.Context(), args[0], &opt); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(warehouseCmd)

	warehouseCmd.AddCommand(addWarehouseCmd)
	warehouseCmd.AddCommand(deleteWarehouseCmd)
	warehouseCmd.AddCommand(listWarehouseCmd)

	deleteWarehouseCmd.Flags().BoolP("force", "f", false, "delete protected warehouses")
	_ = viper.BindPFlag("warehouse_delete_force", deleteWarehouseCmd.Flags().Lookup("force"))

	listWarehouseCmd.Flags().StringSlice("status", []string{}, "filter by status, values can be [active inactive]")
	_ = viper.BindPFlag("warehouse_list_status", listWarehouseCmd.Flags().Lookup("status"))

}
