package handlers

import (
	"net/http"

	"db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ResponseMessage struct {
	Code    int
	Message string
}

func LogoutHandler(ctx *gin.Context) {
	resp := map[string]string{"Logout": "OK"}
	ctx.JSON(http.StatusOK, resp)
}

func RegisterHandler(ctx *gin.Context) {
	registerUser(ctx)
}
func registerUser(ctx *gin.Context) {
	var userRegister Register
	var user db.Users

	if err := ctx.ShouldBindJSON(&userRegister); err == nil {
		if userRegister.Password != userRegister.ConfirmPassword {
			respondWithMessage(http.StatusInternalServerError, "Password and confirm password not equals", ctx)
			return
		}

		if !db.CheckUserByUserName(userRegister.Username) {

			hash, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
			if err != nil {
				return
			}

			user.Username = userRegister.Username
			user.Password = string(hash)
			db.AddUser(&user)
			respondWithMessage(http.StatusCreated, "User registered", ctx)
		} else {
			respondWithMessage(http.StatusInternalServerError, "User name is exist", ctx)
			return
		}

	} else {
		respondWithMessage(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
}

func respondWithMessage(code int, message string, ctx *gin.Context) {
	response := ResponseMessage{
		code,
		message,
	}

	ctx.JSON(code, &response)
	ctx.Abort()
}
