// Package handlers contains the HTTP handlers. They read JSON, call the
// validation service and write JSON. No validation rules live here.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/api-validator/dto"
	"github.com/parthvinchhi/api-validator/services"
)

// Health is GET /health.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// NotFound handles unknown routes.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, dto.ValidateResponse{Success: false, Message: "Route not found"})
}

// MethodNotAllowed handles known routes called with the wrong HTTP method.
func MethodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, dto.ValidateResponse{Success: false, Message: "Method not allowed"})
}

// Validate is POST /api/v1/validate. The validator comes from the "type" field.
func Validate(c *gin.Context) {
	req, ok := readRequest(c)
	if !ok {
		return
	}
	if req.Type == "" {
		badRequest(c, "Invalid request: type is required")
		return
	}
	runValidation(c, req.Type, req.Value)
}

// ValidateByType returns the handler for POST /api/v1/validate/<type>.
// The type is fixed by the route, so the body only needs "value".
func ValidateByType(validationType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, ok := readRequest(c)
		if !ok {
			return
		}
		runValidation(c, validationType, req.Value)
	}
}

// readRequest decodes the JSON body. On failure it writes the error response
// itself and returns false.
func readRequest(c *gin.Context) (dto.ValidateRequest, bool) {
	var req dto.ValidateRequest

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, dto.ValidateResponse{
				Success: false,
				Message: "Request body too large",
			})
			return req, false
		}
		badRequest(c, "Invalid request")
		return req, false
	}

	if req.Value == "" {
		badRequest(c, "Invalid request: value is required")
		return req, false
	}

	return req, true
}

// runValidation is used by BOTH the generic and the dedicated endpoints,
// so the behaviour (and the response shape) can never drift apart.
func runValidation(c *gin.Context, validationType, value string) {
	result, err := services.Validate(validationType, value)

	switch {
	case errors.Is(err, services.ErrUnsupportedType):
		c.JSON(http.StatusUnprocessableEntity, dto.ValidateResponse{
			Success: false,
			Valid:   false,
			Type:    "unknown",
			Message: "Unsupported validation type",
		})
	case errors.Is(err, services.ErrValueTooLong):
		c.JSON(http.StatusUnprocessableEntity, dto.ValidateResponse{
			Success: false,
			Valid:   false,
			Type:    services.NormalizeType(validationType),
			Message: "Value is too long",
		})
	case err != nil:
		c.JSON(http.StatusInternalServerError, dto.ValidateResponse{
			Success: false,
			Valid:   false,
			Message: "Internal server error",
		})
	default:
		// HTTP 200 even when valid=false: the API worked, the input was just invalid.
		c.JSON(http.StatusOK, dto.ValidateResponse{
			Success: true,
			Valid:   result.Valid,
			Type:    result.Type,
			Message: result.Message,
		})
	}
}

func badRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, dto.ValidateResponse{
		Success: false,
		Valid:   false,
		Message: message,
	})
}
