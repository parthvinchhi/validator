// Package routes builds the Gin router. It is the one place where you can see
// every URL the service exposes.
package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/api-validator/config"
	"github.com/parthvinchhi/api-validator/handlers"
	"github.com/parthvinchhi/api-validator/middleware"
	"github.com/parthvinchhi/api-validator/validators"
)

func Setup(cfg config.Config) *gin.Engine {
	router := gin.New() // gin.New() has no default logger, so nothing logs request data.

	// nil = trust no proxy. Client IP is then the direct connection address.
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		log.Fatalf("invalid TRUSTED_PROXIES: %v", err)
	}
	router.HandleMethodNotAllowed = true

	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	if cfg.RateLimitPerMinute > 0 {
		router.Use(middleware.RateLimit(cfg.RateLimitPerMinute))
	}

	// Not protected by the API key, so load balancers can probe it.
	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1")
	api.Use(middleware.BodyLimit(cfg.MaxBodyBytes))
	api.Use(middleware.APIKey(cfg.APIKey))
	{
		api.POST("/validate", handlers.Validate)

		// One dedicated endpoint per registered validator:
		// /validate/aadhaar, /validate/email, /validate/mobile, ...
		for _, validationType := range validators.Types() {
			api.POST("/validate/"+validationType, handlers.ValidateByType(validationType))
		}
	}

	router.NoRoute(handlers.NotFound)
	router.NoMethod(handlers.MethodNotAllowed)

	return router
}
