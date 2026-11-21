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
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
)

const defaultRequestTimeout = 30 * time.Second

type recordCmdSpec struct {
	use             string
	short           string
	long            string
	example         string
	recordType      string
	valueFieldJSON  string
	successVerb     string
	errorHintFormat string
}

type recordCreateResult struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	Type      string `json:"type"`
	Content   string `json:"content"`
}

func newAddRecordCmd(app *cliApp, spec recordCmdSpec) *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:     spec.use,
		Short:   spec.short,
		GroupID: GroupManagement,
		Long:    spec.long,
		Example: spec.example,
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domain := args[0]
			subdomain := args[1]
			content := args[2]
			fqdn := formatFQDN(subdomain, domain)

			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			recordID, err := client.CreateRecord(ctx, domain, porkbun.CreateRecordRequest{
				Name:    subdomain,
				Type:    spec.recordType,
				Content: content,
			})
			if err != nil {
				if spec.errorHintFormat != "" {
					return fmt.Errorf("creating %s record: %w - %s", spec.recordType, err, fmt.Sprintf(spec.errorHintFormat, fqdn))
				}
				return fmt.Errorf("creating %s record: %w", spec.recordType, err)
			}

			if outputJSON {
				if spec.valueFieldJSON != "" && spec.valueFieldJSON != "content" {
					out := map[string]string{
						"id":                recordID,
						"domain":            domain,
						"subdomain":         subdomain,
						"type":              spec.recordType,
						spec.valueFieldJSON: content,
					}
					return printJSON(cmd.OutOrStdout(), out)
				}
				return printJSON(cmd.OutOrStdout(), recordCreateResult{
					ID:        recordID,
					Domain:    domain,
					Subdomain: subdomain,
					Type:      spec.recordType,
					Content:   content,
				})
			}

			verb := spec.successVerb
			if verb == "" {
				verb = "pointing to"
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), theme.Pass.Render(fmt.Sprintf(
				"Successfully created %s record for %s %s %s (ID: %s)",
				spec.recordType, fqdn, verb, content, theme.ID.Render(recordID),
			)))
			return err
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output result in JSON format")
	return cmd
}

func printJSON(w io.Writer, value any) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	_, err = fmt.Fprintln(w, string(encoded))
	return err
}

func formatFQDN(subdomain, domain string) string {
	if subdomain == "" {
		return domain
	}
	return subdomain + "." + domain
}

func confirmAction(cmd *cobra.Command, prompt string, skipConfirm bool) (bool, error) {
	if skipConfirm {
		return true, nil
	}

	in := cmd.InOrStdin()
	if file, ok := in.(*os.File); ok {
		if !isatty.IsTerminal(file.Fd()) && !isatty.IsCygwinTerminal(file.Fd()) {
			return false, errors.New("refusing destructive operation in non-interactive mode; re-run with --yes to confirm")
		}
	}

	if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "%s [y/N]: ", theme.Warn.Render(prompt)); err != nil {
		return false, err
	}

	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("reading confirmation: %w", err)
	}

	reply := strings.ToLower(strings.TrimSpace(line))
	if reply == "y" || reply == "yes" {
		return true, nil
	}

	_, _ = fmt.Fprintln(cmd.ErrOrStderr(), theme.Muted.Render("Aborted."))
	return false, nil
}
