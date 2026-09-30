// Package services sits between the HTTP handlers and the validators.
// It picks the right validator for a type and applies API-level rules
// (such as the maximum input length). It knows nothing about HTTP.
package services

import (
	"errors"
	"strings"

	"github.com/parthvinchhi/api-validator/validators"
)

// MaxValueLength is the longest "value" the API will process (in bytes).
const MaxValueLength = 256

var (
	ErrUnsupportedType = errors.New("unsupported validation type")
	ErrValueTooLong    = errors.New("value too long")
)

// NormalizeType makes type names case-insensitive: " PAN " becomes "pan".
func NormalizeType(validationType string) string {
	return strings.ToLower(strings.TrimSpace(validationType))
}

// Validate runs the validator registered for validationType.
func Validate(validationType, value string) (validators.Result, error) {
	fn, ok := validators.Get(NormalizeType(validationType))
	if !ok {
		return validators.Result{}, ErrUnsupportedType
	}
	if len(value) > MaxValueLength {
		return validators.Result{}, ErrValueTooLong
	}
	return fn(value), nil
}
