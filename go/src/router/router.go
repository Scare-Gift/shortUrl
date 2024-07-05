package router

import (
	"github.com/gin-gonic/gin"
	"shorturl/go/src/controller"
	"shorturl/go/src/middleware"
)

func Router() *gin.Engine {
	router := gin.Default()

	urlApi := router.Group("/api")
	{
		urlApi.GET("/:shortUrl", controller.ResolveShortUrl)
		urlApi.Use(middleware.JWTAuth()).POST("/short", controller.CreateShortUrl)
	}
	userApi := router.Group("/user")
	{
		userApi.POST("/login", controller.Login)
		userApi.Use(middleware.JWTAuth()).POST("/logout", controller.Logout)
	}
	return router
}
