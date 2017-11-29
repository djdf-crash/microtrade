package main

import (
	"handlers"
	"middlewares"
	"os"

	"db"

	"github.com/gin-gonic/gin"
)

func main() {

	port := "8080"

	db.InitDB()
	defer db.DB.Close()

	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
		if ginMode == "release" {
			port = "80"
		}
	}

	router := gin.Default()

	api := router.Group("/api")
	v1 := api.Group("/v1")

	rest := v1.Group("/rest")
	rest.Static("/public", "./public")
	rest.StaticFile("/favicon.ico", "./public/favicon.ico")

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
