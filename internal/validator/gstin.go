package validator

import (
	"regexp"
	"strings"
)

var gstinRe = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z][1-9A-Z]Z[0-9A-Z]$`)

const gstinChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type GSTIN struct{}

func (GSTIN) Type() string { return "gstin" }

func (GSTIN) Validate(in string) Result {
	s := strings.ToUpper(strings.TrimSpace(in))
	if s == "" {
		return fail("gstin", CodeEmpty, "GSTIN is required")
	}
	if !gstinRe.MatchString(s) {
		return fail("gstin", CodeInvalidFormat, "GSTIN must be 15 characters like 27AAPFU0939F1ZV")
	}
	if gstinCheckChar(s[:14]) != s[14] {
		return fail("gstin", CodeInvalidChecksum, "GSTIN checksum is invalid")
	}
	return ok("gstin", s)
}

func gstinCheckChar(first14 string) byte {
	sum := 0
	for i := 0; i < len(first14); i++ {
		v := strings.IndexByte(gstinChars, first14[i])
		f := 1
		if i%2 == 1 {
			f = 2
		}
		p := v * f
		sum += p/36 + p%36
	}
	return gstinChars[(36-sum%36)%36]
}
