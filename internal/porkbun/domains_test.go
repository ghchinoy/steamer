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
	"strings"
	"testing"
)

func TestFlexibleString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    FlexibleString
		wantErr bool
	}{
		{name: "string", input: `"586562066"`, want: "586562066"},
		{name: "large_number", input: `586562066`, want: "586562066"},
		{name: "zero", input: `0`, want: "0"},
		{name: "null", input: `null`, want: ""},
		{name: "invalid_string", input: `"unterminated`, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got FlexibleString
			err := got.UnmarshalJSON([]byte(tc.input))
			if tc.wantErr {
				if err == nil {
					t.Errorf("FlexibleString.UnmarshalJSON(%q) = nil error, want error", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("FlexibleString.UnmarshalJSON(%q) unexpected error: %v", tc.input, err)
			}
			if got != tc.want {
				t.Errorf("FlexibleString.UnmarshalJSON(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestListDomainsResponse_Unmarshal(t *testing.T) {
	payload := `{
		"status": "SUCCESS",
		"domains": [
			{
				"domain": "example.com",
				"status": "ACTIVE",
				"tld": "com",
				"createDate": "2025-01-01 00:00:00",
				"expireDate": "2026-01-01 00:00:00",
				"securityLock": 1,
				"whoisPrivacy": 1,
				"autoRenew": 1,
				"notLocal": 0,
				"labels": [{"id": 123, "title": "prod", "color": "#ff0000"}]
			},
			{
				"domain": "example.dev",
				"status": "ACTIVE",
				"tld": "dev",
				"createDate": "2025-01-01 00:00:00",
				"expireDate": "2026-01-01 00:00:00",
				"securityLock": "1",
				"whoisPrivacy": "1",
				"autoRenew": "1",
				"notLocal": "0",
				"labels": [{"id": "456", "title": "dev", "color": "#00ff00"}]
			}
		]
	}`

	var res ListDomainsResponse
	if err := json.Unmarshal([]byte(payload), &res); err != nil {
		t.Fatalf("json.Unmarshal(ListDomainsResponse) unexpected error: %v", err)
	}
	if got, want := len(res.Domains), 2; got != want {
		t.Fatalf("len(ListDomainsResponse.Domains) = %d, want %d", got, want)
	}
	if got, want := string(res.Domains[0].Labels[0].ID), "123"; got != want {
		t.Errorf("Domains[0].Labels[0].ID = %q, want %q", got, want)
	}
	if got, want := string(res.Domains[1].SecurityLock), "1"; got != want {
		t.Errorf("Domains[1].SecurityLock = %q, want %q", got, want)
	}
}

func TestCreateRecordResponse_IDFormatting(t *testing.T) {
	payload := `{"status":"SUCCESS","id":586562066}`
	var res CreateRecordResponse
	if err := json.Unmarshal([]byte(payload), &res); err != nil {
		t.Fatalf("json.Unmarshal(%q) unexpected error: %v", payload, err)
	}
	if got, want := string(res.ID), "586562066"; got != want {
		t.Errorf("CreateRecordResponse.ID for %q = %q, want %q", payload, got, want)
	}
}

func TestClient_DomainsEndpoints(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/domain/listAll":
			_, _ = w.Write([]byte(`{"status":"SUCCESS","domains":[{"domain":"example.com","status":"ACTIVE","tld":"com","expireDate":"2027-01-01"}]}`))
		case "/domain/checkDomain/example.com":
			_, _ = w.Write([]byte(`{"status":"SUCCESS","response":{"avail":"yes","premium":"no","price":"9.73"}}`))
		case "/domain/checkDomain/api-error.com":
			_, _ = w.Write([]byte(`{"status":"ERROR","message":"rate limit exceeded"}`))
		case "/pricing/get":
			_, _ = w.Write([]byte(`{"status":"SUCCESS","pricing":{"com":{"registration":"9.73","renewal":"10.37","transfer":"9.73"}}}`))
		default:
			http.NotFound(w, r)
		}
	})

	client := newTestClient(handler)
	ctx := t.Context()

	domains, err := client.ListDomains(ctx)
	if err != nil {
		t.Fatalf("ListDomains(ctx) unexpected error: %v", err)
	}
	if got, want := len(domains), 1; got != want {
		t.Fatalf("len(ListDomains(ctx)) = %d, want %d", got, want)
	}
	if got, want := domains[0].Domain, "example.com"; got != want {
		t.Errorf("ListDomains(ctx)[0].Domain = %q, want %q", got, want)
	}

	check, err := client.CheckDomain(ctx, "example.com")
	if err != nil {
		t.Fatalf("CheckDomain(ctx, %q) unexpected error: %v", "example.com", err)
	}
	if got, want := check.Response.Avail, "yes"; got != want {
		t.Errorf("CheckDomain(ctx, %q).Response.Avail = %q, want %q", "example.com", got, want)
	}

	if _, err := client.CheckDomain(ctx, "api-error.com"); err == nil || !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Errorf("CheckDomain(ctx, %q) error = %v, want error containing %q", "api-error.com", err, "rate limit exceeded")
	}

	pricing, err := client.FetchPricing(ctx)
	if err != nil {
		t.Fatalf("FetchPricing(ctx) unexpected error: %v", err)
	}
	if got, want := pricing.Pricing["com"].Registration, "9.73"; got != want {
		t.Errorf("FetchPricing(ctx).Pricing[%q].Registration = %q, want %q", "com", got, want)
	}
}

type handlerTransport struct {
	handler http.Handler
}

func (ht handlerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	ht.handler.ServeHTTP(rec, req)
	return rec.Result(), nil
}

func newTestClient(handler http.Handler) *Client {
	return &Client{
		APIKey:       "test-key",
		SecretAPIKey: "test-secret",
		BaseURL:      "https://porkbun.test",
		HTTPClient: &http.Client{
			Transport: handlerTransport{handler: handler},
		},
	}
}
