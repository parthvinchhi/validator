package validator

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Name validates a person's name: letters (any script), spaces, . ' -
type Name struct{}

func (Name) Type() string { return "name" }

func (Name) Validate(in string) Result {
	in = strings.Join(strings.Fields(in), " ") // trim + collapse spaces
	if in == "" {
		return fail("name", CodeEmpty, "name is required")
	}
	if n := utf8.RuneCountInString(in); n < 2 || n > 100 {
		return fail("name", CodeInvalidLength, "name must be 2 to 100 characters")
	}
	letters := 0
	for _, r := range in {
		switch {
		case unicode.IsLetter(r):
			letters++
		case unicode.IsMark(r), r == ' ', r == '.', r == '\'', r == '-':
		default:
			return fail("name", CodeInvalidFormat, "name contains invalid characters")
		}
	}
	if letters < 2 {
		return fail("name", CodeInvalidFormat, "name must contain at least 2 letters")
	}
	return ok("name", in)
}
