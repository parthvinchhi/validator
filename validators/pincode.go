package validators

import "strings"

// ValidatePincode validates an Indian PIN code.
//
// Normalization: leading/trailing whitespace is trimmed. Nothing else is
// changed, so "380 001" is rejected.
//
// Rules: exactly 6 digits, first digit 1-9.
//
// This checks format only. It does not check that the PIN is actually in use.
func ValidatePincode(value string) Result {
	const t = "pincode"

	pin := strings.TrimSpace(value)

	if pin == "" {
		return invalid(t, "PIN code is required")
	}
	if !isDigits(pin) {
		return invalid(t, "PIN code must contain digits only")
	}
	if len(pin) != 6 {
		return invalid(t, "PIN code must be 6 digits")
	}
	if pin[0] == '0' {
		return invalid(t, "PIN code cannot start with 0")
	}

	return valid(t, "PIN code format is valid")
}
