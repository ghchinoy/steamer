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
	"testing"
)

func TestListDomainsResponseUnmarshal(t *testing.T) {
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
				"labels": [
					{
						"id": 123,
						"title": "prod",
						"color": "#ff0000"
					}
				]
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
				"labels": [
					{
						"id": "456",
						"title": "dev",
						"color": "#00ff00"
					}
				]
			}
		]
	}`

	var res ListDomainsResponse
	if err := json.Unmarshal([]byte(payload), &res); err != nil {
		t.Fatalf("failed to unmarshal ListDomainsResponse with mixed types: %v", err)
	}
	if len(res.Domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(res.Domains))
	}
}
