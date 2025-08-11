package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary      Ping
// @Description  Do ping
// @Tags         ping
// @Success      200  {object}  map[string]string
// @Router       /api/v1/ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
