package handlers

import (
	"middlewares"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *gin.Context) {
	tokenString, expire, err := middlewares.AuthMiddleware.LoginHandler(ctx)
	if err != nil {
		respondWithMessage(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	respondWithMessage(http.StatusCreated, "token:"+tokenString+"; expire:"+expire.Format(time.RFC3339), ctx)
}

func ResetPasswordReqHandler(ctx *gin.Context) {
	resetPasswordReq(ctx)
}

func ConfirmPasswordReqHandler(ctx *gin.Context) {
	confirmPasswordReq(ctx)
}

func ChangePasswordHandler(ctx *gin.Context) {
	changePasswordUser(ctx)
}

func LogoutHandler(ctx *gin.Context) {
	resp := map[string]string{"Logout": "OK"}
	ctx.JSON(http.StatusOK, resp)
}

func RegisterHandler(ctx *gin.Context) {
	registerUser(ctx)
}

func RefreshTokenHandler(ctx *gin.Context) {
	middlewares.AuthMiddleware.RefreshHandler(ctx)
}
