package validator

import (
	"regexp"
	"strings"
)

// Format: 5 letters, 4 digits, 1 letter. 4th letter is the holder type.
var panRe = regexp.MustCompile(`^[A-Z]{3}[ABCFGHLJPT][A-Z][0-9]{4}[A-Z]$`)

type PAN struct{}

func (PAN) Type() string { return "pan" }

func (PAN) Validate(in string) Result {
	s := strings.ToUpper(strings.TrimSpace(in))
	if s == "" {
		return fail("pan", CodeEmpty, "PAN is required")
	}
	if !panRe.MatchString(s) {
		return fail("pan", CodeInvalidFormat, "PAN must look like ABCPE1234F")
	}
	return ok("pan", s)
}
