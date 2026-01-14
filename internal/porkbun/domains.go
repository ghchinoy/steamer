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

type Domain struct {
	Domain       string      `json:"domain"`
	Status       string      `json:"status"`
	TLD          string      `json:"tld"`
	CreateDate   string      `json:"createDate"`
	ExpireDate   string      `json:"expireDate"`
	SecurityLock string      `json:"securityLock"`
	WhoisPrivacy string      `json:"whoisPrivacy"`
	AutoRenew    interface{} `json:"autoRenew"`
	NotLocal     interface{} `json:"notLocal"`
	Labels       []Label     `json:"labels"`
}

type Label struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type ListDomainsRequest struct {
	BaseRequest
	Start         string `json:"start,omitempty"`
	IncludeLabels string `json:"includeLabels,omitempty"`
}

type ListDomainsResponse struct {
	APIResponse
	Domains []Domain `json:"domains"`
}

func (c *Client) ListDomains() ([]Domain, error) {
	req := ListDomainsRequest{
		BaseRequest: BaseRequest{
			APIKey:       c.APIKey,
			SecretAPIKey: c.SecretAPIKey,
		},
		IncludeLabels: "yes",
	}
	var res ListDomainsResponse
	err := c.post("domain/listAll", req, &res)
	if err != nil {
		return nil, err
	}
	return res.Domains, nil
}
