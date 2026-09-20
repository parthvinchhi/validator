package handler_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/parthvinchhi/api-validator/internal/config"
	"github.com/parthvinchhi/api-validator/internal/server"
	"github.com/parthvinchhi/api-validator/internal/validator"
)

func newServer() http.Handler {
	cfg := config.Config{APIKeys: []string{"secret"}, RateRPS: 1000, RateBurst: 1000, MaxBatch: 2}
	return server.NewRouter(cfg, validator.DefaultRegistry(), slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func do(h http.Handler, method, path, key, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if key != "" {
		req.Header.Set("X-API-Key", key)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestEndpoints(t *testing.T) {
	h := newServer()
	tests := []struct {
		name, method, path, key, body string
		want                          int
	}{
		{"health needs no key", "GET", "/healthz", "", "", 200},
		{"no key", "POST", "/v1/validate/pan", "", `{"value":"ABCPE1234F"}`, 401},
		{"valid pan", "POST", "/v1/validate/pan", "secret", `{"value":"ABCPE1234F"}`, 200},
		{"invalid pan is still 200", "POST", "/v1/validate/pan", "secret", `{"value":"bad"}`, 200},
		{"unknown type", "POST", "/v1/validate/xyz", "secret", `{"value":"a"}`, 404},
		{"malformed json", "POST", "/v1/validate/pan", "secret", `{oops`, 400},
		{"batch ok", "POST", "/v1/validate/batch", "secret", `{"items":[{"type":"pan","value":"x"},{"type":"nope","value":"y"}]}`, 200},
		{"batch too big", "POST", "/v1/validate/batch", "secret", `{"items":[{"type":"pan","value":"x"},{"type":"pan","value":"x"},{"type":"pan","value":"x"}]}`, 400},
		{"list types", "GET", "/v1/validators", "secret", "", 200},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := do(h, tc.method, tc.path, tc.key, tc.body).Code; got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}
