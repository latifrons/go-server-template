package controllers

import (
	"github.com/avct/uasurfer"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (rpc *RpcController) DebugUA(c *gin.Context) {
	ua := c.GetHeader("User-Agent")
	if ua == "" {
		return
	}
	uag := rpc.GetUA(ua)

	Response(c, http.StatusOK, 0, "", DebugUAResponse{
		BrowserName: uag.Browser.Name.String(),
		DeviceType:  uag.DeviceType.String(),
		OsName:      uag.OS.Name.String(),
		OsPlatform:  uag.OS.Platform.String(),
		DbPlatform:  string(rpc.GetPlatform(c)),
	})
}

func (rpc *RpcController) GetPlatform(c *gin.Context) Platform {
	ua := c.GetHeader("User-Agent")
	if ua == "" {
		return Platform_PC
	}
	uag := rpc.GetUA(ua)
	switch uag.DeviceType {
	case uasurfer.DevicePhone:
		return Platform_M
	default:
		switch uag.OS.Name {
		case uasurfer.OSAndroid:
			return Platform_M
		case uasurfer.OSiOS:
			return Platform_M
		}
	}
	return Platform_PC
}

func (rpc *RpcController) GetUA(uas string) *uasurfer.UserAgent {
	ua := uasurfer.Parse(uas)
	return ua
}

func (rpc *RpcController) Panic(c *gin.Context) {
	time.Sleep(10 * time.Second) // 模拟接口超时
	c.JSON(200, gin.H{"message": "this is panic"})
}
