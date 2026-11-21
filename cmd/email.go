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
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/theme"
)

func newEmailCmd(app *cliApp) *cobra.Command {
	emailCmd := &cobra.Command{
		Use:     "email",
		Short:   "Manage email forwarding rules for domains",
		GroupID: GroupEmail,
		Long:    `View, create, and delete email forwarding rules for domains in your Porkbun account.`,
		Example: `  # List all email forwards for a domain
  steamer email list example.com

  # Add an email forwarding rule
  steamer email add example.com hello user@example.org

  # Remove an email forwarding rule
  steamer email rm example.com hello user@example.org`,
	}

	emailCmd.AddCommand(
		newEmailListCmd(app),
		newEmailAddCmd(app),
		newEmailRmCmd(app),
	)
	return emailCmd
}

func newEmailListCmd(app *cliApp) *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "list <domain>",
		Short: "List email forwards for a specific domain",
		Long:  `Retrieves and displays all configured email forwarding rules, forward quotas, and DNS status for the domain.`,
		Example: `  # List email forwards for example.com
  steamer email list example.com

  # Output email forwards as JSON
  steamer email list example.com --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			res, err := client.FetchEmailForwards(ctx, domain)
			if err != nil {
				return fmt.Errorf("retrieving email forwards for %q: %w", domain, err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), res)
			}

			out := cmd.OutOrStdout()
			if len(res.Forwards) == 0 {
				_, err = fmt.Fprintf(out, "No email forwards found for %s.\n", domain)
				return err
			}

			_, _ = fmt.Fprintf(out, "%s %s %s\n",
				theme.Accent.Render(fmt.Sprintf("%-32s", "ADDRESS")),
				theme.Accent.Render(fmt.Sprintf("%-32s", "FORWARD TO")),
				theme.Accent.Render(fmt.Sprintf("%-20s", "CREATED")),
			)

			for _, f := range res.Forwards {
				_, _ = fmt.Fprintf(out, "%s %-32s %s\n",
					theme.ID.Render(fmt.Sprintf("%-32s", f.Address)),
					f.ForwardTo,
					theme.Muted.Render(f.Created),
				)
			}

			dnsStatus := theme.Pass.Render("Configured")
			if !res.DNSConfigured {
				dnsStatus = theme.Warn.Render("Not Configured")
			}

			_, _ = fmt.Fprintln(out)
			_, err = fmt.Fprintf(out, "%s %d/%d used  |  %s %s\n",
				theme.Muted.Render("Quota:"),
				res.Limits.Used,
				res.Limits.Max,
				theme.Muted.Render("MX/SPF DNS:"),
				dnsStatus,
			)
			return err
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output forwards in JSON format")
	return cmd
}

func newEmailAddCmd(app *cliApp) *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "add <domain> <address> <forward-to>",
		Short: "Add an email forwarding rule to a domain",
		Long:  `Creates a new email forwarding rule in Porkbun for the specified domain, routing emails from <address> to <forward-to>.`,
		Example: `  # Forward hello@example.com to your destination email
  steamer email add example.com hello user@example.org

  # Forward using full address format
  steamer email add example.com hello@example.com user@example.org`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			address := args[1]
			forwardTo := args[2]
			cleanAddr := strings.TrimSuffix(address, "@"+domain)

			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			if err := client.AddEmailForward(ctx, domain, address, forwardTo); err != nil {
				return fmt.Errorf("creating email forward: %w", err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), map[string]string{
					"domain":    domain,
					"address":   fmt.Sprintf("%s@%s", cleanAddr, domain),
					"forwardTo": forwardTo,
					"status":    "created",
				})
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), theme.Pass.Render(fmt.Sprintf(
				"Successfully created email forward: %s@%s -> %s",
				cleanAddr, domain, forwardTo,
			)))
			return err
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output result in JSON format")
	return cmd
}

func newEmailRmCmd(app *cliApp) *cobra.Command {
	var (
		outputJSON bool
		yes        bool
	)

	cmd := &cobra.Command{
		Use:   "rm <domain> <address> <forward-to>",
		Short: "Remove an email forwarding rule from a domain",
		Long:  `Deletes an existing email forwarding rule in Porkbun for the specified domain, address, and destination. Prompts for confirmation unless --yes is provided.`,
		Example: `  # Remove forwarding from hello@example.com to user@example.org
  steamer email rm example.com hello user@example.org

  # Remove without interactive confirmation
  steamer email rm example.com hello user@example.org --yes`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			address := args[1]
			forwardTo := args[2]
			cleanAddr := strings.TrimSuffix(address, "@"+domain)

			prompt := fmt.Sprintf("Delete email forward %s@%s -> %s?", cleanAddr, domain, forwardTo)
			confirmed, err := confirmAction(cmd, prompt, yes)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			if err := client.DeleteEmailForward(ctx, domain, address, forwardTo); err != nil {
				return fmt.Errorf("deleting email forward: %w", err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), map[string]string{
					"domain":    domain,
					"address":   fmt.Sprintf("%s@%s", cleanAddr, domain),
					"forwardTo": forwardTo,
					"status":    "deleted",
				})
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), theme.Pass.Render(fmt.Sprintf(
				"Successfully deleted email forward: %s@%s -> %s",
				cleanAddr, domain, forwardTo,
			)))
			return err
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output result in JSON format")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip interactive confirmation prompt")
	return cmd
}
