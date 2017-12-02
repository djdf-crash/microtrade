package main

import (
	"handlers"
	"middlewares"
	"os"

	"db"

	"validators"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

	binding.Validator.RegisterValidation("emailValidator", validators.EmailValidator)

	router.Static("/", "./public/static")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	router.StaticFile("/index.html", "./public/index.html")

	v1 := router.Group("/api/v1")

	rest := v1.Group("/rest")

	rest.Use(middlewares.AuthMiddleware().MiddlewareFunc())

	{
		rest.POST("/logout", handlers.LogoutHandler)
		rest.POST("/refreshToken", handlers.RefreshTokenHandler)
	}

	authGroup := v1.Group("/auth")

	{
		authGroup.POST("/login", handlers.LoginHandler)
		authGroup.POST("/register", handlers.RegisterHandler)
	}

	router.Run(":" + port)
}
