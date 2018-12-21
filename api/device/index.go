package device

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
)

// DeviceControlAPI : sets up device control API
func DeviceControlAPI(r *gin.RouterGroup) {
	deviceAPI := r.Group("/device")
	userAuth := middleware.UserAuth.MiddlewareFunc()
	ownerCheck := middleware.OwnerCheck
	deviceAPI.GET("/set", userAuth, ownerCheck, setState)
	deviceAPI.GET("/fetch", userAuth, ownerCheck, fetchInfo)

	deviceAPI.GET("/greeting", greeting)
}
