package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/api-validator/config"
	"github.com/parthvinchhi/api-validator/dto"
	"github.com/parthvinchhi/api-validator/routes"
)

func testConfig() config.Config {
	return config.Config{Port: "8080", Environment: "test", MaxBodyBytes: 4096}
}

func newRouter(cfg config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return routes.Setup(cfg)
}

func send(router http.Handler, method, path, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) dto.ValidateResponse {
	t.Helper()
	var resp dto.ValidateResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v (body: %s)", err, rec.Body.String())
	}
	return resp
}

func TestHealth(t *testing.T) {
	rec := send(newRouter(testConfig()), http.MethodGet, "/health", "", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestGenericEndpoint(t *testing.T) {
	router := newRouter(testConfig())

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantOK     bool // success
		wantValid  bool
		wantType   string
	}{
		{"valid pan", `{"type":"pan","value":"ABCDE1234F"}`, 200, true, true, "pan"},
		{"invalid pan", `{"type":"pan","value":"12345"}`, 200, true, false, "pan"},
		{"type is case-insensitive", `{"type":" PAN ","value":"ABCDE1234F"}`, 200, true, true, "pan"},
		{"valid email", `{"type":"email","value":"user@example.com"}`, 200, true, true, "email"},
		{"valid mobile", `{"type":"mobile","value":"9876543210"}`, 200, true, true, "mobile"},
		{"valid aadhaar", `{"type":"aadhaar","value":"234567890124"}`, 200, true, true, "aadhaar"},
		{"valid pincode", `{"type":"pincode","value":"380001"}`, 200, true, true, "pincode"},
		{"whitespace value is processed, not a bad request", `{"type":"pan","value":"   "}`, 200, true, false, "pan"},
		{"unsupported type", `{"type":"passport","value":"X"}`, 422, false, false, "unknown"},
		{"missing type", `{"value":"ABCDE1234F"}`, 400, false, false, ""},
		{"missing value", `{"type":"pan"}`, 400, false, false, ""},
		{"empty value", `{"type":"pan","value":""}`, 400, false, false, ""},
		{"malformed json", `{"type":"pan",`, 400, false, false, ""},
		{"empty body", ``, 400, false, false, ""},
		{"json array instead of object", `["pan"]`, 400, false, false, ""},
		{"value is a number", `{"type":"pan","value":123}`, 400, false, false, ""},
		{"value too long", `{"type":"pan","value":"` + strings.Repeat("A", 300) + `"}`, 422, false, false, "pan"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := send(router, http.MethodPost, "/api/v1/validate", tc.body, nil)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.wantStatus, rec.Body.String())
			}
			resp := decode(t, rec)
			if resp.Success != tc.wantOK || resp.Valid != tc.wantValid || resp.Type != tc.wantType {
				t.Errorf("got %+v, want success=%v valid=%v type=%q", resp, tc.wantOK, tc.wantValid, tc.wantType)
			}
			if resp.Message == "" {
				t.Error("message must not be empty")
			}
		})
	}
}

func TestDedicatedEndpoints(t *testing.T) {
	router := newRouter(testConfig())

	tests := []struct {
		path      string
		body      string
		wantValid bool
	}{
		{"/api/v1/validate/pan", `{"value":"ABCDE1234F"}`, true},
		{"/api/v1/validate/pan", `{"value":"bad"}`, false},
		{"/api/v1/validate/email", `{"value":"user@example.com"}`, true},
		{"/api/v1/validate/email", `{"value":"user@invalid"}`, false},
		{"/api/v1/validate/mobile", `{"value":"9876543210"}`, true},
		{"/api/v1/validate/mobile", `{"value":"1234567890"}`, false},
		{"/api/v1/validate/aadhaar", `{"value":"234567890124"}`, true},
		{"/api/v1/validate/aadhaar", `{"value":"234567890123"}`, false},
		{"/api/v1/validate/pincode", `{"value":"380001"}`, true},
		{"/api/v1/validate/pincode", `{"value":"012345"}`, false},
	}

	for _, tc := range tests {
		t.Run(tc.path+" "+tc.body, func(t *testing.T) {
			rec := send(router, http.MethodPost, tc.path, tc.body, nil)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			resp := decode(t, rec)
			if !resp.Success || resp.Valid != tc.wantValid {
				t.Errorf("got %+v, want success=true valid=%v", resp, tc.wantValid)
			}
		})
	}
}

func TestDedicatedAndGenericGiveSameAnswer(t *testing.T) {
	router := newRouter(testConfig())

	dedicated := decode(t, send(router, http.MethodPost, "/api/v1/validate/pan", `{"value":"abcde1234f"}`, nil))
	generic := decode(t, send(router, http.MethodPost, "/api/v1/validate", `{"type":"pan","value":"abcde1234f"}`, nil))

	if dedicated != generic {
		t.Errorf("dedicated %+v differs from generic %+v", dedicated, generic)
	}
}

func TestDedicatedEndpointMalformedAndMissing(t *testing.T) {
	router := newRouter(testConfig())

	for _, body := range []string{`{`, `{}`, `{"value":""}`, `{"value":5}`} {
		rec := send(router, http.MethodPost, "/api/v1/validate/pan", body, nil)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body %q: status = %d, want 400", body, rec.Code)
		}
	}
}

func TestBodyTooLarge(t *testing.T) {
	router := newRouter(testConfig())
	body := `{"value":"` + strings.Repeat("A", 5000) + `"}`

	rec := send(router, http.MethodPost, "/api/v1/validate/pan", body, nil)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want 413", rec.Code)
	}
}

func TestUnknownRouteAndWrongMethod(t *testing.T) {
	router := newRouter(testConfig())

	rec := send(router, http.MethodPost, "/api/v1/nothing", `{}`, nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown route: status = %d, want 404", rec.Code)
	}
	if resp := decode(t, rec); resp.Success {
		t.Error("unknown route should have success=false")
	}

	rec = send(router, http.MethodGet, "/api/v1/validate", "", nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong method: status = %d, want 405", rec.Code)
	}
}

func TestAPIKey(t *testing.T) {
	cfg := testConfig()
	cfg.APIKey = "secret-key"
	router := newRouter(cfg)
	body := `{"value":"ABCDE1234F"}`

	if rec := send(router, http.MethodPost, "/api/v1/validate/pan", body, nil); rec.Code != http.StatusUnauthorized {
		t.Errorf("no key: status = %d, want 401", rec.Code)
	}
	if rec := send(router, http.MethodPost, "/api/v1/validate/pan", body, map[string]string{"X-API-Key": "wrong"}); rec.Code != http.StatusUnauthorized {
		t.Errorf("wrong key: status = %d, want 401", rec.Code)
	}
	if rec := send(router, http.MethodPost, "/api/v1/validate/pan", body, map[string]string{"X-API-Key": "secret-key"}); rec.Code != http.StatusOK {
		t.Errorf("right key: status = %d, want 200", rec.Code)
	}
	if rec := send(router, http.MethodGet, "/health", "", nil); rec.Code != http.StatusOK {
		t.Errorf("health must not need a key: status = %d", rec.Code)
	}
}

func TestRateLimit(t *testing.T) {
	cfg := testConfig()
	cfg.RateLimitPerMinute = 2
	router := newRouter(cfg)
	body := `{"value":"ABCDE1234F"}`

	for i := 1; i <= 2; i++ {
		if rec := send(router, http.MethodPost, "/api/v1/validate/pan", body, nil); rec.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d, want 200", i, rec.Code)
		}
	}
	if rec := send(router, http.MethodPost, "/api/v1/validate/pan", body, nil); rec.Code != http.StatusTooManyRequests {
		t.Errorf("third request: status = %d, want 429", rec.Code)
	}
}

func TestCORS(t *testing.T) {
	cfg := testConfig()
	cfg.CORSAllowedOrigins = []string{"https://app.example.com"}
	router := newRouter(cfg)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/validate", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("preflight status = %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Errorf("allow-origin = %q", got)
	}

	req = httptest.NewRequest(http.MethodOptions, "/api/v1/validate", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("unlisted origin got allow-origin %q", got)
	}
}
