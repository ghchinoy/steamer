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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.porkbun.com/api/json/v3"

type Client struct {
	APIKey       string
	SecretAPIKey string
	HTTPClient   *http.Client
}

type BaseRequest struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:       apiKey,
		SecretAPIKey: secretKey,
		HTTPClient:   &http.Client{},
	}
}

func (c *Client) post(endpoint string, body interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)

	// Inject credentials into the body if it's a map or a struct that embeds BaseRequest
	// For simplicity in this implementation, we'll assume the caller passes a struct
	// or we can use a more generic approach.
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIResponse
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return fmt.Errorf("api error: %s", apiErr.Message)
		}
		return fmt.Errorf("api error: %s (status %d)", string(respBody), resp.StatusCode)
	}

	return json.Unmarshal(respBody, result)
}

type PingResponse struct {
	APIResponse
	YourIP string `json:"yourIp"`
}

func (c *Client) Ping() (*PingResponse, error) {
	req := BaseRequest{
		APIKey:       c.APIKey,
		SecretAPIKey: c.SecretAPIKey,
	}
	var res PingResponse
	err := c.post("ping", req, &res)
	if err != nil {
		return nil, err
	}
	if res.Status != "SUCCESS" {
		return nil, fmt.Errorf("ping failed: %s", res.Message)
	}
	return &res, nil
}
