package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS allows browsers on the listed origins to call the API.
// With an empty list no CORS headers are sent, which is the right default for
// server-to-server use (CORS only matters for calls made from browser JavaScript).
// Use "*" in the list to allow every origin.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowOrigin := matchOrigin(origin, allowedOrigins)

		if allowOrigin != "" {
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", allowOrigin)
			h.Set("Vary", "Origin")
			h.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")
			h.Set("Access-Control-Max-Age", "600")
		}

		// Answer browser preflight requests here so they never reach the API-key check.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func matchOrigin(origin string, allowed []string) string {
	for _, a := range allowed {
		if a == "*" {
			return "*"
		}
		if origin != "" && a == origin {
			return origin
		}
	}
	return ""
}
