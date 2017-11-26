package handlers

import (
	"github.com/gin-gonic/gin"
)

func LoginHandler(ctx *gin.Context) {
	resp := map[string]string{"hello":"world"}
	ctx.JSON(200, resp)
}

func respondWithError(code int, message string, ctx *gin.Context) {
	resp := map[string]string{"error": message}

	ctx.JSON(code, resp)
	ctx.Abort()
}