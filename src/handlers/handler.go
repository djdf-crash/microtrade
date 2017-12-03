package handlers

import (
	"db"
	"middlewares"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	Email           string `json:"email" binding:"required,emailValidator"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type ChangePasswordReq struct {
	Email       string `json:"email" binding:"required,emailValidator"`
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ResponseMessage struct {
	Error Message
}

type Message struct {
	Code    int
	Message string
}

func LoginHandler(ctx *gin.Context) {
	tokenString, expire, err := middlewares.AuthMiddleware.LoginHandler(ctx)
	if err != nil {
		respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	respondWithMessage(http.StatusCreated, "token:"+tokenString+"; expire:"+expire.Format(time.RFC3339), ctx)
}

func ChangePasswordHandler(ctx *gin.Context) {
	changePasswordUser(ctx)
}
func changePasswordUser(ctx *gin.Context) {
	var changePassword ChangePasswordReq

	if err := ctx.ShouldBindJSON(&changePassword); err == nil {

		user := db.FindUserByName(changePassword.Email)

		if changePassword.NewPassword == "" {
			respondWithMessage(http.StatusBadRequest, "New password is empty", ctx)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePassword.Password)); err != nil {
			respondWithMessage(http.StatusBadRequest, "Password is not correctly", ctx)
			return
		}

		newHash, _ := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), bcrypt.DefaultCost)

		user.Password = string(newHash)

		db.UpdateUser(&user)

	} else {
		respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
		return
	}
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
	var user db.User

	if err := ctx.ShouldBindJSON(&userRegister); err == nil {
		if userRegister.Password != userRegister.ConfirmPassword {
			respondWithMessage(http.StatusBadRequest, "Password and confirm password not equals", ctx)
			return
		}

		if !db.CheckUserByEmail(userRegister.Email) {

			hash, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
			if err != nil {
				return
			}

			user.Email = userRegister.Email
			user.Password = string(hash)

			if err := db.AddUser(&user); err != nil {
				respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
			} else {

				jwtToken, expire, _ := middlewares.AuthMiddleware.TokenGenerator(user.Email)

				respondWithMessage(http.StatusCreated, "token:"+jwtToken+"; expire:"+expire.Format(time.RFC3339), ctx)
			}
		} else {
			respondWithMessage(http.StatusBadRequest, "User name is exist", ctx)
			return
		}

	} else {
		respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
		return
	}
}

func RefreshTokenHandler(ctx *gin.Context) {
	middlewares.AuthMiddleware.RefreshHandler(ctx)
}

func respondWithMessage(code int, message string, ctx *gin.Context) {
	response := ResponseMessage{
		Message{
			code,
			message,
		},
	}

	ctx.JSON(code, &response)
	ctx.Abort()
}
