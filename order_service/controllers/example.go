package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(context *gin.Context) {
	context.String(http.StatusOK, "pong")
}
