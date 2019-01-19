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
		group.POST("/editDevice", editDevice)
		group.POST("/getDeviceInfo", getDeviceInfo)
		group.POST("/getDeviceLog", getDeviceLog)
		group.POST("/editProfile", editProfile)
		group.POST("/changePassword", changePassword)

		group.GET("/myDevices", getDevicesList)
		group.GET("/getProfile", getProfile)
		group.GET("/checkStatus", checkStatus)
	}
}
