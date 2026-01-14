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

package tui

import (
	"fmt"
	"steamer/internal/porkbun"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)
)

type state int

const (
	viewDomains state = iota
	viewRecords
)

type Model struct {
	client    *porkbun.Client
	state     state
	domains   []porkbun.Domain
	records   []porkbun.DNSRecord
	cursor    int
	selected  int
	err       error
	loading   bool
	domain    string // currently viewed domain
}

func NewModel(client *porkbun.Client, initialDomain string) Model {
	return Model{
		client:  client,
		state:   viewDomains,
		cursor:  0,
		loading: true,
		domain:  initialDomain,
	}
}

func (m Model) Init() tea.Cmd {
	if m.domain != "" {
		return m.fetchRecords(m.domain)
	}
	return m.fetchDomains
}

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
			} else if m.state == viewRecords && m.cursor < len(m.records)-1 {
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

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}

	if m.loading {
		return "Loading..."
	}

	s := titleStyle.Render("STEAMER - Porkbun Manager") + "\n\n"

	if m.state == viewDomains {
		s += "Your Domains:\n\n"
		for i, d := range m.domains {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s += selectedStyle.Render(fmt.Sprintf("%s %s", cursor, d.Domain)) + "\n"
			} else {
				s += fmt.Sprintf("%s %s", cursor, d.Domain) + "\n"
			}
		}
		s += "\n(j/k: navigate, enter: view records, q: quit)"
	} else {
		s += fmt.Sprintf("DNS Records for %s:\n\n", m.domain)
		s += fmt.Sprintf("%-10s %-25s %-10s %-30s\n", "ID", "NAME", "TYPE", "CONTENT")
		for i, r := range m.records {
			line := fmt.Sprintf("%-10v %-25s %-10s %-30s", r.ID, r.Name, r.Type, r.Content)
			if m.cursor == i {
				s += selectedStyle.Render("> "+line) + "\n"
			} else {
				s += "  " + line + "\n"
			}
		}
		s += "\n(j/k: navigate, esc: back to domains, q: quit)"
	}

	return s
}

// Commands and Messages
type domainsMsg []porkbun.Domain
type recordsMsg []porkbun.DNSRecord
type errorMsg error

func (m Model) fetchDomains() tea.Msg {
	domains, err := m.client.ListDomains()
	if err != nil {
		return errorMsg(err)
	}
	return domainsMsg(domains)
}

func (m Model) fetchRecords(domain string) tea.Cmd {
	return func() tea.Msg {
		records, err := m.client.RetrieveRecords(domain)
		if err != nil {
			return errorMsg(err)
		}
		return recordsMsg(records)
	}
}
