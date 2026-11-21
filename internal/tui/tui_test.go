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
	"context"
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ghchinoy/steamer/internal/porkbun"
)

func TestModel_NavigationAndViews(t *testing.T) {
	fake := &fakeDomainClient{
		domains: []porkbun.Domain{
			{Domain: "alpha.example.com"},
			{Domain: "beta.example.com"},
		},
		records: map[string][]porkbun.DNSRecord{
			"beta.example.com": {
				{ID: "42", Name: "www.beta.example.com", Type: "A", Content: "192.0.2.1"},
			},
		},
	}

	ctxFactory := func() (context.Context, context.CancelFunc) {
		return context.WithCancel(t.Context())
	}

	m := NewModel(fake, "", ctxFactory)
	if !strings.Contains(m.View(), "Loading...") {
		t.Errorf("Model.View() while loading = %q, want %q", m.View(), "Loading...")
	}

	// Simulate Init fetching domains.
	initCmd := m.Init()
	if initCmd == nil {
		t.Fatalf("Model.Init() = nil, want fetchDomains command")
	}
	msg := initCmd()
	updated, _ := m.Update(msg)
	m = updated.(Model)

	view := m.View()
	if !strings.Contains(view, "alpha.example.com") || !strings.Contains(view, "beta.example.com") {
		t.Errorf("Model.View() after domainsMsg = %q, want both domains listed", view)
	}

	// Move cursor down to beta.example.com and press enter.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if got, want := m.cursor, 1; got != want {
		t.Errorf("Model.cursor after 'j' = %d, want %d", got, want)
	}

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if cmd == nil {
		t.Fatalf("Model.Update(enter) returned nil cmd, want fetchRecords cmd")
	}

	updated, _ = m.Update(cmd())
	m = updated.(Model)
	recordsView := m.View()
	if !strings.Contains(recordsView, "DNS Records for beta.example.com") || !strings.Contains(recordsView, "192.0.2.1") {
		t.Errorf("Model.View() after recordsMsg = %q, want records for beta.example.com", recordsView)
	}

	// Press esc to return to domains view.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if got, want := m.state, viewDomains; got != want {
		t.Errorf("Model.state after esc = %v, want %v", got, want)
	}

	// Press q to quit.
	_, quitCmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if quitCmd == nil {
		t.Errorf("Model.Update('q') returned nil cmd, want tea.Quit")
	}
}

func TestModel_ErrorState(t *testing.T) {
	fake := &fakeDomainClient{err: errors.New("api unavailable")}
	m := NewModel(fake, "", func() (context.Context, context.CancelFunc) {
		return context.WithCancel(t.Context())
	})

	updated, _ := m.Update(m.Init()())
	m = updated.(Model)
	if !strings.Contains(m.View(), "api unavailable") {
		t.Errorf("Model.View() after errorMsg = %q, want error containing %q", m.View(), "api unavailable")
	}
}

type fakeDomainClient struct {
	domains []porkbun.Domain
	records map[string][]porkbun.DNSRecord
	err     error
}

func (f *fakeDomainClient) ListDomains(_ context.Context) ([]porkbun.Domain, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.domains, nil
}

func (f *fakeDomainClient) RetrieveRecords(_ context.Context, domain string) ([]porkbun.DNSRecord, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.records[domain], nil
}
