package validator

import "testing"

func TestValidators(t *testing.T) {
	// Synthetic Aadhaar-shaped numbers built with the Verhoeff generator.
	base := "23456789012"
	cd := verhoeffCheckDigit(base)
	good := base + string(cd)
	bad := base + string('0'+(cd-'0'+1)%10)

	tests := []struct {
		typ, in string
		valid   bool
		code    string
	}{
		{"pan", "ABCPE1234F", true, ""},
		{"pan", "abcpe1234f", true, ""},
		{"pan", "ABCDE1234F", false, CodeInvalidFormat}, // 4th letter D is not a holder type
		{"pan", "", false, CodeEmpty},
		{"aadhaar", good, true, ""},
		{"aadhaar", bad, false, CodeInvalidChecksum},
		{"aadhaar", "123456789012", false, CodeInvalidFormat},
		{"mobile", "9876543210", true, ""},
		{"mobile", "+91 98765-43210", true, ""},
		{"mobile", "09876543210", true, ""},
		{"mobile", "5876543210", false, CodeInvalidFormat},
		{"email", "a.b@example.com", true, ""},
		{"email", "Bob <a@b.com>", false, CodeInvalidFormat},
		{"email", "nodomain@", false, CodeInvalidFormat},
		{"name", "Ravi Kumar", true, ""},
		{"name", "D'Souza-Rao", true, ""},
		{"name", "R2D2", false, CodeInvalidFormat},
		{"name", "A", false, CodeInvalidLength},
		{"gstin", "27AAPFU0939F1ZV", true, ""},
		{"gstin", "27AAPFU0939F1ZX", false, CodeInvalidChecksum},
		{"ifsc", "SBIN0001234", true, ""},
		{"ifsc", "SBIN1001234", false, CodeInvalidFormat},
	}
	reg := DefaultRegistry()
	for _, tc := range tests {
		t.Run(tc.typ+"/"+tc.in, func(t *testing.T) {
			r, err := reg.Validate(tc.typ, tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if r.Valid != tc.valid || r.Code != tc.code {
				t.Fatalf("got valid=%v code=%q, want valid=%v code=%q", r.Valid, r.Code, tc.valid, tc.code)
			}
		})
	}
}

func TestUnknownType(t *testing.T) {
	if _, err := DefaultRegistry().Validate("nope", "x"); err != ErrUnknownType {
		t.Fatalf("want ErrUnknownType, got %v", err)
	}
}
