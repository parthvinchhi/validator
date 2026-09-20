package validator

import (
	"net/mail"
	"strings"
)

// Email validates syntax only (no MX lookup).
type Email struct{}

func (Email) Type() string { return "email" }

func (Email) Validate(in string) Result {
	s := strings.TrimSpace(in)
	if s == "" {
		return fail("email", CodeEmpty, "email is required")
	}
	if len(s) > 254 {
		return fail("email", CodeInvalidLength, "email is too long")
	}
	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Address != s { // rejects "Name <a@b.c>" forms
		return fail("email", CodeInvalidFormat, "invalid email address")
	}
	at := strings.LastIndex(s, "@")
	local, domain := s[:at], s[at+1:]
	if len(local) > 64 || !strings.Contains(domain, ".") ||
		strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") ||
		strings.Contains(domain, "..") {
		return fail("email", CodeInvalidFormat, "invalid email address")
	}
	return ok("email", strings.ToLower(s))
}
