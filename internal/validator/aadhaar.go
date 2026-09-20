package validator

import (
	"regexp"
	"strings"
)

var aadhaarRe = regexp.MustCompile(`^[2-9][0-9]{11}$`)

// Aadhaar checks 12 digits, first digit 2-9, and the Verhoeff checksum.
// This proves the number is well-formed, NOT that it exists or is issued.
type Aadhaar struct{}

func (Aadhaar) Type() string { return "aadhaar" }

func (Aadhaar) Validate(in string) Result {
	s := strings.NewReplacer(" ", "", "-", "").Replace(in)
	if s == "" {
		return fail("aadhaar", CodeEmpty, "Aadhaar number is required")
	}
	if !aadhaarRe.MatchString(s) {
		return fail("aadhaar", CodeInvalidFormat, "Aadhaar must be 12 digits and not start with 0 or 1")
	}
	if !verhoeffValid(s) {
		return fail("aadhaar", CodeInvalidChecksum, "Aadhaar checksum is invalid")
	}
	return ok("aadhaar", s)
}
