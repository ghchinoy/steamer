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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ghchinoy/steamer/internal/porkbun"
)

func TestNewRootCmd_SubcommandsRegistered(t *testing.T) {
	root := NewRootCmd()
	wantPaths := [][]string{
		{"add-a"},
		{"add-aaaa"},
		{"add-cname"},
		{"add-ns"},
		{"add-txt"},
		{"list-domains"},
		{"list-records"},
		{"list-tlds"},
		{"rm"},
		{"search"},
		{"tui"},
		{"version"},
		{"email", "list"},
		{"email", "add"},
		{"email", "rm"},
	}

	for _, path := range wantPaths {
		cmd, _, err := root.Find(path)
		if err != nil {
			t.Errorf("root.Find(%v) unexpected error: %v", path, err)
			continue
		}
		if cmd == nil || cmd == root {
			t.Errorf("root.Find(%v) = %v, want registered subcommand", path, cmd)
		}
	}
}

func TestVersionCmd(t *testing.T) {
	root := NewRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"version", "--json"})

	if err := root.ExecuteContext(t.Context()); err != nil {
		t.Fatalf("version --json unexpected error: %v", err)
	}

	var info map[string]string
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		t.Fatalf("json.Unmarshal(version output %q) unexpected error: %v", out.String(), err)
	}
	if info["version"] == "" {
		t.Errorf("version --json output = %+v, want non-empty version field", info)
	}

	flagRoot := NewRootCmd()
	var flagOut bytes.Buffer
	flagRoot.SetOut(&flagOut)
	flagRoot.SetErr(&bytes.Buffer{})
	flagRoot.SetArgs([]string{"--version"})
	if err := flagRoot.ExecuteContext(t.Context()); err != nil {
		t.Fatalf("--version unexpected error: %v", err)
	}
	if !strings.Contains(flagOut.String(), "steamer version") {
		t.Errorf("--version output = %q, want prefix %q", flagOut.String(), "steamer version")
	}
}

func TestRootCmd_MissingExplicitConfigReturnsError(t *testing.T) {
	root := NewRootCmd()
	missingPath := filepath.Join(t.TempDir(), "does-not-exist.yaml")
	root.SetArgs([]string{"--config", missingPath, "list-domains"})
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})

	err := root.ExecuteContext(t.Context())
	if err == nil {
		t.Fatalf("ExecuteContext(--config %q) = nil, want error", missingPath)
	}
	if !strings.Contains(err.Error(), "reading config file") {
		t.Errorf("ExecuteContext(--config %q) error = %v, want error mentioning reading config file", missingPath, err)
	}
}

func TestAddRecordCmd_RootDomainFormattingAndJSON(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dns/create/example.com" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"SUCCESS","id":586562066}`))
	})

	// Root domain text output should not contain a leading dot (".example.com").
	out, _, err := executeWithMockHandler(t, handler, nil, "add-a", "example.com", "", "192.0.2.1")
	if err != nil {
		t.Fatalf("add-a root domain unexpected error: %v", err)
	}
	if strings.Contains(out, ".example.com") {
		t.Errorf("add-a root domain output = %q, must not have leading dot", out)
	}
	if !strings.Contains(out, "example.com") || !strings.Contains(out, "586562066") {
		t.Errorf("add-a root domain output = %q, want domain and record ID", out)
	}

	// NS JSON output preserves nameserver field.
	jsonOut, _, err := executeWithMockHandler(t, handler, nil, "add-ns", "example.com", "sub", "ns1.example.net", "--json")
	if err != nil {
		t.Fatalf("add-ns --json unexpected error: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal([]byte(jsonOut), &parsed); err != nil {
		t.Fatalf("json.Unmarshal(add-ns output %q) unexpected error: %v", jsonOut, err)
	}
	if got, want := parsed["nameserver"], "ns1.example.net"; got != want {
		t.Errorf("add-ns --json nameserver = %q, want %q", got, want)
	}
	if got, want := parsed["id"], "586562066"; got != want {
		t.Errorf("add-ns --json id = %q, want %q", got, want)
	}
}

func TestRmCmd_ConfirmationGuard(t *testing.T) {
	deleted := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dns/delete/example.com/123" {
			deleted++
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"SUCCESS"}`))
			return
		}
		http.NotFound(w, r)
	})

	// Declining confirmation should abort without calling DELETE.
	_, errOut, err := executeWithMockHandler(t, handler, strings.NewReader("n\n"), "rm", "example.com", "123")
	if err != nil {
		t.Fatalf("rm declined unexpected error: %v", err)
	}
	if deleted != 0 {
		t.Errorf("rm declined deleted count = %d, want 0", deleted)
	}
	if !strings.Contains(errOut, "Aborted") {
		t.Errorf("rm declined stderr = %q, want %q", errOut, "Aborted")
	}

	// Answering yes executes deletion.
	out, _, err := executeWithMockHandler(t, handler, strings.NewReader("y\n"), "rm", "example.com", "123")
	if err != nil {
		t.Fatalf("rm confirmed unexpected error: %v", err)
	}
	if got, want := deleted, 1; got != want {
		t.Errorf("rm confirmed deleted count = %d, want %d", got, want)
	}
	if !strings.Contains(out, "123") {
		t.Errorf("rm confirmed stdout = %q, want record ID %q", out, "123")
	}

	// --yes bypasses prompt and supports --json.
	jsonOut, _, err := executeWithMockHandler(t, handler, nil, "rm", "example.com", "123", "--yes", "--json")
	if err != nil {
		t.Fatalf("rm --yes --json unexpected error: %v", err)
	}
	if got, want := deleted, 2; got != want {
		t.Errorf("rm --yes deleted count = %d, want %d", got, want)
	}
	if !strings.Contains(jsonOut, `"status": "deleted"`) {
		t.Errorf("rm --yes --json stdout = %q, want status deleted", jsonOut)
	}
}

func TestSearchCmd_JSONSerializesErrors(t *testing.T) {
	origWait := searchRateLimitWait
	searchRateLimitWait = 0
	t.Cleanup(func() { searchRateLimitWait = origWait })

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/domain/checkDomain/good.com":
			_, _ = w.Write([]byte(`{"status":"SUCCESS","response":{"avail":"yes","premium":"no","price":"9.73"}}`))
		case "/domain/checkDomain/good.net":
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"status":"ERROR","message":"too many requests"}`))
		default:
			http.NotFound(w, r)
		}
	})

	out, _, err := executeWithMockHandler(t, handler, nil, "search", "good", "--tlds", "com,net", "--json")
	if err != nil {
		t.Fatalf("search --json unexpected error: %v", err)
	}

	var results []searchResult
	if err := json.Unmarshal([]byte(out), &results); err != nil {
		t.Fatalf("json.Unmarshal(search output %q) unexpected error: %v", out, err)
	}
	if got, want := len(results), 2; got != want {
		t.Fatalf("len(search results) = %d, want %d", got, want)
	}
	if results[0].Error != "" || results[0].Response == nil || results[0].Response.Response.Avail != "yes" {
		t.Errorf("results[0] = %+v, want available domain without error", results[0])
	}
	if !strings.Contains(results[1].Error, "too many requests") {
		t.Errorf("results[1].Error = %q, want error containing %q", results[1].Error, "too many requests")
	}
}

func TestEmailCmd_ListAndRmConfirmation(t *testing.T) {
	deleted := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/email/getForwards/example.com":
			_, _ = w.Write([]byte(`{
				"status":"SUCCESS",
				"domain":"example.com",
				"forwards":[{"address":"hello@example.com","forwardTo":"alice@example.org","created":"2026-01-01"}],
				"limits":{"used":1,"max":20},
				"dnsConfigured":true
			}`))
		case "/email/deleteForward/example.com":
			deleted++
			_, _ = w.Write([]byte(`{"status":"SUCCESS"}`))
		default:
			http.NotFound(w, r)
		}
	})

	out, _, err := executeWithMockHandler(t, handler, nil, "email", "list", "example.com")
	if err != nil {
		t.Fatalf("email list unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello@example.com") || !strings.Contains(out, "1/20 used") {
		t.Errorf("email list stdout = %q, want forward and quota", out)
	}

	_, _, err = executeWithMockHandler(t, handler, strings.NewReader("n\n"), "email", "rm", "example.com", "hello", "alice@example.org")
	if err != nil {
		t.Fatalf("email rm declined unexpected error: %v", err)
	}
	if deleted != 0 {
		t.Errorf("email rm declined deleted count = %d, want 0", deleted)
	}

	_, _, err = executeWithMockHandler(t, handler, nil, "email", "rm", "example.com", "hello", "alice@example.org", "--yes")
	if err != nil {
		t.Fatalf("email rm --yes unexpected error: %v", err)
	}
	if got, want := deleted, 1; got != want {
		t.Errorf("email rm --yes deleted count = %d, want %d", got, want)
	}
}

type testRoundTripper struct {
	handler http.Handler
}

func (rt testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	rt.handler.ServeHTTP(rec, req)
	return rec.Result(), nil
}

func executeWithMockHandler(t *testing.T, handler http.Handler, stdin *strings.Reader, args ...string) (string, string, error) {
	t.Helper()

	app := &cliApp{
		v: viper.New(),
		newClient: func() (*porkbun.Client, error) {
			return &porkbun.Client{
				APIKey:       "test-key",
				SecretAPIKey: "test-secret",
				BaseURL:      "https://porkbun.test",
				HTTPClient: &http.Client{
					Transport: testRoundTripper{handler: handler},
				},
			}, nil
		},
	}

	root := &cobra.Command{
		Use:           "steamer",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddGroup(
		&cobra.Group{ID: GroupInfo, Title: "Information Commands:"},
		&cobra.Group{ID: GroupManagement, Title: "Record Management Commands:"},
		&cobra.Group{ID: GroupTUI, Title: "Interactive Commands:"},
		&cobra.Group{ID: GroupEmail, Title: "Email Forwarding Commands:"},
	)
	root.AddCommand(
		newAddACmd(app),
		newAddAAAACmd(app),
		newAddCNAMECmd(app),
		newAddNSCmd(app),
		newAddTXTCmd(app),
		newEmailCmd(app),
		newListDomainsCmd(app),
		newListRecordsCmd(app),
		newListTLDsCmd(app),
		newRmCmd(app),
		newSearchCmd(app),
		newTUICmd(app),
	)

	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	if stdin != nil {
		root.SetIn(stdin)
	}
	root.SetArgs(args)

	err := root.ExecuteContext(t.Context())
	return stdout.String(), stderr.String(), err
}
