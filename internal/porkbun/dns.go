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
	"context"
	"fmt"
	"net/url"
)

// DNSRecord represents a single DNS record in Porkbun.
type DNSRecord struct {
	ID      FlexibleString `json:"id"`
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Content string         `json:"content"`
	TTL     string         `json:"ttl"`
	Prio    string         `json:"prio"`
	Notes   string         `json:"notes"`
}

// RetrieveDNSResponse is the response from the DNS retrieve endpoint.
type RetrieveDNSResponse struct {
	APIResponse
	Records []DNSRecord `json:"records"`
}

// CreateRecordRequest is the request body for creating a DNS record.
type CreateRecordRequest struct {
	BaseRequest
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl,omitempty"`
	Prio    string `json:"prio,omitempty"`
}

// CreateRecordResponse is the response from the DNS create endpoint.
type CreateRecordResponse struct {
	APIResponse
	ID FlexibleString `json:"id"`
}

// RetrieveRecords fetches all DNS records for the given domain.
func (c *Client) RetrieveRecords(ctx context.Context, domain string) ([]DNSRecord, error) {
	if err := validatePathSegment("domain", domain); err != nil {
		return nil, err
	}
	var res RetrieveDNSResponse
	endpoint := fmt.Sprintf("dns/retrieve/%s", url.PathEscape(domain))
	if err := c.post(ctx, endpoint, c.baseRequest(), &res); err != nil {
		return nil, err
	}
	if err := checkAPIStatus("api error", res.APIResponse); err != nil {
		return nil, err
	}
	return res.Records, nil
}

// CreateRecord creates a new DNS record for the given domain.
func (c *Client) CreateRecord(ctx context.Context, domain string, record CreateRecordRequest) (string, error) {
	if err := validatePathSegment("domain", domain); err != nil {
		return "", err
	}
	record.BaseRequest = c.baseRequest()

	var res CreateRecordResponse
	endpoint := fmt.Sprintf("dns/create/%s", url.PathEscape(domain))
	if err := c.post(ctx, endpoint, record, &res); err != nil {
		return "", err
	}
	if err := checkAPIStatus("api error", res.APIResponse); err != nil {
		return "", err
	}
	return string(res.ID), nil
}

// DeleteRecord deletes the specified DNS record from the given domain.
func (c *Client) DeleteRecord(ctx context.Context, domain, recordID string) error {
	if err := validatePathSegment("domain", domain); err != nil {
		return err
	}
	if err := validatePathSegment("record ID", recordID); err != nil {
		return err
	}
	var res APIResponse
	endpoint := fmt.Sprintf("dns/delete/%s/%s", url.PathEscape(domain), url.PathEscape(recordID))
	if err := c.post(ctx, endpoint, c.baseRequest(), &res); err != nil {
		return err
	}
	return checkAPIStatus("delete record failed", res)
}
