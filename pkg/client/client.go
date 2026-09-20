// Package client is a small Go SDK for calling the validator API from other projects.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	Type       string `json:"type"`
	Valid      bool   `json:"valid"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Normalized string `json:"normalized,omitempty"`
}

type Item struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{BaseURL: baseURL, APIKey: apiKey, HTTP: &http.Client{Timeout: 5 * time.Second}}
}

func (c *Client) post(ctx context.Context, path string, in, out any) error {
	b, _ := json.Marshal(in)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.APIKey)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("validator api: unexpected status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// Validate validates a single value, e.g. Validate(ctx, "pan", "ABCPE1234F").
func (c *Client) Validate(ctx context.Context, typ, value string) (Result, error) {
	var res Result
	err := c.post(ctx, "/v1/validate/"+typ, map[string]string{"value": value}, &res)
	return res, err
}

// Batch validates many values in one request.
func (c *Client) Batch(ctx context.Context, items []Item) ([]Result, error) {
	var out struct {
		Results []Result `json:"results"`
	}
	err := c.post(ctx, "/v1/validate/batch", map[string]any{"items": items}, &out)
	return out.Results, err
}
