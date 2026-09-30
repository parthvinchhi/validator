package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit rejects request bodies larger than maxBytes. Because it takes the
// limit as a parameter, a future file-upload route group can use a bigger limit.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
