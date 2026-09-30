// Package middleware contains small, single-purpose Gin middleware.
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger writes one line per request: method, path, status, duration.
// It never logs the request body or the query string, so submitted values
// (PAN, Aadhaar, mobile, email...) never reach the application log.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s -> %d (%s)", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}
