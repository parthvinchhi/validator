package validators

import "testing"

func TestValidateMobile(t *testing.T) {
	runCases(t, ValidateMobile, "mobile", []testCase{
		{"valid", "9876543210", true},
		{"valid starts with 6", "6123456789", true},
		{"valid with spaces", " 98765 43210 ", true},
		{"valid with hyphen", "98765-43210", true},
		{"valid +91 prefix", "+919876543210", true},
		{"valid 91 prefix", "919876543210", true},
		{"valid 0 prefix", "09876543210", true},
		{"invalid starts with 5", "5876543210", false},
		{"invalid starts with 0", "0123456789", false},
		{"empty", "", false},
		{"only spaces", "   ", false},
		{"too short", "987654321", false},
		{"too long", "98765432101", false},
		{"letters", "98765abcde", false},
		{"all letters", "abcdefghij", false},
		{"special characters", "98765@3210", false},
		{"other country code", "+449876543210", false},
	})
}
