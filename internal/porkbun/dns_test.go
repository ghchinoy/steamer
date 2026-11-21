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

package porkbun

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestClient_PingAndDNSEndpoints(t *testing.T) {
	var createdReq CreateRecordRequest

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/ping":
			_, _ = w.Write([]byte(`{"status":"SUCCESS","yourIp":"192.0.2.10"}`))
		case "/dns/retrieve/example.com":
			_, _ = w.Write([]byte(`{"status":"SUCCESS","records":[{"id":586562066,"name":"www.example.com","type":"A","content":"192.0.2.1"}]}`))
		case "/dns/retrieve/error.com":
			_, _ = w.Write([]byte(`{"status":"ERROR","message":"domain not found"}`))
		case "/dns/create/example.com":
			if err := json.NewDecoder(r.Body).Decode(&createdReq); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`{"status":"SUCCESS","id":586562066}`))
		case "/dns/delete/example.com/586562066":
			_, _ = w.Write([]byte(`{"status":"SUCCESS"}`))
		case "/dns/delete/example.com/999":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"status":"ERROR","message":"invalid record id"}`))
		case "/dns/retrieve/badgateway.com":
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(strings.Repeat("<html>bad gateway</html>", 30)))
		default:
			http.NotFound(w, r)
		}
	})

	client := newTestClient(handler)
	ctx := t.Context()

	ping, err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping(ctx) unexpected error: %v", err)
	}
	if got, want := ping.YourIP, "192.0.2.10"; got != want {
		t.Errorf("Ping(ctx).YourIP = %q, want %q", got, want)
	}

	records, err := client.RetrieveRecords(ctx, "example.com")
	if err != nil {
		t.Fatalf("RetrieveRecords(ctx, %q) unexpected error: %v", "example.com", err)
	}
	if got, want := len(records), 1; got != want {
		t.Fatalf("len(RetrieveRecords(ctx, %q)) = %d, want %d", "example.com", got, want)
	}
	if got, want := string(records[0].ID), "586562066"; got != want {
		t.Errorf("RetrieveRecords(ctx, %q)[0].ID = %q, want %q", "example.com", got, want)
	}

	if _, err := client.RetrieveRecords(ctx, "error.com"); err == nil || !strings.Contains(err.Error(), "domain not found") {
		t.Errorf("RetrieveRecords(ctx, %q) error = %v, want error containing %q", "error.com", err, "domain not found")
	}

	id, err := client.CreateRecord(ctx, "example.com", CreateRecordRequest{
		Name:    "www",
		Type:    "A",
		Content: "192.0.2.1",
	})
	if err != nil {
		t.Fatalf("CreateRecord(ctx, %q, ...) unexpected error: %v", "example.com", err)
	}
	if got, want := id, "586562066"; got != want {
		t.Errorf("CreateRecord(ctx, %q, ...) = %q, want %q", "example.com", got, want)
	}
	if got, want := createdReq.APIKey, "test-key"; got != want {
		t.Errorf("CreateRecord injected APIKey = %q, want %q", got, want)
	}

	if err := client.DeleteRecord(ctx, "example.com", "586562066"); err != nil {
		t.Fatalf("DeleteRecord(ctx, %q, %q) unexpected error: %v", "example.com", "586562066", err)
	}

	if err := client.DeleteRecord(ctx, "example.com", "999"); err == nil || !strings.Contains(err.Error(), "invalid record id") {
		t.Errorf("DeleteRecord(ctx, %q, %q) error = %v, want error containing %q", "example.com", "999", err, "invalid record id")
	}

	_, err = client.RetrieveRecords(ctx, "badgateway.com")
	if err == nil {
		t.Fatalf("RetrieveRecords(ctx, %q) = nil error, want error", "badgateway.com")
	}
	if len(err.Error()) > maxErrorBodyBytes+64 {
		t.Errorf("RetrieveRecords(ctx, %q) error length = %d, want truncated body <= %d", "badgateway.com", len(err.Error()), maxErrorBodyBytes+64)
	}
}

func TestValidatePathSegment(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid_domain", value: "example.com", wantErr: false},
		{name: "valid_numeric_id", value: "586562066", wantErr: false},
		{name: "empty", value: "", wantErr: true},
		{name: "path_traversal", value: "../dns/delete/other.com/1", wantErr: true},
		{name: "slash", value: "example.com/extra", wantErr: true},
		{name: "query_injection", value: "example.com?foo=bar", wantErr: true},
		{name: "fragment", value: "example.com#frag", wantErr: true},
		{name: "control_char", value: "example.com\n", wantErr: true},
		{name: "space", value: "example .com", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePathSegment("domain", tc.value)
			if tc.wantErr && err == nil {
				t.Errorf("validatePathSegment(%q, %q) = nil, want error", "domain", tc.value)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("validatePathSegment(%q, %q) = %v, want nil", "domain", tc.value, err)
			}
		})
	}
}
