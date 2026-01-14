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

import "fmt"

type DNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
	Prio    string `json:"prio"`
	Notes   string `json:"notes"`
}

type RetrieveDNSResponse struct {
	APIResponse
	Records []DNSRecord `json:"records"`
}

type CreateRecordRequest struct {
	BaseRequest
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl,omitempty"`
	Prio    string `json:"prio,omitempty"`
}

type CreateRecordResponse struct {
	APIResponse
	ID string `json:"id"`
}

func (c *Client) RetrieveRecords(domain string) ([]DNSRecord, error) {
	req := BaseRequest{
		APIKey:       c.APIKey,
		SecretAPIKey: c.SecretAPIKey,
	}
	var res RetrieveDNSResponse
	endpoint := fmt.Sprintf("dns/retrieve/%s", domain)
	err := c.post(endpoint, req, &res)
	if err != nil {
		return nil, err
	}
	return res.Records, nil
}

func (c *Client) CreateRecord(domain string, record CreateRecordRequest) (string, error) {
	record.APIKey = c.APIKey
	record.SecretAPIKey = c.SecretAPIKey
	
	var res CreateRecordResponse
	endpoint := fmt.Sprintf("dns/create/%s", domain)
	err := c.post(endpoint, record, &res)
	if err != nil {
		return "", err
	}
	if res.Status != "SUCCESS" {
		return "", fmt.Errorf("create record failed: %s", res.Message)
	}
	return res.ID, nil
}

func (c *Client) DeleteRecord(domain, id string) error {
	req := BaseRequest{
		APIKey:       c.APIKey,
		SecretAPIKey: c.SecretAPIKey,
	}
	var res APIResponse
	endpoint := fmt.Sprintf("dns/delete/%s/%s", domain, id)
	err := c.post(endpoint, req, &res)
	if err != nil {
		return err
	}
	if res.Status != "SUCCESS" {
		return fmt.Errorf("delete record failed: %s", res.Message)
	}
	return nil
}
