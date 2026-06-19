package routes

import (
	"SWE-Live/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	memberReadHandler *handler.MemberReadHandler,
	memberWriteHadler *handler.MemberWriteHandler,
) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	rest := router.Group("/rest")
	{
		memberReadHandler.RegisterRoutes(rest)
		memberWriteHadler.RegisterRoutes(rest)
	}
}
