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

	"github.com/spf13/cobra"
)

var listDomainsJSON bool

var listDomainsCmd = &cobra.Command{
	Use:     "list-domains",
	Short:   "List all domains in your Porkbun account",
	GroupID: GroupInfo,
	Long:    `Retrieves and displays a list of all domains associated with your Porkbun account, including their current status and expiration dates.`,
	Example: `  # List domains in a table
  steamer list-domains

  # Output domains as JSON for scripting
  steamer list-domains --json`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client := porkbun.NewClient(apiKey, secretKey)
		domains, err := client.ListDomains()
		if err != nil {
			fmt.Printf("Error listing domains: %v\n", err)
			os.Exit(1)
		}

		if listDomainsJSON {
			b, err := json.MarshalIndent(domains, "", "  ")
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
			return
		}

		fmt.Printf("%-25s %-10s %-10s %-20s\n", "DOMAIN", "STATUS", "TLD", "EXPIRATION")
		for _, d := range domains {
			fmt.Printf("%-25s %-10s %-10s %-20s\n", d.Domain, d.Status, d.TLD, d.ExpireDate)
		}
	},
}

func init() {
	listDomainsCmd.Flags().BoolVar(&listDomainsJSON, "json", false, "Output results in JSON format")
	rootCmd.AddCommand(listDomainsCmd)
}
