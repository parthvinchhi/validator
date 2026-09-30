package validators

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	runCases(t, ValidateEmail, "email", []testCase{
		{"valid", "user@example.com", true},
		{"valid uppercase and spaces around", "  USER@EXAMPLE.COM ", true},
		{"valid plus and dots", "user.name+tag@example.co.in", true},
		{"valid subdomain", "user@mail.example.com", true},
		{"empty", "", false},
		{"only spaces", "   ", false},
		{"missing at", "userexample.com", false},
		{"missing domain", "user@", false},
		{"missing local part", "@example.com", false},
		{"domain without dot", "user@invalid", false},
		{"empty domain label", "user@example..com", false},
		{"domain label starts with hyphen", "user@-example.com", false},
		{"one letter tld", "user@example.c", false},
		{"numeric tld", "user@example.123", false},
		{"space inside", "us er@example.com", false},
		{"display name", "Name <user@example.com>", false},
		{"too long", strings.Repeat("a", 250) + "@example.com", false},
	})
}
