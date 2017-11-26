package main

import (
	"github.com/gin-gonic/gin"
	"middlewares"
	"handlers"
)

func main() {
	router := gin.Default()
	router.Static("/public", "./public")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")

	v1 := router.Group("/v1")

	v1.Use(middlewares.AuthMiddleware)

	{
		v1.GET("/login",handlers.LoginHandler)
		v1.POST("/logout",)
		v1.POST("/register",)
	}

	router.Run() // listen and serve on 0.0.0.0:8080
}
