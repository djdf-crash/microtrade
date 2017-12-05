package main

import (
	"db"
	"handlers"
	"middlewares"

	"validators"

	"config"
	"log"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {

	err := config.InitConfig("./config.json")
	if err != nil {
		log.Panic(err.Error())
	}

	db.InitDB()
	defer db.DB.Close()

	port := config.AppConfig.Port

	gin.SetMode(config.AppConfig.ModeStart)

	router := gin.Default()

	binding.Validator.RegisterValidation("emailValidator", validators.EmailValidator)

	router.Use(handlers.StaticHandler("/", static.LocalFile("./public", true)))

	router.GET("/token/:token", handlers.ConfirmPasswordReqHandler)

	v1 := router.Group("/api/v1")

	rest := v1.Group("/rest")

	rest.Use(middlewares.AuthMiddleware.MiddlewareFunc())
	{
		rest.POST("/password/change", handlers.ChangePasswordHandler)
		rest.POST("/logout", handlers.LogoutHandler)
		rest.POST("/refreshToken", handlers.RefreshTokenHandler)
	}

	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", handlers.LoginHandler)
		authGroup.POST("/register", handlers.RegisterHandler)
	}

	passwordGroup := authGroup.Group("/password")
	{
		passwordGroup.POST("/reset", handlers.ResetPasswordReqHandler)
	}

	router.Run(port)
}
