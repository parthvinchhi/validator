# Validation API

A small, stateless HTTP service (Go + Gin) that validates common user inputs:
mobile number, email, PAN, Aadhaar and PIN code. Your other applications call it
over HTTP and get a JSON answer.

```text
My Application --HTTP--> Validation API --> JSON result --> My Application
```

## 1. What it does

Receive input, validate its **format/rules locally**, return a result, discard the input.

## 2. What it deliberately does NOT do

No users, signup, login, JWT, sessions, profiles, password storage, validation
history, database or any other persistent storage. Your applications handle their
own users. The optional API key (section 16) protects the service between your
own apps; it is not user authentication.

It also does **not verify** anything with an outside authority (see "Validation vs
verification" below).

## 3. Stateless architecture

```text
INPUT -> VALIDATOR -> RESULT -> RESPONSE -> INPUT DISCARDED
```

There is no database package and no file writing. The only in-memory state is the
rate limiter's per-IP request counters, which hold IP addresses and counts (never
submitted values) and disappear on restart.

### Validation vs verification

| | Validation (this project) | Verification (future) |
|---|---|---|
| Question | "Does this look like a PAN?" | "Does this PAN really exist and belong to this person?" |
| How | Local rules: length, characters, pattern, checksum | Call an external authority/provider, or send an OTP |
| Example | `ABCDE1234F` matches `AAAAA9999A` | Income-tax/KYC provider lookup, mobile OTP |

A response of `"valid": true` means "correct format", **not** "this is a real, active
identity".

## 4. Folder structure

```text
validation-api/
├── main.go                  start-up: load config, build router, run server
├── go.mod
├── .env.example             list of environment variables
├── config/config.go         reads environment variables
├── dto/validation.go        JSON request/response shapes
├── validators/              ALL validation rules (one file per validator)
│   ├── validator.go         Result type + the registry (type name -> function)
│   ├── mobile.go  email.go  pan.go  aadhaar.go  pincode.go
│   └── *_test.go            unit tests, next to the code (Go convention)
├── services/                picks a validator by type, enforces max input length
├── handlers/                HTTP only: read JSON, call service, write JSON
├── routes/routes.go         every URL in one place
├── middleware/              logger, recovery, CORS, rate limit, API key, body limit
├── client/client.go         optional helper for other Go projects
├── examples/http_client/    copy-paste net/http example
└── tests/api_test.go        HTTP-level tests for the endpoints
```

Why each folder exists: `validators` is the only place rules live, so adding a
rule never touches HTTP code. `services` exists so handlers do not need to know how
validators are looked up (it is also the natural place for future verification
calls). `handlers`, `routes`, `middleware`, `dto` and `config` each do one HTTP or
setup job. Unit tests sit next to validators because Go tests in the same package
can be run and found easily; `tests/` is only for end-to-end API tests.

## 5. Validators and normalization

Nothing is changed silently: every change below is intentional and happens only
inside the validator, on a copy. The submitted value is never echoed back.

| Type | Normalization | Rules |
|---|---|---|
| `mobile` | trim; remove spaces and `-`; strip `+91`, `91` (12 digits) or `0` (11 digits) prefix | 10 digits, first digit 6-9 |
| `email` | trim only (case unchanged) | valid syntax (`net/mail`), plain address only, length limits, domain with 2+ valid labels, TLD at least 2 chars and not all digits. Format only; ASCII domains only |
| `pan` | trim; **convert to uppercase** | `AAAAA9999A` (5 letters, 4 digits, 1 letter). 4th-letter holder type not enforced. No public checksum exists |
| `aadhaar` | trim; remove spaces and `-` | 12 digits, first digit 2-9, Verhoeff checksum |
| `pincode` | trim only | 6 digits, first digit 1-9 |

The `type` field is case-insensitive (`PAN` and `pan` both work). Maximum `value`
length is 256 bytes.

## 6. API endpoints

| Method | Path | Purpose |
|---|---|---|
| GET | `/health` | liveness check |
| POST | `/api/v1/validate` | generic, validator chosen by `type` |
| POST | `/api/v1/validate/mobile` | dedicated |
| POST | `/api/v1/validate/email` | dedicated |
| POST | `/api/v1/validate/pan` | dedicated |
| POST | `/api/v1/validate/aadhaar` | dedicated |
| POST | `/api/v1/validate/pincode` | dedicated |

Dedicated routes are created automatically from the validator registry, and they
call the same code path as the generic route (`runValidation` in the handler), so
logic is never duplicated.

## 7. Request formats

Dedicated: `{"value": "ABCDE1234F"}`
Generic: `{"type": "pan", "value": "ABCDE1234F"}`

## 8. Response formats and status codes

Valid: `{"success":true,"valid":true,"type":"pan","message":"PAN format is valid"}`
Invalid: `{"success":true,"valid":false,"type":"pan","message":"Invalid PAN format"}`
Unsupported type: `{"success":false,"valid":false,"type":"unknown","message":"Unsupported validation type"}`
Malformed request: `{"success":false,"valid":false,"message":"Invalid request"}`

| Status | When |
|---|---|
| 200 | request processed, **including** `valid:false` |
| 400 | invalid JSON, wrong field types, missing/empty `value`, missing `type` (generic) |
| 401 | API key enabled and missing/wrong |
| 404 | unknown route |
| 405 | wrong HTTP method |
| 413 | request body larger than 4 KB |
| 422 | well-formed JSON the API cannot process: unsupported `type`, or `value` longer than 256 bytes |
| 429 | rate limit exceeded |
| 500 | unexpected server error (never for merely invalid input) |

A `value` that is only whitespace is processed normally (200, `valid:false`). A
missing or `""` value is a malformed request (400).

## 9. Running locally

You need Go 1.22 or newer.

### Windows (PowerShell)

```powershell
cd validation-api
go mod tidy          # downloads Gin and creates go.sum
go test ./...
go run main.go
```

Optional settings (same PowerShell window, before `go run`):

```powershell
$env:PORT = "8080"
$env:API_KEY = "my-secret-key"
```

### Linux / macOS

```bash
cd validation-api
go mod tidy
go test ./...
go run main.go
# optional: PORT=8080 API_KEY=my-secret-key go run main.go
```

Test with curl or Postman: `POST http://localhost:8080/api/v1/validate/pan`, header
`Content-Type: application/json`, raw JSON body `{"value":"ABCDE1234F"}`.

## 10. Testing

```bash
go test ./...          # everything
go test ./validators   # validator rules only
go vet ./...
```

## 11. Curl examples

Linux/macOS/Git Bash (single quotes):

```bash
curl -X POST http://localhost:8080/api/v1/validate/pan -H "Content-Type: application/json" -d '{"value":"ABCDE1234F"}'
curl -X POST http://localhost:8080/api/v1/validate/email -H "Content-Type: application/json" -d '{"value":"user@example.com"}'
curl -X POST http://localhost:8080/api/v1/validate/mobile -H "Content-Type: application/json" -d '{"value":"9876543210"}'
curl -X POST http://localhost:8080/api/v1/validate/aadhaar -H "Content-Type: application/json" -d '{"value":"234567890124"}'
curl -X POST http://localhost:8080/api/v1/validate/pincode -H "Content-Type: application/json" -d '{"value":"380001"}'
curl -X POST http://localhost:8080/api/v1/validate -H "Content-Type: application/json" -d '{"type":"pan","value":"ABCDE1234F"}'
```

Windows CMD uses escaped double quotes: `-d "{\"value\":\"ABCDE1234F\"}"`. In PowerShell
use `curl.exe` (not `curl`, which is an alias) with `-d '{\"value\":\"ABCDE1234F\"}'`.

If an API key is set, add `-H "X-API-Key: my-secret-key"`.

Note: `123456789012` is **invalid** Aadhaar (starts with 1 and fails the checksum), so
use `234567890124` for a valid example.

## 12. Go client integration

Plain `net/http`: see `examples/http_client/main.go` (run with `go run ./examples/http_client`).

Helper package: `client/client.go` is worth using because it removes about 30
lines of repeated request code from each of your 5 projects. The simplest way to
share it is to **copy the file** into each project (it uses only the standard
library). Importing it as a Go module also works once you rename `module validation-api`
in `go.mod` to a real path such as `github.com/you/validation-api`.

```go
v := client.New("http://localhost:8080", os.Getenv("VALIDATION_API_KEY"))

result, err := v.ValidatePAN("ABCDE1234F")
if err != nil {
    // the call failed: service down, bad API key, unsupported type...
    return err
}
if !result.Valid {
    // the value is invalid: show result.Message to the user
}
```

Decide what your app does when the service is unreachable (fail closed or fall
back). Always keep the client timeout (it defaults to 5 seconds).

## 13. Adding a validator (example: IFSC)

1. Create `validators/ifsc.go`:

```go
package validators

import (
	"regexp"
	"strings"
)

var ifscPattern = regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)

// ValidateIFSC: trims and uppercases; 4 letters, the digit 0, then 6 letters/digits.
func ValidateIFSC(value string) Result {
	const t = "ifsc"
	ifsc := strings.ToUpper(strings.TrimSpace(value))
	if ifsc == "" {
		return invalid(t, "IFSC is required")
	}
	if !ifscPattern.MatchString(ifsc) {
		return invalid(t, "Invalid IFSC format")
	}
	return valid(t, "IFSC format is valid")
}
```

2. Register it: add one line to the `registry` map in `validators/validator.go`:
   `"ifsc": ValidateIFSC,`
3. Add `validators/ifsc_test.go` (copy `pincode_test.go` and change the cases).
4. Run `go test ./...`. That is all: `POST /api/v1/validate` with `"type":"ifsc"` and
   `POST /api/v1/validate/ifsc` now work. No handler, route or other validator changes.
5. Optional: add a `Validate...` convenience method to `client/client.go`.

The same steps work for GSTIN, passport, vehicle number, URL, date of birth, username, etc.
Anything that is a single text value fits this design. (Strong password checks
work technically, but sending real passwords to another service is a risky
habit; prefer checking password strength inside your own app.)

## 14. Adding file validation later (images, PDFs)

Files do not fit `func(value string) Result` (they are bytes, a filename, and a
size, sent as `multipart/form-data`). Do **not** force them into the registry.
Add a parallel, separate path:

1. Create a package `filevalidators/` with functions such as
   `func ValidatePDF(filename string, data []byte) Result`.
2. Create `handlers/file_handler.go` that reads `c.FormFile("file")`, opens it,
   reads it with a size cap, and calls the file validator.
3. In `routes.go`, add a separate group with its **own larger body limit**
   (`BodyLimit` already takes the size as a parameter):

```go
files := router.Group("/api/v1/validate")
files.Use(middleware.BodyLimit(10 << 20), middleware.APIKey(cfg.APIKey))
files.POST("/pdf", handlers.ValidatePDF)
files.POST("/image", handlers.ValidateImage)
```

   Register these explicit routes before or alongside the loop of text routes;
   the names (`pdf`, `image`, `file`) must not collide with text validator names.
4. Detect the MIME type from the file contents (`http.DetectContentType`), not from
   the filename or client-provided header. Keep the file in memory or a temp file
   that is deleted immediately, consistent with the stateless rule.
5. Keep the same response shape so clients see the same JSON.

## 15. Adding external verification later

Keep it in a separate package so "format is OK" and "this really exists" never mix:

```text
validators/     local format validation   (exists today)
verification/   external provider calls   (future)
```

```go
package verification

type Result struct {
	Verified bool
	Message  string
}

type PANVerificationProvider interface {
	VerifyPAN(pan string) (Result, error)
}
```

Suggested flow: a new route such as `POST /api/v1/verify/pan` runs the local
validator first (no point calling a paid provider with a malformed PAN), then calls
the provider. Provider credentials come from environment variables, calls need
timeouts, provider failures map to 502/503 (not 500 for bad input), and the
submitted value must still never be logged. Interfaces are not used today because
there is only one implementation of everything; add the interface when the first
provider arrives, and a second provider gives you a reason to keep it.
Regulated identity checks such as Aadhaar need an authorized provider and legal
review.

## 16. Security

Implemented:

- Request body limit (4 KB) and 256-byte value limit.
- Panic recovery returning JSON 500 (logs panic type and stack, not the value).
- Log line contains only method, path, status, duration, never the body or query string.
  Validators and handlers contain no logging at all.
- CORS: off by default (server-to-server). Set `CORS_ALLOWED_ORIGINS` for browsers.
- Per-IP rate limit, in memory (`RATE_LIMIT_PER_MINUTE`, default 300, 0 = off).
  All your apps calling from one server share one IP, so keep the limit generous.
- Optional API key via `X-API-Key` (constant-time comparison), `/health` excluded.
  Service-to-service protection only, not user authentication.
- HTTP server timeouts.

Not implemented: TLS (terminate it at a reverse proxy), per-key rate limits, key rotation,
any legal or privacy certification. This project makes no compliance claims (for example
DPDP Act). Aadhaar numbers in particular are sensitive and legally regulated in India,
so get proper advice before handling them in production.

Environment variables: `PORT`, `ENVIRONMENT` (`production` enables Gin release
mode), `API_KEY`, `CORS_ALLOWED_ORIGINS`, `RATE_LIMIT_PER_MINUTE`, `TRUSTED_PROXIES`.

## 17. Privacy

The application is stateless: it receives input, validates it, returns the result
and discards the input. It has no database and keeps no validation history.

That covers **the application only**. Reverse proxies (nginx, Caddy, cloud load
balancers), hosting providers, access logs, APM/monitoring tools, WAFs and server
or container logs can still record request data, including bodies if
configured that way, unless you separately configure them not to. Check each
layer you deploy. Also remember that callers send values over the network, so use HTTPS.

## 18. Production deployment

- Run behind a reverse proxy that terminates HTTPS; do not expose plain HTTP publicly.
- Prefer a private network/VPC so only your apps can reach the service.
- Set `ENVIRONMENT=production` and a long random `API_KEY` (for example
  `openssl rand -hex 32`); store it in your secret manager, not in git.
- Set `TRUSTED_PROXIES` to your proxy's IP so rate limiting sees real client IPs.
- Turn off request-body logging in the proxy; keep access logs short-lived.
- Build a binary: `go build -o validation-api .` and run it with systemd or in a
  container (a scratch/distroless image is enough; no files or database needed).
- Run 2 or more instances if you need availability. They do not need to share
  anything, but each rate-limits separately.
- Use `GET /health` for load balancer checks.
