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

func newListRecordsCmd(app *cliApp) *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:     "list-records <domain>",
		Short:   "List DNS records for a specific domain",
		GroupID: GroupInfo,
		Long:    `Retrieves and displays all DNS records (A, CNAME, TXT, etc.) for the specified domain from your Porkbun account.`,
		Example: `  # List records for example.com
  steamer list-records example.com

  # Output records as JSON
  steamer list-records example.com --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			records, err := client.RetrieveRecords(ctx, domain)
			if err != nil {
				return fmt.Errorf("retrieving records for %q: %w", domain, err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), records)
			}

			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(out, "%s %s %s %s\n",
				theme.Accent.Render(fmt.Sprintf("%-10s", "ID")),
				theme.Accent.Render(fmt.Sprintf("%-25s", "NAME")),
				theme.Accent.Render(fmt.Sprintf("%-10s", "TYPE")),
				theme.Accent.Render(fmt.Sprintf("%-30s", "CONTENT")),
			)
			for _, r := range records {
				_, _ = fmt.Fprintf(out, "%s %-25s %s %-30s\n",
					theme.ID.Render(fmt.Sprintf("%-10s", string(r.ID))),
					r.Name,
					theme.Muted.Render(fmt.Sprintf("%-10s", r.Type)),
					r.Content,
				)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output results in JSON format")
	return cmd
}
