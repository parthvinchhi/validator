// Package server wires routes and middleware together.
package server

import (
	"log/slog"
	"net/http"

	"github.com/parthvinchhi/api-validator/internal/config"
	"github.com/parthvinchhi/api-validator/internal/handler"
	"github.com/parthvinchhi/api-validator/internal/middleware"
	"github.com/parthvinchhi/api-validator/internal/validator"
)

func NewRouter(cfg config.Config, reg *validator.Registry, log *slog.Logger) http.Handler {
	h := &handler.Handler{Reg: reg, MaxBatch: cfg.MaxBatch}

	// Protected routes: auth first, then rate limiting per key.
	api := http.NewServeMux()
	api.HandleFunc("POST /v1/validate/batch", h.Batch)
	api.HandleFunc("POST /v1/validate/{type}", h.Validate)
	api.HandleFunc("GET /v1/validators", h.Types)
	protected := middleware.Chain(api,
		middleware.Auth(cfg.APIKeys),
		middleware.RateLimit(cfg.RateRPS, cfg.RateBurst),
	)

	root := http.NewServeMux()
	root.HandleFunc("GET /healthz", handler.Health)
	root.Handle("/v1/", protected)

	return middleware.Chain(root,
		middleware.Recovery(log),
		middleware.Logging(log),
	)
}
