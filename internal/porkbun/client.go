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

// Package porkbun provides a client for the Porkbun API v3.
package porkbun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL     = "https://api.porkbun.com/api/json/v3"
	defaultHTTPTimeout = 30 * time.Second
	maxResponseBytes   = 10 << 20 // 10 MiB
	maxErrorBodyBytes  = 256
)

// Client is a Porkbun API client.
type Client struct {
	APIKey       string
	SecretAPIKey string
	BaseURL      string
	HTTPClient   *http.Client
}

// BaseRequest contains the credentials required for every Porkbun API request.
type BaseRequest struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
}

// FlexibleString is a string that can unmarshal from either a JSON string or a JSON number
// without float64 conversion or scientific notation formatting.
type FlexibleString string

// UnmarshalJSON implements json.Unmarshaler for FlexibleString.
func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*fs = ""
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*fs = FlexibleString(s)
		return nil
	}
	*fs = FlexibleString(string(data))
	return nil
}

// APIResponse is the common response structure from the Porkbun API.
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// PingResponse is the response from the ping endpoint.
type PingResponse struct {
	APIResponse
	YourIP string `json:"yourIp"`
}

// NewClient creates a new Porkbun API client with sensible defaults.
func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:       apiKey,
		SecretAPIKey: secretKey,
		BaseURL:      defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}
}

// Ping checks the connection to the Porkbun API.
func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	req := c.baseRequest()
	var res PingResponse
	if err := c.post(ctx, "ping", req, &res); err != nil {
		return nil, err
	}
	if err := checkAPIStatus("ping failed", res.APIResponse); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) baseRequest() BaseRequest {
	return BaseRequest{
		APIKey:       c.APIKey,
		SecretAPIKey: c.SecretAPIKey,
	}
}

func (c *Client) baseURL() string {
	if trimmed := strings.TrimRight(c.BaseURL, "/"); trimmed != "" {
		return trimmed
	}
	return defaultBaseURL
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: defaultHTTPTimeout}
}

func (c *Client) post(ctx context.Context, endpoint string, body, result any) error {
	reqURL := fmt.Sprintf("%s/%s", c.baseURL(), strings.TrimLeft(endpoint, "/"))

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("creating http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("sending http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	if len(respBody) > maxResponseBytes {
		return fmt.Errorf("response body exceeded %d bytes", maxResponseBytes)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIResponse
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return fmt.Errorf("api error: %s", apiErr.Message)
		}
		snippet := strings.TrimSpace(string(respBody))
		if len(snippet) > maxErrorBodyBytes {
			snippet = snippet[:maxErrorBodyBytes] + "..."
		}
		return fmt.Errorf("api error: %s (status %d)", snippet, resp.StatusCode)
	}

	dec := json.NewDecoder(bytes.NewReader(respBody))
	dec.UseNumber()
	if err := dec.Decode(result); err != nil {
		return fmt.Errorf("decoding response body: %w", err)
	}
	return nil
}

func checkAPIStatus(prefix string, res APIResponse) error {
	if res.Status == "SUCCESS" {
		return nil
	}
	msg := res.Message
	if msg == "" {
		msg = "unknown error"
	}
	return fmt.Errorf("%s: %s", prefix, msg)
}

func validatePathSegment(label, value string) error {
	if value == "" {
		return fmt.Errorf("invalid %s: must not be empty", label)
	}
	if strings.Contains(value, "..") {
		return fmt.Errorf("invalid %s %q: must not contain path traversal", label, value)
	}
	for _, r := range value {
		if r <= 0x1F || r == 0x7F || r == '/' || r == '\\' || r == '?' || r == '#' || r == ' ' || r == '\t' {
			return fmt.Errorf("invalid %s %q: contains disallowed character", label, value)
		}
	}
	return nil
}
