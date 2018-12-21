package user

import (
	"../middleware"
	"github.com/gin-gonic/gin"
)

func SetUpUserAPI(r *gin.RouterGroup) {
	userAPI := r.Group("/user")
	// define short name
	// userAuth := middleware.UserAuth.MiddlewareFunc()
	// ownerCheck := middleware.OwnerCheck
	{
		userAPI.POST("/login", middleware.UserAuth.LoginHandler)
		userAPI.POST("/register", register)
		// userAPI.POST("/ws", register)
		// userAPI.POST("addDevice", userAuth, addDevice)
		// userAPI.POST("removeDevice", userAuth, ownerCheck, removeDevice)
	}
}
