package validator

import (
	"regexp"
	"strings"
)

var ifscRe = regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)

type IFSC struct{}

func (IFSC) Type() string { return "ifsc" }

func (IFSC) Validate(in string) Result {
	s := strings.ToUpper(strings.TrimSpace(in))
	if s == "" {
		return fail("ifsc", CodeEmpty, "IFSC is required")
	}
	if !ifscRe.MatchString(s) {
		return fail("ifsc", CodeInvalidFormat, "IFSC must look like SBIN0001234")
	}
	return ok("ifsc", s)
}
