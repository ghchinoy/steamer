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
	"time"

	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
)

var searchRateLimitWait = 10 * time.Second

type searchResult struct {
	Domain   string                       `json:"domain"`
	Response *porkbun.DomainCheckResponse `json:"response,omitempty"`
	Error    string                       `json:"error,omitempty"`
}

func newSearchCmd(app *cliApp) *cobra.Command {
	var (
		outputJSON bool
		searchTLDs []string
	)

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			domainsToCheck := expandSearchDomains(args[0], searchTLDs)
			results := make([]searchResult, 0, len(domainsToCheck))
			out := cmd.OutOrStdout()

			for i, domain := range domainsToCheck {
				if i > 0 && searchRateLimitWait > 0 {
					if !outputJSON {
						_, _ = fmt.Fprintf(out, "%s Waiting %s for Porkbun rate limits...\n", theme.Warn.Render("⏳"), searchRateLimitWait)
					}
					if err := waitWithContext(cmd.Context(), searchRateLimitWait); err != nil {
						return err
					}
				}

				ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
				res, checkErr := client.CheckDomain(ctx, domain)
				cancel()

				item := searchResult{Domain: domain, Response: res}
				if checkErr != nil {
					item.Error = checkErr.Error()
				}
				results = append(results, item)
			}

			if outputJSON {
				return printJSON(out, results)
			}

			for _, r := range results {
				_, _ = fmt.Fprintf(out, "%s %s\n", theme.Accent.Render("Domain:"), r.Domain)
				if r.Error != "" {
					_, _ = fmt.Fprintf(out, "%s %s\n\n", theme.Fail.Render("Error:"), r.Error)
					continue
				}

				if r.Response == nil || r.Response.Response.Avail != "yes" {
					_, _ = fmt.Fprintf(out, "%s %s\n\n", theme.Fail.Render("Status:"), "Unavailable (Already registered)")
					continue
				}

				_, _ = fmt.Fprintf(out, "%s %s\n", theme.Pass.Render("Status:"), "Available")
				if r.Response.Response.Price != "" {
					_, _ = fmt.Fprintf(out, "%s %s\n", theme.Accent.Render("Price:"), "$"+r.Response.Response.Price)
				}
				isPremium := "No"
				if r.Response.Response.Premium == "yes" {
					isPremium = "Yes"
				}
				_, _ = fmt.Fprintf(out, "%s %s\n\n", theme.Accent.Render("Premium:"), isPremium)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output results in JSON format")
	cmd.Flags().StringSliceVar(&searchTLDs, "tlds", []string{"com", "net", "org", "co", "io", "dev"}, "Comma-separated list of TLDs to check when a phrase is provided")
	return cmd
}

func expandSearchDomains(query string, tlds []string) []string {
	if strings.Contains(query, ".") {
		return []string{query}
	}
	domains := make([]string, 0, len(tlds))
	for _, tld := range tlds {
		cleanTLD := strings.TrimPrefix(tld, ".")
		domains = append(domains, fmt.Sprintf("%s.%s", query, cleanTLD))
	}
	return domains
}

func waitWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
