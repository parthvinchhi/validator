package validators

import (
	"net/mail"
	"strings"
)

// ValidateEmail validates email syntax.
//
// Normalization: leading/trailing whitespace is trimmed. The case is NOT changed.
//
// Rules:
//   - total length at most 254, local part at most 64
//   - valid address syntax (standard library net/mail), a plain address only
//     (display names such as "Name <a@b.com>" are rejected)
//   - domain has at least two dot-separated labels, each 1-63 characters of
//     letters, digits or hyphens (not starting/ending with a hyphen)
//   - last label is at least 2 characters and not all digits
//
// This checks format only. It does not prove the mailbox exists
// (that needs verification, e.g. an email OTP).
// Internationalized (non-ASCII) domain names are rejected.
func ValidateEmail(value string) Result {
	const t = "email"

	email := strings.TrimSpace(value)

	if email == "" {
		return invalid(t, "Email is required")
	}
	if len(email) > 254 {
		return invalid(t, "Email is too long")
	}
	if strings.ContainsAny(email, " \t\r\n") {
		return invalid(t, "Email must not contain spaces")
	}

	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return invalid(t, "Invalid email format")
	}

	at := strings.LastIndex(email, "@")
	if at < 1 {
		return invalid(t, "Invalid email format")
	}
	local, domain := email[:at], email[at+1:]

	if len(local) > 64 {
		return invalid(t, "Email local part is too long")
	}
	if !isValidDomain(domain) {
		return invalid(t, "Invalid email domain")
	}

	return valid(t, "Email format is valid")
}

func isValidDomain(domain string) bool {
	if len(domain) > 253 {
		return false
	}

	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}

	for _, label := range labels {
		if label == "" || len(label) > 63 {
			return false
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for i := 0; i < len(label); i++ {
			if !isLetterDigitOrHyphen(label[i]) {
				return false
			}
		}
	}

	tld := labels[len(labels)-1]
	return len(tld) >= 2 && !isDigits(tld)
}

func isLetterDigitOrHyphen(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '-'
}
