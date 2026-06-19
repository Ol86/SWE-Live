package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	rest := router.Group("/rest")
	{

	}
}
