package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"SWE-Live/internal/config"
	"SWE-Live/internal/routes"
	"SWE-Live/pkg/logger"
)

func main() {
	config, error := config.LoadAndValidate()
	if error != nil {
		log.Fatalf("Critical error while loading the config: %v", error)
	}
	appLogger := logger.InitLogger(config.Environment)

	appLogger.Info("Logger initiallized succresfully",
		"env", config.Environment,
		"port", config.Port,
	)

	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	routes.SetupRoutes(router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router,
	}

	if error := server.ListenAndServe(); error != nil && error != http.ErrServerClosed {
		appLogger.Error("Server-Error: %v", error)
	}
}
