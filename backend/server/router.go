package server

import (
	"github.com/ip-05/quizzus/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	auth := new(controllers.AuthController)

	authGroup := router.Group("auth")

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	authGroup.GET("/google", auth.GoogleLogin)
	authGroup.GET("/google/callback", auth.GoogleCallback)

	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/logout", auth.Logout)
	authGroup.GET("/me", auth.Me)

	return router

}
