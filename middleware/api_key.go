package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/api-validator/dto"
)

// APIKey requires the header "X-API-Key" to match the configured key.
// If key is empty, the check is disabled.
//
// This is service-to-service protection so only your own applications can call
// the service. It is NOT user authentication.
func APIKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key == "" {
			c.Next()
			return
		}

		given := c.GetHeader("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(given), []byte(key)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ValidateResponse{
				Success: false,
				Valid:   false,
				Message: "Unauthorized",
			})
			return
		}

		c.Next()
	}
}
