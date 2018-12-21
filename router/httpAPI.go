package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rod41732/cu-smart-farm-backend/api/middleware"
	"github.com/rod41732/cu-smart-farm-backend/api/user"
	"github.com/rod41732/cu-smart-farm-backend/common"
)

func SetUpHttpAPI(r *gin.Engine) {

	common.ShouldPrintDebug = true
	httpAPI := r.Group("api/v1")
	// define short name
	userAuth := middleware.UserAuth.MiddlewareFunc()
	// ownerCheck := middleware.OwnerCheck
	{
		httpAPI.POST("/login", middleware.UserAuth.LoginHandler)
		httpAPI.POST("/register", user.Register)
		httpAPI.POST("/ws", userAuth, user.WebSocket)
	}
}
