package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Ping godoc
// @Summary      Ping pong
// @Description  Ping-Pong
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Router       /api/v1/ping [get]
func Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong!",
	})
}
