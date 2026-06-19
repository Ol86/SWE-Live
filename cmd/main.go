package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"SWE-Live/internal/config"
	"SWE-Live/pkg/logger"
)

func main() {
	config, error := config.LoadAndValidate()
	if error != nil {
		//TODO: Logger
	}
	logger.InitLogger(config.Environment)
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	http.ListenAndServe(config.Port, router)
}
