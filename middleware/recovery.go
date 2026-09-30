package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/parthvinchhi/api-validator/dto"

	"github.com/gin-gonic/gin"
)

// Recovery turns a panic into a JSON 500 response instead of crashing the service.
// Only the panic's type and the stack trace are logged, not the panic value,
// because the value could contain submitted data.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered (%T)\n%s", r, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ValidateResponse{
					Success: false,
					Valid:   false,
					Message: "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
