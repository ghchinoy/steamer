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
	"strings"
	"time"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"

	"github.com/spf13/cobra"
)

var searchJSON bool
var searchTlds []string

var searchCmd = &cobra.Command{
	Use:     "search <domain-or-phrase>",
	Short:   "Check domain availability and pricing",
	GroupID: GroupInfo,
	Args:    cobra.ExactArgs(1),
	Long:    `Queries the Porkbun API to check if a specific domain is available for registration, and retrieves its first-year pricing if available. If a phrase is provided without a TLD (e.g., 'mynewidea'), it will check a predefined list of popular TLDs or the TLDs specified via the --tlds flag.`,
	Example: `  # Check if a specific domain is available
  steamer search mynewidea.com

  # Search a phrase against default TLDs (.com, .net, .org, .co, .io, .dev)
  steamer search mynewidea

  # Search a phrase against specific TLDs
  steamer search mynewidea --tlds ai,app,xyz

  # Check availability and output as JSON
  steamer search mynewidea.com --json`,
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client := porkbun.NewClient(apiKey, secretKey)

		// Determine if it's a phrase or a specific domain
		var domainsToCheck []string
		if strings.Contains(query, ".") {
			domainsToCheck = []string{query}
		} else {
			for _, tld := range searchTlds {
				// Strip leading dot if user provided one in the flag (e.g., --tlds .com)
				cleanTLD := strings.TrimPrefix(tld, ".")
				domainsToCheck = append(domainsToCheck, fmt.Sprintf("%s.%s", query, cleanTLD))
			}
		}

		// Collect results sequentially to respect the 1-check-per-10s rate limit
		finalResults := make([]struct {
			Domain string
			Res    *porkbun.DomainCheckResponse
			Err    error
		}, 0, len(domainsToCheck))

		for i, d := range domainsToCheck {
			if i > 0 {
				if !searchJSON {
					fmt.Printf("%s Waiting 10s for Porkbun rate limits...\n", theme.Warn.Render("⏳"))
				}
				time.Sleep(10 * time.Second)
			}
			res, err := client.CheckDomain(d)
			finalResults = append(finalResults, struct {
				Domain string
				Res    *porkbun.DomainCheckResponse
				Err    error
			}{Domain: d, Res: res, Err: err})
		}

		if searchJSON {
			b, err := json.MarshalIndent(finalResults, "", "  ")
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
			return
		}

		// Output formatting
		for _, r := range finalResults {
			fmt.Printf("%s %s\n", theme.Accent.Render("Domain:"), r.Domain)
			if r.Err != nil {
				fmt.Printf("%s %s\n", theme.Fail.Render("Error:"), r.Err.Error())
				fmt.Println()
				continue
			}

			if r.Res.Response.Avail == "yes" {
				fmt.Printf("%s %s\n", theme.Pass.Render("Status:"), "Available")
				if r.Res.Response.Price != "" {
					fmt.Printf("%s %s\n", theme.Accent.Render("Price:"), "$"+r.Res.Response.Price)
				}
				isPremium := "No"
				if r.Res.Response.Premium == "yes" {
					isPremium = "Yes"
				}
				fmt.Printf("%s %s\n", theme.Accent.Render("Premium:"), isPremium)
			} else {
				fmt.Printf("%s %s\n", theme.Fail.Render("Status:"), "Unavailable (Already registered)")
			}
			fmt.Println()
		}
	},
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output results in JSON format")
	searchCmd.Flags().StringSliceVar(&searchTlds, "tlds", []string{"com", "net", "org", "co", "io", "dev"}, "Comma-separated list of TLDs to check when a phrase is provided")
	rootCmd.AddCommand(searchCmd)
}
