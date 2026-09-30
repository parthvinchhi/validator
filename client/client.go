// Package client is a tiny helper for calling the Validation API from other
// Go projects. It only uses the standard library, so you can simply copy this
// one file into your own project.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Response mirrors the JSON returned by the API.
type Response struct {
	Success bool   `json:"success"`
	Valid   bool   `json:"valid"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// New creates a client. baseURL looks like "http://localhost:8080".
// apiKey may be empty if the service does not use API_KEY.
func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Validate calls POST /api/v1/validate.
//
// A non-nil error means the call itself failed (network problem, bad API key,
// unsupported type...). If err is nil, check resp.Valid to see whether the
// submitted value is valid.
func (c *Client) Validate(validationType, value string) (*Response, error) {
	payload, err := json.Marshal(map[string]string{"type": validationType, "value": value})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/validate", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling validation service: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64*1024)).Decode(&result); err != nil {
		return nil, fmt.Errorf("unexpected response (HTTP %d): %w", resp.StatusCode, err)
	}
	if !result.Success {
		return &result, fmt.Errorf("validation service error (HTTP %d): %s", resp.StatusCode, result.Message)
	}

	return &result, nil
}

func (c *Client) ValidateMobile(value string) (*Response, error)  { return c.Validate("mobile", value) }
func (c *Client) ValidateEmail(value string) (*Response, error)   { return c.Validate("email", value) }
func (c *Client) ValidatePAN(value string) (*Response, error)     { return c.Validate("pan", value) }
func (c *Client) ValidateAadhaar(value string) (*Response, error) { return c.Validate("aadhaar", value) }
func (c *Client) ValidatePincode(value string) (*Response, error) { return c.Validate("pincode", value) }
