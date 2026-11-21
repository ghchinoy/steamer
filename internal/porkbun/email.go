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

// EmailForwardsResponse is the response from the email/getForwards endpoint.
type EmailForwardsResponse struct {
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

// FetchEmailForwards retrieves all email forwarding rules for a domain.
func (c *Client) FetchEmailForwards(ctx context.Context, domain string) (*EmailForwardsResponse, error) {
	if err := validatePathSegment("domain", domain); err != nil {
		return nil, err
	}
	var res EmailForwardsResponse
	endpoint := fmt.Sprintf("email/getForwards/%s", url.PathEscape(domain))
	if err := c.post(ctx, endpoint, c.baseRequest(), &res); err != nil {
		return nil, err
	}
	if err := checkAPIStatus("api error", res.APIResponse); err != nil {
		return nil, err
	}
	return &res, nil
}

// AddEmailForward creates a new email forwarding rule for the domain.
func (c *Client) AddEmailForward(ctx context.Context, domain, address, forwardTo string) error {
	if err := validatePathSegment("domain", domain); err != nil {
		return err
	}
	cleanAddress := strings.TrimSuffix(address, "@"+domain)
	req := AddForwardRequest{
		BaseRequest: c.baseRequest(),
		Address:     cleanAddress,
		ForwardTo:   forwardTo,
	}
	var res APIResponse
	endpoint := fmt.Sprintf("email/addForward/%s", url.PathEscape(domain))
	if err := c.post(ctx, endpoint, req, &res); err != nil {
		return err
	}
	return checkAPIStatus("api error", res)
}

// DeleteEmailForward removes an existing email forwarding rule from the domain.
func (c *Client) DeleteEmailForward(ctx context.Context, domain, address, forwardTo string) error {
	if err := validatePathSegment("domain", domain); err != nil {
		return err
	}
	cleanAddress := strings.TrimSuffix(address, "@"+domain)
	req := DeleteForwardRequest{
		BaseRequest: c.baseRequest(),
		Address:     cleanAddress,
		ForwardTo:   forwardTo,
	}
	var res APIResponse
	endpoint := fmt.Sprintf("email/deleteForward/%s", url.PathEscape(domain))
	if err := c.post(ctx, endpoint, req, &res); err != nil {
		return err
	}
	return checkAPIStatus("api error", res)
}
