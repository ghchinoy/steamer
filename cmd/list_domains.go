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

	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/theme"
)

func newListDomainsCmd(app *cliApp) *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:     "list-domains",
		Short:   "List all domains in your Porkbun account",
		GroupID: GroupInfo,
		Long:    `Retrieves and displays a list of all domains associated with your Porkbun account, including their current status and expiration dates.`,
		Example: `  # List domains in a table
  steamer list-domains

  # Output domains as JSON for scripting
  steamer list-domains --json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			domains, err := client.ListDomains(ctx)
			if err != nil {
				return fmt.Errorf("listing domains: %w", err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), domains)
			}

			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(out, "%s %s %s %s\n",
				theme.Accent.Render(fmt.Sprintf("%-25s", "DOMAIN")),
				theme.Accent.Render(fmt.Sprintf("%-10s", "STATUS")),
				theme.Accent.Render(fmt.Sprintf("%-10s", "TLD")),
				theme.Accent.Render(fmt.Sprintf("%-20s", "EXPIRATION")),
			)
			for _, d := range domains {
				statusColor := theme.Pass
				if d.Status != "ACTIVE" {
					statusColor = theme.Warn
				}
				_, _ = fmt.Fprintf(out, "%-25s %s %s %-20s\n",
					d.Domain,
					statusColor.Render(fmt.Sprintf("%-10s", d.Status)),
					theme.Muted.Render(fmt.Sprintf("%-10s", d.TLD)),
					d.ExpireDate,
				)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output results in JSON format")
	return cmd
}
