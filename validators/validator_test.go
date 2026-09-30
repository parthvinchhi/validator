package validators

import "testing"

type testCase struct {
	name  string
	input string
	want  bool
}

// runCases runs a table of inputs through one validator.
func runCases(t *testing.T, fn ValidatorFunc, validationType string, cases []testCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := fn(tc.input)
			if got.Valid != tc.want {
				t.Errorf("input %q: valid = %v, want %v (message: %s)", tc.input, got.Valid, tc.want, got.Message)
			}
			if got.Type != validationType {
				t.Errorf("type = %q, want %q", got.Type, validationType)
			}
			if got.Message == "" {
				t.Error("message must not be empty")
			}
		})
	}
}

func TestRegistry(t *testing.T) {
	for _, name := range []string{"mobile", "email", "pan", "aadhaar", "pincode"} {
		if _, ok := Get(name); !ok {
			t.Errorf("validator %q is not registered", name)
		}
	}
	if _, ok := Get("does-not-exist"); ok {
		t.Error("unknown type should not be found")
	}
	if len(Types()) != 5 {
		t.Errorf("expected 5 types, got %d", len(Types()))
	}
}
