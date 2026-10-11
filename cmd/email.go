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

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
	"github.com/spf13/cobra"
)

var emailListJSON bool

var emailCmd = &cobra.Command{
	Use:     "email",
	Short:   "Manage email forwarding rules for domains",
	GroupID: GroupEmail,
	Long:    `View, create, and delete email forwarding rules for domains in your Porkbun account.`,
	Example: `  # List all email forwards for a domain
  steamer email list chinoy.si

  # Add an email forwarding rule
  steamer email add chinoy.si hello ghchinoy@gmail.com

  # Remove an email forwarding rule
  steamer email rm chinoy.si hello ghchinoy@gmail.com`,
}

var emailListCmd = &cobra.Command{
	Use:   "list [domain]",
	Short: "List email forwards for a specific domain",
	Long:  `Retrieves and displays all configured email forwarding rules, forward quotas, and DNS status for the domain.`,
	Example: `  # List email forwards for chinoy.si
  steamer email list chinoy.si

  # Output email forwards as JSON
  steamer email list chinoy.si --json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain := args[0]
		client := porkbun.NewClient(apiKey, secretKey)
		res, err := client.GetEmailForwards(domain)
		if err != nil {
			fmt.Println(theme.Fail.Render(fmt.Sprintf("Error retrieving email forwards for %s: %v", domain, err)))
			os.Exit(1)
		}

		if emailListJSON {
			b, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
			return
		}

		if len(res.Forwards) == 0 {
			fmt.Printf("No email forwards found for %s.\n", domain)
			return
		}

		fmt.Printf("%s %s %s\n",
			theme.Accent.Render(fmt.Sprintf("%-32s", "ADDRESS")),
			theme.Accent.Render(fmt.Sprintf("%-32s", "FORWARD TO")),
			theme.Accent.Render(fmt.Sprintf("%-20s", "CREATED")),
		)

		for _, f := range res.Forwards {
			fmt.Printf("%-32s %-32s %s\n",
				theme.ID.Render(f.Address),
				f.ForwardTo,
				theme.Muted.Render(f.Created),
			)
		}

		dnsStatus := theme.Pass.Render("Configured")
		if !res.DNSConfigured {
			dnsStatus = theme.Warn.Render("Not Configured")
		}

		fmt.Println()
		fmt.Printf("%s %d/%d used  |  %s %s\n",
			theme.Muted.Render("Quota:"),
			res.Limits.Used,
			res.Limits.Max,
			theme.Muted.Render("MX/SPF DNS:"),
			dnsStatus,
		)
	},
}

var emailAddCmd = &cobra.Command{
	Use:   "add [domain] [address] [forward-to]",
	Short: "Add an email forwarding rule to a domain",
	Long:  `Creates a new email forwarding rule in Porkbun for the specified domain, routing emails from [address] to [forward-to].`,
	Example: `  # Forward hello@chinoy.si to your personal email
  steamer email add chinoy.si hello ghchinoy@gmail.com

  # Forward using full address format
  steamer email add chinoy.si hello@chinoy.si ghchinoy@gmail.com`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain := args[0]
		address := args[1]
		forwardTo := args[2]

		cleanAddr := strings.TrimSuffix(address, "@"+domain)

		client := porkbun.NewClient(apiKey, secretKey)
		err = client.AddEmailForward(domain, cleanAddr, forwardTo)
		if err != nil {
			fmt.Println(theme.Fail.Render(fmt.Sprintf("Error creating email forward: %v", err)))
			os.Exit(1)
		}

		fmt.Println(theme.Pass.Render(fmt.Sprintf("Successfully created email forward: %s@%s -> %s", cleanAddr, domain, forwardTo)))
	},
}

var emailRmCmd = &cobra.Command{
	Use:   "rm [domain] [address] [forward-to]",
	Short: "Remove an email forwarding rule from a domain",
	Long:  `Deletes an existing email forwarding rule in Porkbun for the specified domain, address, and destination.`,
	Example: `  # Remove forwarding from hello@chinoy.si to ghchinoy@gmail.com
  steamer email rm chinoy.si hello ghchinoy@gmail.com`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domain := args[0]
		address := args[1]
		forwardTo := args[2]

		cleanAddr := strings.TrimSuffix(address, "@"+domain)

		client := porkbun.NewClient(apiKey, secretKey)
		err = client.DeleteEmailForward(domain, cleanAddr, forwardTo)
		if err != nil {
			fmt.Println(theme.Fail.Render(fmt.Sprintf("Error deleting email forward: %v", err)))
			os.Exit(1)
		}

		fmt.Println(theme.Pass.Render(fmt.Sprintf("Successfully deleted email forward: %s@%s -> %s", cleanAddr, domain, forwardTo)))
	},
}

func init() {
	emailListCmd.Flags().BoolVar(&emailListJSON, "json", false, "Output forwards in JSON format")

	emailCmd.AddCommand(emailListCmd)
	emailCmd.AddCommand(emailAddCmd)
	emailCmd.AddCommand(emailRmCmd)

	rootCmd.AddCommand(emailCmd)
}
