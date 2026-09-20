package validator

import (
	"regexp"
	"strings"
)

var mobileRe = regexp.MustCompile(`^[6-9][0-9]{9}$`)

// Mobile validates Indian mobile numbers. Accepts +91 / 91 / 0 prefixes.
type Mobile struct{}

func (Mobile) Type() string { return "mobile" }

func (Mobile) Validate(in string) Result {
	s := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(in)
	if s == "" {
		return fail("mobile", CodeEmpty, "mobile number is required")
	}
	switch {
	case strings.HasPrefix(s, "+91"):
		s = s[3:]
	case len(s) == 12 && strings.HasPrefix(s, "91"):
		s = s[2:]
	case len(s) == 11 && strings.HasPrefix(s, "0"):
		s = s[1:]
	}
	if !mobileRe.MatchString(s) {
		return fail("mobile", CodeInvalidFormat, "mobile must be 10 digits starting with 6-9")
	}
	return ok("mobile", s)
}
