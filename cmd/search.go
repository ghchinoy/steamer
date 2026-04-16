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

var searchJSON bool

var searchCmd = &cobra.Command{
	Use:     "search <domain>",
	Short:   "Check domain availability and pricing",
	GroupID: GroupInfo,
	Args:    cobra.ExactArgs(1),
	Long:    `Queries the Porkbun API to check if a specific domain is available for registration, and retrieves its first-year pricing if available.`,
	Example: `  # Check if a domain is available
  steamer search mynewidea.com

  # Check availability and output as JSON
  steamer search mynewidea.com --json`,
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client := porkbun.NewClient(apiKey, secretKey)
		res, err := client.CheckDomain(domain)
		if err != nil {
			fmt.Printf("Error checking domain: %v\n", err)
			os.Exit(1)
		}

		if searchJSON {
			b, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
			return
		}

		fmt.Printf("%s %s\n", theme.Accent.Render("Domain:"), domain)
		if res.Response.Avail == "yes" {
			fmt.Printf("%s %s\n", theme.Pass.Render("Status:"), "Available")
			if res.Response.Price != "" {
				fmt.Printf("%s %s\n", theme.Accent.Render("Price:"), "$"+res.Response.Price)
			}
			isPremium := "No"
			if res.Response.Premium == "yes" {
				isPremium = "Yes"
			}
			fmt.Printf("%s %s\n", theme.Accent.Render("Premium:"), isPremium)
		} else {
			fmt.Printf("%s %s\n", theme.Fail.Render("Status:"), "Unavailable (Already registered)")
		}
	},
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output results in JSON format")
	rootCmd.AddCommand(searchCmd)
}
