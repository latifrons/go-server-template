package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Health check
// @Description Get to check health
// @Router /health [get]
// @Tags System
// @Success 200 {string} string "ok"
func (rpc *RpcController) Health(c *gin.Context) {
	cors(c)
	Response(c, http.StatusOK, 0, "", "okk")
}
