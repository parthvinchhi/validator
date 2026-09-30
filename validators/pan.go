package validators

import (
	"regexp"
	"strings"
)

var panPattern = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]$`)

// ValidatePAN validates the format of an Indian PAN.
//
// Normalization: whitespace is trimmed and letters are converted to UPPERCASE,
// so "abcde1234f" is accepted as "ABCDE1234F". No other characters are removed.
//
// Rules: 5 letters, 4 digits, 1 letter (AAAAA9999A).
//
// Note: the 4th letter normally shows the holder type (P, C, H, F, A, T, B, L,
// J, G). It is deliberately NOT enforced here, because the rule is about issuing
// practice, not format. PAN has no public checksum. Whether a PAN really exists
// requires verification with an authorized provider.
func ValidatePAN(value string) Result {
	const t = "pan"

	pan := strings.ToUpper(strings.TrimSpace(value))

	if pan == "" {
		return invalid(t, "PAN is required")
	}
	if len(pan) != 10 {
		return invalid(t, "Invalid PAN format: must be 10 characters")
	}
	if !panPattern.MatchString(pan) {
		return invalid(t, "Invalid PAN format")
	}

	return valid(t, "PAN format is valid")
}
