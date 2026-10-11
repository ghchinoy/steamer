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
	"fmt"
	"strings"
)

// EmailForward represents a configured email forwarding rule in Porkbun.
type EmailForward struct {
	ID        FlexibleString `json:"id,omitempty"`
	Address   string         `json:"address"`
	ForwardTo string         `json:"forwardTo"`
	Created   string         `json:"created,omitempty"`
}

// ForwardLimits indicates used and max email forwarding rules allowed on the domain.
type ForwardLimits struct {
	Used int `json:"used"`
	Max  int `json:"max"`
}

// GetForwardsResponse is the response from the email/getForwards endpoint.
type GetForwardsResponse struct {
	APIResponse
	Domain        string         `json:"domain"`
	Forwards      []EmailForward `json:"forwards"`
	Limits        ForwardLimits  `json:"limits"`
	DNSConfigured bool           `json:"dnsConfigured"`
}

// AddForwardRequest is the request body for adding an email forward.
type AddForwardRequest struct {
	BaseRequest
	Address   string `json:"address"`
	ForwardTo string `json:"forwardTo"`
}

// DeleteForwardRequest is the request body for deleting an email forward.
type DeleteForwardRequest struct {
	BaseRequest
	Address   string `json:"address"`
	ForwardTo string `json:"forwardTo"`
}

// GetEmailForwards retrieves all email forwarding rules for a domain.
func (c *Client) GetEmailForwards(domain string) (*GetForwardsResponse, error) {
	req := BaseRequest{
		APIKey:       c.APIKey,
		SecretAPIKey: c.SecretAPIKey,
	}
	var res GetForwardsResponse
	endpoint := fmt.Sprintf("email/getForwards/%s", domain)
	err := c.post(endpoint, req, &res)
	if err != nil {
		return nil, err
	}
	if res.Status != "SUCCESS" {
		msg := res.Message
		if msg == "" {
			msg = "unknown error"
		}
		return nil, fmt.Errorf("api error: %s", msg)
	}
	return &res, nil
}

// AddEmailForward creates a new email forwarding rule for the domain.
func (c *Client) AddEmailForward(domain, address, forwardTo string) error {
	// Normalize address by trimming domain suffix if user supplied it
	cleanAddress := strings.TrimSuffix(address, "@"+domain)

	req := AddForwardRequest{
		BaseRequest: BaseRequest{
			APIKey:       c.APIKey,
			SecretAPIKey: c.SecretAPIKey,
		},
		Address:   cleanAddress,
		ForwardTo: forwardTo,
	}
	var res APIResponse
	endpoint := fmt.Sprintf("email/addForward/%s", domain)
	err := c.post(endpoint, req, &res)
	if err != nil {
		return err
	}
	if res.Status != "SUCCESS" {
		msg := res.Message
		if msg == "" {
			msg = "unknown error"
		}
		return fmt.Errorf("api error: %s", msg)
	}
	return nil
}

// DeleteEmailForward removes an existing email forwarding rule from the domain.
func (c *Client) DeleteEmailForward(domain, address, forwardTo string) error {
	// Normalize address by trimming domain suffix if user supplied it
	cleanAddress := strings.TrimSuffix(address, "@"+domain)

	req := DeleteForwardRequest{
		BaseRequest: BaseRequest{
			APIKey:       c.APIKey,
			SecretAPIKey: c.SecretAPIKey,
		},
		Address:   cleanAddress,
		ForwardTo: forwardTo,
	}
	var res APIResponse
	endpoint := fmt.Sprintf("email/deleteForward/%s", domain)
	err := c.post(endpoint, req, &res)
	if err != nil {
		return err
	}
	if res.Status != "SUCCESS" {
		msg := res.Message
		if msg == "" {
			msg = "unknown error"
		}
		return fmt.Errorf("api error: %s", msg)
	}
	return nil
}
