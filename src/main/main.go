package main

import (
	"handlers"
	"middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	port := "8080"

	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
		if ginMode == "release" {
			port = "80"
		}
	}

	router := gin.Default()
	router.Static("/public", "./public")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")

	api := router.Group("/api")
	v1 := api.Group("/v1")

	rest := v1.Group("/rest")

	rest.Use(middlewares.AuthMiddleware().MiddlewareFunc())

	{
		rest.POST("/logout", handlers.LogoutHandler)
	}

	authGroup := v1.Group("/auth")

	{
		authGroup.POST("/login", middlewares.AuthMiddleware().LoginHandler)
		authGroup.POST("/register", handlers.RegisterHandler)
	}

	router.Run(":" + port) // listen and serve on 0.0.0.0:8080
}
