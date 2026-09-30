package validators

import "strings"

// Verhoeff checksum tables. Aadhaar's last digit is a Verhoeff check digit.
var verhoeffD = [10][10]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 2, 3, 4, 0, 6, 7, 8, 9, 5},
	{2, 3, 4, 0, 1, 7, 8, 9, 5, 6},
	{3, 4, 0, 1, 2, 8, 9, 5, 6, 7},
	{4, 0, 1, 2, 3, 9, 5, 6, 7, 8},
	{5, 9, 8, 7, 6, 0, 4, 3, 2, 1},
	{6, 5, 9, 8, 7, 1, 0, 4, 3, 2},
	{7, 6, 5, 9, 8, 2, 1, 0, 4, 3},
	{8, 7, 6, 5, 9, 3, 2, 1, 0, 4},
	{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
}

var verhoeffP = [8][10]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 5, 7, 6, 2, 8, 3, 0, 9, 4},
	{5, 8, 0, 3, 7, 9, 6, 1, 4, 2},
	{8, 9, 1, 6, 0, 4, 3, 5, 2, 7},
	{9, 4, 5, 3, 1, 2, 6, 8, 7, 0},
	{4, 2, 8, 6, 5, 7, 3, 9, 0, 1},
	{2, 7, 9, 3, 8, 0, 6, 4, 1, 5},
	{7, 0, 4, 6, 9, 1, 3, 2, 5, 8},
}

// ValidateAadhaar validates the format and checksum of an Aadhaar number.
//
// Normalization: whitespace is trimmed and spaces/hyphens are removed, so
// "2345 6789 0124" is accepted.
//
// Rules: exactly 12 digits, first digit 2-9, and the Verhoeff checksum passes.
//
// This does NOT prove the Aadhaar was issued to anyone. That requires an
// authorized verification service.
func ValidateAadhaar(value string) Result {
	const t = "aadhaar"

	number := strings.TrimSpace(value)
	number = strings.NewReplacer(" ", "", "-", "").Replace(number)

	if number == "" {
		return invalid(t, "Aadhaar number is required")
	}
	if !isDigits(number) {
		return invalid(t, "Aadhaar number must contain digits only")
	}
	if len(number) != 12 {
		return invalid(t, "Aadhaar number must be 12 digits")
	}
	if number[0] == '0' || number[0] == '1' {
		return invalid(t, "Aadhaar number cannot start with 0 or 1")
	}
	if !verhoeffValid(number) {
		return invalid(t, "Aadhaar checksum is invalid")
	}

	return valid(t, "Aadhaar format is valid")
}

// verhoeffValid expects digits only.
func verhoeffValid(number string) bool {
	c := 0
	for i := 0; i < len(number); i++ {
		digit := int(number[len(number)-1-i] - '0')
		c = verhoeffD[c][verhoeffP[i%8][digit]]
	}
	return c == 0
}
