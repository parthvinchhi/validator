// Package model holds request/response shapes for the HTTP API.
package model

import "github.com/parthvinchhi/api-validator/internal/validator"

type ValidateRequest struct {
	Value string `json:"value"`
}

type BatchItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type BatchRequest struct {
	Items []BatchItem `json:"items"`
}

type BatchResponse struct {
	Results []validator.Result `json:"results"`
}

type TypesResponse struct {
	Types []string `json:"types"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
