package main

import (
	"SWE-Live/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger("develoment")
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}
