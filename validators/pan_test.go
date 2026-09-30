package validators

import "testing"

func TestValidatePAN(t *testing.T) {
	runCases(t, ValidatePAN, "pan", []testCase{
		{"valid", "ABCDE1234F", true},
		{"lowercase is accepted (normalized)", "abcde1234f", true},
		{"mixed case is accepted (normalized)", "AbCdE1234f", true},
		{"spaces around are trimmed", " ABCDE1234F ", true},
		{"empty", "", false},
		{"too short", "ABCDE1234", false},
		{"too long", "ABCDE12345F", false},
		{"digits first", "12345ABCDE", false},
		{"wrong pattern", "ABCD12345F", false},
		{"last char is digit", "ABCDE12345", false},
		{"special character", "ABCDE1234!", false},
		{"hyphen inside", "ABCDE-1234F", false},
		{"space inside", "ABCDE 1234F", false},
	})
}
