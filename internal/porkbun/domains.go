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

// Domain represents a domain registered with Porkbun.
type Domain struct {
	Domain       string         `json:"domain"`
	Status       string         `json:"status"`
	TLD          string         `json:"tld"`
	CreateDate   string         `json:"createDate"`
	ExpireDate   string         `json:"expireDate"`
	SecurityLock FlexibleString `json:"securityLock"`
	WhoisPrivacy FlexibleString `json:"whoisPrivacy"`
	AutoRenew    FlexibleString `json:"autoRenew"`
	NotLocal     FlexibleString `json:"notLocal"`
	Labels       []Label        `json:"labels"`
}

// Label represents a user-defined label in Porkbun.
type Label struct {
	ID    FlexibleString `json:"id"`
	Title string         `json:"title"`
	Color string         `json:"color"`
}

// ListDomainsRequest is the request body for listing all domains.
type ListDomainsRequest struct {
	BaseRequest
	Start         string `json:"start,omitempty"`
	IncludeLabels string `json:"includeLabels,omitempty"`
}

// ListDomainsResponse is the response from the listAll domains endpoint.
type ListDomainsResponse struct {
	APIResponse
	Domains []Domain `json:"domains"`
}

// DomainPricing is the detailed pricing and availability information.
type DomainPricing struct {
	Avail   string `json:"avail"`
	Premium string `json:"premium"`
	Price   string `json:"price"`
}

// DomainCheckResponse is the response from the checkDomain endpoint.
type DomainCheckResponse struct {
	APIResponse
	Response DomainPricing `json:"response"`
}

// TLDPricing represents the pricing info for a single TLD.
type TLDPricing struct {
	Registration string `json:"registration"`
	Renewal      string `json:"renewal"`
	Transfer     string `json:"transfer"`
}

// PricingResponse is the response from the pricing/get endpoint.
type PricingResponse struct {
	APIResponse
	Pricing map[string]TLDPricing `json:"pricing"`
}

// ListDomains fetches all domains in the user's Porkbun account.
func (c *Client) ListDomains(ctx context.Context) ([]Domain, error) {
	req := ListDomainsRequest{
		BaseRequest:   c.baseRequest(),
		IncludeLabels: "yes",
	}
	var res ListDomainsResponse
	if err := c.post(ctx, "domain/listAll", req, &res); err != nil {
		return nil, err
	}
	if err := checkAPIStatus("api error", res.APIResponse); err != nil {
		return nil, err
	}
	return res.Domains, nil
}

// CheckDomain checks the availability and pricing of a domain.
func (c *Client) CheckDomain(ctx context.Context, domain string) (*DomainCheckResponse, error) {
	if err := validatePathSegment("domain", domain); err != nil {
		return nil, err
	}
	var res DomainCheckResponse
	endpoint := fmt.Sprintf("domain/checkDomain/%s", url.PathEscape(domain))
	if err := c.post(ctx, endpoint, c.baseRequest(), &res); err != nil {
		return nil, err
	}
	if err := checkAPIStatus("api error", res.APIResponse); err != nil {
		return nil, err
	}
	return &res, nil
}

// FetchPricing retrieves pricing information for all supported TLDs.
func (c *Client) FetchPricing(ctx context.Context) (*PricingResponse, error) {
	var res PricingResponse
	if err := c.post(ctx, "pricing/get", c.baseRequest(), &res); err != nil {
		return nil, err
	}
	if err := checkAPIStatus("api error", res.APIResponse); err != nil {
		return nil, err
	}
	return &res, nil
}
