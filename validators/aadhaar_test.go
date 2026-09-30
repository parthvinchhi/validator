package validators

import "testing"

func TestValidateAadhaar(t *testing.T) {
	// 234567890124 and 499900001236 have correct Verhoeff check digits.
	runCases(t, ValidateAadhaar, "aadhaar", []testCase{
		{"valid", "234567890124", true},
		{"valid second number", "499900001236", true},
		{"valid with spaces", "2345 6789 0124", true},
		{"valid with hyphens", "2345-6789-0124", true},
		{"valid with spaces around", " 234567890124 ", true},
		{"wrong checksum", "234567890123", false},
		{"starts with 1", "123456789012", false},
		{"starts with 0", "012345678901", false},
		{"too short", "23456789012", false},
		{"too long", "2345678901245", false},
		{"letters", "23456789012A", false},
		{"special characters", "2345678901@4", false},
		{"empty", "", false},
	})
}
