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
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"

	"github.com/spf13/cobra"
)

var listRecordsJSON bool

var listRecordsCmd = &cobra.Command{
	Use:     "list-records [domain]",
	Short:   "List DNS records for a specific domain",
	GroupID: GroupInfo,
	Long:    `Retrieves and displays all DNS records (A, CNAME, TXT, etc.) for the specified domain from your Porkbun account.`,
	Example: `  # List records for aaie.cloud
  steamer list-records aaie.cloud

  # Output records as JSON
  steamer list-records aaie.cloud --json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain := args[0]
		client := porkbun.NewClient(apiKey, secretKey)
		records, err := client.RetrieveRecords(domain)
		if err != nil {
			fmt.Printf("Error retrieving records for %s: %v\n", domain, err)
			os.Exit(1)
		}

		if listRecordsJSON {
			b, err := json.MarshalIndent(records, "", "  ")
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
			return
		}

		fmt.Printf("%s %s %s %s\n",
			theme.Accent.Render(fmt.Sprintf("%-10s", "ID")),
			theme.Accent.Render(fmt.Sprintf("%-25s", "NAME")),
			theme.Accent.Render(fmt.Sprintf("%-10s", "TYPE")),
			theme.Accent.Render(fmt.Sprintf("%-30s", "CONTENT")),
		)
		for _, r := range records {
			fmt.Printf("%s %s %s %s\n",
				theme.ID.Render(fmt.Sprintf("%-10v", r.ID)),
				fmt.Sprintf("%-25s", r.Name),
				theme.Muted.Render(fmt.Sprintf("%-10s", r.Type)),
				fmt.Sprintf("%-30s", r.Content),
			)
		}
	},
}

func init() {
	listRecordsCmd.Flags().BoolVar(&listRecordsJSON, "json", false, "Output results in JSON format")
	rootCmd.AddCommand(listRecordsCmd)
}
