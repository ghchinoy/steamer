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

func newRmCmd(app *cliApp) *cobra.Command {
	var (
		outputJSON bool
		yes        bool
	)

	cmd := &cobra.Command{
		Use:     "rm <domain> <record-id>",
		Short:   "Remove a DNS record from a domain using its ID",
		GroupID: GroupManagement,
		Long:    `Deletes a specific DNS record from your Porkbun domain. You must provide the exact record ID, which can be found using the 'list-records' command. Prompts for confirmation unless --yes is provided.`,
		Example: `  # Delete record ID 123456789 from example.com
  steamer rm example.com 123456789

  # Delete without interactive confirmation
  steamer rm example.com 123456789 --yes`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			recordID := args[1]

			prompt := fmt.Sprintf("Delete DNS record %s from %s?", recordID, domain)
			confirmed, err := confirmAction(cmd, prompt, yes)
			if err != nil {
				return err
			}
			if !confirmed {
				return nil
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			if err := client.DeleteRecord(ctx, domain, recordID); err != nil {
				return fmt.Errorf("deleting record %q from %q: %w", recordID, domain, err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), map[string]string{
					"id":     recordID,
					"domain": domain,
					"status": "deleted",
				})
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), theme.Pass.Render(fmt.Sprintf(
				"Successfully deleted record %s from %s",
				theme.ID.Render(recordID), domain,
			)))
			return err
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output result in JSON format")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip interactive confirmation prompt")
	return cmd
}
