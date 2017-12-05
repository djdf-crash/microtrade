package handlers

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *gin.Context) {
	loginReq(ctx)
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
	refreshToken(ctx)
}

func StaticHandler(urlPrefix string, fs static.ServeFileSystem) gin.HandlerFunc {
	return staticFilesGet(urlPrefix, fs)
}
