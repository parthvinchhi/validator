// Package handler contains the HTTP handlers.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/parthvinchhi/api-validator/internal/model"
	"github.com/parthvinchhi/api-validator/internal/validator"
)

const maxBody = 1 << 20 // 1 MB

type Handler struct {
	Reg      *validator.Registry
	MaxBatch int
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, model.ErrorResponse{Error: msg, Code: code})
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeErr(w, http.StatusBadRequest, "BAD_REQUEST", "malformed JSON body")
		return false
	}
	return true
}

// Validate handles POST /v1/validate/{type}.
// Invalid input is a normal 200 with valid=false; only bad requests get 4xx.
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	var req model.ValidateRequest
	if !decode(w, r, &req) {
		return
	}
	res, err := h.Reg.Validate(r.PathValue("type"), req.Value)
	if errors.Is(err, validator.ErrUnknownType) {
		writeErr(w, http.StatusNotFound, validator.CodeUnknownType, "unknown validator type")
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// Batch handles POST /v1/validate/batch.
func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	var req model.BatchRequest
	if !decode(w, r, &req) {
		return
	}
	if len(req.Items) == 0 {
		writeErr(w, http.StatusBadRequest, "BAD_REQUEST", "items must not be empty")
		return
	}
	if len(req.Items) > h.MaxBatch {
		writeErr(w, http.StatusBadRequest, "BATCH_TOO_LARGE", "too many items in batch")
		return
	}
	out := make([]validator.Result, len(req.Items))
	for i, it := range req.Items {
		res, err := h.Reg.Validate(it.Type, it.Value)
		if err != nil {
			res = validator.Result{Type: it.Type, Code: validator.CodeUnknownType, Message: "unknown validator type"}
		}
		out[i] = res
	}
	writeJSON(w, http.StatusOK, model.BatchResponse{Results: out})
}

// Types handles GET /v1/validators.
func (h *Handler) Types(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, model.TypesResponse{Types: h.Reg.Types()})
}

// Health handles GET /healthz.
func Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
