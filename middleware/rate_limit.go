package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/parthvinchhi/api-validator/dto"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count       int
	windowStart time.Time
}

// RateLimit allows at most perMinute requests per client IP per minute
// (a simple fixed window). The counters live only in memory and are lost on
// restart; nothing is written anywhere.
//
// Behind a reverse proxy, set TRUSTED_PROXIES so the real client IP is used,
// otherwise every request looks like it comes from the proxy.
func RateLimit(perMinute int) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Remove old entries once a minute so the map cannot grow forever.
	go func() {
		for range time.Tick(time.Minute) {
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.windowStart) >= time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, found := visitors[ip]
		if !found || time.Since(v.windowStart) >= time.Minute {
			v = &visitor{windowStart: time.Now()}
			visitors[ip] = v
		}
		v.count++
		allowed := v.count <= perMinute
		mu.Unlock()

		if !allowed {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ValidateResponse{
				Success: false,
				Valid:   false,
				Message: "Too many requests",
			})
			return
		}

		c.Next()
	}
}
