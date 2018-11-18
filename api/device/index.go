package device

import (
	"../middleware"
	"github.com/gin-gonic/gin"
)

// DeviceControlAPI : sets up device control API
func DeviceControlAPI(r *gin.RouterGroup) {
	deviceAPI := r.Group("/device")
	userAuth := middleware.UserAuth.MiddlewareFunc()

	deviceAPI.GET("/set", userAuth, middleware.OwnerCheck, setState)
	deviceAPI.GET("/greeting", greeting)
}
