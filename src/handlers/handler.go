package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(ctx *gin.Context) {
	resp := map[string]string{"hello":"world"}
	ctx.JSON(http.StatusOK, resp)
}

func respondWithError(code int, message string, ctx *gin.Context) {
	resp := map[string]string{"error": message}

	ctx.JSON(code, resp)
	ctx.Abort()
}