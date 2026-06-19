package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"SWE-Live/internal/config"
	"SWE-Live/internal/handler"
	"SWE-Live/internal/repository"
	"SWE-Live/internal/routes"
	"SWE-Live/internal/service"
	"SWE-Live/pkg/logger"
)

func main() {
	config, err := config.LoadAndValidate()
	if err != nil {
		log.Fatalf("Critical error while loading the config: %v", err)
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

	dbPool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbPool.Close()

	memberRepository := repository.NewMemberRepository(dbPool)
	memberReadService := service.NewMemberReadService(memberRepository)
	memberReadHandler := handler.NewMemberReadHandler(memberReadService)

	routes.SetupRoutes(router, memberReadHandler)

	server := &http.Server{
		Addr:    config.Port,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error("Server-Error: %v", err)
	}
}
