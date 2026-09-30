package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/api-validator/config"
	"github.com/parthvinchhi/api-validator/routes"
)

func main() {
	cfg := config.Load()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		if cfg.APIKey == "" {
			log.Println("WARNING: API_KEY is not set; anyone who can reach this service can call it")
		}
	}

	router := routes.Setup(cfg)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Validation API listening on port %s (environment: %s)", cfg.Port, cfg.Environment)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
