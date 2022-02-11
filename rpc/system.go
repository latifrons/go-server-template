package rpc

import (
	"github.com/gin-gonic/gin"
)

// Health does heartbeat checking.
// @Summary Health check
// @Description Get to check health
// @Router /health [get]
// @Tags System
// @Success 200 {string} string "ok"
func (rpc *RpcController) Health(c *gin.Context) {
	rpc.ResponseOK(c, "okk")
}
