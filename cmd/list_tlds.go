package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"

	"github.com/spf13/cobra"
)

var listTldsJSON bool
var listTldsForce bool

var listTldsCmd = &cobra.Command{
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
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, secretKey, err := getClientConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pricing, err := getCachedOrFetchPricing(apiKey, secretKey, listTldsForce)
		if err != nil {
			fmt.Printf("Error fetching TLD pricing: %v\n", err)
			os.Exit(1)
		}

		if listTldsJSON {
			b, err := json.MarshalIndent(pricing, "", "  ")
			if err != nil {
				fmt.Printf("Error encoding JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(b))
			return
		}

		// Sort TLDs alphabetically
		tlds := make([]string, 0, len(pricing))
		for tld := range pricing {
			tlds = append(tlds, tld)
		}
		sort.Strings(tlds)

		columns := []table.Column{
			{Title: "TLD", Width: 15},
			{Title: "REGISTRATION", Width: 15},
			{Title: "RENEWAL", Width: 15},
			{Title: "TRANSFER", Width: 15},
		}

		var rows []table.Row
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
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		m := tldTableModel{t}
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func getCachedOrFetchPricing(apiKey, secretKey string, force bool) (map[string]porkbun.TLDPricing, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return fetchPricing(apiKey, secretKey)
	}

	cacheDir := filepath.Join(home, ".config", "steamer")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return fetchPricing(apiKey, secretKey)
	}

	cacheFile := filepath.Join(cacheDir, "tlds.json")

	if !force {
		info, err := os.Stat(cacheFile)
		if err == nil {
			// Cache exists, check if it's less than 7 days old
			if time.Since(info.ModTime()) < 7*24*time.Hour {
				data, err := os.ReadFile(cacheFile)
				if err == nil {
					var cached map[string]porkbun.TLDPricing
					if err := json.Unmarshal(data, &cached); err == nil {
						return cached, nil
					}
				}
			}
		}
	}

	// Fetch fresh pricing
	pricing, err := fetchPricing(apiKey, secretKey)
	if err != nil {
		return nil, err
	}

	// Try to cache it, but don't fail if we can't
	if data, err := json.Marshal(pricing); err == nil {
		_ = os.WriteFile(cacheFile, data, 0600)
	}

	return pricing, nil
}

func fetchPricing(apiKey, secretKey string) (map[string]porkbun.TLDPricing, error) {
	client := porkbun.NewClient(apiKey, secretKey)
	res, err := client.GetPricing()
	if err != nil {
		return nil, err
	}
	return res.Pricing, nil
}

func init() {
	listTldsCmd.Flags().BoolVar(&listTldsJSON, "json", false, "Output results in JSON format")
	listTldsCmd.Flags().BoolVar(&listTldsForce, "force", false, "Force refresh the TLD cache")
	rootCmd.AddCommand(listTldsCmd)
}

type tldTableModel struct {
	table table.Model
}

func (m tldTableModel) Init() tea.Cmd { return nil }

func (m tldTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
