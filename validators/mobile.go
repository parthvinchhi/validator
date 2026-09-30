package validators

import "strings"

// ValidateMobile validates an Indian mobile number.
//
// Normalization (applied before checking):
//   - leading/trailing whitespace is trimmed
//   - spaces and hyphens are removed
//   - a "+91" prefix, a "91" prefix (12 digits total) or a "0" prefix
//     (11 digits total) is removed
//
// Rules: exactly 10 digits, first digit 6, 7, 8 or 9.
func ValidateMobile(value string) Result {
	const t = "mobile"

	number := strings.TrimSpace(value)
	number = strings.NewReplacer(" ", "", "-", "").Replace(number)

	if number == "" {
		return invalid(t, "Mobile number is required")
	}

	switch {
	case strings.HasPrefix(number, "+91"):
		number = number[3:]
	case len(number) == 12 && strings.HasPrefix(number, "91"):
		number = number[2:]
	case len(number) == 11 && strings.HasPrefix(number, "0"):
		number = number[1:]
	}

	if !isDigits(number) {
		return invalid(t, "Mobile number must contain digits only")
	}
	if len(number) != 10 {
		return invalid(t, "Mobile number must be 10 digits")
	}
	if number[0] < '6' {
		return invalid(t, "Mobile number must start with 6, 7, 8 or 9")
	}

	return valid(t, "Mobile number format is valid")
}
