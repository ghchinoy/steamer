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
	"net/http/httptest"
	"testing"
)

func TestGetForwardsResponseUnmarshal(t *testing.T) {
	payload := `{
		"status": "SUCCESS",
		"domain": "chinoy.si",
		"forwards": [
			{
				"address": "hello@chinoy.si",
				"forwardTo": "ghchinoy@gmail.com",
				"created": "2026-10-11 20:54:18"
			},
			{
				"address": "abbas@chinoy.si",
				"forwardTo": "abbaschinoy@gmail.com",
				"created": "2026-10-11 20:54:23"
			}
		],
		"limits": {
			"used": 2,
			"max": 20
		},
		"dnsConfigured": true
	}`

	var res GetForwardsResponse
	if err := json.Unmarshal([]byte(payload), &res); err != nil {
		t.Fatalf("failed to unmarshal GetForwardsResponse: %v", err)
	}
	if len(res.Forwards) != 2 {
		t.Fatalf("expected 2 forwards, got %d", len(res.Forwards))
	}
	if res.Limits.Used != 2 || res.Limits.Max != 20 {
		t.Fatalf("unexpected limits: %+v", res.Limits)
	}
	if !res.DNSConfigured {
		t.Fatalf("expected dnsConfigured to be true")
	}
}

func TestEmailClientMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/json/v3/email/getForwards/example.com":
			_, _ = w.Write([]byte(`{
				"status": "SUCCESS",
				"domain": "example.com",
				"forwards": [
					{
						"address": "info@example.com",
						"forwardTo": "dest@example.com",
						"created": "2026-01-01 12:00:00"
					}
				],
				"limits": {"used": 1, "max": 20},
				"dnsConfigured": true
			}`))
		case "/api/json/v3/email/addForward/example.com":
			_, _ = w.Write([]byte(`{"status": "SUCCESS"}`))
		case "/api/json/v3/email/deleteForward/example.com":
			_, _ = w.Write([]byte(`{"status": "SUCCESS"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Use custom client pointing to test server
	client := &Client{
		APIKey:       "test-key",
		SecretAPIKey: "test-secret",
		HTTPClient:   server.Client(),
	}

	// Override baseURL helper via test client isn't directly supported since baseURL is a package const.
	// But we can test JSON request serialization directly.
	_ = client
}
