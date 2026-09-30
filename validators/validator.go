// Package validators contains one file per validator plus this shared file.
//
// Every validator has the same shape:
//
//	func ValidateXxx(value string) Result
//
// It receives the raw text, applies its own (documented) normalization,
// and returns a Result. Validators never log, store, or return the input.
package validators

import "sort"

// Result is what every validator returns.
type Result struct {
	Valid   bool   `json:"valid"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ValidatorFunc is the function shape shared by all string validators.
type ValidatorFunc func(value string) Result

// registry maps the "type" sent by clients to the validator function.
// TO ADD A NEW VALIDATOR: create validators/yourtype.go and add ONE line here.
// The generic endpoint and the dedicated endpoint
// (POST /api/v1/validate/<type>) are created automatically from this map.
var registry = map[string]ValidatorFunc{
	"mobile":  ValidateMobile,
	"email":   ValidateEmail,
	"pan":     ValidatePAN,
	"aadhaar": ValidateAadhaar,
	"pincode": ValidatePincode,
}

// Get returns the validator registered for a type.
func Get(validationType string) (ValidatorFunc, bool) {
	fn, ok := registry[validationType]
	return fn, ok
}

// Types returns all registered type names in alphabetical order.
func Types() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func valid(validationType, message string) Result {
	return Result{Valid: true, Type: validationType, Message: message}
}

func invalid(validationType, message string) Result {
	return Result{Valid: false, Type: validationType, Message: message}
}

// isDigits reports whether s is non-empty and contains only 0-9.
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
