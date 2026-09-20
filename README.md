# validator-api

HTTP API that validates Indian identifiers: name, mobile, email, PAN, Aadhaar, GSTIN, IFSC.
Standard library only, no third-party dependencies (Go 1.22+).

**Format/checksum validation only.** It does not verify that a PAN/Aadhaar exists or belongs to a person.

## Run
    cp .env.example .env && export $(grep -v '^#' .env | xargs)
    make run

    curl -X POST localhost:8080/v1/validate/pan \
      -H 'X-API-Key: projectA-key-change-me' -d '{"value":"ABCPE1234F"}'

## Add a validator
1. Create `internal/validator/yourtype.go` implementing `Type()` and `Validate()`.
2. Add it to `DefaultRegistry()` in `internal/validator/validator.go`.
3. Add a row to `validator_test.go`.

## Notes
- Bodies and query strings are never logged. Run behind HTTPS (reverse proxy / load balancer).
- Rate limiting is in-memory per key; use a shared store if you run multiple replicas.
- Change the module path in `go.mod` to your repo before publishing.
