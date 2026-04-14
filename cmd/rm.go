// Copyright 2026 Google LLC
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

package cmd

import (
	"fmt"
	"os"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm [domain] [record-id]",
	Short:   "Remove a DNS record from a domain using its ID",
	GroupID: GroupManagement,
	Long:    `Deletes a specific DNS record from your Porkbun domain. You must provide the exact record ID, which can be found using the 'list-records' command.`,
	Example: `  # Delete record ID 123456789 from aaie.cloud
  steamer rm aaie.cloud 123456789`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain := args[0]
		id := args[1]

		client := porkbun.NewClient(apiKey, secretKey)
		err = client.DeleteRecord(domain, id)
		if err != nil {
			fmt.Println(theme.Fail.Render(fmt.Sprintf("Error deleting record: %v", err)))
			os.Exit(1)
		}

		fmt.Println(theme.Pass.Render(fmt.Sprintf("Successfully deleted record %s from %s", id, domain)))
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
