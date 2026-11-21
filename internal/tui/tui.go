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

// Package tui provides the terminal user interface for steamer.
package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ghchinoy/steamer/internal/porkbun"
	"github.com/ghchinoy/steamer/internal/theme"
)

// DomainClient fetches domains and DNS records for the TUI.
type DomainClient interface {
	ListDomains(ctx context.Context) ([]porkbun.Domain, error)
	RetrieveRecords(ctx context.Context, domain string) ([]porkbun.DNSRecord, error)
}

type state int

const (
	viewDomains state = iota
	viewRecords
)

// Model is the Bubble Tea model for the steamer TUI.
type Model struct {
	client     DomainClient
	newContext func() (context.Context, context.CancelFunc)
	state      state
	domains    []porkbun.Domain
	records    []porkbun.DNSRecord
	cursor     int
	selected   int
	err        error
	loading    bool
	domain     string // currently viewed domain
}

// NewModel creates a new TUI model.
func NewModel(client DomainClient, initialDomain string, newContext func() (context.Context, context.CancelFunc)) Model {
	if newContext == nil {
		newContext = func() (context.Context, context.CancelFunc) {
			return context.TODO(), func() {}
		}
	}
	return Model{
		client:     client,
		newContext: newContext,
		state:      viewDomains,
		cursor:     0,
		loading:    true,
		domain:     initialDomain,
	}
}

// Init initializes the TUI.
func (m Model) Init() tea.Cmd {
	if m.domain != "" {
		return m.fetchRecords(m.domain)
	}
	return m.fetchDomains
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.state == viewDomains && m.cursor < len(m.domains)-1 {
				m.cursor++
			}
			if m.state == viewRecords && m.cursor < len(m.records)-1 {
				m.cursor++
			}
		case "enter":
			if m.state == viewDomains && len(m.domains) > 0 {
				m.selected = m.cursor
				m.domain = m.domains[m.selected].Domain
				m.state = viewRecords
				m.cursor = 0
				m.loading = true
				return m, m.fetchRecords(m.domain)
			}
		case "esc", "backspace":
			if m.state == viewRecords {
				m.state = viewDomains
				m.cursor = 0
				m.records = nil
			}
		}

	case domainsMsg:
		m.loading = false
		m.domains = msg
	case recordsMsg:
		m.loading = false
		m.records = msg
		m.state = viewRecords
	case errorMsg:
		m.loading = false
		m.err = msg
	}

	return m, nil
}

// View renders the TUI.
func (m Model) View() string {
	if m.err != nil {
		return theme.Fail.Render(fmt.Sprintf("Error: %v", m.err)) + "\n\nPress q to quit."
	}

	if m.loading {
		return theme.Muted.Render("Loading...")
	}

	if m.state == viewDomains {
		return m.renderDomainsView()
	}
	return m.renderRecordsView()
}

func (m Model) renderDomainsView() string {
	var b strings.Builder
	b.WriteString(theme.Title.Render("STEAMER - Porkbun Manager"))
	b.WriteString("\n\nYour Domains:\n\n")

	for i, d := range m.domains {
		line := fmt.Sprintf("  %s", d.Domain)
		if m.cursor == i {
			line = theme.Selected.Render(fmt.Sprintf("> %s", d.Domain))
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}

	b.WriteString("\n")
	b.WriteString(theme.Muted.Render("(j/k: navigate, enter: view records, q: quit)"))
	return b.String()
}

func (m Model) renderRecordsView() string {
	var b strings.Builder
	b.WriteString(theme.Title.Render("STEAMER - Porkbun Manager"))
	fmt.Fprintf(&b, "\n\nDNS Records for %s:\n\n", m.domain)
	fmt.Fprintf(&b, "%-10s %-25s %-10s %-30s\n", "ID", "NAME", "TYPE", "CONTENT")

	for i, r := range m.records {
		row := fmt.Sprintf("%-10s %-25s %-10s %-30s", string(r.ID), r.Name, r.Type, r.Content)
		line := "  " + row
		if m.cursor == i {
			line = theme.Selected.Render("> " + row)
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}

	b.WriteString("\n")
	b.WriteString(theme.Muted.Render("(j/k: navigate, esc: back to domains, q: quit)"))
	return b.String()
}

type (
	domainsMsg []porkbun.Domain
	recordsMsg []porkbun.DNSRecord
	errorMsg   error
)

func (m Model) fetchDomains() tea.Msg {
	ctx, cancel := m.newContext()
	defer cancel()

	domains, err := m.client.ListDomains(ctx)
	if err != nil {
		return errorMsg(err)
	}
	return domainsMsg(domains)
}

func (m Model) fetchRecords(domain string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := m.newContext()
		defer cancel()

		records, err := m.client.RetrieveRecords(ctx, domain)
		if err != nil {
			return errorMsg(err)
		}
		return recordsMsg(records)
	}
}
