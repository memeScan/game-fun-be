package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the service is healthy
// @Tags system
// @Produce json
// @Success 200 {string} string "ok"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
