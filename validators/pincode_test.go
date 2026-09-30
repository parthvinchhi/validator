package validators

import "testing"

func TestValidatePincode(t *testing.T) {
	runCases(t, ValidatePincode, "pincode", []testCase{
		{"valid", "380001", true},
		{"valid with spaces around", " 380001 ", true},
		{"starts with 0", "000001", false},
		{"too short", "38000", false},
		{"too long", "3800011", false},
		{"letters", "38000A", false},
		{"space inside", "380 001", false},
		{"hyphen inside", "380-001", false},
		{"empty", "", false},
	})
}
