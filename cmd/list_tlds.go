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
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
)

const tldCacheTTL = 7 * 24 * time.Hour

func newListTLDsCmd(app *cliApp) *cobra.Command {
	var (
		outputJSON bool
		force      bool
	)

	cmd := &cobra.Command{
		Use:     "list-tlds",
		Short:   "List all supported TLDs and their pricing",
		GroupID: GroupInfo,
		Long:    `Retrieves and displays a list of all Top-Level Domains (TLDs) supported by Porkbun, along with their registration, renewal, and transfer prices. Results are cached locally for 7 days to improve performance.`,
		Example: `  # List all TLDs in a table
  steamer list-tlds

  # Force refresh the cached list of TLDs
  steamer list-tlds --force

  # Output TLDs as JSON
  steamer list-tlds --json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := app.newClient()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), defaultRequestTimeout)
			defer cancel()

			pricing, err := loadTLDPricing(ctx, client, force)
			if err != nil {
				return fmt.Errorf("fetching TLD pricing: %w", err)
			}

			if outputJSON {
				return printJSON(cmd.OutOrStdout(), pricing)
			}

			tlds := slices.Sorted(maps.Keys(pricing))

			columns := []table.Column{
				{Title: "TLD", Width: 15},
				{Title: "REGISTRATION", Width: 15},
				{Title: "RENEWAL", Width: 15},
				{Title: "TRANSFER", Width: 15},
			}

			rows := make([]table.Row, 0, len(tlds))
			for _, tld := range tlds {
				p := pricing[tld]
				rows = append(rows, table.Row{
					"." + tld,
					"$" + p.Registration,
					"$" + p.Renewal,
					"$" + p.Transfer,
				})
			}

			t := table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
				table.WithFocused(true),
				table.WithHeight(20),
			)

			s := table.DefaultStyles()
			s.Header = s.Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(theme.TableHeaderBorder).
				BorderBottom(true).
				Bold(false)
			s.Selected = theme.TableSelected
			t.SetStyles(s)

			m := tldTableModel{table: t}
			if _, err := tea.NewProgram(m, tea.WithContext(cmd.Context())).Run(); err != nil {
				return fmt.Errorf("running TUI: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output results in JSON format")
	cmd.Flags().BoolVar(&force, "force", false, "Force refresh the TLD cache")
	return cmd
}

func loadTLDPricing(ctx context.Context, client *porkbun.Client, force bool) (map[string]porkbun.TLDPricing, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return fetchPricing(ctx, client)
	}

	cacheDir := filepath.Join(home, ".config", "steamer")
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		return fetchPricing(ctx, client)
	}

	root, err := os.OpenRoot(cacheDir)
	if err != nil {
		return fetchPricing(ctx, client)
	}
	defer func() { _ = root.Close() }()

	const cacheFile = "tlds.json"
	if !force {
		if cached, ok := readCachedPricing(root, cacheFile); ok {
			return cached, nil
		}
	}

	pricing, err := fetchPricing(ctx, client)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(pricing); err == nil {
		_ = root.WriteFile(cacheFile, data, 0o600)
	}

	return pricing, nil
}

func readCachedPricing(root *os.Root, filename string) (map[string]porkbun.TLDPricing, bool) {
	info, err := root.Stat(filename)
	if err != nil || time.Since(info.ModTime()) >= tldCacheTTL {
		return nil, false
	}
	data, err := root.ReadFile(filename)
	if err != nil {
		return nil, false
	}
	var cached map[string]porkbun.TLDPricing
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, false
	}
	return cached, true
}

func fetchPricing(ctx context.Context, client *porkbun.Client) (map[string]porkbun.TLDPricing, error) {
	res, err := client.FetchPricing(ctx)
	if err != nil {
		return nil, err
	}
	return res.Pricing, nil
}

type tldTableModel struct {
	table table.Model
}

func (m tldTableModel) Init() tea.Cmd { return nil }

func (m tldTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tldTableModel) View() string {
	return "\n" + theme.Accent.Render("Porkbun TLD Pricing") + "\n\n" + m.table.View() + "\n\n  Press q or esc to quit.\n"
}
