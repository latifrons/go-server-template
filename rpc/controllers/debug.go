package controllers

import (
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/avct/uasurfer"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (rpc *RpcController) GetPlatform(c *gin.Context) dbgorm.Platform {
	ua := c.GetHeader("User-Agent")
	if ua == "" {
		return dbgorm.Platform_PC
	}
	uag := rpc.GetUA(ua)
	switch uag.DeviceType {
	case uasurfer.DevicePhone:
		return dbgorm.Platform_M
	default:
		switch uag.OS.Name {
		case uasurfer.OSAndroid:
			return dbgorm.Platform_M
		case uasurfer.OSiOS:
			return dbgorm.Platform_M
		}
	}
	return dbgorm.Platform_PC
}

func (rpc *RpcController) GetUA(uas string) *uasurfer.UserAgent {
	ua := uasurfer.Parse(uas)
	return ua
}
