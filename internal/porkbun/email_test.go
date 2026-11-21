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
	"testing"
)

func TestEmailForwardsResponse_Unmarshal(t *testing.T) {
	payload := `{
		"status": "SUCCESS",
		"domain": "example.com",
		"forwards": [
			{
				"address": "hello@example.com",
				"forwardTo": "alice@example.org",
				"created": "2026-10-11 20:54:18"
			},
			{
				"address": "support@example.com",
				"forwardTo": "bob@example.org",
				"created": "2026-10-11 20:54:23"
			}
		],
		"limits": {
			"used": 2,
			"max": 20
		},
		"dnsConfigured": true
	}`

	var res EmailForwardsResponse
	if err := json.Unmarshal([]byte(payload), &res); err != nil {
		t.Fatalf("json.Unmarshal(EmailForwardsResponse) unexpected error: %v", err)
	}
	if got, want := len(res.Forwards), 2; got != want {
		t.Fatalf("len(EmailForwardsResponse.Forwards) = %d, want %d", got, want)
	}
	if res.Limits.Used != 2 || res.Limits.Max != 20 {
		t.Errorf("EmailForwardsResponse.Limits = %+v, want {Used:2 Max:20}", res.Limits)
	}
	if !res.DNSConfigured {
		t.Errorf("EmailForwardsResponse.DNSConfigured = %v, want true", res.DNSConfigured)
	}
}

func TestClient_EmailEndpoints(t *testing.T) {
	var lastAddReq AddForwardRequest
	var lastDeleteReq DeleteForwardRequest

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/email/getForwards/example.com":
			_, _ = w.Write([]byte(`{
				"status": "SUCCESS",
				"domain": "example.com",
				"forwards": [
					{
						"address": "info@example.com",
						"forwardTo": "dest@example.org",
						"created": "2026-01-01 12:00:00"
					}
				],
				"limits": {"used": 1, "max": 20},
				"dnsConfigured": true
			}`))
		case "/email/addForward/example.com":
			if err := json.NewDecoder(r.Body).Decode(&lastAddReq); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`{"status": "SUCCESS"}`))
		case "/email/deleteForward/example.com":
			if err := json.NewDecoder(r.Body).Decode(&lastDeleteReq); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`{"status": "SUCCESS"}`))
		default:
			http.NotFound(w, r)
		}
	})

	client := newTestClient(handler)
	ctx := t.Context()

	forwards, err := client.FetchEmailForwards(ctx, "example.com")
	if err != nil {
		t.Fatalf("FetchEmailForwards(ctx, %q) unexpected error: %v", "example.com", err)
	}
	if got, want := len(forwards.Forwards), 1; got != want {
		t.Fatalf("len(FetchEmailForwards(ctx, %q).Forwards) = %d, want %d", "example.com", got, want)
	}
	if got, want := forwards.Forwards[0].Address, "info@example.com"; got != want {
		t.Errorf("FetchEmailForwards(ctx, %q).Forwards[0].Address = %q, want %q", "example.com", got, want)
	}

	if err := client.AddEmailForward(ctx, "example.com", "hello@example.com", "dest@example.org"); err != nil {
		t.Fatalf("AddEmailForward(ctx, %q, %q, %q) unexpected error: %v", "example.com", "hello@example.com", "dest@example.org", err)
	}
	if got, want := lastAddReq.Address, "hello"; got != want {
		t.Errorf("AddEmailForward normalized Address = %q, want %q", got, want)
	}

	if err := client.DeleteEmailForward(ctx, "example.com", "hello@example.com", "dest@example.org"); err != nil {
		t.Fatalf("DeleteEmailForward(ctx, %q, %q, %q) unexpected error: %v", "example.com", "hello@example.com", "dest@example.org", err)
	}
	if got, want := lastDeleteReq.Address, "hello"; got != want {
		t.Errorf("DeleteEmailForward normalized Address = %q, want %q", got, want)
	}
}
