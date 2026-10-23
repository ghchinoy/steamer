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

var addNsJSON bool

var addNsCmd = &cobra.Command{
	Use:     "add-ns [domain] [subdomain] [nameserver]",
	Short:   "Add a new NS record to a domain",
	GroupID: GroupManagement,
	Long:    `Creates a new NS (Name Server) record for the specified domain and subdomain, delegating DNS authority for the subdomain to the target nameserver.`,
	Example: `  # Delegate int.aaie.cloud to a Google Cloud DNS nameserver
  steamer add-ns aaie.cloud int ns-cloud-b1.googledomains.com

  # Delegate and return output as JSON
  steamer add-ns aaie.cloud int ns-cloud-b1.googledomains.com --json`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(theme.Fail.Render(err.Error()))
			os.Exit(1)
		}

		domain := args[0]
		subdomain := args[1]
		nameserver := args[2]

		client := porkbun.NewClient(apiKey, secretKey)
		id, err := client.CreateRecord(domain, porkbun.CreateRecordRequest{
			Name:    subdomain,
			Type:    "NS",
			Content: nameserver,
		})
		if err != nil {
			fmt.Println(theme.Fail.Render(fmt.Sprintf("Error creating NS record: %v - hint: verify the domain exists in your Porkbun account and no conflicting CNAME record exists for '%s.%s'", err, subdomain, domain)))
			os.Exit(1)
		}

		if addNsJSON {
			out := map[string]string{
				"id":         id,
				"domain":     domain,
				"subdomain":  subdomain,
				"type":       "NS",
				"nameserver": nameserver,
			}
			b, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				fmt.Println(theme.Fail.Render(fmt.Sprintf("Error encoding JSON: %v", err)))
				os.Exit(1)
			}
			fmt.Println(string(b))
			return
		}

		fmt.Println(theme.Pass.Render(fmt.Sprintf("Successfully created NS record for %s.%s pointing to %s (ID: %s)", subdomain, domain, nameserver, theme.ID.Render(id))))
	},
}

func init() {
	addNsCmd.Flags().BoolVar(&addNsJSON, "json", false, "Output result in JSON format")
	rootCmd.AddCommand(addNsCmd)
}
