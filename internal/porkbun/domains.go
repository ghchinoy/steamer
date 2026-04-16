package porkbun

import "fmt"

// Domain represents a domain registered with Porkbun.
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

// Label represents a user-defined label in Porkbun.
type Label struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
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

// ListDomains fetches all domains in the user's Porkbun account.
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

// CheckDomain checks the availability and pricing of a domain.
func (c *Client) CheckDomain(domain string) (*DomainCheckResponse, error) {
	req := BaseRequest{
		APIKey:       c.APIKey,
		SecretAPIKey: c.SecretAPIKey,
	}
	var res DomainCheckResponse
	endpoint := fmt.Sprintf("domain/checkDomain/%s", domain)
	err := c.post(endpoint, req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
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

// GetPricing retrieves pricing information for all supported TLDs.
func (c *Client) GetPricing() (*PricingResponse, error) {
	req := BaseRequest{
		APIKey:       c.APIKey,
		SecretAPIKey: c.SecretAPIKey,
	}
	var res PricingResponse
	err := c.post("pricing/get", req, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
