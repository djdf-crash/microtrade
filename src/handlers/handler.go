package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(ctx *gin.Context) {
	resp := map[string]string{"Logout": "OK"}
	ctx.JSON(http.StatusOK, resp)
}

func RegisterHandler(ctx *gin.Context) {
	resp := map[string]string{"Register": "OK"}
	ctx.JSON(http.StatusOK, resp)
}

func respondWithError(code int, message string, ctx *gin.Context) {
	resp := map[string]string{"error": message}

	ctx.JSON(code, resp)
	ctx.Abort()
}
