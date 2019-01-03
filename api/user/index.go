package user

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
)

// UserAPI sets up user API for specified router
func UserAPI(r *gin.RouterGroup) {
	group := r.Group("/user")
	group.Use(middleware.UserAuth.MiddlewareFunc())
	{
		group.POST("/addDevice", addDevice)
		group.POST("/removeDevice", removeDevice)
		group.POST("/setDevice", setDevice)
		group.POST("/renameDevice", renameDevice)
		group.GET("/myDevices", getDevicesList)
		group.POST("/getDeviceInfo", getDeviceInfo)
		group.POST("/getDeviceLog", getDeviceLog)
	}
}
