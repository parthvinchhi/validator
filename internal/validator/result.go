package validator

// Result is the outcome of a single validation.
type Result struct {
	Type       string `json:"type"`
	Valid      bool   `json:"valid"`
	Code       string `json:"code,omitempty"`       // e.g. INVALID_FORMAT, INVALID_CHECKSUM
	Message    string `json:"message,omitempty"`    // human-readable reason
	Normalized string `json:"normalized,omitempty"` // cleaned value, when useful
}

// Error codes.
const (
	CodeEmpty           = "EMPTY"
	CodeInvalidFormat   = "INVALID_FORMAT"
	CodeInvalidChecksum = "INVALID_CHECKSUM"
	CodeInvalidLength   = "INVALID_LENGTH"
	CodeUnknownType     = "UNKNOWN_TYPE"
)

func ok(typ, normalized string) Result {
	return Result{Type: typ, Valid: true, Normalized: normalized}
}

func fail(typ, code, msg string) Result {
	return Result{Type: typ, Valid: false, Code: code, Message: msg}
}
