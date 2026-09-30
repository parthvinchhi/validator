// Package dto holds the JSON request and response shapes of the API.
package dto

// ValidateRequest is the body for every validation endpoint.
// "type" is only used by the generic POST /api/v1/validate endpoint.
type ValidateRequest struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ValidateResponse is returned by every endpoint so clients can rely on one shape.
// "type" is left out when it does not apply (for example a malformed request).
type ValidateResponse struct {
	Success bool   `json:"success"`
	Valid   bool   `json:"valid"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message"`
}
